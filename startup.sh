#!/bin/bash

bash ./scripts/self_signed_certificate.sh "local" "$(pwd)/certs"
bash ./scripts/add_to_local_dns.sh "127.0.0.1" "dashboard.local"

