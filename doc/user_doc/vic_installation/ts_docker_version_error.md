# Docker Commands Fail with a Docker API Version Error #

After a successful deployment of a vSphere Integrated Containers Engine virtual container host (VCH), attempting to run a Docker command fails with a Docker version error.

## Problem ##
When you attempt to run a Docker command from a Docker client that is connecting to a VCH, the command fails with the following error: 
<pre>Error response from daemon: client is newer than server 
(client API version: <i>x.xx</i>, server API version: <i>y.yy</i>)</pre>

## Cause ##

This version of vSphere Integrated Containers Engine supports Docker 1.13, that includes version 1.25 of the Docker API. You are using a more recent version of the Docker client, that includes a version of the Docker API that is incompatible.

## Solution ##
1. Open a terminal on the system on which you run the Docker client.
2. Set the Docker client API to the same version as the one that is used by vSphere Integrated Containers Engine.

    <pre>export DOCKER_API_VERSION=1.25</pre> 

4. Check that your Docker client can now connect to the VCH by running a Docker command. 

   - With TLS authentication: <pre>docker -H <i>virtual_container_host_address</i>:2376 --tls info</pre>
   - Without TLS authentication:<pre>docker -H <i>virtual_container_host_address</i>:2375 info</pre>

   The `docker info` command should succeed and you should see information about the VCH.