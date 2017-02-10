# VCH Deployment Fails with an I/O Timeout Error #
When you use `vic-machine create` to deploy a virtual container host (VCH), deployment fails with an I/O timeout error.

## Problem ##

Deployment fails with the error `Connection failed with error: i/o timeout.`

## Cause ##

During deployment, `vic-machine` checks that the endpoint VM is reachable from Docker clients. If this check fails, `vic-machine create` fails with an error. This error can be caused by the Docker API being slow to start or because it has failed to start.

## Solution ##

The solution to choose depends on whether the API is slow to start or whether it failed to start.

### Docker API is Slow to Start ###

Wait for a few minutes, then run the `docker info` command to test the responsiveness of the Docker API. 

If `docker info` succeeds, it shows information about the VCH, including confirmation that the storage driver is vSphere Integrated Containers. 
<pre>Storage Driver: vSphere Integrated Containers <i>version</i> Backend Engine</pre> 
This output means that the VCH is running correctly and can now accept Docker commands.

If `docker info` times out, it means that the Docker API did not start.

### Docker API Did Not Start ###

If the Docker API was not responsive when you ran `docker info`, download the VCH log bundle and examine the logs to determine why the installation failed. Collecting the vSphere log bundle might also be useful for troubleshooting.

- For information about how to download VCH logs by using the VCH Admin Portal, see [Access the VCH Admin Portal](../vic_admin/access_vicadmin.html) in *vSphere Integrated Containers Engine Administration*.
- For information about how to collect VCH logs manually, see [Access vSphere Integrated Containers Engine Log Bundles](../vic_admin/log_bundles.html) in *vSphere Integrated Containers Engine Administration*.