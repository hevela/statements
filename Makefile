
tidy: ## Cleans the Go module.
	@echo "==> Tidying module"
	@go mod tidy
.PHONY: tidy

mod-download: ## Downloads the Go module.
	@echo "==> Downloading Go module"
	@go mod download -x

lint: mod-download
	@golangci-lint --version
	@golangci-lint run -v ./...

lint-fix:
	@goimports -local "github.com/hevela/statements/" -w .

test: mod-download
	@go test -race -v ./... -coverprofile=coverage.txt -covermode=atomic

build: mod-download
	@go build -o ./. cmd/app/main.go