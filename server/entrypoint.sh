#!/bin/sh
ip route del default
ip route add default via 10.0.11.254
# go build -o server main.go

chmod +x server
./server