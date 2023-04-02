#!/bin/bash
if [ -z "$1" ] ; then
    echo "Usage: <command> <cert_file.pem>"
    exit 1
fi

CERT_FILE=$1
echo $CERT_FILE

filename="$(basename "$CERT_FILE")"
CERT_NAME=${filename%.*}

#echo $name_without_ext
#CERT_NAME=${$(basename "$CERT_FILE")%.*}  # Remove the file extension
echo $CERT_NAME
sudo security remove-trusted-cert -d "$CERT_FILE"
security delete-certificate -c "$CERT_NAME" "$HOME/Library/Keychains/login.keychain"

