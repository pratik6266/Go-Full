ENV ?= dev
BIN_DIR := bin
APP := studentapi

.PHONY: run build run-windows test fmt vet clean help swag vendor tidy precommit db monitor db-server

# Run the app directly (cross-platform)
run:
	go run ./cmd/studentapi

# Build the binary into bin/
build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP) ./cmd/studentapi

# Windows build-and-run helper (uses .exe)
run-windows: build
	$(BIN_DIR)/$(APP).exe || $(BIN_DIR)/$(APP)

test:
	go test -v ./...

fmt:
	gofmt -s -w .

vet:
	go vet ./...

clean:
	rm -rf $(BIN_DIR)

swag:
	swag init -g cmd/studentapi/main.go -d . -o ./docs --parseDependency --parseInternal

vendor:
	go mod vendor

tidy:
	go mod tidy

precommit: tidy vendor vet fmt test

db-server:
	docker compose -f db-server.yaml up --build

monitor:
	docker compose up --build

db:
	docker compose -f db.yaml up --build

help:
	@echo "make run          # Run the API with 'go run'"
	@echo "make build        # Build binary to $(BIN_DIR)/$(APP)"
	@echo "make run-windows  # Build and run binary (tries .exe first)"
	@echo "make test         # Run tests"
	@echo "make fmt          # Format code"
	@echo "make vet          # Go vet"
	@echo "make clean        # Remove build artifacts"