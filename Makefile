.PHONY: test build-cli
test:
	@echo "Running tests..."
	@go test -timeout 300s -p 1 -count=1 -race -cover -v ./...

build-cli:
	@echo "Building..."
	@go build -o bin/cli -v ./cmd/cli/main.go