NAME = go-springfield-bank
PROJECT_PATH ?= github.com/helder-jaspion/go-springfield-bank
COMMAND_HANDLER ?= serverd
VERSION ?= dev
OS ?= linux

.PHONY: dev-up
dev-up:
	@echo "Starting dev deps"
	docker-compose -f deployments/docker-compose-dev.yml up -d

.PHONY: dev-down
dev-down:
	@echo "Shutting dev deps down"
	docker-compose -f deployments/docker-compose-dev.yml down

.PHONY: start
start:
	@echo "Starting"
	docker-compose -f deployments/docker-compose.yml up -d --build

.PHONY: stop
stop:
	@echo "Stopping"
	docker-compose -f deployments/docker-compose.yml down

.PHONY: setup
setup:
	@go mod download
	@echo "Installing tools"
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/rakyll/gotest@latest
	@go install github.com/resotto/gochk/cmd/gochk@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go mod tidy

.PHONY: clean
clean:
	@echo "Cleaning releases"
	GOOS=${OS} go clean -i -x ./...
	rm -f build/${COMMAND_HANDLER}
	rm -f dist/${COMMAND_HANDLER}
	rm -f coverage.txt

.PHONY: test
test:
	@echo "Running tests"
	gotest -race -failfast -v ./...

.PHONY: test-coverage
test-coverage:
	@echo "Running coverage tests"
	gotest -race -failfast -coverprofile=coverage.out -covermode=atomic ./...

.PHONY: compile
compile: clean
	@echo "Building "${COMMAND_HANDLER}
	env GOOS=${OS} GOARCH=amd64 go build -v -o dist/${COMMAND_HANDLER} cmd/${COMMAND_HANDLER}/main.go
	@echo "Binary generated at dist/"${COMMAND_HANDLER}

.PHONY: build
build: clean
	@echo "Building Docker image"
	docker build -t ${NAME}-${COMMAND_HANDLER}:${VERSION} build -f build/Dockerfile

.PHONY: goformat
goformat:
	go mod tidy
	gci -local ${PROJECT_PATH} -w .
	gofumpt -w -extra .
	go fmt ./...

.PHONY: generate
generate:
	@echo "Generating Go files"
	go generate ./...
	swag init -g cmd/${COMMAND_HANDLER}/main.go -o api
	$(MAKE) goformat

.PHONY: lint
lint:
	@echo "Running golangci-lint"
	golangci-lint run ./...

.PHONY: archlint
archlint:
	@echo "Running architecture linter(gochk)"
	gochk -e -c ./gochk-arch-lint.json

.PHONY: pre-push
pre-push: lint archlint test

