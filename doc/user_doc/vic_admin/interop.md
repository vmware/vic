# Interoperability of vSphere Integrated Containers Engine with Other VMware Software
IT administrators use vCenter Server to view and manage containers. vSphere Integrated Containers Engine works seamlessly with VMware products. 

## VMware vRealize&reg; Suite 
vRealize Suite is available for health monitoring, performance analysis, and compliance across private and public clouds to move businesses faster.

## VMware vSphere vMotion&reg;  

IT teams can assure service-level agreements for container workloads with VMware vSphere Distributed Resource Scheduler&trade; (DRS) as well as reduce planned and unplanned downtime with VMware vSphere vMotion.

You can restart the virtual container host (VCH) without needing to take the containers offline during the process. You do not require a native agent on the ESXi host. The appliance VM does not need to be running for vMotion to occur. Clusters with non-container VMs can also vMotion with fully automated DRS.


## Maintenance Mode ##
Hosts with container VMs can enter maintenance mode without manual intervention, with these exceptions:

- For a standalone ESXi host, you must power down the VCH and any container VMs before entering maintenance mode. 
- In a clustered vSphere environment with DRS set to automatic, DRS migrates the VCH to another host in the cluster before the host enters maintenance mode. 
- For a host with container VMs running, DRS migrates the container VMs to another host in the cluster before the host enters maintenance mode.

## VMware vSAN&trade;
The VCH maintains filesystem layers inherent in container images by mapping to discrete VMDK files, all of which can be which can be housed in shared vSphere datastores, including vSAN and NFS datastores.

## vCenter Linked Mode Environments
You can deploy VCHs in vCenter Linked Mode environments.

## Multiple VLAN
You can deploy and run Docker commands against a VCH on a vCenter appliance that has multiple portgroups on different VLANs within the datacenter.

<!--
## vSphere Instant Clone
vSphere Integrated Containers Engine allows you to create and run multiple containers rapidly with minimal overhead using vSphere 6 Instant Clone technology, which provisions child VMs forked directly from a parent VM template running a Linux kernel. vSphere Integrated Containers Engine creates the kernel and a few supporting resources to run containers using Photon OS technology.
-->

## vSphere Features Not Supported in This Release
vSphere Integrated Containers Engine does not currently support the following vSphere features:

- vSphere Storage DRS&trade;: You cannot configure VCHs to use datastores that are in Storage DRS clusters.
- vSphere High Availability: You can deploy VCHs to systems that are configured with High Availability, but you cannot use High Availability to fail over the VCHs themselves.
- vSphere Fault Tolerance: You cannot configure Fault Tolerance on VCHs.
- vSphere Virtual Volumes&trade;: You cannot use Virtual Volumes as the target datastores for image stores or volume stores.
- Snapshots: Creating snapshots of the VCH endpoint VM or container VMs can cause vSphere Integrated Containers Engine not to function correctly.
- Powering on, powering off, or deleting the VCH endpoint VM or container VMs can cause vSphere Integrated Containers Engine not to function correctly. Use vic-machine to delete VCHs and use Docker commands to perform operations on container VMs.
