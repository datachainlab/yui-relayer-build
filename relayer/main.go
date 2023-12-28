package main

import (
	"log"

	ethereum "github.com/datachainlab/ethereum-ibc-relay-chain/pkg/relay/ethereum/module"
	"github.com/datachainlab/ethereum-ibc-relay-chain/pkg/relay/ethereum/signers/hd"
	tendermint "github.com/hyperledger-labs/yui-relayer/chains/tendermint/module"
	"github.com/hyperledger-labs/yui-relayer/cmd"
	mock "github.com/hyperledger-labs/yui-relayer/provers/mock/module"
)

func main() {
	if err := cmd.Execute(
		tendermint.Module{},
		mock.Module{},
		ethereum.Module{},
		hd.Module{},
	); err != nil {
		log.Fatal(err)
	}
}
