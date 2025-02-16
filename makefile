.PHONY: build build-server build-client run run-server run-client run-router stop stop-server stop-client stop-router logs-server restart-server

# Build all images
build: build-server build-client build-router

build-server-base:
	docker build -f server/server_base.Dockerfile . -t server:base
# Build server image
build-server:
	docker-compose -f server/docker-compose.yml build
# docker build -f server/Dockerfile . -t server

build-client-base:
	docker build -f client/client_base.Dockerfile . -t client:base

# Build client image
build-client:
	docker build -f client/client.Dockerfile . -t client

# Build router image
build-router:
	docker build -f router/router.Dockerfile router

setup-infra:
	sh setup_infra.sh

# Clean up images
clean:
	docker rmi server client router

clean-client:
	docker stop client || true
	docker rm client || true

# Run containers
run: run-router run-server run-client

# Run individual containers
run-server:
	docker compose -p $(name) -f server/docker-compose.yml up -d 
# docker run -d --name $(name) --cap-add NET_ADMIN --network servers server

run-client:
	docker run -d --name client --cap-add NET_ADMIN --network clients client

run-router:
	docker run -d --name router router

# Stop all containers
stop: stop-server stop-client stop-router

# Stop individual containers
stop-server:
	docker-compose -f server/docker-compose.yml down

stop-client:
	docker stop client || true
	docker rm client || true

stop-router:
	docker stop router || true
	docker rm router || true

# Show server logs
logs-server:
	docker-compose -f server/docker-compose.yml logs -f

# Restart server
restart-server: stop-server run-server