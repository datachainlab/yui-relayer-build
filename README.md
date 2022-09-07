# yui-relayer-build

This repository provides a table of [IBC](https://github.com/cosmos/ibc) module combinations supported by [yui-relayer](https://github.com/hyperledger-labs/yui-relayer) and their E2E tests. The E2E tests are performed for each version of the yui-relayer.

In addition, the yui-relayer version corresponds to the [ibc-go](https://github.com/cosmos/ibc-go) version. So please check here for the version compatibility table: https://github.com/hyperledger-labs/yui-relayer#compatibility-with-ibc

## Supported combinations

### v0.3

#### Tendermint

| counterparty | e2e-test                                                                                                                                                                                                      |
|--------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Tendermint   | [![tm2tm@v0.3](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.3-tm2tm.yml/badge.svg?branch=v0.3)](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.3-tm2tm.yml) |

### v0.2

#### Tendermint

| counterparty | e2e-test                                                                                                                                                                                                         |
|--------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Ethereum     | [![tm2eth@v0.2](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.2-tm2eth.yml/badge.svg?branch=v0.2)](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.2-tm2eth.yml) |
| Fabric       | [![tm2fab@v0.2](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.2-tm2fab.yml/badge.svg?branch=v0.2)](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.2-tm2fab.yml) |
| Tendermint   | [![tm2tm@v0.2](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.2-tm2tm.yml/badge.svg?branch=v0.2)](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.2-tm2tm.yml)    |


#### Ethereum

| counterparty | e2e-test                                                                                                                                                                                                         |
|--------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Tendermint   | [![tm2eth@v0.2](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.2-tm2eth.yml/badge.svg?branch=v0.2)](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.2-tm2eth.yml) |


#### Fabric

| counterparty | e2e-test                                                                                                                                                                                                                  |
|--------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Corda        | [![corda2fab@v0.2](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.2-corda2fab.yml/badge.svg?branch=v0.2)](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.2-corda2fab.yml) |
| Fabric       | [![fab2fab@v0.2](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.2-fab2fab.yml/badge.svg)](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.2-fab2fab.yml)                   |
| Tendermint   | [![tm2fab@v0.2](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.2-tm2fab.yml/badge.svg?branch=v0.2)](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.2-tm2fab.yml)          |

#### Corda

| counterparty | e2e-test                                                                                                                                                                                                                        |
|--------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Corda        | [![corda2corda@v0.2](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.2-corda2corda.yml/badge.svg?branch=v0.2)](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.2-corda2corda.yml) |
| Fabric       | [![corda2fab@v0.2](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.2-corda2fab.yml/badge.svg?branch=v0.2)](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.2-corda2fab.yml)       |
