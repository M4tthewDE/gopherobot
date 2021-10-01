.PHONY: build lint test clean

build:
	go build -o target/gopherobot .
lint:
	golangci-lint run . bot/... cmd/... config/... provider/... util/...
test:
	go test ./...
clean:
	rm -rf target/
