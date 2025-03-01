OS := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
ARCH := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))


BINARY := gwolf
BUILD_DIR := ./bin

build: $(BUILD_DIR)
		go build -o $(BUILD_DIR)/$(BINARY) main.go
		@echo "Building $(BINARY) main.go"

test:
		go test -short ./...
		@echo "Running Tests"

$(BUILD_DIR):
		@mkdir -p $@

.PHONY: build