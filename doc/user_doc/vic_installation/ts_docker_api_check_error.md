# VCH Deployment Fails with Docker API Endpoint Check Failed Error #
When you use `vic-machine create` to deploy a virtual container host (VCH), deployment fails because
vic-machine cannot contact the Docker API endpoint.

## Problem ##

Deployment fails with Docker API endpoint check failed:

<pre>
Docker API endpoint check failed:
API may be slow to start - try to connect to API after a few minutes:
   Run docker -H 192.168.218.160:2376 --tls info
   If command succeeds, VCH is started. If command fails, VCH failed to install - see documentation for troubleshooting.
</pre>


## Cause ##

During install, vic-machine checks that the endpoint that Docker clients connect to is reachable. If
this check fails, vic-machine create fails with an error. The Docker API may be slow to start or it
may have failed to start.

## Solution ##

### Docker API is slow to start ###

After waiting a few minutes, run the `docker info` command shown in the error message to test the
API responsiveness. If this succeeds, output similar to below will be shown:

<pre>
$ docker -H 192.168.218.160:2376 --tls info
Containers: 0
 Running: 0
 Paused: 0
 Stopped: 0
Images: 0
Server Version: v0.8.0-0-1652a8b
Storage Driver: vSphere Integrated Containers v0.8.0-0-1652a8b Backend Engine
VolumeStores:
vSphere Integrated Containers v0.8.0-0-1652a8b Backend Engine: RUNNING
 VCH mhz limit: 2944 Mhz
 VCH memory limit: 1.259 GiB
 VMware Product: VMware ESXi
 VMware OS: vmnix-x86
 VMware OS version: 6.5.0
Plugins:
 Volume:
 Network: bridge
Swarm:
 NodeID:
 Is Manager: false
 Node Address:
Security Options:
Operating System: vmnix-x86
OSType: vmnix-x86
Architecture: x86_64
CPUs: 2944
Total Memory: 1.259 GiB
Name: oceanlab
ID: vSphere Integrated Containers
Docker Root Dir:
Debug Mode (client): false
Debug Mode (server): false
Registry: registry-1.docker.io
</pre>

This output means that the VCH is running as expected and can accept Docker commands.

If the command times out, proceed to `Docker API failed to start` for troubleshooting.

### Docker API failed to start ###

If the Docker API is not responsive, run `vic-machine inspect`, login to the VCH Admin Portal, and
download the Log Bundle. Follow troubleshooting steps to determine why the installation failed.

If the VCH Admin Portal is not responsive, run `vic-machine debug` to enable the VCH debugging mode.
SSH to the VCH to collect the VIC logs. Follow troubleshooting steps to determine why the installation
failed. Collecting the vSphere log bundle may also be useful in troubleshooting.
