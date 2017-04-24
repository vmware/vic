#!/usr/bin/env bash

# Install CA into root store

CERT_FILE="/root/ca/certs/STARK_ENTERPRISES_ROOT_CA.crt"

while getopts ":f:" opt; do
  case $opt in
    f) CERT_FILE="$OPTARG"
    ;;
    \?) echo "Invalid option -$OPTARG" >&2
    ;;
  esac
done

cp $CERT_FILE /usr/local/share/ca-certificates

dpkg-reconfigure --frontend=noninteractive ca-certificates
