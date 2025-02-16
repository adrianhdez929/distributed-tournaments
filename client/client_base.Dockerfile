FROM golang:1.22-alpine

WORKDIR /app
COPY shared shared
COPY client client

WORKDIR /app/client
RUN go get