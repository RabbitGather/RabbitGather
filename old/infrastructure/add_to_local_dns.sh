#!/bin/bash
if [ -z "$1" ] || [ -z "$2" ]; then
    echo "Usage: <command> <IP> <host_name>"
    exit 1
fi

HOSTS_FILE="/etc/hosts"
HOST_ENTRY="$1       $2"
#echo $HOST_ENTRY
if ! grep -qxF "$HOST_ENTRY" "$HOSTS_FILE"; then
    echo "$HOST_ENTRY" | sudo tee -a "$HOSTS_FILE"
fi
