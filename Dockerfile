FROM golang:1.14.2 as build

RUN mkdir /opt/app
WORKDIR /opt/app

COPY ./go.mod ./go.mod
RUN go mod download
RUN apt update && apt install -y protobuf-compiler
COPY ./bin/* /usr/bin/
COPY . .
RUN protoc --go_out=. --go_opt=paths=source_relative  --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/*.proto
RUN go build -o server server.go
