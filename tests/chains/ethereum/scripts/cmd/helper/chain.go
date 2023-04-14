package helper

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/client"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/consts"
	ibccommitment "github.com/hyperledger-labs/yui-ibc-solidity/pkg/contract/ibccommitmenttesthelper"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/contract/ibchandler"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/contract/ics20bank"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/contract/ics20transferbank"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/contract/simpletoken"
	ibcclient "github.com/hyperledger-labs/yui-ibc-solidity/pkg/ibc/core/client"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/wallet"
)

const MnemonicPhrase = "math razor capable expose worth grape metal sunset metal sudden usage scheme"

type Chain struct {
	chainID        int64
	client         *client.ETHClient
	lc             *LightClient
	mnemonicPhrase string
	keys           map[uint32]*ecdsa.PrivateKey

	ContractConfig ContractConfig

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
}

func NewChain(chainID int64, client *client.ETHClient, lc *LightClient, config ContractConfig, mnemonicPhrase string, ibcID uint64) *Chain {
	ibcHandler, err := ibchandler.NewIbchandler(config.GetIBCHandlerAddress(), client)
	if err != nil {
		log.Print(err)
		return nil
	}
	ibcCommitment, err := ibccommitment.NewIbccommitmenttesthelper(config.GetIBCCommitmentTestHelperAddress(), client)
	if err != nil {
		log.Print(err)
		return nil
	}
	simpletoken, err := simpletoken.NewSimpletoken(config.GetSimpleTokenAddress(), client)
	if err != nil {
		log.Print(err)
		return nil
	}
	ics20transfer, err := ics20transferbank.NewIcs20transferbank(config.GetICS20TransferBankAddress(), client)
	if err != nil {
		log.Print(err)
		return nil
	}
	ics20bank, err := ics20bank.NewIcs20bank(config.GetICS20BankAddress(), client)
	if err != nil {
		log.Print(err)
		return nil
	}

	return &Chain{
		client:         client,
		chainID:        chainID,
		lc:             lc,
		ContractConfig: config,
		mnemonicPhrase: mnemonicPhrase,
		keys:           make(map[uint32]*ecdsa.PrivateKey),
		IBCID:          ibcID,

		IBCHandler:    *ibcHandler,
		IBCCommitment: *ibcCommitment,

		SimpleToken:   *simpletoken,
		ICS20Transfer: *ics20transfer,
		ICS20Bank:     *ics20bank,
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
		state, err := chain.lc.GetState(ctx, chain.ContractConfig.GetIBCHandlerAddress(), nil, nil)
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

type ContractConfig interface {
	GetIBCHandlerAddress() common.Address
	GetIBCCommitmentTestHelperAddress() common.Address

	GetSimpleTokenAddress() common.Address
	GetICS20TransferBankAddress() common.Address
	GetICS20BankAddress() common.Address
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

func CreateChannel() Channel {
	return Channel{
		PortID:               "transfer",
		ID:                   "channel-0",
		ClientID:             "mock-client-0",
		CounterpartyClientID: "mock-client-0",
		Version:              "transfer-1",
	}
}

func InitializeChains() (*Chain, *Chain, error) {
	ethClientA, err := client.NewETHClient("http://127.0.0.1:8645")
	if err != nil {
		return nil, nil, err
	}

	ethClientB, err := client.NewETHClient("http://127.0.0.1:8745")
	if err != nil {
		return nil, nil, err
	}

	chainA := NewChain(2018, ethClientA, NewLightClient(ethClientA, ibcclient.MockClient), consts.Contract, MnemonicPhrase, uint64(time.Now().UnixNano()))
	chainB := NewChain(2019, ethClientB, NewLightClient(ethClientB, ibcclient.MockClient), consts.Contract, MnemonicPhrase, uint64(time.Now().UnixNano()))

	chainA.UpdateHeader()
	chainB.UpdateHeader()
	return chainA, chainB, nil
}
