FROM golang:latest
MAINTAINER m4tthewde github.com/m4tthewde
WORKDIR /go/src/github.com/m4tthewde/gopherobot
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o target/gopherobot . 

FROM alpine:latest
RUN apk add --no-cache git
COPY --from=0 /go/src/github.com/m4tthewde/gopherobot/ .
CMD ["./target/gopherobot"]