FROM golang:1.22-alpine AS build-server

WORKDIR /app
COPY shared shared
COPY server/go.mod server/go.mod
COPY server/go.sum server/go.sum
# To acknowledge go there is a module here
RUN touch server/main.go
RUN echo "package main; func main() {}" > server/main.go

WORKDIR /app/server
RUN go get

FROM build-server

WORKDIR /app
COPY shared /app/shared
COPY server /app/server

WORKDIR /app/server
RUN go build -o server /app/server/main.go

EXPOSE 50053

RUN chmod +x /app/server/entrypoint.sh

CMD [ "/app/server/entrypoint.sh" ]