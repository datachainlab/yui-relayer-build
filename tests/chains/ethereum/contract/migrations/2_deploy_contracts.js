const IBFT2Client = artifacts.require("IBFT2Client");
const MockClient = artifacts.require("MockClient");
const IBCClient = artifacts.require("IBCClient");
const IBCConnection = artifacts.require("IBCConnection");
const IBCChannelHandshake = artifacts.require("IBCChannelHandshake");
const IBCPacket = artifacts.require("IBCPacket");
const IBCHandler = artifacts.require("OwnableIBCHandler");
const SimpleToken = artifacts.require("SimpleToken");
const ICS20TransferBank = artifacts.require("ICS20TransferBank");
const ICS20Bank = artifacts.require("ICS20Bank");

const deployCore = async (deployer) => {
  await deployer.deploy(IBCClient);
  await deployer.deploy(IBCConnection);
  await deployer.deploy(IBCChannelHandshake);
  await deployer.deploy(IBCPacket);
  await deployer.deploy(IBCHandler, IBCClient.address, IBCConnection.address, IBCChannelHandshake.address, IBCPacket.address);

  await deployer.deploy(IBFT2Client, IBCHandler.address);
  await deployer.deploy(MockClient, IBCHandler.address);

};

const deployApp = async (deployer) => {
  console.log("deploying app contracts");

  await deployer.deploy(SimpleToken, "simple", "simple", 1000000);
  await deployer.deploy(ICS20Bank)
  await deployer.deploy(ICS20TransferBank, IBCHandler.address, ICS20Bank.address);
};

module.exports = async function (deployer) {
  await deployCore(deployer);
  await deployApp(deployer);
};
