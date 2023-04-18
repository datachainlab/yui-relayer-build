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
	"time"

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
	ibcclient "github.com/hyperledger-labs/yui-ibc-solidity/pkg/ibc/core/client"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/wallet"
)

type PathConfig struct {
	Src      PathInfo `json:"src"`
	Dst      PathInfo `json:"dst"`
	Mnemonic string   `json:"mnemonic"`
}

type PathInfo struct {
	ChainID      string `json:"chain-id"`
	ClientID     string `json:"client-id"`
	ConnectionID string `json:"connection-id"`
	ChannelID    string `json:"channel-id"`
	PortID       string `json:"port-id"`
	Order        string `json:"order"`
	Version      string `json:"version"`
}

type ChainConfig struct {
	Chain     ChainInfo    `json:"chain"`
	Prover    ProverConfig `json:"prover"`
	Addresses struct {
		SimpleTokenAddress       string `json:"simple_token_address"`
		ICS20TransferBankAddress string `json:"ics20_transfer_bank_address"`
		ICS20BankAddress         string `json:"ics20_bank_address"`
	}
}

type ChainInfo struct {
	Type        string `json:"@type"`
	ChainID     string `json:"chain_id"`
	EthChainID  int    `json:"eth_chain_id"`
	RPCAddr     string `json:"rpc_addr"`
	HDWMnemonic string `json:"hdw_mnemonic"`
	HDWPath     string `json:"hdw_path"`
	IBCAddress  string `json:"ibc_address"`
}

type ProverConfig struct {
	Type string `json:"@type"`
}

type Connection struct {
	ID                   string
	ClientID             string
	CounterpartyClientID string
	NextChannelVersion   string
	Channels             []Channel
}

type Channel struct {
	PortID               string
	ID                   string
	ClientID             string
	CounterpartyClientID string
	Version              string
}

type Chain struct {
	chainID        int64
	client         *client.ETHClient
	lc             *LightClient
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

	// State
	LastLCState LightClientState

	// IBC specific helpers
	ClientIDs   []string      // ClientID's used on this chain
	Connections []*Connection // track connectionID's created for this chain
	IBCID       uint64

	// Channel specific helpers
	Channel Channel
}

func NewChain(pathInfo PathInfo, chainConfig ChainConfig, client *client.ETHClient, lc *LightClient, mnemonicPhrase string, ibcID uint64) *Chain {
	ibcHandler, err := ibchandler.NewIbchandler(common.HexToAddress(chainConfig.Chain.IBCAddress), client)
	if err != nil {
		log.Print(err)
		return nil
	}
	simpletoken, err := simpletoken.NewSimpletoken(common.HexToAddress(chainConfig.Addresses.SimpleTokenAddress), client)
	if err != nil {
		log.Print(err)
		return nil
	}
	ics20transfer, err := ics20transferbank.NewIcs20transferbank(common.HexToAddress(chainConfig.Addresses.ICS20TransferBankAddress), client)
	if err != nil {
		log.Print(err)
		return nil
	}
	ics20bank, err := ics20bank.NewIcs20bank(common.HexToAddress(chainConfig.Addresses.ICS20BankAddress), client)
	if err != nil {
		log.Print(err)
		return nil
	}

	return &Chain{
		client:         client,
		chainID:        int64(chainConfig.Chain.EthChainID),
		lc:             lc,
		ChainConfig:    chainConfig,
		mnemonicPhrase: mnemonicPhrase,
		keys:           make(map[uint32]*ecdsa.PrivateKey),
		IBCID:          ibcID,

		IBCHandler: *ibcHandler,

		SimpleToken:   *simpletoken,
		ICS20Transfer: *ics20transfer,
		ICS20Bank:     *ics20bank,
		Channel: Channel{
			PortID:               pathInfo.PortID,
			ID:                   pathInfo.ChannelID,
			ClientID:             pathInfo.ClientID,
			CounterpartyClientID: pathInfo.ClientID,
			Version:              pathInfo.Version,
		},
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

func (chain *Chain) LastHeader() *gethtypes.Header {
	return chain.LastLCState.Header()
}

func (chain *Chain) UpdateHeader() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	for {
		state, err := chain.lc.GetState(ctx, common.HexToAddress(chain.ChainConfig.Chain.IBCAddress), nil, nil)
		if err != nil {
			panic(err)
		}
		if chain.LastLCState == nil || state.Header().Number.Cmp(chain.LastHeader().Number) == 1 {
			chain.LastLCState = state
			return
		} else {
			continue
		}
	}
}

func (chain *Chain) GetChannel() Channel {
	return chain.Channel
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

func InitializeChains(configDir string) (*Chain, *Chain, error) {
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
	ethClientA, err := client.NewETHClient(src.Chain.RPCAddr)
	if err != nil {
		return nil, nil, err
	}
	ethClientB, err := client.NewETHClient(dst.Chain.RPCAddr)
	if err != nil {
		return nil, nil, err
	}
	chainA := NewChain(pathConfig.Src, *src, ethClientA, NewLightClient(ethClientA, ibcclient.MockClient), src.Chain.HDWMnemonic, uint64(time.Now().UnixNano()))
	chainB := NewChain(pathConfig.Dst, *dst, ethClientB, NewLightClient(ethClientB, ibcclient.MockClient), dst.Chain.HDWMnemonic, uint64(time.Now().UnixNano()))

	chainA.UpdateHeader()
	chainB.UpdateHeader()
	return chainA, chainB, nil
}

func parsePathConfig(configDir string) (*PathConfig, error) {
	files, err := os.ReadDir(configDir)
	if err != nil {
		return nil, err
	}
	var pathConfig PathConfig
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
