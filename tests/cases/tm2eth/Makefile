include ../../docker.mk

export USE_FAKELOST ?= yes

.PHONY: network
network:
	TAG=${DOCKER_TAG} $(DOCKER_COMPOSE) \
		-f ../docker-compose-test.yaml \
		up -d \
		tendermint-chain0-mock ethereum-geth0

.PHONY: test
test:
	./scripts/fixture
	./scripts/init-rly
	./scripts/handshake
	./scripts/test-channel-upgrade
	./scripts/test-mockapp-packet-relay
	./scripts/test-ics20-packet-relay
	./scripts/test-mockapp-packet-timeout
	./scripts/test-service

.PHONY: network-down
network-down:
	TAG=${DOCKER_TAG} $(DOCKER_COMPOSE) \
		-f ../docker-compose-test.yaml \
		down -v --remove-orphans
