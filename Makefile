.PHONY: docs build test test-v cover cover-check lint fmt vet tidy install-tools

build:
	go build ./...

test:
	go test -race ./...

test-v:
	go test -race -v ./...

cover:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

cover-check:
	go test -race -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out | grep total | awk '{print $$3}' | \
		awk -F. '{if ($$1 < 70) {print "Coverage " $$0 " is below 70%"; exit 1}}'

lint:
	golangci-lint run ./...

fmt:
	gofmt -w .
	goimports -w -local github.com/Edge-Center/edgecenteredgemon-go .

vet:
	go vet ./...

tidy:
	go mod tidy

install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

docs:
	go list ./... | xargs -n1 go doc -all
