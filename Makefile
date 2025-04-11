OS := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
ARCH := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))


BINARY := gwolf
BUILD_DIR := ./bin

env:
	export DB_PASSWD=adminPass123
	export DB_USER=local_user

build: $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY) main.go
	@echo "Building $(BINARY) main.go"

run: $(BUILD_DIR)
	$(BUILD_DIR)/$(BINARY)

test:
	go test -short ./...
	@echo "Running Tests"

clean:
	@echo "Cleaning"
	rm -rf $(BUILD_DIR)

$(BUILD_DIR):
		@mkdir -p $@

.PHONY: