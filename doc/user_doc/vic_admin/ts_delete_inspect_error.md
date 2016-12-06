# Deleting or Inspecting a VCH Fails with a Not a VCH or Resource Pool Not Found Error #

When you use `vic-machine delete` or `vic-machine inspect` to delete or inspect a virtual container host (VCH) and you specify the address of an ESXi host in the `target` option, the operation fails with "an error stating that the target is not a VCH or that the resource pool cannot be found".

## Problem ##
Deleting or inspecting a VCH fails with one of the following error messages: 

<pre>### Inspecting VCH ####                      
Not a VCH                                    
Failed to get Virtual Container Host <i>vch_name</i>   
Not a VCH                                    
--------------------                         
vic-machine-<i>os</i> failed: inspect failed</pre>

<pre>### Inspecting VCH ####                      
Failed to get VCH resource pool "<i>path_to_resource_pool</i>": 
resource pool '<i>path_to_resource_pool</i>' not found
Failed to get Virtual Container Host <i>vch_name</i>   
resource pool '<i>path_to_resource_pool</i>' not found                                   
--------------------                         
vic-machine-<i>os</i> failed: inspect failed</pre>

<pre>### Removing VCH ####                      
Not a VCH                                    
Failed to get Virtual Container Host <i>vch_name</i>    
Not a VCH                                    
--------------------                         
vic-machine-<i>os</i> failed: delete failed</pre>

<pre>### Removing VCH ####                      
Failed to get VCH resource pool "<i>path_to_resource_pool</i>": 
resource pool '<i>path_to_resource_pool</i>' not found
Failed to get Virtual Container Host <i>vch_name</i>   
resource pool '<i>path_to_resource_pool</i>' not found                                   
--------------------                         
vic-machine-<i>os</i> failed: delete failed</pre>

## Cause ##
You set the `target` option to the address of an ESXi host that is managed by a vCenter Server instance. If there are multiple ESXi hosts in a cluster, the error that you see depends on the host that you specify in the `target` option. 

- If you set the `target` option to the ESXi host on which the VCH is running, you see the error `Not a VCH, Failed to get Virtual Container Host`.
- If you set the `target` option to an ESXi host in the cluster that is not the one on which the VCH is running, you see the error `Not a VCH, Failed to get VCH resource pool`.

## Solution ##


1. Run `vic-machine ls` with the `target` option set to the same ESXi host.

  The `vic-machine ls` operation fails but informs you of the address of the  vCenter Server instance that manages the ESXi host.
2. Run `vic-machine delete` or `vic-machine inspect` again, setting the `target` option to the address of the vCenter Server instance that was returned by `vic-machine ls`.
