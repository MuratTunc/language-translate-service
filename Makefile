BACK_END_SERVICE_BINARY=backendServiceApp
GOFULLPATH := /usr/local/go/bin/go

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images in the background ..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_back_end_service
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_back_end_service: builds the build_back_end_service binary as a linux executable
build_back_end_service:
	@echo "Building binary..."
	cd ../back-end-service && env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOFULLPATH) build -o ${BACK_END_SERVICE_BINARY} ./cmd/api
	@echo "Done!"