NAME = go-springfield-bank
PROJECT_PATH ?= github.com/helder-jaspion/go-springfield-bank
COMMAND_HANDLER ?= serverd
VERSION ?= dev
OS ?= linux

.PHONY: setup
setup:
	@echo "  > Getting deps tools..."
	go mod tidy
	GO111MODULE=on go install \
	github.com/resotto/gochk/cmd/gochk \
	github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: clean
clean:
	@echo "  >  Cleaning releases..."
	GOOS=${OS} go clean -i -x ./...
	rm -f build/${COMMAND_HANDLER}

.PHONY: test
test:
	@echo "  >  Running Tests..."
	go test -v ./...

.PHONY: compile
compile: clean
	@echo "  >  Building "${COMMAND_HANDLER}"..."
	env GOOS=${OS} GOARCH=amd64 go build -v -o build/${COMMAND_HANDLER} ${PROJECT_PATH}/cmd/${COMMAND_HANDLER}
	echo "Binary generated at build/"${COMMAND_HANDLER}

.PHONY: build
build: clean
	@echo "  >  Building Docker  image..."
	docker build -t ${NAME}-${COMMAND_HANDLER}:${VERSION} build -f build/Dockerfile

.PHONY: lint
lint:
	@echo "  >  Running linters..."
	golangci-lint run ./...

.PHONY: archlint
archlint:
	@echo "  >  Running architecture linter(gochk)..."
	gochk -c ./gochk-arch-lint.json

.PHONY: test-coverage
test-coverage:
	@echo "  >  Running tests..."
	go test -failfast -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: go-fmt
go-fmt:
	@echo "  >  Formatting go files..."
	go fmt ./...

.PHONY: pre-push
pre-push: go-fmt test lint archlint

