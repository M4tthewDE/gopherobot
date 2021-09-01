build:
	go build -o target/gopherobot .
lint:
	golangci-lint run --enable-all . bot/... cmd/... config/... provider/... util/...