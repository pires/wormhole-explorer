## help: print this help message
.PHONY: help build doc test
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

build:
	go build -o bin/analytics cmd/main.go
	
doc:
	swag init -pd


test:
	go test -v -cover ./...
