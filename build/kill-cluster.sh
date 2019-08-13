#!/usr/bin/env bash

set -e

echo "killing all container..."

for i in $(docker ps -q)
do
   echo "killing container $i"
   docker kill $i
done
