.PHONY: build lint test clean

build:
	go build -ldflags \
	"-X de.com.fdm/gopherobot/config.GitCommit=$(shell git rev-parse --short HEAD) -X de.com.fdm/gopherobot/config.GitBranch=$(shell git branch --show-current)" \
	-o target/gopherobot .

lint:
	golangci-lint run . bot/... commands/... config/... util/...

test:
	go test ./...

run:
	go run -ldflags \
	"-X de.com.fdm/gopherobot/config.GitCommit=$(shell git rev-parse --short HEAD) -X de.com.fdm/gopherobot/config.GitBranch=$(shell git branch --show-current)" \
	.

docker:
	docker build --tag gopherobot .

docker-run:
	docker build --tag gopherobot .
	docker container run -d --name gopherobot gopherobot

docker-clean:
	docker container stop gopherobot
	docker container rm gopherobot

