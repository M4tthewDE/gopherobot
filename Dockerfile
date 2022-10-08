FROM golang:1.19-alpine
LABEL author="github.com/m4tthewde"

WORKDIR /go/src/github.com/m4tthewde/gopherobot
COPY . .
RUN go mod download
RUN apk add --no-cache git make
RUN make build

FROM alpine:latest
COPY --from=0 /go/src/github.com/m4tthewde/gopherobot/target/gopherobot .
CMD ["./gopherobot"]