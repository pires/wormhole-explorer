build:
	CGO_ENABLED=0 GOOS=linux go build -o pipeline cmd/main.go

test:
	go test -v -cover ./...

.PHONY: build test
