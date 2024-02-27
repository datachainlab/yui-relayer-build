
async function deploy(deployer, contractName, args = []) {
  const factory = await hre.ethers.getContractFactory(contractName);
  const contract = await factory.connect(deployer).deploy(...args);
  await contract.waitForDeployment();
  return contract;
}

async function deployIBC(deployer) {
  const logicNames = [
    "IBCClient",
    "IBCConnectionSelfStateNoValidation",
    "IBCChannelHandshake",
    "IBCChannelPacketSendRecv",
    "IBCChannelPacketTimeout"
  ];
  const logics = [];
  for (const name of logicNames) {
    const logic = await deploy(deployer, name);
    logics.push(logic);
  }
  return deploy(deployer, "OwnableIBCHandler", logics.map(l => l.target));
}

async function main() {
  // This is just a convenience check
  if (network.name === "hardhat") {
    console.warn(
      "You are trying to deploy a contract to the Hardhat Network, which" +
        "gets automatically created and destroyed every time. Use the Hardhat" +
        " option '--network localhost'"
    );
  }

  // ethers is available in the global scope
  const [deployer] = await hre.ethers.getSigners();
  console.log(
    "Deploying the contracts with the account:",
    await deployer.getAddress()
  );
  console.log("Account balance:", (await hre.ethers.provider.getBalance(deployer.getAddress())).toString());

  const ibcHandler = await deployIBC(deployer);
  console.log("IBCHandler address:", ibcHandler.target);

  const erc20token = await deploy(deployer, "ERC20Token", ["simple", "simple", 1000000]);
  console.log("ERC20Token address:", erc20token.target);

  const ics20bank = await deploy(deployer, "ICS20Bank");
  console.log("ICS20Bank address:", ics20bank.target);

  const ics20transferbank = await deploy(deployer, "ICS20TransferBank", [ibcHandler.target, ics20bank.target]);
  console.log("ICS20TransferBank address:", ics20transferbank.target);

  const mockClient = await deploy(deployer, "MockClient", [ibcHandler.target]);
  console.log("MockClient address:", mockClient.target);

  await ibcHandler.bindPort("transfer", ics20transferbank.target);
  await ibcHandler.registerClient("mock-client", mockClient.target);
  await ics20bank.setOperator(ics20transferbank.target);

}

if (require.main === module) {
  main()
    .then(() => process.exit(0))
    .catch((error) => {
      console.error(error);
      process.exit(1);
    });
}

exports.deployIBC = deployIBC;
