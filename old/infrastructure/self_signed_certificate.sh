#!/bin/bash

HOST_NAME=$1
CERT_DIR=${2:-$(pwd)/certs}

if [ -z "$1" ]; then
    echo "Usage: <command> <hostname>"
    echo "Usage: <command> <hostname> <cert_output_dir>"
    echo "if <cert_output_dir> is not specified, it will be default as \"\$(pwd)/certs\" directory"
    exit 1
fi

mkdir -p "${CERT_DIR}"

openssl genrsa -aes256 -passout pass:password -out "${CERT_DIR}/rootCA.key" 4096

openssl req -x509 -new -nodes -key "${CERT_DIR}/rootCA.key" -sha256 -days 365000 -subj "/C=TW/ST=Taipei/L=Taipei/O=meowalien/OU=meowalien/CN=meowalien" -passin pass:password -out "${CERT_DIR}/rootCA.crt"

sudo security add-certificates "${CERT_DIR}/rootCA.crt"
sudo security add-trusted-cert -r trustRoot -d "${CERT_DIR}/rootCA.crt"

openssl genrsa -out "${CERT_DIR}/${HOST_NAME}.key" 4096

echo "[req]
default_bits = 4096
prompt = no
default_md = sha512
distinguished_name = dn

[dn]
C=TW
ST=Taiwan
L=Taipei
O=meowalien
OU=meowalien
emailAddress=meowalien@meowalien.com
CN=*.${HOST_NAME}" > "/tmp/${HOST_NAME}.cnf"

openssl req -new -sha256 -nodes -key "${CERT_DIR}/${HOST_NAME}.key" -out "${CERT_DIR}/${HOST_NAME}.csr" -config "/tmp/${HOST_NAME}.cnf"
echo "authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment, keyAgreement, keyCertSign, cRLSign
subjectAltName = @alt_names

[alt_names]
DNS.1 = *.${HOST_NAME}" > "/tmp/${HOST_NAME}-v3.ext"

openssl x509 -req -passin pass:password -in "${CERT_DIR}/${HOST_NAME}.csr" -CA "${CERT_DIR}/rootCA.crt" -CAkey "${CERT_DIR}/rootCA.key" -CAcreateserial -out "${CERT_DIR}/${HOST_NAME}.crt" -days 36500 -sha256 -extfile "/tmp/${HOST_NAME}-v3.ext"

cat "${CERT_DIR}/${HOST_NAME}.crt" "${CERT_DIR}/rootCA.crt" > "${CERT_DIR}/${HOST_NAME}_fullchain.crt"
sudo security add-certificates "${CERT_DIR}/${HOST_NAME}_fullchain.crt"
sudo security add-trusted-cert -r trustRoot -d "${CERT_DIR}/${HOST_NAME}_fullchain.crt"
