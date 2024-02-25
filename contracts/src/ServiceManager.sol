// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

import {BytesLib} from "eigenlayer/contracts/libraries/BytesLib.sol";

import {
    ServiceManagerBase,
    IAVSDirectory,
    IRegistryCoordinator,
    IStakeRegistry
} from "eigenlayer-middleware/src/ServiceManagerBase.sol";

import {ChallengeManager} from "src/ChallengeManager.sol";

/// @title ServiceManager
/// @author @cairoeth
/// @notice Primary entrypoint for procuring services for preconfirmations.
contract ServiceManager is ServiceManagerBase {
    using BytesLib for bytes;

    /*//////////////////////////////////////////////////////////////
                                VARIABLES
    //////////////////////////////////////////////////////////////*/

    ChallengeManager public immutable challengeManager;

    /*//////////////////////////////////////////////////////////////
                                 ERRORS
    //////////////////////////////////////////////////////////////*/

    /// @notice Emitted when the caller is not the ChallengeManager.
    error NotChallengeManager();

    /*//////////////////////////////////////////////////////////////
                                MODIFIERS
    //////////////////////////////////////////////////////////////*/

    /// @notice Ensures a function is only callable by the TaskManager.
    modifier onlyChallengeManager() {
        if (msg.sender != address(challengeManager)) {
            revert NotChallengeManager();
        }
        _;
    }

    /*//////////////////////////////////////////////////////////////
                               CONSTRUCTOR
    //////////////////////////////////////////////////////////////*/

    constructor(
        IAVSDirectory _avsDirectory,
        IRegistryCoordinator _registryCoordinator,
        IStakeRegistry _stakeRegistry,
        ChallengeManager _challengeManager
    ) ServiceManagerBase(_avsDirectory, _registryCoordinator, _stakeRegistry) {
        challengeManager = _challengeManager;
    }

    /*//////////////////////////////////////////////////////////////
                            EXTERNAL FUNCTIONS
    //////////////////////////////////////////////////////////////*/

    /// @notice Called by the ChallengeManager to slash an operator that violated a preconfirmation.
    /// @param operatorAddr The address of the operator to be slashed.
    function freezeOperator(address operatorAddr) external onlyChallengeManager {
        // TODO
        // slasher.freezeOperator(operatorAddr);
    }
}
