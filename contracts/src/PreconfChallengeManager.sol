// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

import {Initializable} from "openzeppelin-upgrades/proxy/utils/Initializable.sol";
import {OwnableUpgradeable} from "openzeppelin-upgrades/access/OwnableUpgradeable.sol";

import {ECDSA} from "openzeppelin/utils/cryptography/ECDSA.sol";
import {EIP712Upgradeable} from "openzeppelin-upgrades/utils/cryptography/draft-EIP712Upgradeable.sol";

import {IProver} from "relic-sdk/packages/contracts/interfaces/IProver.sol";
import {Fact} from "relic-sdk/packages/contracts/lib/Facts.sol";

import {PreconfServiceManager} from "src/PreconfServiceManager.sol";

/// @title PreconfChallengeManager
/// @author @cairoeth
/// @notice Manages challenges for preconfirmations.
contract PreconfChallengeManager is Initializable, OwnableUpgradeable, EIP712Upgradeable {
    /*//////////////////////////////////////////////////////////////
                                VARIABLES
    //////////////////////////////////////////////////////////////*/

    /// @notice The PreconfServiceManager contract address.
    address public serviceManager;

    /// @notice The Transaction Prover address from Relic.
    address public prover;

    /// @notice Number of blocks that the challenger has to raise and resolve a challenge.
    uint32 public constant CHALLENGE_WINDOW_BLOCK = 100;

    /// @notice Mapping of hashed preconfirmations that have been challenged.
    mapping(bytes32 => bool) public challenges;

    /// @notice The preconfirmation struct.
    /// @param hashes The hashes of the transactions to be included.
    /// @param blockTarget The block target to include the transactions.
    struct Preconf {
        bytes32[] hashes;
        uint256 blockTarget;
    }

    /*//////////////////////////////////////////////////////////////
                                 EVENTS
    //////////////////////////////////////////////////////////////*/

    event Challenged(Preconf preconf, address challenger);

    /*//////////////////////////////////////////////////////////////
                                 ERRORS
    //////////////////////////////////////////////////////////////*/

    /// @notice Emitted when the block target of the preconfirmation has not been reached yet.
    error NotReached();

    /// @notice Emitted when the bundle contains no transaction hashes.
    error EmptyBundle();

    /// @notice Emitted when the preconfirmation has already been challenged.
    error AlreadyChallenged();

    /// @notice Emitted when the challenge window has expired.
    error WindowExpired();

    /*//////////////////////////////////////////////////////////////
                               CONSTRUCTOR
    //////////////////////////////////////////////////////////////*/

    /// @notice Prevent the implementation contract from being initialized.
    /// @dev The proxy contract state will still be able to call this function because the constructor does not affect the proxy state.
    constructor() {
        _disableInitializers();
    }

    /*//////////////////////////////////////////////////////////////
                               INITIALIZER
    //////////////////////////////////////////////////////////////*/

    /// @notice Initializes the contract.
    /// @param owner The owner of the contract.
    /// @param _serviceManager The address of the PreconfServiceManager contract.
    /// @param _prover The address of the transaction prover.
    function initialize(address owner, address _serviceManager, address _prover) external initializer {
        __Ownable_init();
        __EIP712_init("PreconfChallengeManager", "0.1.0");
        _transferOwnership(owner);

        serviceManager = _serviceManager;
        prover = _prover;
    }

    /*//////////////////////////////////////////////////////////////
                            EXTERNAL FUNCTIONS
    //////////////////////////////////////////////////////////////*/

    /// @notice Raises and resolves a challenge for a signed preconfirmation.
    /// @param preconf The preconfirmation data.
    /// @param preconfSigned The signed preconfirmation data.
    /// @param proof The proof of the preconfirmation.
    function raiseAndResolveChallenge(Preconf calldata preconf, bytes calldata preconfSigned, bytes calldata proof)
        external
    {
        if (block.number < preconf.blockTarget) revert NotReached();
        if (block.number > preconf.blockTarget + CHALLENGE_WINDOW_BLOCK) revert WindowExpired();
        if (preconf.hashes.length == 0) revert EmptyBundle();

        // Get the hash of the preconfirmation.
        bytes32 preconfHash = keccak256(abi.encode(preconf));

        if (challenges[preconfHash]) revert AlreadyChallenged();

        // TODO: require challenger to pay a bond?

        // Check proof with prover.
        Fact memory fact = IProver(prover).prove(proof, false);

        // Check the response. True if transaction was included in block, false otherwise.
        bool response = abi.decode(fact.data, (bool));

        // If preconfirmation was violated, slash the operator.
        if (response != true) {
            // Get address of operator from the signed preconfirmation.
            address operator = ECDSA.recover(preconfHash, preconfSigned);

            // Slash via the PreconfServiceManager.
            PreconfServiceManager(serviceManager).freezeOperator(operator);
        }

        challenges[preconfHash] = true;

        emit Challenged(preconf, msg.sender);
    }
}
