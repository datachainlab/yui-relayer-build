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
	Use:   "wallet",
	Short: "wallet command",
	Long:  "wallet command",
	Run: func(cmd *cobra.Command, args []string) {
		walletIndex, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		balanceOf(walletIndex)
	},
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

const (
	relayer  = 0
	deployer = 0
)

func balanceOf(index int64) {
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
	ctx := context.Background()

	chainA := NewChain(2018, ethClientA, NewLightClient(ethClientA, ibcclient.MockClient), consts.Contract, mnemonicPhrase, uint64(time.Now().UnixNano()))
	chainB := NewChain(2019, ethClientB, NewLightClient(ethClientB, ibcclient.MockClient), consts.Contract, mnemonicPhrase, uint64(time.Now().UnixNano()))

	baseDenom := strings.ToLower(chainA.ContractConfig.GetSimpleTokenAddress().String())

	bankA, err := chainA.ICS20Bank.BalanceOf(chainA.CallOpts(ctx, relayer), chainA.CallOpts(ctx, uint32(index)).From, baseDenom)
	if err != nil {
		log.Println("BalanceOf Error: ", err)
		os.Exit(1)
	}
	fmt.Println("ChainA: ", bankA)

	bankB, err := chainB.ICS20Bank.BalanceOf(chainB.CallOpts(ctx, relayer), chainB.CallOpts(ctx, uint32(index)).From, baseDenom)
	if err != nil {
		log.Println("BalanceOf Error: ", err)
		os.Exit(1)
	}
	fmt.Println("ChainB: ", bankB)

}

type ContractConfig interface {
	GetIBCHandlerAddress() common.Address
	GetIBCCommitmentTestHelperAddress() common.Address

	GetSimpleTokenAddress() common.Address
	GetICS20TransferBankAddress() common.Address
	GetICS20BankAddress() common.Address
}

type LightClient struct {
	client     *client.ETHClient
	clientType string
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
	// LastLCState LightClientState

	// IBC specific helpers
	ClientIDs []string // ClientID's used on this chain
	// Connections []*TestConnection // track connectionID's created for this chain
	IBCID uint64
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

func (chain *Chain) CallOpts(ctx context.Context, index uint32) *bind.CallOpts {
	opts := chain.TxOpts(ctx, index)
	return &bind.CallOpts{
		From:    opts.From,
		Context: opts.Context,
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

func NewLightClient(cl *client.ETHClient, clientType string) *LightClient {
	return &LightClient{client: cl, clientType: clientType}
}
