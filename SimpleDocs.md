## Install

### Golang

[https://go.dev/dl/](https://go.dev/dl/)

Golang MacOS - amd64

Set GOROOT và GOPATH (arcording to computer)

```bash
curl -o golang.pkg https://dl.google.com/go/go1.20.4.darwin-amd64.pkg
sudo open golang.pkg
```

```bash
export GOROOT=$HOME/usr/local/go
export GOPATH=$HOME/go
```

### Tomo

Create tomo folder

```bash
mkdir tomo
cd tomo
```

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
- Create account or import account (with your private key) (at least 2 account)
  - Create Keystore directory: `mkdir keystore`
  ```bash
  export PRIVATE_KEY_DIRECTORY=[DIRECTORY TO STORE PRIVATE KEY]
  export PASSWORD_DIRECTORY=[DIRECTORY TO STORE PASSWORD OF KEYSTORE FILE]
  ```
  - Create new account:
    ```bash
    touch $PASSWORD_DIRECTORY/pw1.txt
    echo [YOUR_PASSWORD_1] >> $PASSWORD_DIRECTORY/pw1.txt
    touch $PASSWORD_DIRECTORY/pw2.txt
    echo [YOUR_PASSWORD_2] >> $PASSWORD_DIRECTORY/pw2.txt
    ```
    ```bash
    tomo account new \
          --password $PASSWORD_DIRECTORY/pw1.txt \
          --keystore $PWD/keystore/1
    ```
    ```bash
    tomo account new \
          --password $PASSWORD_DIRECTORY/pw2.txt \
          --keystore $PWD/keystore/2
    ```
  - Import account:
    ```bash
    touch $PRIVATE_KEY_DIRECTORY/pk1.txt
    echo [YOUR_PRIVATE_KEY_1] >> $PRIVATE_KEY_DIRECTORY/pk1.txt
    touch $PRIVATE_KEY_DIRECTORY/pk2.txt
    echo [YOUR_PRIVATE_KEY_2] >> $PRIVATE_KEY_DIRECTORY/pk2.txt
    ```
    ```bash
    tomo  account import $PRIVATE_KEY_DIRECTORY/pk1.txt \
        --keystore $PWD/keystore/1 \
        --password $PRIVATE_KEY_DIRECTORY/pk1.txt
    ```
    ```bash
    tomo  account import $PRIVATE_KEY_DIRECTORY/pk2.txt \
        --keystore $PWD/keystore/2 \
        --password $PRIVATE_KEY_DIRECTORY/pk2.txt
    ```

## Create genesis file with `puppeth`

- Run `puppeth`
  ```bash
  puppeth
  ```
  - Set chain name: `c98chain`
  - Configure new genesis: `2`
  - Select `POSV` consensus: `3`
  - Set blocktime (default 2 seconds): `Enter`
  - Set reward of each epoch: `250`
  - Set addresses to be initial masternodes: Address created before
  - Set account to seal: Account 1
  - Set number of blocks of each epoch (default 900): `Enter`
  - Set gap: `5`
  - Set foundation wallet address: `Enter`
  - Account confirm Foundation MultiSignWallet: Account 1 & 2
  - Require for confirm tx in Foudation MultiSignWallet: `2`
  - Account confirm Team MultiSignWallet: Account 1 & 2
  - Require for confirm tx in Team MultiSignWallet: `2`
  - Enter swap wallet address for fund 55 million TOMO: Account 1
  - Enter account be pre-funded: Account 2
  - Enter network ID: `3172`
- Export genesis file
  - Select `2. Manage existing genesis`
  - Select `2. Export genesis configuration`
  - Enter genesis filename: `genesis.json`
- `Ctrl + C` to end

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

  Example: `"enode://372853cfc9cc509bdd79db961cf791e8b2c8fdbadd5b4a25b0e59187f3be9a6e1d26e381f8ed4ae71d81c72ad7f53430af605955293df66660232ad235633880@[::]:30301"`

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
  --mine --gasprice 2500 \ --bootnodesv5 "enode://372853cfc9cc509bdd79db961cf791e8b2c8fdbadd5b4a25b0e59187f3be9a6e1d26e381f8ed4ae71d81c72ad7f53430af605955293df66660232ad235633880@[::]:30301" \
  console
  ```

- Node 1 can Commit and seal block

## Connect node to sync and execute

- Open IPC node 1
  ```bash
  tomo attach nodes/1/tomo.ipc
  ```

→ Successful 3 master nodes network
