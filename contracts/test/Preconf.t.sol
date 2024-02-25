// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

import {TransparentUpgradeableProxy} from "openzeppelin/proxy/transparent/TransparentUpgradeableProxy.sol";

import {BLSMockAVSDeployer} from "lib/eigenlayer-middleware/test/utils/BLSMockAVSDeployer.sol";

import {PreconfChallengeManager} from "src/PreconfChallengeManager.sol";
import {PreconfServiceManager} from "src/PreconfServiceManager.sol";

contract TestPreconf is BLSMockAVSDeployer {
    /*//////////////////////////////////////////////////////////////
                                CONTRACTS
    //////////////////////////////////////////////////////////////*/

    PreconfChallengeManager public challenge;
    PreconfServiceManager public service;

    /*//////////////////////////////////////////////////////////////
                                 HELPERS
    //////////////////////////////////////////////////////////////*/

    function setUp() public {
        _setUpBLSMockAVSDeployer();

        PreconfChallengeManager challengeImplementation = new PreconfChallengeManager();

        challenge = PreconfChallengeManager(
            address(
                new TransparentUpgradeableProxy(
                    address(challengeImplementation),
                    address(proxyAdmin),
                    abi.encodeWithSelector(
                        PreconfChallengeManager.initialize.selector,
                        address(this),
                        computeCreateAddress(address(this), vm.getNonce(address(this)) + 2),
                        address(this)
                    )
                )
            )
        );

        PreconfServiceManager serviceImplementation =
            new PreconfServiceManager(avsDirectory, registryCoordinator, stakeRegistry);

        service = PreconfServiceManager(
            address(
                new TransparentUpgradeableProxy(
                    address(serviceImplementation),
                    address(proxyAdmin),
                    abi.encodeWithSelector(
                        PreconfServiceManager.initialize.selector, pauserRegistry, 0, address(this), address(challenge)
                    )
                )
            )
        );
    }

    /*//////////////////////////////////////////////////////////////
                                TESTS
    //////////////////////////////////////////////////////////////*/

    /// @notice Unit test to check that the PreconfChallengeManager has the variables set correctly.
    function testChallengeManager() public {
        assertEq(challenge.prover(), address(this));
        assertEq(challenge.serviceManager(), address(service));
    }

    /// @notice Unit test to check that the PreconfServiceManager has the variables set correctly.
    function testServiceManager() public {
        assertEq(address(service.challengeManager()), address(challenge));
    }

    /// @notice Unit test to check that only PreconfChallengeManager can call the slash function.
    function testUnauthorizedSlash() public {
        vm.expectRevert(abi.encodeWithSelector(PreconfServiceManager.NotChallengeManager.selector));
        service.freezeOperator(address(0));
    }
}
