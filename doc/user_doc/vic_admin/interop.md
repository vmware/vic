# Interoperability of vSphere Integrated Containers Engine with Other VMware Software
IT administrators use vCenter Server to view and manage containers. vSphere Integrated Containers Engine works seamlessly with VMware products. 

## VMware vRealize&reg; Suite 
vRealize Suite is available for health monitoring, performance analysis, and compliance across private and public clouds to move businesses faster.

## VMware vSphere vMotion&reg;  

IT teams can assure service-level agreements for container workloads with VMware vSphere Distributed Resource Scheduler&trade; (DRS) as well as reduce planned and unplanned downtime with VMware vSphere vMotion.

You can restart or upgrade the virtual container host without needing to take the containers offline during the process. You do not require a native agent on the ESXi host. The appliance VM does not need to be running for vMotion to occur. Clusters with non-container VMs can also vMotion with fully automated DRS.

Hosts with container VMs can enter maintenance mode without manual intervention, with these exceptions:

- For a standalone host, you must first power down the virtual container host or any other VMs before entering maintenance mode. 
- In a clustered vSphere environment with DRS set to automatic, DRS first migrates the virtual container host to another host in the cluster and then enters maintenance mode. 
- For a host with a container VM running, DRS first migrates the VM to another host in the cluster and then enters maintenance mode.

## VMware vSAN&trade;
The virtual container host maintains filesystem layers inherent in container images by mapping to discrete VMDK files, all of which can be which can be housed in shared vSphere datastores, including vSAN and NFS datastores.

## Enhanced Link Mode Environments
You can deploy virtual container hosts in enhanced linked mode environments.

## Multiple VLAN
You can deploy and run Docker commands against a virtual container host on a vCenter appliance that has multiple portgroups on different VLANs within the datacenter.

<!--
## vSphere Instant Clone
vSphere Integrated Containers Engine allows you to create and run multiple containers rapidly with minimal overhead using vSphere 6 Instant Clone technology, which provisions child VMs forked directly from a parent VM template running a Linux kernel. vSphere Integrated Containers Engine creates the kernel and a few supporting resources to run containers using Photon OS technology.
-->

## vSphere Features Not Supported in This Release
vSphere Integrated Containers Engine does not currently support the following vSphere features:

- vSphere Storage DRS&trade;: You cannot configure virtual container hosts to use datastores that are in Storage DRS clusters.
- vSphere High Availability: You can deploy virtual container hosts to systems that are configured with High Availability, but you cannot use High Availability to fail over the virtual container hosts themselves.
- vSphere Fault Tolerance: You cannot configure Fault Tolerance on virtual container hosts.
- vSphere Virtual Volumes&trade;: You cannot use Virtual Volumes as the target datastores for image stores or volume stores.
- Snapshots: Creating snapshots of the virtual container host endpoint VM or container VMs can cause vSphere Integrated Containers Engine not to function correctly.
- Powering on, powering off, or deleting the virtual container host endpoint VM or container VMs can cause vSphere Integrated Containers Engine not to function correctly.
