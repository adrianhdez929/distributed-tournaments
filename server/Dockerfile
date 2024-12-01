FROM golang:1.22-alpine


WORKDIR /app

COPY shared /app/shared
COPY server /app/server

WORKDIR /app/server
RUN go get

RUN go build -o server /app/server/main.go

EXPOSE 50053

# RUN ip route del default
# RUN ip route add default via 10.0.11.254

# CMD [ "./server" ]
RUN chmod +x /app/server/entrypoint.sh

CMD [ "/app/server/entrypoint.sh" ]