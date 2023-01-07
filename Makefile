.PHONY: test
test:
	@echo "Running tests..."
	@go test -timeout 300s -p 1 -count=1 -race -cover -v ./...