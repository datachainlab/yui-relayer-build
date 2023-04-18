package main

import (
	"context"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/datachainlab/yui-relayer-build/tests/chains/ethereum/scripts/cmd/helper"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "transfer",
	Short: "transfer command",
	Long:  "transfer command fromIndex toIndex amount",
	Run: func(cmd *cobra.Command, args []string) {
		pathFile := args[0]
		fromIndex, err := strconv.ParseInt(args[1], 10, 32)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		toIndex, err := strconv.ParseInt(args[2], 10, 32)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		amount, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		Transfer(pathFile, uint32(fromIndex), uint32(toIndex), amount)
	},
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func Transfer(pathFile string, fromIndex, toIndex uint32, amount int64) error {
	chainA, chainB, err := helper.InitializeChains(pathFile)
	if err != nil {
		log.Println("InitializeChains Error: ", err)
		os.Exit(1)
	}
	ctx := context.Background()
	const (
		relayer  = 0
		deployer = 0
	)
	chanA := chainA.GetChannel()
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
