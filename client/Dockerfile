FROM golang:1.22-alpine AS build-client

WORKDIR /app
COPY shared shared
COPY client/go.mod client/go.mod
COPY client/go.sum client/go.sum
# To acknowledge go there is a module here
RUN touch client/main.go
RUN echo "package main; func main() {}" > client/main.go

WORKDIR /app/client
RUN go get

FROM build-client

WORKDIR /app
COPY shared /app/shared
COPY client /app/client

WORKDIR /app/client
RUN go get

RUN go build -o client /app/client/main.go

RUN chmod +x /app/client/entrypoint.sh

# CMD while true; do sleep 1; done
# ENTRYPOINT [ "entrypoint.sh" ]


CMD [ "/app/client/entrypoint.sh" ]
