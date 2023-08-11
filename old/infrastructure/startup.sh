#!/bin/bash

cd ./infrastructure || exit 1
CURRENT_DIR=$(pwd)

#bash ./scripts/self_signed_certificate.sh "meowalien.local" "$(pwd)/certs"
#bash ./scripts/add_to_local_dns.sh "127.0.0.1" "dashboard.meowalien.local"
#bash ./scripts/add_to_local_dns.sh "127.0.0.1" "registry.meowalien.local"


cd ./traefik && docker-compose up -d && cd "$CURRENT_DIR" || exit 1
cd ./harbor || exit 1 && docker-compose up -d && cd "$CURRENT_DIR" || exit 1


#docker compose up traefik