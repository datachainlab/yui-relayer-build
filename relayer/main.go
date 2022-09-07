package main

import (
	"log"

	corda "github.com/hyperledger-labs/yui-corda-ibc/go/relay/module"
	fabric "github.com/hyperledger-labs/yui-fabric-ibc/relay/module"
	ethereum "github.com/hyperledger-labs/yui-ibc-solidity/pkg/relay/module"
	tendermint "github.com/hyperledger-labs/yui-relayer/chains/tendermint/module"
	"github.com/hyperledger-labs/yui-relayer/cmd"
	mock "github.com/hyperledger-labs/yui-relayer/provers/mock/module"
)

func main() {
	if err := cmd.Execute(
		tendermint.Module{},
		ethereum.Module{},
		fabric.Module{},
		corda.Module{},
		mock.Module{},
	); err != nil {
		log.Fatal(err)
	}
}
