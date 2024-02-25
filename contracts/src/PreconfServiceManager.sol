// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

import {Pausable} from "eigenlayer-core/contracts/permissions/Pausable.sol";
import {IAVSDirectory} from "eigenlayer-core/contracts/interfaces/IAVSDirectory.sol";
import {IPauserRegistry} from "eigenlayer-core/contracts/interfaces/IPauserRegistry.sol";

import {PreconfChallengeManager} from "src/PreconfChallengeManager.sol";

import {ServiceManagerBase} from "eigenlayer-middleware/ServiceManagerBase.sol";
import {IRegistryCoordinator} from "eigenlayer-middleware/interfaces/IRegistryCoordinator.sol";
import {IStakeRegistry} from "eigenlayer-middleware/interfaces/IStakeRegistry.sol";

/// @title PreconfServiceManager
/// @author @cairoeth
/// @notice Primary entrypoint for procuring services for preconfirmations.
contract PreconfServiceManager is ServiceManagerBase, Pausable {
    /*//////////////////////////////////////////////////////////////
                                VARIABLES
    //////////////////////////////////////////////////////////////*/

    PreconfChallengeManager public challengeManager;

    /*//////////////////////////////////////////////////////////////
                                 ERRORS
    //////////////////////////////////////////////////////////////*/

    /// @notice Emitted when the caller is not the PreconfChallengeManager.
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

    /// @notice Prevent the implementation contract from being initialized.
    /// @dev The proxy contract state will still be able to call this function because the constructor does not affect the proxy state.
    constructor(IAVSDirectory _avsDirectory, IRegistryCoordinator _registryCoordinator, IStakeRegistry _stakeRegistry)
        ServiceManagerBase(_avsDirectory, _registryCoordinator, _stakeRegistry)
    {
        _disableInitializers();
    }

    /*//////////////////////////////////////////////////////////////
                               INITIALIZER
    //////////////////////////////////////////////////////////////*/

    /// @notice Initializes the contract.
    /// @param _pauserRegistry The PauserRegistry contract.
    /// @param _initialPausedStatus The initial paused status.
    /// @param _initialOwner The initial owner of the contract.
    /// @param _challengeManager The PreconfChallengeManager contract.
    function initialize(
        IPauserRegistry _pauserRegistry,
        uint256 _initialPausedStatus,
        address _initialOwner,
        PreconfChallengeManager _challengeManager
    ) external initializer {
        _initializePauser(_pauserRegistry, _initialPausedStatus);
        _transferOwnership(_initialOwner);

        challengeManager = _challengeManager;
    }

    /*//////////////////////////////////////////////////////////////
                            EXTERNAL FUNCTIONS
    //////////////////////////////////////////////////////////////*/

    /// @notice Called by the PreconfChallengeManager to slash an operator that violated a preconfirmation.
    /// @param operatorAddr The address of the operator to be slashed.
    function freezeOperator(address operatorAddr) external onlyChallengeManager {
        // TODO
        // slasher.freezeOperator(operatorAddr);
    }
}
