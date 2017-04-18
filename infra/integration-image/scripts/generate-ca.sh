#!/usr/bin/env bash

# Generate a CA

set -euf -o pipefail

OUTDIR="/root/ca"
CA_NAME="STARK_ENTERPRISES_ROOT_CA"

while getopts ":c:d:" opt; do
  case $opt in
    c) CA_NAME="$OPTARG"
    ;;
    d) OUTDIR="$OPTARG"
    ;;
    \?) echo "Invalid option -$OPTARG" >&2
    ;;
  esac
done


### Create root CA
mkdir -p $OUTDIR
cp openssl.cnf $OUTDIR
cd $OUTDIR
mkdir certs crl csr newcerts private
chmod 700 private
touch index.txt
echo 1000 > serial

# Generate root CA key
# Private key is not encrypted - use -aes256 to specify a password
openssl genrsa -out private/ca.key.pem 4096
chmod 400 private/ca.key.pem

# Generate root CA CSR
openssl req -config openssl.cnf \
    -new -sha256 \
    -key private/ca.key.pem \
    -out csr/ca.csr.pem \
    -extensions v3_ca \
    -subj "/C=US/ST=California/L=Los Angeles/O=Stark Enterprises/OU=Stark Enterprises Certificate Authority/CN=Stark Enterprises Global CA"

# Self sign for root CA certificate
openssl x509 -req -extfile openssl.cnf \
    -extensions v3_ca \
    -days 7300 -in csr/ca.csr.pem -signkey private/ca.key.pem -out certs/ca.cert.pem

chmod 444 certs/ca.cert.pem
openssl x509 -noout -text -in certs/ca.cert.pem

# Output CRT format
openssl x509 -in certs/ca.cert.pem -inform PEM -out certs/$CA_NAME.crt
