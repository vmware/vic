# Connection to Docker Daemon Fails  #
After a successful installation of vSphere Integrated Containers to a vCenter Server 6.0 instance, connecting to the Docker daemon fails.

## Problem ##
After you have set the `DOCKER_HOST` variable to point to your virtual container host, when you attempt a Docker operation in your Docker client, the connection to the Docker daemon fails with the error `Cannot connect to the Docker daemon. Is the docker daemon running on this host?`.

## Cause ##
You have deployed vSphere Integrated Containers onto an ESXi 5.5 host that is managed by vCenter Server 6.0.

## Solution ##
None. If you are deploying vSphere Integrated Containers to a vCenter Server instance, both vCenter Server and the ESXi host on which you are deploying the appliance must be version 6.0. 

You can only deploy the vSphere Integrated Containers appliance to an ESXi 5.5 host if that host is not managed by vCenter Server.