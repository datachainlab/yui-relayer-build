# yui-relayer-build

This repository provides a table of [IBC](https://github.com/cosmos/ibc) module combinations supported by [yui-relayer](https://github.com/hyperledger-labs/yui-relayer) and their E2E tests. The E2E tests are performed for each version of the yui-relayer.

In addition, the yui-relayer version corresponds to the [ibc-go](https://github.com/cosmos/ibc-go) version. So please check here for the version compatibility table: https://github.com/hyperledger-labs/yui-relayer#compatibility-with-ibc

## Supported combinations

### v0.5

#### Tendermint

| counterparty | e2e-test                                                                                                                                                                                                                | directory                                                                             |
|--------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------|
| Ethereum     | [![tm2eth@v0.5](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.5-tm2eth.yml/badge.svg?branch=v0.5)](https://github.com/datachainlab/yui-relayer-build/blob/v0.5/.github/workflows/v0.5-tm2eth.yml) | [Link](https://github.com/datachainlab/yui-relayer-build/tree/v0.5/tests/cases/tm2eth) |
| Tendermint   | [![tm2tm@v0.5](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.5-tm2tm.yml/badge.svg?branch=v0.5)](https://github.com/datachainlab/yui-relayer-build/blob/v0.5/.github/workflows/v0.5-tm2tm.yml) | [Link](https://github.com/datachainlab/yui-relayer-build/tree/v0.5/tests/cases/tm2tm) |

#### Ethereum

| counterparty | e2e-test                                                                                                                                                                                                                | directory                                                                             |
|--------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------|
| Ethereum     | [![eth2eth@v0.5](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.5-eth2eth.yml/badge.svg?branch=v0.5)](https://github.com/datachainlab/yui-relayer-build/blob/v0.5/.github/workflows/v0.5-eth2eth.yml) | [Link](https://github.com/datachainlab/yui-relayer-build/tree/v0.5/tests/cases/eth2eth) |
| Tendermint   | [![tm2eth@v0.5](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.5-tm2eth.yml/badge.svg?branch=v0.5)](https://github.com/datachainlab/yui-relayer-build/blob/v0.5/.github/workflows/v0.5-tm2eth.yml) | [Link](https://github.com/datachainlab/yui-relayer-build/tree/v0.5/tests/cases/tm2eth) |


### v0.4

#### Tendermint

| counterparty | e2e-test                                                                                                                                                                                                                | directory                                                                             |
|--------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------|
| Ethereum     | [![tm2eth@v0.4](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.4-tm2eth.yml/badge.svg?branch=v0.4)](https://github.com/datachainlab/yui-relayer-build/blob/v0.4/.github/workflows/v0.4-tm2eth.yml) | [Link](https://github.com/datachainlab/yui-relayer-build/tree/v0.4/tests/cases/tm2eth) |
| Tendermint   | [![tm2tm@v0.4](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.4-tm2tm.yml/badge.svg?branch=v0.4)](https://github.com/datachainlab/yui-relayer-build/blob/v0.4/.github/workflows/v0.4-tm2tm.yml) | [Link](https://github.com/datachainlab/yui-relayer-build/tree/v0.4/tests/cases/tm2tm) |

#### Ethereum

| counterparty | e2e-test                                                                                                                                                                                                                | directory                                                                             |
|--------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------|
| Ethereum     | [![eth2eth@v0.4](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.4-eth2eth.yml/badge.svg?branch=v0.4)](https://github.com/datachainlab/yui-relayer-build/blob/v0.4/.github/workflows/v0.4-eth2eth.yml) | [Link](https://github.com/datachainlab/yui-relayer-build/tree/v0.4/tests/cases/eth2eth) |
| Tendermint   | [![tm2eth@v0.4](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.4-tm2eth.yml/badge.svg?branch=v0.4)](https://github.com/datachainlab/yui-relayer-build/blob/v0.4/.github/workflows/v0.4-tm2eth.yml) | [Link](https://github.com/datachainlab/yui-relayer-build/tree/v0.4/tests/cases/tm2eth) |


### v0.3

#### Tendermint

| counterparty | e2e-test                                                                                                                                                                                                                | directory                                                                             |
|--------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------|
| Tendermint   | [![tm2tm@v0.3](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.3-tm2tm.yml/badge.svg?branch=v0.3)](https://github.com/datachainlab/yui-relayer-build/blob/v0.3/.github/workflows/v0.3-tm2tm.yml) | [Link](https://github.com/datachainlab/yui-relayer-build/tree/v0.3/tests/cases/tm2tm) |

#### Ethereum

| counterparty | e2e-test                                                                                                                                                                                                                | directory                                                                             |
|--------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------|
| Ethereum   | [![eth2eth@v0.3](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.3-eth2eth.yml/badge.svg?branch=v0.3)](https://github.com/datachainlab/yui-relayer-build/blob/v0.3/.github/workflows/v0.3-eth2eth.yml) | [Link](https://github.com/datachainlab/yui-relayer-build/tree/v0.3/tests/cases/eth2eth) |

### v0.2

#### Tendermint

| counterparty | e2e-test                                                                                                                                                                                                                   | directory                                                                              |
|--------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------|
| Ethereum     | [![tm2eth@v0.2](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.2-tm2eth.yml/badge.svg?branch=v0.2)](https://github.com/datachainlab/yui-relayer-build/blob/v0.2/.github/workflows/v0.2-tm2eth.yml) | [Link](https://github.com/datachainlab/yui-relayer-build/tree/v0.2/tests/cases/tm2eth) |
| Fabric       | [![tm2fab@v0.2](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.2-tm2fab.yml/badge.svg?branch=v0.2)](https://github.com/datachainlab/yui-relayer-build/blob/v0.2/.github/workflows/v0.2-tm2fab.yml) | [Link](https://github.com/datachainlab/yui-relayer-build/tree/v0.2/tests/cases/tm2fab) |
| Tendermint   | [![tm2tm@v0.2](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.2-tm2tm.yml/badge.svg?branch=v0.2)](https://github.com/datachainlab/yui-relayer-build/blob/v0.2/.github/workflows/v0.2-tm2tm.yml)    | [Link](https://github.com/datachainlab/yui-relayer-build/tree/v0.2/tests/cases/tm2tm)  |


#### Ethereum

| counterparty | e2e-test                                                                                                                                                                                                                   | directory                                                                              |
|--------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------|
| Tendermint   | [![tm2eth@v0.2](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.2-tm2eth.yml/badge.svg?branch=v0.2)](https://github.com/datachainlab/yui-relayer-build/blob/v0.2/.github/workflows/v0.2-tm2eth.yml) | [Link](https://github.com/datachainlab/yui-relayer-build/tree/v0.2/tests/cases/tm2eth) |


#### Fabric

| counterparty | e2e-test                                                                                                                                                                                                                            | directory                                                                                 |
|--------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------|
| Corda        | [![corda2fab@v0.2](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.2-corda2fab.yml/badge.svg?branch=v0.2)](https://github.com/datachainlab/yui-relayer-build/blob/v0.2/.github/workflows/v0.2-corda2fab.yml) | [Link](https://github.com/datachainlab/yui-relayer-build/tree/v0.2/tests/cases/corda2fab) |
| Fabric       | [![fab2fab@v0.2](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.2-fab2fab.yml/badge.svg)](https://github.com/datachainlab/yui-relayer-build/blob/v0.2/.github/workflows/v0.2-fab2fab.yml)                   | [Link](https://github.com/datachainlab/yui-relayer-build/tree/v0.2/tests/cases/fab2fab)   |
| Tendermint   | [![tm2fab@v0.2](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.2-tm2fab.yml/badge.svg?branch=v0.2)](https://github.com/datachainlab/yui-relayer-build/blob/v0.2/.github/workflows/v0.2-tm2fab.yml)          | [Link](https://github.com/datachainlab/yui-relayer-build/tree/v0.2/tests/cases/tm2fab)    |

#### Corda

| counterparty | e2e-test                                                                                                                                                                                                                                  | directory                                                                                   |
|--------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------|
| Corda        | [![corda2corda@v0.2](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.2-corda2corda.yml/badge.svg?branch=v0.2)](https://github.com/datachainlab/yui-relayer-build/blob/v0.2/.github/workflows/v0.2-corda2corda.yml) | [Link](https://github.com/datachainlab/yui-relayer-build/tree/v0.2/tests/cases/corda2corda) |
| Fabric       | [![corda2fab@v0.2](https://github.com/datachainlab/yui-relayer-build/actions/workflows/v0.2-corda2fab.yml/badge.svg?branch=v0.2)](https://github.com/datachainlab/yui-relayer-build/blob/v0.2/.github/workflows/v0.2-corda2fab.yml)       | [Link](https://github.com/datachainlab/yui-relayer-build/tree/v0.2/tests/cases/corda2fab)   |
