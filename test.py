import requests
from web3 import Web3
import redis


web3 = Web3(Web3.HTTPProvider("http://127.0.0.1:8545"))
private_key = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
account = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
nonce = web3.eth.get_transaction_count(account)

tx = {'nonce': nonce, 'to': account, 'value': web3.to_wei(10, 'ether'), 'gas': 25000, 'gasPrice': web3.to_wei(50, 'gwei')}
signed_tx = web3.eth.account.sign_transaction(tx, private_key).rawTransaction.hex()
block = str(hex(web3.eth.block_number + 1))
print(block)


red = redis.StrictRedis('localhost', 6379, charset="utf-8", decode_responses=True)

def user_counter():
    sub = red.pubsub()
    sub.subscribe('hints')
    for message in sub.listen():
         print(message)
         if message is not None and isinstance(message, dict):
            data = message.get('data')
            print(data)

while True:
    count = 0
    if count == 0:
        r = requests.post('http://localhost:8080', json={
            "params": [
                {
                    "version": "v0.1",
                    "inclusion": {
                        "block": block,
                        "maxBlock": block
                    },
                    "body": [
                        {
                            "tx": str(signed_tx),
                            "canRevert": False
                        }
                    ],
                    "validity": {
                        "refund": [],
                        "refundConfig": []
                    },
                    "privacy": {
                        "hints": [
                            # "calldata"
                            # "contract_address",
                            # "logs",
                            # "function_selector",
                            "hash"
                            # "tx_hash"
                        ]
                    }
                }
            ],
            "method": "mev_sendBundle",
            "id": 1,
            "jsonrpc": "2.0"
        })

        print(f"Status Code: {r.status_code}, Response: {r.json()}")

        count += 1

    user_counter()