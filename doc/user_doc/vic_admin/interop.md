# Interoperability of vSphere Integrated Containers Engine with Other VMware Software
vSphere administrators can use vSphere to view and manage virtual container hosts (VCHs) and container VMs. vSphere Integrated Containers Engine works with the following vSphere features and VMware products. 

## Performing Operations on VCHs and Container VMs in vSphere ##

- If you restart a VCH endpoint VM, it comes back up in the same state that it was in when it shut down. 
- If you shut down the VCH vApp, the VCH endpoint VM and all of the container VMs that it manages are shut down. 
- If you use DHCP on the client network, the IP address of the VCH endpoint VM might change after a restart. Use `vic-machine inspect` to obtain the new IP address.
- Do not manually delete a VCH vApp or VCH endpoint VM. Use `vic-machine delete` to delete VCHs.
- Manually powering off container VMs can result in incorrect end-times for container opeations. Do not manually delete a contaienr VM. Use Docker commands to perform operations on container VMs.

## VMware vRealize&reg; Suite 
You can use VMware vRealize Automation to provide a self-provisioning service for VCHs, by using the vRealize Automation interface or APIs to request VCHs. At the end of such a provisioning process, vRealize Automation can communicate the VCH endpoint VM address to the requester.

## VMware vSphere vMotion&reg;  

You can use vMotion to move VCHs without needing to take the container VMs offline. The VCH endpoint VM does not need to be running for vMotion to occur on the container VMs. Clusters with a mix of container VMs and non-container VMs can use vMotion with fully automated DRS. 

## Maintenance Mode ##
Hosts with container VMs can enter maintenance mode without manual intervention, with these exceptions:

- For a standalone ESXi host, you must power down VCHs and any container VMs before entering maintenance mode.
- In a clustered vSphere environment with DRS set to automatic, DRS migrates VCHs to another host in the cluster before the host enters maintenance mode.
- For a host with running container VMs, DRS migrates the container VMs to another host in the cluster before the host enters maintenance mode.
- When you put a host into maintenance mode, vSphere Integrated Containers Engine powers off container VMs with open sessions. Container VMs that do not have open sessions are migrated to other hosts by vMotion. 

## VMware vSAN&trade;
VCHs maintain filesystem layers inherent in container images by mapping to discrete VMDK files, all of which can be housed in shared vSphere datastores, including vSAN, NFS, and iSCSI datastores.

## vCenter Linked Mode Environments
You can deploy VCHs in vCenter Linked Mode environments.

## vSphere Features Not Supported in This Release
vSphere Integrated Containers Engine does not currently support the following vSphere features:

- vSphere Storage DRS&trade;: You cannot configure VCHs to use Storage DRS datastore clusters. However, you can specify the path to a specific datastore within a Storage DRS datastore cluster by specifying the full inventory path to the datastore in the `vic-machine create --image-store` option. For example, `--image-store /dc1/datastore/my-storage-pod/datastore1`. You can also specify the relative path from a datastore folder in a datacenter, for example `--image-store my-storage-pod/datastore1`.
- vSphere High Availability: You can deploy VCHs to systems that are configured with High Availability, but you cannot use High Availability to fail over the VCHs themselves.
- vSphere Fault Tolerance: You cannot configure Fault Tolerance on VCHs.
- vSphere Virtual Volumes&trade;: You cannot use Virtual Volumes as the target datastores for image stores or volume stores.
- Snapshots: Creating and reverting to snapshots of the VCH endpoint VM or container VMs can cause vSphere Integrated Containers Engine not to function correctly.