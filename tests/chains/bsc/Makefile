.PHONY:docker-images
docker-images:
	docker compose -f docker-compose.image.yml build bsc-rpc
	docker compose -f docker-compose.image.yml build bsc-bootstrap

.PHONY:bootstrap
bootstrap:
	docker compose -f docker-compose.simple.yml run --rm bootstrap-simple
	docker compose -f docker-compose.simple.yml run --rm bootstrap-simple2

.PHONY:network
network:
	docker compose -f docker-compose.simple.yml up -d bsc-rpc bsc-validator1-1
	docker compose -f docker-compose.simple.yml up -d bsc-rpc2 bsc-validator2-1

.PHONY: deploy-contract
deploy-contract:
	cd contract && npm install && npm run migrate

.PHONY:network-down
network-down:
	docker compose -f docker-compose.simple.yml down -v
