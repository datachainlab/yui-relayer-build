#!/usr/bin/env bash

set -eu

DOCKER_BUILD="docker build --rm --no-cache --pull"

DOCKER_REPO=$1
DOCKER_TAG=$2
DOCKER_IMAGE=$3
SCAFFOLD_CONTAINER=$4
CHAIN_ID=$5

docker cp ./contracts/addresses/${CHAIN_ID} ${SCAFFOLD_CONTAINER}:/root/addresses
docker cp ./contracts/abis ${SCAFFOLD_CONTAINER}:/root/abis
docker commit --pause=true ${SCAFFOLD_CONTAINER} ${DOCKER_REPO}${DOCKER_IMAGE}:${DOCKER_TAG}
