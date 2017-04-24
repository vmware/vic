#!/usr/bin/env bash

# Sign a CSR with specified CA

set -euf -o pipefail

CA_NAME="STARK_ENTERPRISES_ROOT_CA" # used to verify cert signature
OUTDIR="/root/ca"
SERVER_CERT_CN="starkenterprises.io"

while getopts ":c:d:n:" opt; do
  case $opt in
    c) CA_NAME="$OPTARG"
    ;;
    d) OUTDIR="$OPTARG"
    ;;
    n) SERVER_CERT_CN="$OPTARG"
    ;;
    \?) echo "Invalid option -$OPTARG" >&2
    ;;
  esac
done


CONF_DIR=`dirname $0`
cd $OUTDIR
openssl ca -config $CONF_DIR/openssl.cnf \
    -batch \
    -extensions server_cert \
    -days 365 -notext -md sha256 \
    -in csr/${SERVER_CERT_CN}.csr.pem \
    -out certs/${SERVER_CERT_CN}.cert.pem

chmod 444 certs/${SERVER_CERT_CN}.cert.pem
openssl x509 -noout -text -in certs/${SERVER_CERT_CN}.cert.pem

# Test certificate
openssl verify -CAfile certs/$CA_NAME.crt certs/${SERVER_CERT_CN}.cert.pem
