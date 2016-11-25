# Docker Commands Fail with a Docker API Version Error #

After a successful deployment of a vSphere Integrated Containers Engine virtual container host (VCH), attempting to run a Docker command fails with a Docker version error.

## Problem ##
When you attempt to run a Docker command from a Docker client that is connecting to a VCH, the command fails with the error `Error response from daemon: client is newer than server (client API version: 1.24, server API version: 1.23)`.

## Cause ##

vSphere Integrated Containers Engine supports Docker 1.11, that includes version 1.23 of the Docker API. You are using version 1.12 of the Docker client, that uses version 1.24 of the Docker API, which is incompatible.

## Solution ##
1. Open a Docker client terminal.
2. Set the Docker client API to the same version as is used by vSphere Integrated Containers Engine.

 <pre>export DOCKER_API_VERSION=1.23</pre> 
4. Check that your Docker client can now connect to the VCH by running a Docker command. 

 <pre>docker -H <i>virtual_container_host_address</i>:2376 --tls info</pre>

 The `docker info` command should succeed and you should see information about the VCH.