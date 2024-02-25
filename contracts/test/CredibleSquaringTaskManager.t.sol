// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.12;

import {TransparentUpgradeableProxy} from "openzeppelin/proxy/transparent/TransparentUpgradeableProxy.sol";

import {BLSMockAVSDeployer} from "eigenlayer-middleware/test/utils/BLSMockAVSDeployer.sol";

import {ChallengeManager} from "src/ChallengeManager.sol";
import {ServiceManager} from "src/ServiceManager.sol";

contract TestPreconfirmations is BLSMockAVSDeployer {
    /*//////////////////////////////////////////////////////////////
                                CONTRACTS
    //////////////////////////////////////////////////////////////*/

    ChallengeManager public challenge;
    ServiceManager public service;

    /*//////////////////////////////////////////////////////////////
                                 HELPERS
    //////////////////////////////////////////////////////////////*/

    function setUp() public {
        _setUpBLSMockAVSDeployer();

        ChallengeManager challengeImplementation = new ChallengeManager();

        challenge = ChallengeManager(
            address(
                new TransparentUpgradeableProxy(
                    address(challengeImplementation),
                    address(proxyAdmin),
                    abi.encodeWithSelector(
                        ChallengeManager.initialize.selector,
                        address(this),
                        computeCreateAddress(address(this), vm.getNonce(address(this)) + 1),
                        address(this)
                    )
                )
            )
        );

        service = new ServiceManager(avsDirectory, registryCoordinator, stakeRegistry, challenge);
    }

    /*//////////////////////////////////////////////////////////////
                                TESTS
    //////////////////////////////////////////////////////////////*/

    /// @notice Unit test to check that the ChallengeManager has the variables set correctly.
    function testChallengeManager() public {
        assertEq(challenge.prover(), address(this));
        assertEq(challenge.serviceManager(), address(service));
    }

    /// @notice Unit test to check that the ServiceManager has the variables set correctly.
    function testServiceManager() public {
        assertEq(address(service.challengeManager()), address(challenge));
    }

    /// @notice Unit test to check that only ChallengeManager can call the slash function.
    function testUnauthorizedSlash() public {
        vm.expectRevert(abi.encodeWithSelector(ServiceManager.NotChallengeManager.selector));
        service.freezeOperator(address(0));
    }
}
