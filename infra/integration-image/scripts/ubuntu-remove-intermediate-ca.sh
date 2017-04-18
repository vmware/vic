#!/usr/bin/env bash

rm /usr/local/share/ca-certificates/intermediate.cert.crt

dpkg-reconfigure --frontend=noninteractive ca-certificates
