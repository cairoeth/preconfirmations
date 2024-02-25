// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

import {PauserRegistry} from "eigenlayer-core/contracts/permissions/PauserRegistry.sol";
import {EmptyContract} from "eigenlayer-core/test/mocks/EmptyContract.sol";

import {BLSApkRegistry} from "eigenlayer-middleware/BLSApkRegistry.sol";
import {RegistryCoordinator} from "eigenlayer-middleware/RegistryCoordinator.sol";
import {OperatorStateRetriever} from "eigenlayer-middleware/OperatorStateRetriever.sol";
import {IRegistryCoordinator} from "eigenlayer-middleware/interfaces/IRegistryCoordinator.sol";
import {IndexRegistry} from "eigenlayer-middleware/IndexRegistry.sol";
import {IIndexRegistry} from "eigenlayer-middleware/interfaces/IIndexRegistry.sol";
import {StakeRegistry} from "eigenlayer-middleware/StakeRegistry.sol";
import {IStakeRegistry} from "eigenlayer-middleware/interfaces/IStakeRegistry.sol";
import {IServiceManager} from "eigenlayer-middleware/interfaces/IServiceManager.sol";
import {IBLSApkRegistry} from "eigenlayer-middleware/interfaces/IBLSApkRegistry.sol";

import {PreconfChallengeManager} from "src/PreconfChallengeManager.sol";
import {PreconfServiceManager} from "src/PreconfServiceManager.sol";
import {ERC20Mock} from "src/ERC20Mock.sol";

import "eigenlayer-scripts/middleware/DeployOpenEigenLayer.s.sol";
import "forge-std/Test.sol";
import "forge-std/Script.sol";
import "forge-std/StdJson.sol";

