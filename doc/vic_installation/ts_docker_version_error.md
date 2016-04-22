# Connection to Docker Client Fails with a Docker Version Error #

After a successful installation of vSphere Integrated Containers, connecting a Docker client to the virtual container host fails with a Docker version error.

## Problem ##
After you have set the `DOCKER_HOST` variable to point to your virtual container host, when you attempt a Docker operation in your Docker client, the connection fails with the error `Error response from daemon: client is newer than server (client API version: 1.22, server API cersion: 1.21)`.

## Cause ##

The vSphere Integrated Containers technical preview uses version 1.21 of the Docker API. Your Docker client uses a more recent version of the Docker API that is incompatible.

## Solution ##
1. Open a Docker client terminal.
2. Set the Docker client API to the same version as is used by vSphere Integrated Containers.

 `export DOCKER_API_VERSION=1.21` 
4. Check that your Docker client can now connect to the virtual container host by running a Docker command. 

 `docker info` 

 You should see information about the virtual container host that is running in vSphere.