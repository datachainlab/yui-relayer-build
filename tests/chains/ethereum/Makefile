include ../../docker.mk

HARDHAT ?= npx hardhat
SCRIPT_DIR ?= $(CURDIR)/scripts

.PHONY: docker-images
docker-images:
	make up-scaffold
	make deploy-contract
	make docker-commit
	make down-scaffold

.PHONY: deploy-contract
deploy-contract:
	cd contracts && npm install && $(HARDHAT) run ./scripts/deploy.js --network ibc0 && $(HARDHAT) run ./scripts/deploy.js --network ibc1

.PHONY: up-scaffold
up-scaffold:
	$(DOCKER_COMPOSE) up --build -d geth0-scaffold geth1-scaffold

.PHONY: down-scaffold
down-scaffold:
	$(DOCKER_COMPOSE) down --volumes --remove-orphans

.PHONY: docker-commit
docker-commit:
	$(SCRIPT_DIR)/docker/commitImage.sh  $(DOCKER_REPO) $(DOCKER_TAG) ethereum-geth0 ethereum-geth0-scaffold
	$(SCRIPT_DIR)/docker/commitImage.sh  $(DOCKER_REPO) $(DOCKER_TAG) ethereum-geth1 ethereum-geth1-scaffold
