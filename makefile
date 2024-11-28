.PHONY: build build-server build-client run run-server run-client run-router stop stop-server stop-client stop-router logs-server restart-server

# Build all images
build: build-server build-client build-router

# Build server image
build-server:
	docker-compose -f server/docker-compose.yml build .

# Build client image
build-client:
	docker build -f client/Dockerfile .

# Build router image
build-router:
	docker build -f router/Dockerfile .

# Clean up images
clean:
	docker rmi distributed-tournaments-server distributed-tournaments-client router

# Run containers
run: run-router run-server run-client

# Run individual containers
run-server:
	docker-compose -f server/docker-compose.yml up -d

run-client:
	docker run -d --name distributed-tournaments-client distributed-tournaments-client

run-router:
	docker run -d --name router router

# Stop all containers
stop: stop-server stop-client stop-router

# Stop individual containers
stop-server:
	docker-compose -f server/docker-compose.yml down

stop-client:
	docker stop distributed-tournaments-client || true
	docker rm distributed-tournaments-client || true

stop-router:
	docker stop router || true
	docker rm router || true

# Show server logs
logs-server:
	docker-compose -f server/docker-compose.yml logs -f

# Restart server
restart-server: stop-server run-server