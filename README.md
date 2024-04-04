# 🔌 Preconfirmations

[![All protocol components](https://github.com/cairoeth/preconfirmations/actions/workflows/protocol.yaml/badge.svg)](https://github.com/cairoeth/preconfirmations/actions/workflows/protocol.yaml)
[![Goreport status](https://goreportcard.com/badge/github.com/flashbots/mev-boost)](https://goreportcard.com/report/github.com/flashbots/mev-boost)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL_v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

Preconfirmation protocol allowing users to get sub-second transaction confirmations on Ethereum.

> This is **experimental software** and is provided on an "as is" and "as available" basis. We **do not give any warranties** and **will not be liable for any losses** incurred through any use of this code base.

## 🌱 Background

Ethereum preconfirmations, or "preconfs" for short, are a proposed mechanism to enable faster transaction confirmation and improve the user experience on Ethereum, and subsequently, layer 2 rollups and validiums. The key idea is to have Ethereum block proposers, known as "preconfers," issue signed promises to users guaranteeing that their transactions will be included and executed within a certain timeframe. Users pay a tip to preconfers for this service. By acquiring preconf promises from upcoming block proposers, users can get assurance of speedy execution, with latencies as low as 100ms.

This project leverages [EigenLayer](https://www.eigenlayer.xyz/) as the key piece of infrastructure to secure preconfer guarantees via a slashing mechanism. Operators that run the AVS software (`preconf-operator`) on top of EigenLayer must be block proposers that can include preconfirmation transactions. If preconfirmation promises by operators are not fulfilled, they are subject to on-chain slashing with proofs, verified with [Relic Protocol](https://relicprotocol.com/).

In order to match users with preconfers, the `preconf-share` middleware acts as a matchmaker, receiving user requests with hints that are gossiped to operators. Operators can decide to promise a preconfirmation or not based on hints and tip. Promises are ranked by `preconf-share` according to the user's validation preferences (desired block, etc).

As `preconf-share` requires users to pass arguments and extra parameters, `rpc` is a JSON-RPC wrapper that simplifies this interaction for usage on wallets (like `Metamask`).

Read more [here](https://ethresear.ch/t/towards-an-implementation-of-based-preconfirmations-leveraging-restaking/19211).

## 🧱 Structure

```ml
preconf-share — "Matchmaker middleware that connects user requests with block proposers"
├─ preconshare — "Main backend to receive requests, gossip hints, rank preconfirmations, and publish signed preconfirmations"
preconf-operator — "AVS (Actively Validated Services) infrastructure that receives and reponds to preconfirmation requests"
├─ core — "Configuration and utilities to interact with on-chain smart contracts"
├─ receiverapi — "Endpoint for preconf-share callback to receive and include raw transactions"
rpc — "JSON-RPC wrapper on top of preconf-share for wallets"
├─ cmd — "Launcher of components to run the RPC server and decode transactions"
├─ server — "Logic for request handlers, processors, and forwarders"
contracts
├─ PreconfChallengeManager — "Manager that allows to raise and resolves challenges for signed preconfirmations"
├─ PreconfServiceManager — "Primary entrypoint for procuring services for preconfirmations"
```

## 👷 Setup

Requires [Go 1.18+](https://go.dev/doc/install), [Docker](https://www.docker.com/products/docker-desktop/) and [Foundry](https://getfoundry.sh/).

```bash
git clone https://github.com/cairoeth/preconfirmations.git

cd preconfirmations

make build-all
```

## 🔧 Usage

To run the protocol, you must start the following components:
- `anvil`: Local network with EigenLayer and preconfirmation contracts.
- `preconf-share`: Middleware software. Requires network.
- `preconf-operator`: AVS software. Requires network and `preconf-share`.
- `rpc`: JSON-RPC wrapper (optional). Requires `preconf-share`.

You can run all the protocol components concurrently with:

```bash
make run-all
```

Alternatively, you can run each component individually:

### Anvil

Anvil (part of Foundry), creates a local testnet node for deploying and testing smart contracts. You can start it with the [state](contracts/anvil-state.json) that already contains the required contracts, or manually deploy them using the [script](contracts/script/Preconf.s.sol).

To run with the pre-defined state:

```bash
make run-anvil
```

### Preconf-Share

Preconf-Share is the middleware that connects users with block proposers. It receives user requests, gossips hints, ranks preconfirmations, and publishes signed preconfirmations. You can read more here.

```bash
make run-share
```

### Preconf-Operator

Preconf-Operator is the AVS infrastructure that receives and responds to preconfirmation requests. It interacts with on-chain smart contracts to provide preconfirmation services. You can read more here.

```bash
make run-operator
```

### RPC

RPC is a JSON-RPC wrapper on top of Preconf-Share for wallets. It simplifies the interaction for users to pass arguments and extra parameters.

```bash
make run-rpc
```

## 🧪 Test

After running all the components, make sure to wait around 10 seconds for the AVS sofware to register the operator. Once everything is ready, you can test the protocol with the JSON-RPC wrapper, or by running the [Python example](test_tx.py). You must have Python installed:

```
# Install dependencies
pip install -r requirements.txt

python test_tx.py
```

This example generates a random transaction that transfers Ether to itself. Check the protocol logs to see how the transaction is preconfirmed and included in the next block.

## 🙏🏼 Acknowledgements

This repository is inspired by or directly modified from many sources, primarily:

- [rpc-endpoint](https://github.com/flashbots/rpc-endpoint): `rpc` is a fork of Flashbot's RPC endpoint.
- [mev-share-node](https://github.com/flashbots/mev-share-node): `preconf-share` is a fork of Flashbot's MEV-Share node.
- [Based preconfirmations](https://ethresear.ch/t/based-preconfirmations/17353)
- [Analyzing BFT & Proposer-Promised Preconfirmations](https://ethresear.ch/t/analyzing-bft-proposer-promised-preconfirmations/17963)

## 🫡 Contributing

Check out the [Contribution Guidelines](./CONTRIBUTING.md)!
