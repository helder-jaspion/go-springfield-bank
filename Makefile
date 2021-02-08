NAME = go-springfield-bank
COMMAND_HANDLER ?= serverd
VERSION ?= dev
OS ?= linux

.PHONY: dev-up
dev-up:
	@echo "  > Starting dev deps..."
	docker-compose -f deployments/docker-compose-dev.yml up -d

.PHONY: dev-down
dev-down:
	@echo "  > Shutting dev deps down..."
	docker-compose -f deployments/docker-compose-dev.yml down

.PHONY: start
start:
	@echo "  > Starting..."
	docker-compose -f deployments/docker-compose.yml up -d --build

.PHONY: stop
stop:
	@echo "  > Stopping..."
	docker-compose -f deployments/docker-compose.yml down

.PHONY: setup
setup:
	@echo "  > Getting deps..."
	go mod tidy
	GO111MODULE=on go install \
	github.com/kevinburke/go-bindata/go-bindata \
	golang.org/x/tools/cmd/goimports \
	github.com/resotto/gochk/cmd/gochk \
	golang.org/x/lint/golint \
	github.com/golangci/golangci-lint/cmd/golangci-lint \
	github.com/swaggo/swag/cmd/swag

.PHONY: clean
clean:
	@echo "  >  Cleaning releases..."
	GOOS=${OS} go clean -i -x ./...
	rm -f build/${COMMAND_HANDLER}
	rm -f dist/${COMMAND_HANDLER}
	rm -f coverage.txt

.PHONY: test
test:
	@echo "  >  Running Tests..."
	go test -v -race ./...

.PHONY: compile
compile: clean
	@echo "  >  Building "${COMMAND_HANDLER}"..."
	env GOOS=${OS} GOARCH=amd64 go build -v -o dist/${COMMAND_HANDLER} cmd/${COMMAND_HANDLER}/main.go
	@echo "Binary generated at dist/"${COMMAND_HANDLER}

.PHONY: build
build: clean
	@echo "  >  Building Docker image..."
	docker build -t ${NAME}-${COMMAND_HANDLER}:${VERSION} build -f build/Dockerfile

.PHONY: generate
generate:
	@echo "  >  Generating Go files..."
	go get -u github.com/kevinburke/go-bindata/...
	go generate ./...
	swag init -g cmd/${COMMAND_HANDLER}/main.go -o api

.PHONY: lint
lint:
	@echo "  >  Running linters..."
	go get -u golang.org/x/lint/golint github.com/golangci/golangci-lint/cmd/golangci-lint
	golint ./...
	golangci-lint run ./...

.PHONY: archlint
archlint:
	@echo "  >  Running architecture linter(gochk)..."
	go get -u github.com/resotto/gochk/cmd/gochk
	gochk -c ./gochk-arch-lint.json

.PHONY: test-coverage
test-coverage:
	@echo "  >  Running tests..."
	go test -v ./... -coverprofile=coverage.out -covermode=atomic

.PHONY: go-fmt
go-fmt:
	@echo "  >  Formatting Go files..."
	go fmt ./...

.PHONY: pre-push
pre-push: go-fmt test lint archlint

