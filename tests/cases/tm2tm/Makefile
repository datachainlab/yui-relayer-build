include ../../docker.mk

.PHONY: network
network:
	TAG=${DOCKER_TAG} $(DOCKER_COMPOSE) \
		-f ../docker-compose-test.yaml \
		up -d \
		tendermint-chain0 tendermint-chain1

.PHONY: test test-unordered test-ordered
test:
	$(MAKE) test-unordered test-ordered

test-unordered:
	./scripts/fixture
	./scripts/init-rly
	./scripts/handshake
	./scripts/test-tx
	./scripts/test-timeout transfer

test-ordered:
	./scripts/fixture
	./scripts/init-rly
	./scripts/handshake path-ordered.json
	./scripts/test-timeout mockapp

.PHONY: network-down
network-down:
	TAG=${DOCKER_TAG} $(DOCKER_COMPOSE) \
		-f ../docker-compose-test.yaml \
		down --volume --remove-orphans
