#!/usr/bin/env bash

# Reload default CAs to remove ALL user installed root certificates

update-ca-certificates --fresh
