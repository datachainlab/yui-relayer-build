#!/usr/bin/env bash

workspace=$(
	cd $(dirname $0)
	pwd
)/..

function prepare() {
	if ! [[ -f /usr/local/bin/geth ]]; then
		echo "geth do not exist!"
		exit 1
	fi
	# run `make clean` to remove files
	cd ${workspace}/genesis
	rm -rf validators.conf
}

function init_validator() {
	node_id=$1

  # set validator address
	mkdir -p ${workspace}/storage/${node_id}/keystore
	cp ${workspace}/validators/keystore/${node_id} ${workspace}/storage/${node_id}/keystore/${node_id}
  validatorAddr=0x$(cat ${workspace}/storage/${node_id}/keystore/${node_id} | jq .address | sed 's/"//g')
	echo ${validatorAddr} >${workspace}/storage/${node_id}/address

  # import BLS vote address
	echo password01 > ${workspace}/storage/${node_id}/blspassword.txt
	geth --datadir ${workspace}/storage/${node_id} bls account import ${workspace}/validators/bls/${node_id} --blspassword ${workspace}/storage/${node_id}/blspassword.txt --blsaccountpassword ${workspace}/storage/${node_id}/blspassword.txt
	geth --datadir ${workspace}/storage/${node_id} bls account list --blspassword ${workspace}/storage/${node_id}/blspassword.txt
  voteAddr=0x$(cat ${workspace}/validators/bls/${node_id} | jq .pubkey | sed 's/"//g')

	echo "${validatorAddr},${validatorAddr},${validatorAddr},0x0000000010000000,${voteAddr}" >>${workspace}/genesis/validators.conf
}

function generate_genesis() {
	INIT_HOLDER_ADDRESSES=$(ls ${workspace}/init-holders | tr '\n' ',')
	INIT_HOLDER_ADDRESSES=${INIT_HOLDER_ADDRESSES/%,/}
	echo "blocks per epoch = ${BLOCKS_PER_EPOCH}"
	sed "s/{{BLOCKS_PER_EPOCH}}/${BLOCKS_PER_EPOCH}/g" ${workspace}/genesis/genesis-template.template >${workspace}/genesis/genesis-template.json
	sed "s/{{INIT_HOLDER_ADDRESSES}}/${INIT_HOLDER_ADDRESSES}/g" ${workspace}/genesis/init_holders.template | sed "s/{{INIT_HOLDER_BALANCE}}/${INIT_HOLDER_BALANCE}/g" >${workspace}/genesis/init_holders.js
	node generate-validator.js
	chainIDHex=$(printf '%04x\n' ${BSC_CHAIN_ID})
	node generate-genesis.js --chainid ${BSC_CHAIN_ID} --bscChainId ${chainIDHex}
}

function init_genesis_data() {
	node_type=$1
	node_id=$2
	geth --datadir ${workspace}/storage/${node_id} init ${workspace}/genesis/genesis.json
	cp ${workspace}/config/config-${node_type}.toml ${workspace}/storage/${node_id}/config.toml
	sed -i -e "s/{{NetworkId}}/${BSC_CHAIN_ID}/g" ${workspace}/storage/${node_id}/config.toml
	if [ "${node_id}" == "bsc-rpc" ]; then
		cp ${workspace}/init-holders/* ${workspace}/storage/${node_id}/keystore
		cp ${workspace}/genesis/genesis.json ${workspace}/storage/${node_id}
		cp ${workspace}/config/bootstrap.key ${workspace}/storage/${node_id}/geth/nodekey
	fi
}

prepare

# First, generate config for each validator
for ((i = 1; i <= ${NUMS_OF_VALIDATOR}; i++)); do
	init_validator "bsc-validator${i}"
  generate_validator_conf "bsc-validator${i}"
done

# Then, use validator configs to generate genesis file
generate_genesis

# Finally, use genesis file to init cluster data
init_genesis_data bsc-rpc bsc-rpc

for ((i = 1; i <= ${NUMS_OF_VALIDATOR}; i++)); do
	init_genesis_data validator "bsc-validator${i}"
done
