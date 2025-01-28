install:
	go mod tidy

clean:
	go mod tidy
	go fmt ./...
	goimports -w .

start:
	go run cmd/main.go

test:
	${GOPATH}/bin/golangci-lint run