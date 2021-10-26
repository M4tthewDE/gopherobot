.PHONY: build lint test clean

build:
	go build -o target/gopherobot .
lint:
	golangci-lint run . bot/... cmd/... config/... provider/... util/...
test:
	go test ./...

docker:
	sudo docker build --tag gopherobot .

docker-run:
	sudo docker build --tag gopherobot .
	sudo docker container run -d --name gopherobot gopherobot

docker-clean:
	sudo docker container stop gopherobot
	sudo docker container rm gopherobot

