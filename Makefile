BINARY  := postal_server
CMD     := ./cmd/server
CLICMD     := ./cmd/http-server
BIN_DIR := bin

ADDR             ?= 9444
REGISTRATION_KEY ?= dev-registration-key

.PHONY: build run test generate install-tools clean lint

## build: compile the server binary to bin/postal_server
build:
	@mkdir -p $(BIN_DIR)
	( cd ./lib/version ; ../../bin/generate-git-commit.sh )
	go build -o $(BIN_DIR)/$(BINARY) $(CMD)

## run: run the server (set REGISTRATION_KEY env var for non-dev use)
run:
	echo REGISTRATION_KEY=$(REGISTRATION_KEY) ADDR=$(ADDR) go run $(CMD)
	./bin/postal_server --port $(ADDR)

## test: run all integration tests
test:
	go test ./tests/... -v -count=1

## generate: regenerate server stubs from openapi.yaml using oapi-codegen
generate: install-tools
	oapi-codegen -config api/oapi-codegen.yaml api/openapi.yaml > api/api.gen.go

## install-tools: install required code-generation tools
install-tools:
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

## tidy: tidy and verify module dependencies
tidy:
	go mod tidy

## lint: run golangci-lint
lint:
	golangci-lint run ./...

## clean: remove build artefacts and database
clean:
	rm -rf $(BIN_DIR)/postal_server 

# git push origin v1.0.0
git_set_tag:
	git tag v0.0.10
	git push origin --tags

.DEFAULT_GOAL := build