contract DeployPreconf is DeployOpenEigenLayer {
    // Preconf contracts
    ProxyAdmin public eigenDAProxyAdmin;
    PauserRegistry public eigenDAPauserReg;

    BLSApkRegistry public apkRegistry;
    PreconfServiceManager public preconfServiceManager;
    PreconfChallengeManager public preconfChallengeManager;
    RegistryCoordinator public registryCoordinator;
    IIndexRegistry public indexRegistry;
    IStakeRegistry public stakeRegistry;
    OperatorStateRetriever public operatorStateRetriever;

    ERC20Mock token;

    struct AddressConfig {
        address eigenLayerCommunityMultisig;
        address eigenLayerOperationsMultisig;
        address eigenLayerPauserMultisig;
        address eigenDACommunityMultisig;
        address eigenDAPauser;
        address churner;
        address ejector;
        address confirmer;
    }

    // deploy all the EigenDA contracts. Relies on many EL contracts having already been deployed.
    function run() external {
        vm.startBroadcast(vm.envUint("PRIVATE_KEY"));
        address deployer = vm.addr(vm.envUint("PRIVATE_KEY"));

        _deployEigenDAAndEigenLayerContracts(deployer, 1, type(uint128).max);

        // Fund operator
        payable(address(0x860B6912C2d0337ef05bbC89b0C2CB6CbAEAB4A5)).transfer(500 ether);

        vm.stopBroadcast();

        // Write JSON output
        string memory parent_object = "parent object";

        string memory deployed_addresses = "addresses";
        vm.serializeAddress(deployed_addresses, "erc20Mock", address(token));
        vm.serializeAddress(deployed_addresses, "erc20MockStrategy", address(deployedStrategyArray[0]));
        vm.serializeAddress(deployed_addresses, "serviceManager", address(preconfServiceManager));
        vm.serializeAddress(deployed_addresses, "challengeManager", address(preconfChallengeManager));
        vm.serializeAddress(deployed_addresses, "registryCoordinator", address(registryCoordinator));
        string memory deployed_addresses_output =
            vm.serializeAddress(deployed_addresses, "operatorStateRetriever", address(operatorStateRetriever));

        // serialize all the data
        string memory finalJson = vm.serializeString(parent_object, deployed_addresses, deployed_addresses_output);

        writeOutput(finalJson, "preconf_avs_deployment_output");
    }

    function _deployEigenDAAndEigenLayerContracts(address deployer, uint8 numStrategies, uint256 maxOperatorCount)
        internal
    {
        token = new ERC20Mock();

        StrategyConfig[] memory strategyConfigs = new StrategyConfig[](1);

        strategyConfigs[0] = StrategyConfig({
            maxDeposits: type(uint256).max,
            maxPerDeposit: type(uint256).max,
            tokenAddress: address(token),
            tokenSymbol: token.symbol()
        });

        _deployEigenLayer(deployer, deployer, deployer, strategyConfigs);

        // deploy proxy admin for ability to upgrade proxy contracts
        eigenDAProxyAdmin = new ProxyAdmin();

        // deploy pauser registry
        {
            address[] memory pausers = new address[](2);
            pausers[0] = deployer;
            pausers[1] = deployer;
            eigenDAPauserReg = new PauserRegistry(pausers, deployer);
        }

        emptyContract = new EmptyContract();

        // hard-coded inputs

        /**
         * First, deploy upgradeable proxy contracts that **will point** to the implementations. Since the implementation contracts are
         * not yet deployed, we give these proxies an empty contract as the initial implementation, to act as if they have no code.
         */
        preconfServiceManager = PreconfServiceManager(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        preconfChallengeManager = PreconfChallengeManager(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        registryCoordinator = RegistryCoordinator(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        indexRegistry = IIndexRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        stakeRegistry = IStakeRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );
        apkRegistry = BLSApkRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenDAProxyAdmin), ""))
        );

        eigenDAProxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(indexRegistry))),
            address(new IndexRegistry(registryCoordinator))
        );

        eigenDAProxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(stakeRegistry))),
            address(new StakeRegistry(registryCoordinator, IDelegationManager(delegation)))
        );

        eigenDAProxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(apkRegistry))), address(new BLSApkRegistry(registryCoordinator))
        );

        {
            IRegistryCoordinator.OperatorSetParam[] memory operatorSetParams =
                new IRegistryCoordinator.OperatorSetParam[](numStrategies);
            for (uint256 i = 0; i < numStrategies; i++) {
                // hard code these for now
                operatorSetParams[i] = IRegistryCoordinator.OperatorSetParam({
                    maxOperatorCount: uint32(maxOperatorCount),
                    kickBIPsOfOperatorStake: 11000, // an operator needs to have kickBIPsOfOperatorStake / 10000 times the stake of the operator with the least stake to kick them out
                    kickBIPsOfTotalStake: 1001 // an operator needs to have less than kickBIPsOfTotalStake / 10000 of the total stake to be kicked out
                });
            }

            uint96[] memory minimumStakeForQuourm = new uint96[](numStrategies);
            IStakeRegistry.StrategyParams[][] memory strategyAndWeightingMultipliers =
                new IStakeRegistry.StrategyParams[][](numStrategies);
            for (uint256 i = 0; i < numStrategies; i++) {
                strategyAndWeightingMultipliers[i] = new IStakeRegistry.StrategyParams[](1);
                strategyAndWeightingMultipliers[i][0] =
                    IStakeRegistry.StrategyParams({strategy: deployedStrategyArray[i], multiplier: 1 ether});
            }

            eigenDAProxyAdmin.upgradeAndCall(
                TransparentUpgradeableProxy(payable(address(registryCoordinator))),
                address(
                    new RegistryCoordinator(
                        IServiceManager(address(preconfServiceManager)), stakeRegistry, apkRegistry, indexRegistry
                    )
                ),
                abi.encodeWithSelector(
                    RegistryCoordinator.initialize.selector,
                    deployer,
                    deployer,
                    deployer,
                    IPauserRegistry(address(eigenDAPauserReg)),
                    0, // initial paused status is nothing paused
                    operatorSetParams,
                    minimumStakeForQuourm,
                    strategyAndWeightingMultipliers
                )
            );
        }

        // Third, upgrade the proxy contracts to use the correct implementation contracts and initialize them.
        eigenDAProxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(preconfServiceManager))),
            address(new PreconfServiceManager(IAVSDirectory(address(avsDirectory)), registryCoordinator, stakeRegistry)),
            abi.encodeWithSelector(
                PreconfServiceManager.initialize.selector,
                eigenDAPauserReg,
                0,
                deployer,
                address(preconfChallengeManager)
            )
        );

        eigenDAProxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(preconfChallengeManager))),
            address(new PreconfChallengeManager()),
            abi.encodeWithSelector(
                PreconfChallengeManager.initialize.selector, deployer, address(preconfServiceManager), address(0)
            )
        );

        operatorStateRetriever = new OperatorStateRetriever();

        {
            IStrategy[] memory strategies = new IStrategy[](numStrategies);
            bool[] memory transferLocks = new bool[](numStrategies);
            for (uint8 i = 0; i < numStrategies; i++) {
                strategies[i] = deployedStrategyArray[i];
            }
            strategyManager.addStrategiesToDepositWhitelist(strategies, transferLocks);
        }
    }

    function writeOutput(string memory outputJson, string memory outputFileName) internal {
        string memory outputDir = string.concat(vm.projectRoot(), "/script/output/");
        string memory chainDir = string.concat(vm.toString(block.chainid), "/");
        string memory outputFilePath = string.concat(outputDir, chainDir, outputFileName, ".json");
        vm.writeJson(outputJson, outputFilePath);
    }
}
