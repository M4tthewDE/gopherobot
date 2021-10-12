.PHONY: build lint test clean

build:
	go build -ldflags \
	"-X de.com.fdm/gopherobot/config.GitCommit=$(shell git rev-parse --short HEAD) -X de.com.fdm/gopherobot/config.GitBranch=$(shell git branch --show-current)" \
	-o target/gopherobot .

lint:
	golangci-lint run . bot/... cmd/... config/... provider/... util/...

test:
	go test ./...

run:
	go run -ldflags \
	"-X de.com.fdm/gopherobot/config.GitCommit=$(shell git rev-parse --short HEAD) -X de.com.fdm/gopherobot/config.GitBranch=$(shell git branch --show-current)" \
	.

docker:
	sudo docker build --tag gopherobot .