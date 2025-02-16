FROM golang:1.22-alpine

WORKDIR /app
COPY shared shared
COPY server server

WORKDIR /app/server
RUN go get