#!/bin/bash
set -a
[ -f .env ] && . .env
set +a

# docker run -v ~/.arbitrum/arb1:/home/user/.arbitrum/arb1 -v /tmp/.arbitrum:/tmp/.arbitrum -p 0.0.0.0:8547:8547 -p 0.0.0.0:8548:8548 nitro-node --parent-chain.connection.url=$L1_RPC_URL --parent-chain.blob-client.beacon-url=$BEACON_URL --chain.id=42161 --chain.name=arb1 --init.latest=genesis  --http.api=net,web3,eth --http.corsdomain=* --http.addr=0.0.0.0 --http.vhosts=* 
# ./target/bin/nitro --parent-chain.connection.url=$L1_RPC_URL --parent-chain.blob-client.beacon-url=$BEACON_URL --chain.id=42161 --chain.name=arb1 --init.latest=genesis  --http.api=net,web3,eth --http.corsdomain=* --http.addr=0.0.0.0 --http.vhosts=* 

docker run -v ~/.arbitrum/arb1:/home/user/.arbitrum/arb1 -v /tmp/.arbitrum:/tmp/.arbitrum -p 0.0.0.0:8547:8547 -p 0.0.0.0:8548:8548 nitro-node --parent-chain.connection.url=$L1_RPC_URL --parent-chain.blob-client.beacon-url=$BEACON_URL --chain.id=42161 --chain.name=arb1 --init.latest=genesis  --http.api=net,web3,eth,lightclient,debug --http.corsdomain=* --http.addr=0.0.0.0 --http.vhosts=* 