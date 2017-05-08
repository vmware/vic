#!/bin/bash

data_dir=/data/kov
cert_dir=${data_dir}/cert
cert=${cert_dir}/server.crt
key=${cert_dir}/server.key

echo KOVD_EXPOSED_PORT="$(ovfenv -k cluster_manager.port)"
echo KOVD_KEY_LOCATION=$key
echo KOVD_CERT_LOCATION=$cert
