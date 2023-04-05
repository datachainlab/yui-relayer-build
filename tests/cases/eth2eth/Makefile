TRUFFLE ?= npx truffle
TEST_OPTS ?= --compile-none --migrate-none
include ../../docker.mk

.PHONY: network
network:
	TAG=${DOCKER_TAG} $(DOCKER_COMPOSE) \
		-f ../docker-compose-test.yaml \
		up -d \
		ethereum-ganache0 ethereum-ganache1

.PHONY: test
test:
	./scripts/init-rly

.PHONY: network-down
network-down:
	TAG=${DOCKER_TAG} $(DOCKER_COMPOSE) \
		-f ../docker-compose-test.yaml \
		down --volume --remove-orphans
