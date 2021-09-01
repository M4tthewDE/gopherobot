build:
	go build -o target/gopherobot .
lint:
	golangci-lint run --enable-all --disable exkhaustivestruct --disable tagliatelle --disable maligned. bot/... cmd/... config/... provider/... util/...