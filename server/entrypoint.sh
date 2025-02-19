#!/bin/sh
ip route del default
ip route add default via 10.0.11.254
# go build -o server main.go

chmod +x server
# while true; do
# echo ""
# done
# Split the host ip into parts
First=$(ip addr show eth0 | grep "inet\b" | awk '{print $2}' | cut -d/ -f1 | cut -d . -f1)
Second=$(ip addr show eth0 | grep "inet\b" | awk '{print $2}' | cut -d/ -f1 | cut -d . -f2)
Third=$(ip addr show eth0 | grep "inet\b" | awk '{print $2}' | cut -d/ -f1 | cut -d . -f3)
Fourth=$(ip addr show eth0 | grep "inet\b" | awk '{print $2}' | cut -d/ -f1 | cut -d . -f4)

# Host ip
Ip=$First"."$Second"."$Third"."$Fourth
# Peer id it assumes that the node with ip X.X.X.2 is the first node in the ring
# PeerId=$First"."$Second"."$Third".2"

./server --ip=$Ip # --peer_id=$PeerId
