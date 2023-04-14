package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/chains"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/client"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/consts"
	ibccommitment "github.com/hyperledger-labs/yui-ibc-solidity/pkg/contract/ibccommitmenttesthelper"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/contract/ibchandler"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/contract/ics20bank"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/contract/ics20transferbank"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/contract/simpletoken"
	ibcclient "github.com/hyperledger-labs/yui-ibc-solidity/pkg/ibc/core/client"
	"github.com/hyperledger-labs/yui-ibc-solidity/pkg/wallet"
	"github.com/spf13/cobra"
)

const mnemonicPhrase = "math razor capable expose worth grape metal sunset metal sudden usage scheme"

var rootCmd = &cobra.Command{
	Use:   "transfer",
	Short: "transfer command",
	Long:  "transfer command",
	Run: func(cmd *cobra.Command, args []string) {
		fromIndex, err := strconv.ParseInt(args[0], 10, 32)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		toIndex, err := strconv.ParseInt(args[1], 10, 32)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		amount, err := strconv.ParseInt(args[2], 10, 64)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		Transfer(uint32(fromIndex), uint32(toIndex), amount)
	},
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// ============== LightClientState ==============
type LightClientState interface {
	Header() *gethtypes.Header
	Proof() *client.StateProof
}

type ETHState struct {
	header     *gethtypes.Header
	StateProof *client.StateProof
}

var _ LightClientState = (*ETHState)(nil)

func (cs ETHState) Header() *gethtypes.Header {
	return cs.header
}

func (cs ETHState) Proof() *client.StateProof {
	return cs.StateProof
}

// ============== IBFT2State ==============
type IBFT2State struct {
	ParsedHeader *chains.ParsedHeader
	StateProof   *client.StateProof
	CommitSeals  [][]byte
}

func (cs IBFT2State) Header() *gethtypes.Header {
	return cs.ParsedHeader.Base
}

func (cs IBFT2State) Proof() *client.StateProof {
	return cs.StateProof
}

// ============== LightClietn ==============
type LightClient struct {
	client     *client.ETHClient
	clientType string
}

func NewLightClient(cl *client.ETHClient, clientType string) *LightClient {
	return &LightClient{client: cl, clientType: clientType}
}

func (lc LightClient) GetState(ctx context.Context, address common.Address, storageKeys [][]byte, bn *big.Int) (LightClientState, error) {
	switch lc.clientType {
	case ibcclient.BesuIBFT2Client:
		return lc.GetIBFT2State(ctx, address, storageKeys, bn)
	case ibcclient.MockClient:
		return lc.GetMockContractState(ctx, address, storageKeys, bn)
	default:
		panic(fmt.Sprintf("unknown client type '%v'", lc.clientType))
	}
}

func (lc LightClient) GetMockContractState(ctx context.Context, address common.Address, storageKeys [][]byte, bn *big.Int) (LightClientState, error) {
	block, err := lc.client.BlockByNumber(ctx, bn)
	if err != nil {
		return nil, err
	}
	proof := &client.StateProof{
		StorageProofRLP: make([][]byte, len(storageKeys)),
	}
	return ETHState{header: block.Header(), StateProof: proof}, nil
}

func (lc LightClient) GetIBFT2State(ctx context.Context, address common.Address, storageKeys [][]byte, bn *big.Int) (LightClientState, error) {
	var state IBFT2State
	block, err := lc.client.BlockByNumber(ctx, bn)
	if err != nil {
		return nil, err
	}
	proof, err := lc.client.GetProof(address, storageKeys, block.Number())
	if err != nil {
		return nil, err
	}
	state.StateProof = proof
	state.ParsedHeader, err = chains.ParseHeader(block.Header())
	if err != nil {
		return nil, err
	}
	state.CommitSeals, err = state.ParsedHeader.ValidateAndGetCommitSeals()
	if err != nil {
		return nil, err
	}
	return state, nil
}

// ============== ContractConfig ==============
type ContractConfig interface {
	GetIBCHandlerAddress() common.Address
	GetIBCCommitmentTestHelperAddress() common.Address

	GetSimpleTokenAddress() common.Address
	GetICS20TransferBankAddress() common.Address
	GetICS20BankAddress() common.Address
}

type TestConnection struct {
	ID                   string
	ClientID             string
	CounterpartyClientID string
	NextChannelVersion   string
	Channels             []TestChannel
}

type TestChannel struct {
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
	ClientIDs   []string          // ClientID's used on this chain
	Connections []*TestConnection // track connectionID's created for this chain
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

func Transfer(fromIndex, toIndex uint32, amount int64) error {
	ethClientA, err := client.NewETHClient("http://127.0.0.1:8645")
	if err != nil {
		log.Println("NewETHClient Error: ", err)
		os.Exit(1)
	}
	ethClientB, err := client.NewETHClient("http://127.0.0.1:8745")
	if err != nil {
		log.Println("NewETHClient Error: ", err)
		os.Exit(1)
	}
	chainA := NewChain(2018, ethClientA, NewLightClient(ethClientA, ibcclient.MockClient), consts.Contract, mnemonicPhrase, uint64(time.Now().UnixNano()))
	chainB := NewChain(2019, ethClientB, NewLightClient(ethClientB, ibcclient.MockClient), consts.Contract, mnemonicPhrase, uint64(time.Now().UnixNano()))

	chainA.UpdateHeader()
	chainB.UpdateHeader()

	ctx := context.Background()

	const (
		relayer  = 0
		deployer = 0
	)

	chanA := TestChannel{
		PortID:               "transfer",
		ID:                   "channel-0",
		ClientID:             "mock-client-0",
		CounterpartyClientID: "mock-client-0",
		Version:              "transfer-1",
	}

	_, err = chainA.SimpleToken.Approve(chainA.TxOpts(ctx, deployer), chainA.ContractConfig.GetICS20BankAddress(), big.NewInt(amount))
	if err != nil {
		log.Println("token approve error: ", err)
		os.Exit(1)
	}
	log.Println("1. token approve success")

	_, err = chainA.ICS20Bank.Deposit(
		chainA.TxOpts(ctx, deployer),
		chainA.ContractConfig.GetSimpleTokenAddress(),
		big.NewInt(amount),
		chainA.CallOpts(ctx, fromIndex).From,
	)
	if err != nil {
		log.Println("deposit error: ", err)
		os.Exit(1)
	}
	log.Println("2. deposit success")

	baseDenom := strings.ToLower(chainA.ContractConfig.GetSimpleTokenAddress().String())

	_, err = chainA.ICS20Transfer.SendTransfer(
		chainA.TxOpts(ctx, fromIndex),
		baseDenom,
		uint64(amount),
		chainB.CallOpts(ctx, toIndex).From,
		chanA.PortID, chanA.ID,
		uint64(chainB.LastHeader().Number.Int64())+1000,
	)
	if err != nil {
		log.Println("sendTransfer error: ", err)
		os.Exit(1)
	}
	log.Println("3. sendTransfer success")

	return nil
}
