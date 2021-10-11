.PHONY: build lint test clean

build:
	go build -o target/gopherobot .
lint:
	golangci-lint run . bot/... cmd/... config/... provider/... util/...
test:
	go test ./...

docker:
	sudo docker build --tag gopherobot .