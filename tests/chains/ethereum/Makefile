include ../../docker.mk

SCRIPT_DIR ?= $(CURDIR)/scripts

CHAIN_ID0 ?= 2018
CHAIN_ID1 ?= 2019

.PHONY: docker-images
docker-images:
	make up-scaffold
	make deploy-contract
	make extract-abi
	make docker-commit
	make down-scaffold

.PHONY: deploy-contract
deploy-contract:
	$(DOCKER_COMPOSE) run contracts sh -c 'npm install &&  npx hardhat run ./scripts/deploy.js --network ibc0 && npx hardhat run ./scripts/deploy.js --network ibc1'

.PHONY: extract-abi
extract-abi:
	$(SCRIPT_DIR)/extractABI.sh

.PHONY: up-scaffold
up-scaffold:
	$(DOCKER_COMPOSE) up --build -d geth0-scaffold geth1-scaffold

.PHONY: down-scaffold
down-scaffold:
	$(DOCKER_COMPOSE) down --volumes --remove-orphans

.PHONY: docker-commit
docker-commit:
	$(SCRIPT_DIR)/docker/commitImage.sh  $(DOCKER_REPO) $(DOCKER_TAG) ethereum-geth0 ethereum-geth0-scaffold $(CHAIN_ID0)
	$(SCRIPT_DIR)/docker/commitImage.sh  $(DOCKER_REPO) $(DOCKER_TAG) ethereum-geth1 ethereum-geth1-scaffold $(CHAIN_ID1)
