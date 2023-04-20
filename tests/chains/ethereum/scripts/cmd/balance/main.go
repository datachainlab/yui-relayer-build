package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/datachainlab/yui-relayer-build/tests/chains/ethereum/scripts/cmd/helper"
	"github.com/spf13/cobra"
)

const (
	relayer = 0
)

var rootCmd = &cobra.Command{
	Use:   "wallet",
	Short: "wallet command",
	Long:  "wallet command walletIndex",
	Run: func(cmd *cobra.Command, args []string) {
		configDir := args[0]
		walletIndex, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		simpleTokenAddress := args[2]
		ics20TransferBankAddress := args[3]
		ics20BankAddress := args[4]
		balanceA, balanceB, err := balanceOf(configDir, walletIndex, simpleTokenAddress, ics20TransferBankAddress, ics20BankAddress)
		if err != nil {
			log.Fatalln("balanceOf Error: ", err)
		}
		fmt.Printf("%d,%d\n", balanceA, balanceB)
	},
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func balanceOf(configDir string, index int64, simpleTokenAddress, ics20TransferBankAddress, ics20BankAddress string) (*big.Int, *big.Int, error) {
	chainA, chainB, err := helper.InitializeChains(configDir, simpleTokenAddress, ics20TransferBankAddress, ics20BankAddress)
	if err != nil {
		return big.NewInt(0), big.NewInt(0), err
	}
	ctx := context.Background()
	baseDenom := strings.ToLower(simpleTokenAddress)
	bankA, err := chainA.ICS20Bank.BalanceOf(chainA.CallOpts(ctx, relayer), chainA.CallOpts(ctx, uint32(index)).From, baseDenom)
	if err != nil {
		return bankA, big.NewInt(0), err
	}
	chanB := chainB.GetChannel()
	expectedDenom := fmt.Sprintf("%v/%v/%v", chanB.PortID, chanB.ID, baseDenom)
	bankB, err := chainB.ICS20Bank.BalanceOf(chainB.CallOpts(ctx, relayer), chainB.CallOpts(ctx, uint32(index)).From, expectedDenom)
	if err != nil {
		return bankA, bankB, err
	}
	return bankA, bankB, nil
}
