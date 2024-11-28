FROM golang:1.22.5-alpine


WORKDIR /app
COPY shared /app/shared
COPY client /app/client

WORKDIR /app/client
RUN go get

ENTRYPOINT [ "client/entrypoint.sh" ]