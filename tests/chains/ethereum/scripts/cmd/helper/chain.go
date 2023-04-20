package helper

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/client"
	ibccommitment "github.com/hyperledger-labs/yui-ibc-solidity/pkg/contract/ibccommitmenttesthelper"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/contract/ibchandler"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/contract/ics20bank"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/contract/ics20transferbank"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/contract/simpletoken"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/relay/ethereum"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/wallet"
	"github.com/hyperledger-labs/yui-relayer/core"
)

type ChainConfig struct {
	Chain  ethereum.ChainConfig `json:"chain"`
	Prover ProverConfig         `json:"prover"`
}

type ProverConfig struct {
	Type string `json:"@type"`
}

type Chain struct {
	chainID        int64
	client         *client.ETHClient
	mnemonicPhrase string
	keys           map[uint32]*ecdsa.PrivateKey

	ChainConfig ChainConfig

	// Core Modules
	IBCHandler    ibchandler.Ibchandler
	IBCCommitment ibccommitment.Ibccommitmenttesthelper

	// App Modules
	SimpleToken   simpletoken.Simpletoken
	ICS20Transfer ics20transferbank.Ics20transferbank
	ICS20Bank     ics20bank.Ics20bank

	// PathEnd specific helpers
	PathEnd core.PathEnd
}

func NewChain(pathEnd *core.PathEnd, chainConfig ChainConfig, client *client.ETHClient, mnemonicPhrase string, simpleTokenAddress, ics20TransferBankAddress, ics20BankAddress string) *Chain {
	ibcHandler, err := ibchandler.NewIbchandler(chainConfig.Chain.IBCAddress(), client)
	if err != nil {
		log.Print(err)
		return nil
	}
	simpletoken, err := simpletoken.NewSimpletoken(common.HexToAddress(simpleTokenAddress), client)
	if err != nil {
		log.Print(err)
		return nil
	}
	ics20transfer, err := ics20transferbank.NewIcs20transferbank(common.HexToAddress(ics20TransferBankAddress), client)
	if err != nil {
		log.Print(err)
		return nil
	}
	ics20bank, err := ics20bank.NewIcs20bank(common.HexToAddress(ics20BankAddress), client)
	if err != nil {
		log.Print(err)
		return nil
	}

	return &Chain{
		client:         client,
		chainID:        chainConfig.Chain.EthChainId,
		ChainConfig:    chainConfig,
		mnemonicPhrase: mnemonicPhrase,
		keys:           make(map[uint32]*ecdsa.PrivateKey),

		IBCHandler: *ibcHandler,

		SimpleToken:   *simpletoken,
		ICS20Transfer: *ics20transfer,
		ICS20Bank:     *ics20bank,
		PathEnd:       *pathEnd,
	}
}

func (chain *Chain) TxOpts(ctx context.Context, index uint32) *bind.TransactOpts {
	return makeGenTxOpts(big.NewInt(chain.chainID), chain.prvKey(index))(ctx)
}

func (chain *Chain) prvKey(index uint32) *ecdsa.PrivateKey {
	key, ok := chain.keys[index]
	if ok {
		return key
	}
	key, err := wallet.GetPrvKeyFromMnemonicAndHDWPath(chain.mnemonicPhrase, fmt.Sprintf("m/44'/60'/0'/0/%v", index))
	if err != nil {
		panic(err)
	}
	chain.keys[index] = key
	return key
}

func (chain *Chain) CallOpts(ctx context.Context, index uint32) *bind.CallOpts {
	opts := chain.TxOpts(ctx, index)
	return &bind.CallOpts{
		From:    opts.From,
		Context: opts.Context,
	}
}

func (chain *Chain) LastHeader(ctx context.Context) (*gethtypes.Header, error) {
	return chain.client.HeaderByNumber(ctx, nil)
}

func makeGenTxOpts(chainID *big.Int, prv *ecdsa.PrivateKey) func(ctx context.Context) *bind.TransactOpts {
	signer := gethtypes.LatestSignerForChainID(chainID)
	addr := gethcrypto.PubkeyToAddress(prv.PublicKey)
	return func(ctx context.Context) *bind.TransactOpts {
		return &bind.TransactOpts{
			From:     addr,
			GasLimit: 6382056,
			Signer: func(address common.Address, tx *gethtypes.Transaction) (*gethtypes.Transaction, error) {
				if address != addr {
					return nil, errors.New("not authorized to sign this account")
				}
				signature, err := gethcrypto.Sign(signer.Hash(tx).Bytes(), prv)
				if err != nil {
					return nil, err
				}
				return tx.WithSignature(signer, signature)
			},
		}
	}
}

func InitializeChains(configDir, simpleTokenAddress, ics20TransferBankAddress, ics20BankAddress string) (*Chain, *Chain, error) {
	pathConfig, err := parsePathConfig(configDir)
	if err != nil {
		return nil, nil, err
	}
	chainConfigs, err := parseChainConfigs(configDir + "/chains")
	if err != nil {
		return nil, nil, err
	}
	src := chainConfigs[0]
	dst := chainConfigs[1]
	ethClientA, err := client.NewETHClient(src.Chain.RpcAddr)
	if err != nil {
		return nil, nil, err
	}
	ethClientB, err := client.NewETHClient(dst.Chain.RpcAddr)
	if err != nil {
		return nil, nil, err
	}
	chainA := NewChain(pathConfig.Src, *src, ethClientA, src.Chain.HdwMnemonic, simpleTokenAddress, ics20TransferBankAddress, ics20BankAddress)
	chainB := NewChain(pathConfig.Dst, *dst, ethClientB, dst.Chain.HdwMnemonic, simpleTokenAddress, ics20TransferBankAddress, ics20BankAddress)

	return chainA, chainB, nil
}

func parsePathConfig(configDir string) (*core.Path, error) {
	files, err := os.ReadDir(configDir)
	if err != nil {
		return nil, err
	}
	var pathConfig core.Path
	for _, f := range files {
		pth := fmt.Sprintf("%s/%s", configDir, f.Name())
		if f.IsDir() {
			continue
		}
		byt, err := os.ReadFile(pth)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(byt, &pathConfig); err != nil {
			return nil, err
		}
	}
	return &pathConfig, nil
}

func parseChainConfigs(configDir string) ([]*ChainConfig, error) {
	files, err := os.ReadDir(configDir)
	if err != nil {
		return nil, err
	}
	var chainConfigs []*ChainConfig
	for _, f := range files {
		var chainConfig ChainConfig
		pth := fmt.Sprintf("%s/%s", configDir, f.Name())
		byt, err := os.ReadFile(pth)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(byt, &chainConfig); err != nil {
			return nil, err
		}
		chainConfigs = append(chainConfigs, &chainConfig)
	}
	return chainConfigs, nil
}
