FROM golang:1.16-alpine 
MAINTAINER m4tthewde github.com/m4tthewde
WORKDIR /app
COPY . ./
RUN go mod download
RUN go build -o target/gopherobot . 
RUN apk add --no-cache git
CMD ["./target/gopherobot"]