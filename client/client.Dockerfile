FROM client:base

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
