#!/bin/sh

FILE="/root/.ssh"

if [ -d "$FILE" ]; then
    echo "exists pass"
else
    echo "not found"
    mkdir -p /root/.ssh
    chmod 700 /root/.ssh
    echo "" > /root/.ssh/authorized_keys
    chmod 600 /root/.ssh/authorized_keys
fi

cd /

dev_address_file="/localnet/dev_address"
dev_address_pwd_file="/localnet/dev_address_password"
dev_address=""

if [ "$(find /localnet -mindepth 1 | wc -l)" -eq 0 ]; then
    if [ ! -d "localnet" ]; then
        mkdir /localnet
    fi
    
    # 마이닝을 위한 초기 개발자 계정 생성
    (echo "windowshyun"; echo "windowshyun") | geth account new --datadir /localnet > output.txt
    dev_address=$(cat output.txt | awk '/Public address of the key:/ {print $6}')
    echo $dev_address > "$dev_address_file"
    echo "windowshyun" > "$dev_address_pwd_file"
    echo "개발자 주소 생성 완료: $dev_address"

    # Create genesis.json file
    genesis_file="/bnbsmartchain/genesis.json"
    modified_genesis_file="/bnbsmartchain/modified_genesis.json"

    old_value="f3a43bb136a7ec60c8fd6d99a36e8139d6dce117"
    new_value="${dev_address#*x}"
    sed -i "s/$old_value/$new_value/g" "$genesis_file"

    old_value="change_admin_wallet"
    new_value="$dev_address"
    sed -i "s/$old_value/$new_value/g" "$genesis_file"


    # Write the updated genesis.json file
    chmod +x "$genesis_file"

    geth --datadir /localnet init /bnbsmartchain/genesis.json
fi

if [ -f "$dev_address_file" ]; then
    # Read developer address if the file exists
    dev_address=$(cat "$dev_address_file")
    echo "개발자 주소 읽어오기 완료: $dev_address"
fi

service ssh start
cd /localnet

# Execute geth command (unlock the 0th account)
geth --http --config /bnbsmartchain/config.toml --datadir /localnet --cache 8000 --rpc.allow-unprotected-txs --txlookuplimit 0 --allow-insecure-unlock --unlock "$dev_address" --password "$dev_address_pwd_file" --mine --miner.etherbase "$dev_address"


#/bin/bash

# Open Port 8545 8546 30303 8548 30304 8575
# docker build --no-cache . -t binance-smartchain:latest
# DockerRun docker run --name bnbsmartchain -it -p 8545:8545 -p 8575:8575 -p 8546:8546 -p 30303:30303 -p 8548:8548 -p 30304:30304 -v /Users/windowshyun/Documents/DockerData/bnbsmartchain:/localnet binance-smartchain:latest