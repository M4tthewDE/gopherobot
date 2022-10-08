FROM golang:1.17-alpine
LABEL author="github.com/m4tthewde"

WORKDIR /go/src/github.com/m4tthewde/gopherobot
COPY . .
RUN go mod download
RUN apk add --no-cache git make
#RUN CGO_ENABLED=0 go build -o target/gopherobot . 
RUN make build

FROM alpine:latest
COPY --from=0 /go/src/github.com/m4tthewde/gopherobot/target/gopherobot .
COPY --from=0 /go/src/github.com/m4tthewde/gopherobot/config.yml .
CMD ["./gopherobot"]