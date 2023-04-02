#!/bin/bash

if [ -z "$1" ] ; then
    echo "Usage: <command> <host_name>"
    exit 1
fi

HOST_ENTRY="$1"
sudo sed -i '' "/$HOST_ENTRY/d" /etc/hosts
