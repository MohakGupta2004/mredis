BINARY_NAME=redis-server
BUILD_DIR=bin

.PHONY: all
all: clean build-server

# Build the Go binary
.PHONY: build-server
build-server:
	@echo "Building the binary..."
	@mkdir -p $(BUILD_DIR)
	cd server && go build -o ../$(BUILD_DIR)/$(BINARY_NAME) ./cmd