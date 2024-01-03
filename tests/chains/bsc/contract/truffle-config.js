const HDWalletProvider = require("@truffle/hdwallet-provider");
const mnemonic =
  "math razor capable expose worth grape metal sunset metal sudden usage scheme";

module.exports = {
  networks: {
    chain0: {
      network_id: "*",
      maxFeePerGas: 100000000000,
      maxPriorityFeePerGas: 10000000000,
      provider: () =>
        new HDWalletProvider({
          mnemonic: {
            phrase: mnemonic,
          },
          providerOrUrl: "http://127.0.0.1:8845",
          addressIndex: 0,
          numberOfAddresses: 10,
        }),
    },
    chain1: {
      network_id: "*",
      maxFeePerGas: 100000000000,
      maxPriorityFeePerGas: 10000000000,
      provider: () =>
        new HDWalletProvider({
          mnemonic: {
            phrase: mnemonic,
          },
          providerOrUrl: "http://127.0.0.1:8945",
          addressIndex: 0,
          numberOfAddresses: 10,
        }),
    },
  },

  compilers: {
    solc: {
      version: "0.8.12",
      settings: {
        optimizer: {
          enabled: true,
          runs: 1000,
        },
      },
    },
  },
  plugins: ["truffle-contract-size"],
};
