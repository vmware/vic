# Deleting or Inspecting a VCH Fails with a Resource Pool Error #

When you use `vic-machine delete` or `vic-machine inspect` to delete or inspect a virtual container host and you specify the address of an ESXi host in the `target` option, the operation fails with a resource pool creation error.

## Problem ##
Deleting or inspecting a virtual container host fails with the error message: 

<pre>Failed to get VCH resource pool "/ha-datacenter/host/localhost./Resources/<i>vch_name</i>": 
resource pool '/ha-datacenter/host/localhost./Resources/<i>vch_name</i>' not found</pre>

## Cause ##
You set the `target` option to the address of an ESXi host that is managed by a vCenter Server instance.

## Solution ##


1. Run `vic-machine ls` with the `target` option set to the same ESXi host.

  The `vic-machine ls` operation fails but informs you of the address of the  vCenter Server instance that manages the ESXi host.
2. Run `vic-machine delete` or `vic-machine inspect` again, setting the `target` option to the address of the vCenter Server instance that was returned by `vic-machine ls`.
