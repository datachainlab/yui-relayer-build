require("@nomicfoundation/hardhat-toolbox");

const mnemonic =
  "math razor capable expose worth grape metal sunset metal sudden usage scheme";

/**
 * @type import('hardhat/config').HardhatUserConfig
 */
module.exports = {
  solidity: {
    version: "0.8.20",
    settings: {
      optimizer: {
        enabled: true,
        runs: 9_999_999
      }
    },
  },
  networks: {
    ibc0: {
      url: "http://127.0.0.1:8645",
      accounts: {
        mnemonic: mnemonic,
      },
      chainId: 2018,
    },
    ibc1: {
      url: "http://127.0.0.1:8745",
      accounts: {
        mnemonic: mnemonic,
      },
      chainId: 2019,
    },
  }
}
