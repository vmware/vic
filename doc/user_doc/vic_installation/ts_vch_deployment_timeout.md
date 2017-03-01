# VCH Deployment Times Out #

When you use `vic-machine create` to deploy a virtual container host (VCH), the operation times out.

## Problem ##
Deployment fails with a timeout error that states that the context deadline has been exceeded. 

<pre>Failed to power on appliance context deadline exceeded. Exiting...
vic-machine-linux failed: Create timed out: 
if slow connection, increase timeout with --timeout</pre>

## Causes ##

This error can have different causes:

- The connection between the system on which you are running `vic-machine` and vCenter Server is slow. The upload of the ISO files exceeds the default 3 minute timeout.
- The upload of the ISO files succeeds but the VCH fails to obtain an IP address.
- The VCH obtained an IP address, but the VCH service does not start or the VCH cannot connect to the Docker API.

## Solutions ##

1. Set the `vic-machine --timeout` option to allow more time for the ISOs to upload. 

   For example, set `--timeout 10m` or `--timeout 20m`. 
2. If the ISO upload succeeds with a longer timeout period but the operation still times out, check the DHCP service to make sure than an IP address is available for the VCH.
3. If the DHCP service is working and the operation still times out, see [VCH Deployment Fails with Docker API Endpoint Check Failed Error](ts_docker_api_check_error.md)  