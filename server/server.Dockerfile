FROM server:base

WORKDIR /app
COPY shared /app/shared
COPY server /app/server

WORKDIR /app/server
RUN go build -o server /app/server/main.go

EXPOSE 50053
EXPOSE 50054
EXPOSE 50055

RUN chmod +x /app/server/entrypoint.sh

CMD [ "/app/server/entrypoint.sh" ]