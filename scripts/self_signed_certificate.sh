#!/bin/bash
if [ -z "$1" ] || [ -z "$2" ]; then
    echo "Usage: <command> <hostname> <cert_output_dir>"
    exit 1
fi

HOST_NAME=$1
CERT_DIR=${2:-$(pwd)/certs}
FULL_CHAIN_NAME="${HOST_NAME}_fullchain.pem"

mkdir -p "${CERT_DIR}"

function create_and_trust_root_cert() {
    openssl genrsa -aes256 -passout pass:password -out "${CERT_DIR}/root.key" 4096

    openssl req -x509 -new -sha512 -days 365 \
        -subj "/C=TW/ST=Taipei/L=Taipei/O=test/OU=lab/CN=root" \
        -passin pass:password \
        -key "${CERT_DIR}/root.key" \
        -out "${CERT_DIR}/root.pem"
    echo "Adding root.pem to keychain, need permission for sudo:"
    sudo security add-certificates ${CERT_DIR}/root.pem
    sudo security add-trusted-cert -r trustRoot -d ${CERT_DIR}/root.pem
}

if ! [ -f "${CERT_DIR}/root.pem" ]; then
    create_and_trust_root_cert || exist 1
fi

function generate_certificate() {
mkdir -p "${CERT_DIR}"

  openssl genrsa -out "${CERT_DIR}/${HOST_NAME}.key" 4096

  openssl req -sha512 -new \
      -subj "/C=TW/ST=Taiwan/L=Taipei/O=test/OU=lab/CN=*.${HOST_NAME}" \
      -key "${CERT_DIR}/${HOST_NAME}.key" \
      -out "${CERT_DIR}/${HOST_NAME}.csr"

  echo "authorityKeyIdentifier=keyid,issuer
        basicConstraints=CA:FALSE
        keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment, keyAgreement, keyCertSign, cRLSign
        extendedKeyUsage = serverAuth
        subjectAltName = DNS:${HOST_NAME},DNS:*.${HOST_NAME}" > /tmp/v3.ext
#  ALTNAMES="DNS:*.local,DNS:local"

#  openssl x509 -req -days 365 -in "${CERT_DIR}/${HOST_NAME}.csr" -signkey "${CERT_DIR}/${HOST_NAME}.key" -out "${CERT_DIR}/${HOST_NAME}.pem"  -extensions v3_req   -extfile <(echo "[v3_req]
#                                                                                   subjectAltName=$ALTNAMES")

  openssl x509 -req -sha512 -days 365 \
      -extfile /tmp/v3.ext \
      -passin pass:password \
      -CA ${CERT_DIR}/root.pem -CAkey ${CERT_DIR}/root.key -CAcreateserial \
      -in "${CERT_DIR}/${HOST_NAME}.csr" \
      -out "${CERT_DIR}/${HOST_NAME}.pem"


  cat ${CERT_DIR}/${HOST_NAME}.pem ${CERT_DIR}/root.pem > ${CERT_DIR}/${FULL_CHAIN_NAME}

  sudo security add-certificates ${CERT_DIR}/${FULL_CHAIN_NAME}
}

# check if directory exists
if ! [ -d "${CERT_DIR}/${HOST_NAME}.pem" ]; then
    generate_certificate || exit 1
fi


