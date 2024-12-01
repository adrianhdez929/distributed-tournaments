#!/bin/sh

# remove hardcoded router ip
ip route del default
ip route add default via 10.0.10.254

# go build -o client main.go
chmod +x client
./client