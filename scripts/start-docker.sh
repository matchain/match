#!/bin/bash

KEY="dev0"
CHAINID="match_699-1"
MONIKER="mymoniker"
DATA_DIR=$(mktemp -d -t match-datadir.XXXXX)

echo "create and add new keys"
./matchd keys add $KEY --home $DATA_DIR --no-backup --chain-id $CHAINID --algo "eth_secp256k1" --keyring-backend test
echo "init Match with moniker=$MONIKER and chain-id=$CHAINID"
./matchd init $MONIKER --chain-id $CHAINID --home $DATA_DIR
echo "prepare genesis: Allocate genesis accounts"
./matchd add-genesis-account \
"$(./matchd keys show $KEY -a --home $DATA_DIR --keyring-backend test)" 1000000000000000000amatch,1000000000000000000stake \
--home $DATA_DIR --keyring-backend test
echo "prepare genesis: Sign genesis transaction"
./matchd gentx $KEY 1000000000000000000stake --keyring-backend test --home $DATA_DIR --keyring-backend test --chain-id $CHAINID
echo "prepare genesis: Collect genesis tx"
./matchd collect-gentxs --home $DATA_DIR
echo "prepare genesis: Run validate-genesis to ensure everything worked and that the genesis file is setup correctly"
./matchd validate-genesis --home $DATA_DIR

echo "starting match node $i in background ..."
./matchd start --pruning=nothing --rpc.unsafe \
--keyring-backend test --home $DATA_DIR \
>$DATA_DIR/node.log 2>&1 & disown

echo "started match node"
tail -f /dev/null