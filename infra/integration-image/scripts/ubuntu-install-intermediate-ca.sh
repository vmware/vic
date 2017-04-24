#!/usr/bin/env bash

cp /root/ca/intermediate/certs/intermediate.cert.crt /usr/local/share/ca-certificates

dpkg-reconfigure --frontend=noninteractive ca-certificates
