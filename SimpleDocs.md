## Install

### Golang

[https://go.dev/dl/](https://go.dev/dl/)

Golang MacOS - amd64

Set GOROOT và GOPATH (tuỳ máy)

`curl -o golang.pkg https://dl.google.com/go/go1.20.4.darwin-amd64.pkg`
`sudo open golang.pkg`

```bash
export GOROOT=$HOME/usr/local/go
export GOPATH=$HOME/go
```

### Tomo

Create tomo folder
`mkdir tomo`
`cd tomo`

Install Tomo and library

```bash
git clone https://github.com/tomochain/tomochain/
cd tomochain
go mod tidy -e
make all
cd ..
```

Alias for tomo:

```bash
alias tomo=$PWD/tomo/build/bin/tomo
alias puppeth=$PWD/tomo/build/bin/puppeth
alias bootnode=$PWD/tomo/build/bin/bootnode
```

## Setup node and account

- Create nodes file
  ```bash
  mkdir nodes
  cd nodes
  mkdir 1
  mkdir 2
  mkdir 3
  cd ..
  ```
- Create account or import account (with your private key)
  ```bash
  mkdir keystore
  touch [YOUR_PASSWORD_FILE_TO_LOCK_YOUR_ACCOUNT]
  echo "abc" >> [YOUR_PASSWORD_FILE_TO_LOCK_YOUR_ACCOUNT]
  ```
  - Create new account:
    ```bash
    tomo account new \
          --password [YOUR_PASSWORD_FILE_TO_LOCK_YOUR_ACCOUNT] \
          --keystore [YOUR_KEYSTORE_FILE_TO_STORE_YOUR_ACCOUNT]
    ```
  - Import account:
    ```bash
    tomo  account import [PRIVATE_KEY_FILE_OF_YOUR_ACCOUNT] \
        --keystore [YOUR_KEYSTORE_FILE_TO_STORE_YOUR_ACCOUNT] \
        --password [YOUR_PASSWORD_FILE_TO_LOCK_YOUR_ACCOUNT]
    ```

## Genesis file

    ```bash
    curl https://gist.githubusercontent.com/terryyyz/968c7d336615a207996b32ea57be120e/raw/b75106da76d2bc8dca898745de3351389d9b4ef3/tomo_genesis.json >> genesis.json
    ```

## Init node with genesis file

```bash
tomo --datadir nodes/1 init genesis.json
tomo --datadir nodes/2 init genesis.json
tomo --datadir nodes/3 init genesis.json
```

## Setup bootnode

- Init bootnode key
  ```bash
  bootnode -genkey bootnode.key
  ```
- Start bootnode

  ```bash
  bootnode -nodekey ./bootnode.key
  ```

  Get bootnode info

  Example: `"enode://d2bb804ef44d29fa98a422d2cebaded916641f6fc78cb8f5bb666748ac7c22cc8019b7f4ce19aac76b89d9943686d1cebd34fe2230063fa1ffdb82ce5b939bb5@[::]:30301"`

## Run node

`YOUR_ACCOUNT_ADDRESS` example: `"0x79d3620f9379d043eaea262f1cac689fc906d5a1"`

- Node 1
  ```bash
  tomo --syncmode "full" \
  --datadir nodes/1 --networkid 3172 --port 10303 \
  --keystore keystore/1 --password pw.json \
  --rpc --rpccorsdomain "*" --rpcaddr 0.0.0.0 --rpcport 1545 --rpcvhosts "*" \
  --rpcapi "admin,db,eth,net,web3,personal,debug" \
  --gcmode "archive" \
  --ws --wsaddr 0.0.0.0 --wsport 1546 --wsorigins "*" --unlock [YOUR_ACCOUNT_ADDRESS] \
  --identity "NODE1" \
  --mine --gasprice 2500 \ --bootnodesv5 [BOOTNODE_INFO] \
  console
  ```
- Node 2
  ```bash
  tomo --syncmode "full" \
  --datadir nodes/2 --networkid 3172 --port 20303 --nodiscover \
  --keystore keystore/2 --password pw.json \
  --rpc --rpccorsdomain "*" --rpcaddr 0.0.0.0 --rpcport 2545 --rpcvhosts "*" \
  --rpcapi "admin,db,eth,net,web3,personal,debug" \
  --gcmode "archive" \
  --ws --wsaddr 0.0.0.0 --wsport 2546 --wsorigins "*" --unlock [YOUR_ACCOUNT_ADDRESS] \
  --identity "NODE2" \
  --mine --gasprice 2500 \ --bootnodesv5 [BOOTNODE_INFO] \
  console
  ```
- Node 3
  ```bash
  tomo --syncmode "full" \
  --datadir nodes/3 --networkid 3172 --port 30303 --nodiscover \
  --keystore keystore/3 --password pw.json \
  --rpc --rpccorsdomain "*" --rpcaddr 0.0.0.0 --rpcport 3545 --rpcvhosts "*" \
  --rpcapi "admin,db,eth,net,web3,personal,debug" \
  --gcmode "archive" \
  --ws --wsaddr 0.0.0.0 --wsport 3546 --wsorigins "*" --unlock [YOUR_ACCOUNT_ADDRESS] \
  --identity "NODE3" \
  --mine --gasprice 2500 \ --bootnodesv5 [BOOTNODE_INFO] \
  console
  ```

## Connect node to sync and execute

- Open IPC node 1
  ```bash
  tomo attach nodes/1/tomo.ipc
  ```
- Get Node information → Get encode
  ```bash
  admin.nodeInfo
  ```
  ![Untitled](docs-image/Untitled%2012.png)
- Open IPC of node 2:
  ```bash
  tomo attach nodes/2/tomo.ipc
  ```
- Add node 1 to peers
  ```bash
  admin.addPeer("enode://3499d673e80770534444f3ef51db4bd13e1778f3c31d4f295d04234fa06acc33b8deabc588791ebebdaee4fe7cf750167e40e825dce7639c61cae4237e636ee0@[::]:10303")
  ```
- Check peers:
  ```bash
  admin.peers
  ```
  If array isn't empty, we're successful
  ![Untitled](docs-image/Untitled%2013.png)

→ Khi đó node 1 và node 2 đã kết nối sync thành công và có thể commit block

- Làm tương tự đưa node 3 vào mạng.

→ Ta được mạng với 3 master node
