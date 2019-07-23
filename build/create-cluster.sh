#!/usr/bin/env bash

set -e

echo "creating $1 containers..."

docker network create meerkat || true
docker run -d --network meerkat meerkat:latest
sleep 3
SEED="$(docker inspect --format '{{ .NetworkSettings.Networks.meerkat.IPAddress }}' $(docker ps -q))"

echo "using seed $SEED"
START=1
END=$1

for i in $(eval echo "{$START..$END}")
do
   echo "creating container $i"
   docker run -d --network meerkat -e SEEDS=$SEED meerkat:latest
done
