# Contracts

## Background

Protocol contracts to power the AVS service with registration, challenge creation/resolution, and slashing.

### `PreconfServiceManager`

`PreconfChallengeManager` is the central contract of the protocol, which inherits `ServiceManagerBase` from the EigenLayer middleware for operator registration and deregistrating in communication with the EigenLayer core contracts. Additionally, it implements the logic to freeze/slash an operator (needs to be implemented as it depends on the latest EigenLayer core contracts). Read more about `ServiceManagerBase` [here](https://github.com/Layr-Labs/eigenlayer-middleware/blob/dev/docs/ServiceManagerBase.md).

#### `initialize`

```solidity
function initialize(
    IPauserRegistry _pauserRegistry,
    uint256 _initialPausedStatus,
    address _initialOwner,
    PreconfChallengeManager _challengeManager
) external initializer
```

The `initialize` function is used to initialize the contract with the following parameters:

* `_pauserRegistry`: The address of the `PauserRegistry` contract.
* `_initialPausedStatus`: The initial paused status of the contract.
* `_initialOwner`: The initial owner of the contract.
* `_challengeManager`: The address of the `PreconfChallengeManager` contract.

*Effects*:
* Initialization of the AVS service contract

*Requirements*:
* Can only be called once.

#### `freezeOperator`

```solidity
function freezeOperator(address operatorAddr) external onlyChallengeManager
```

Freezes/slashes an operator due to violation of a preconfirmation promise. Needs to be fully implemented.

* `operatorAddr`: The address of the operator to penalize.

*Effects*:
* Penalizes the operator by freezing/slashing their stake.

*Requirements*:
* Can only be called by the `PreconfChallengeManager` contract.

### `PreconfChallengeManager`

`PreconfChallengeManager` is a periphery contract that allows to raise and resolve challenges by integrating with Relic Protocol's [Transaction Prover](https://explorer.relicprotocol.com/mainnet/prover/0xbb1Aa14e41a67360973cca8cab17DE5200876A24). 

#### `raiseAndResolveChallenge`

```solidity
function raiseAndResolveChallenge(Preconf calldata preconf, bytes calldata preconfSigned, bytes calldata proof)
    external
```

Allows anyone to raise and resolve a challenge for a signed preconfirmation. The function takes three parameters:

* `preconf`: the preconfirmation data
* `preconfSigned`: the preconfirmation data signed by the operator/preconfer
* `proof`: the proof that the preconfirmation's transaction hash was not included in the promised block.

Using these parameters, the function verifies the proof with the prover to determine if the preconfirmation promise was violated. If the proof is valid, the contract penalizes the operator by calling the `freezeOperator` function from the `PreconfServiceManager` contract.

*Effects*:
* Raises a challenge to penalize an operator.

*Requirements*:
* Can only be called with valid proof parameters.

## Usage

Requires [Foundry](https://github.com/foundry-rs/foundry).

To build the contracts:

```sh
git clone https://github.com/cairoeth/preconfirmations.git
cd contracts
forge install
forge build
```

### Tests

In order to run unit tests, run:

```sh
forge test
```

### Coverage

To check the test coverage, run:

```sh
forge coverage
```

### Deploy

In order to deploy the contracts, you need to set the `PRIVATE_KEY` variable in the `.env` file, then run the script to deploy on the network (assuming Anvil is running):

```sh
forge script script/Preconf.s.sol:DeployPreconf --rpc-url http://127.0.0.1:8545 --broadcast
```