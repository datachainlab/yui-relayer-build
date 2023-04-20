package main

import (
	"context"
	"log"
	"math/big"
	"strconv"
	"strings"

	"github.com/datachainlab/yui-relayer-build/tests/chains/ethereum/scripts/cmd/helper"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "transfer",
	Short: "transfer command",
	Long:  "transfer command fromIndex toIndex amount",
	Run: func(cmd *cobra.Command, args []string) {
		configDir := args[0]
		fromIndex, err := strconv.ParseInt(args[1], 10, 32)
		if err != nil {
			log.Fatalln(err)
		}
		toIndex, err := strconv.ParseInt(args[2], 10, 32)
		if err != nil {
			log.Fatalln(err)
		}
		amount, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			log.Fatalln(err)
		}
		simpleTokenAddress := args[4]
		ics20TransferBankAddress := args[5]
		ics20BankAddress := args[6]
		err = Transfer(configDir, uint32(fromIndex), uint32(toIndex), amount, simpleTokenAddress, ics20TransferBankAddress, ics20BankAddress)
		if err != nil {
			log.Fatalln("transfer Error: ", err)
		}
	},
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}

func Transfer(configDir string, fromIndex, toIndex uint32, amount int64, simpleTokenAddress, ics20TransferBankAddress, ics20BankAddress string) error {
	chainA, chainB, err := helper.InitializeChains(configDir, simpleTokenAddress, ics20TransferBankAddress, ics20BankAddress)
	if err != nil {
		return err
	}
	ctx := context.Background()
	const (
		relayer  = 0
		deployer = 0
	)
	chanA := chainA.GetChannel()
	_, err = chainA.SimpleToken.Approve(chainA.TxOpts(ctx, deployer), common.HexToAddress(ics20BankAddress), big.NewInt(amount))
	if err != nil {
		return err
	}
	log.Println("1. token approve success")

	_, err = chainA.ICS20Bank.Deposit(
		chainA.TxOpts(ctx, deployer),
		common.HexToAddress(simpleTokenAddress),
		big.NewInt(amount),
		chainA.CallOpts(ctx, fromIndex).From,
	)
	if err != nil {
		return err
	}
	log.Println("2. deposit success")

	var heightB uint64
	if header, err := chainB.LastHeader(ctx); err != nil {
		log.Fatalf("failed to the latest header of chain B: %v", err)
	} else {
		heightB = header.Number.Uint64()
	}
	baseDenom := strings.ToLower(simpleTokenAddress)
	_, err = chainA.ICS20Transfer.SendTransfer(
		chainA.TxOpts(ctx, fromIndex),
		baseDenom,
		uint64(amount),
		chainB.CallOpts(ctx, toIndex).From,
		chanA.PortID, chanA.ChannelID,
		heightB+1000,
	)
	if err != nil {
		return err
	}
	log.Println("3. sendTransfer success")

	return nil
}
