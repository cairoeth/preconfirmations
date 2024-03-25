from threading import Thread
import requests
from web3 import Web3
import random

web3 = Web3(Web3.HTTPProvider("http://127.0.0.1:8545"))
private_key = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
account = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
nonce = web3.eth.get_transaction_count(account)

tx = {'nonce': nonce, 'to': account, 'value': web3.to_wei(
    random.uniform(0.1, 10), 'ether'), 'gas': 25000, 'gasPrice': web3.to_wei(50, 'gwei')}
signed_tx = web3.eth.account.sign_transaction(
    tx, private_key).rawTransaction.hex()
block = str(hex(web3.eth.block_number + 1))


def request():
    r = requests.get('http://localhost:8080', json={
        "params": [
            {
                "version": "v0.1",
                "inclusion": {
                    "desiredBlock": block,
                    "maxBlock": block,
                    "tip": block
                },
                "body": [{"tx": str(signed_tx)}],
                "privacy": {
                    "hints": [
                        "hash"
                    ]
                }
            }
        ],
        "method": "preconf_sendRequest",
        "id": 1,
        "jsonrpc": "2.0"
    })

    json_response = r.json()

    if json_response["result"]["preconfSignature"] is None:
        raise Exception("No preconfirmation")

    print(f"Status Code (request): {r.status_code}, Response: {r.json()}")

# def threaded_process_range():
#     store = {}
#     threads = []

#     # create the threads
#     threads.append(Thread(target=request))

#     # start the threads
#     [t.start() for t in threads]
#     # wait for the threads to finish
#     [t.join() for t in threads]
#     return store


# threaded_process_range()

request()
