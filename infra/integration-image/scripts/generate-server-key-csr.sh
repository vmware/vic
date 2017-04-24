#!/usr/bin/env bash

# Generate a server key and CSR for specified CN

set -euf -o pipefail

OUTDIR="/root/ca"
SERVER_CERT_CN="starkenterprises.io"

while getopts ":d:n:" opt; do
  case $opt in
    d) OUTDIR="$OPTARG"
    ;;
    n) SERVER_CERT_CN="$OPTARG"
    ;;
    \?) echo "Invalid option -$OPTARG" >&2
    ;;
  esac
done


cd $OUTDIR
# Generate server key
# Private key is not encrypted - use -aes256 to specify a password
openssl genrsa -out private/${SERVER_CERT_CN}.key.pem 4096
chmod 400 private/${SERVER_CERT_CN}.key.pem

# Generate server CSR
openssl req -config openssl.cnf \
    -new -sha256 \
    -key private/${SERVER_CERT_CN}.key.pem \
    -out csr/${SERVER_CERT_CN}.csr.pem \
    -subj "/C=US/ST=California/L=Los Angeles/O=Stark Enterprises/OU=Stark Enterprises Web Services/CN=${SERVER_CERT_CN}"
