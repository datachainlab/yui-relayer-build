include ../../docker.mk

.PHONY: network
network:
	TAG=${DOCKER_TAG} $(DOCKER_COMPOSE) \
		-f ../docker-compose-test.yaml \
		up -d \
		ethereum-geth0 ethereum-geth1

.PHONY: test
test:
	$(MAKE) test-single
	$(MAKE) test-single-ordered
	$(MAKE) test-multi

.PHONY: test-multi
test-multi:
	./scripts/fixture multi
	./scripts/init-rly
	./scripts/handshake
	./scripts/test-channel-upgrade multi
	./scripts/test-tx

.PHONY: test-single
test-single:
	./scripts/fixture single
	./scripts/init-rly
	./scripts/handshake
	./scripts/test-channel-upgrade single
	./scripts/test-tx
	./scripts/test-timeout

.PHONY: test-single-ordered
test-single-ordered:
	./scripts/fixture single
	./scripts/init-rly path-ordered.json
	./scripts/handshake
	./scripts/test-tx
	./scripts/test-timeout

.PHONY: network-down
network-down:
	TAG=${DOCKER_TAG} $(DOCKER_COMPOSE) \
		-f ../docker-compose-test.yaml \
		down --volume --remove-orphans
