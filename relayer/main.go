package main

import (
	"log"

	tendermint "github.com/hyperledger-labs/yui-relayer/chains/tendermint/module"
	"github.com/hyperledger-labs/yui-relayer/cmd"
)

func main() {
	if err := cmd.Execute(
		tendermint.Module{},
	); err != nil {
		log.Fatal(err)
	}
}
