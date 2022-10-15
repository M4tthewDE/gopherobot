FROM golang:1.19-alpine
LABEL author="github.com/m4tthewde"

WORKDIR /go/src/github.com/m4tthewde/gopherobot
COPY . .
RUN apk add --update --no-cache --repository http://dl-3.alpinelinux.org/alpine/edge/community vips-dev git make gcc musl-dev
RUN make build

FROM alpine:latest
RUN apk add --update --no-cache --repository http://dl-3.alpinelinux.org/alpine/edge/community vips-dev gcc musl-dev
COPY --from=0 /go/src/github.com/m4tthewde/gopherobot/target/gopherobot .
CMD ["./gopherobot"]