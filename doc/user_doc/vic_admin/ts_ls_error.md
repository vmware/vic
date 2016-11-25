# Running `vic-machine ls` on an ESXi Host Fails with an Error #

When you use `vic-machine ls` to list virtual container hosts (VCHs) and you specify the address of an ESXi host in the `target` option, the operation fails with an error.

## Problem ##
Listing VCHs fails with the error message: 

<pre>Target is managed by vCenter server "<i>vcenter_server_address</i>", 
please change --target to vCenter server address or select a standalone ESXi</pre>

## Cause ##
You set the `target` option to the address of an ESXi host that is managed by a vCenter Server instance.

## Solution ##
Set the `target` option to the address of the vCenter Server instance that manages the ESXi host on which the VCHs are running. 