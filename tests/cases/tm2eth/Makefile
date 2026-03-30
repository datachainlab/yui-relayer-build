include ../../docker.mk

export USE_FAKELOST ?= yes

.PHONY: network
network:
	TAG=${DOCKER_TAG} $(DOCKER_COMPOSE) \
		-f ../docker-compose-test.yaml \
		up -d \
		tendermint-chain0-mock ethereum-geth0

.PHONY: test test-normal test-timeout
test:
	$(MAKE) test-normal test-timeout

test-normal:
	./scripts/fixture
	./scripts/init-rly
	./scripts/handshake
	./scripts/test-channel-upgrade
	./scripts/test-mockapp-packet-relay
	./scripts/test-ics20-packet-relay
	./scripts/test-service

test-timeout:
	./scripts/fixture
	./scripts/init-rly
	./scripts/handshake
	./scripts/test-mockapp-packet-timeout

.PHONY: network-down
network-down:
	TAG=${DOCKER_TAG} $(DOCKER_COMPOSE) \
		-f ../docker-compose-test.yaml \
		down -v --remove-orphans
