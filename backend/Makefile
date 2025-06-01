# Project configuration
APP_NAME = Webcross
CMD_DIR = ./cmd/server
OPENAPI_SPEC = ./api/openapi.yaml
GEN_SPECS_CONFIG = ./api/specs.gen.yaml
GEN_STUBS_CONFIG = ./api/stubs.gen.yaml
GEN_SPECS_OUTPUT = ./internal/api/specs.gen.go
GEN_STUBS_OUTPUT = ./internal/api/stubs.gen.go

# Tools
OAPI_CODEGEN = oapi-codegen

# Commands
.PHONY: all build run clean fmt lint generate

all: generate build

generate: $(OPENAPI_SPEC) $(GEN_CONFIG)
	@echo "Generating OpenAPI code..."
	$(OAPI_CODEGEN) -config $(GEN_SPECS_CONFIG) $(OPENAPI_SPEC)
	$(OAPI_CODEGEN) -config $(GEN_STUBS_CONFIG) $(OPENAPI_SPEC)

build:
	@echo "Building..."
	go build -o bin/$(APP_NAME) $(CMD_DIR)

run: generate
	@echo "Running server..."
	@command -v air >/dev/null 2>&1 && air || go run ./cmd/server


clean:
	@echo "Cleaning..."
	rm -rf bin/

fmt:
	@echo "Formatting code..."
	go fmt ./...

lint:
	@echo "Linting code..."
	golangci-lint run

install-tools:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
