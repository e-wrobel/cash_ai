APP_NAME = processor
PKG = ./cmd/$(APP_NAME)
BIN_DIR = bin
BIN_PATH = $(BIN_DIR)/$(APP_NAME)

MOCK_PROVIDER_PATH = ./mock-provider/main.go

.PHONY: run build test mock clean

run: ## Run the processor service
	go run $(PKG)

build: ## Build the binary
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_PATH) $(PKG)

test: ## Run tests
	go test ./...

mock: ## Run the mock transaction provider
	go run $(MOCK_PROVIDER_PATH)

clean: ## Remove generated binaries
	rm -rf $(BIN_DIR)

help: ## Show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | awk 'BEGIN {FS = \":.*?## \"}; {printf \"\033[36m%-12s\033[0m %s\\n\", $$1, $$2}'