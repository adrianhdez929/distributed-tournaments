FROM golang:1.22.5-alpine


WORKDIR /app

COPY shared /app/shared
COPY server /app/server

WORKDIR /app/server
RUN go get

ENTRYPOINT [ "server/entrypoint.sh" ]