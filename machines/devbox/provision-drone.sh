#!/bin/bash

# This depends on docker being on the box and running

docker pull drone/drone

curl http://downloads.drone.io/drone-cli/drone_linux_amd64.tar.gz | tar zx
sudo install -t /usr/local/bin drone
