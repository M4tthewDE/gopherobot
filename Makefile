build:
	go build -o target/gopherobot .
lint:
	golangci-lint run --enable-all --disable exhaustivestruct --disable tagliatelle --disable maligned --disable gomnd . bot/... cmd/... config/... provider/... util/...