#!/usr/bin/env bash

rm /usr/local/share/ca-certificates/ca.cert.crt

dpkg-reconfigure --frontend=noninteractive ca-certificates
