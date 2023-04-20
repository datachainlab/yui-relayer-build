const IBCHandler = artifacts.require("OwnableIBCHandler");
const MockClient = artifacts.require("MockClient");
const ICS20TransferBank = artifacts.require("ICS20TransferBank");
const ICS20Bank = artifacts.require("ICS20Bank");

const PortTransfer = "transfer"
const MockClientType = "mock-client"

module.exports = async function (deployer) {
  const ibcHandler = await IBCHandler.deployed();
  const ics20Bank = await ICS20Bank.deployed();

  for(const promise of [
    () => ibcHandler.bindPort(PortTransfer, ICS20TransferBank.address),
    () => ibcHandler.registerClient(MockClientType, MockClient.address),
    () => ics20Bank.setOperator(ICS20TransferBank.address),
  ]) {
    const result = await promise();
    console.log(result);
    if(!result.receipt.status) {
      throw new Error(`transaction failed to execute. ${result.tx}`);
    }
  }
};
