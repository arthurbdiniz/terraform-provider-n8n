TEST_ARGS ?= -v -cover -coverprofile=coverage.out -timeout=120s -parallel=10


default: fmt lint install generate

build:
	go build -v ./...

install: build
	go install -v ./...

lint:
	golangci-lint run

generate:
	cd tools; go generate ./...

fmt:
	gofmt -s -w -e .

coverage-html:
	go tool cover -html=coverage.out -o coverage.html

test:
	@if [ "$(ACC)" = "1" ]; then \
		echo "Running acceptance tests..."; \
		TF_ACC=1 go test $(TEST_ARGS) ./...; \
	else \
		echo "Running unit tests..."; \
		go test $(TEST_ARGS) ./...; \
	fi

.PHONY: fmt lint test build install generate
