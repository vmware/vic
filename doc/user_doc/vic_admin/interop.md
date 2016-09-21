# Interoperability of vSphere Integrated Containers Engine with Other VMware Software
IT administrators use vCenter Server to view and manage containers. vSphere Integrated Containers Engine works seamlessly with VMware products. 

## vRealize Suite 
The vRealize Suite is available for health monitoring, performance analysis, and compliance across private and public clouds to move businesses faster.

## vSphere vMotion  

IT teams can assure service-level agreements for container workloads with VMware vSphere Distributed Resource Scheduler as well as reduce planned and unplanned downtime with VMware vSphere vMotion.

**NOTE**: In the current builds of vSphere Integrated Containers Engine, support for vSphere vMotion is not yet fully implemented.
<!--
You can apply vSphere High Availability and Fault Tolerance to both the container VMs and the virtual container host, so that containers and the virtual container host can power on or off independently of each other.
--> 

You can restart or upgrade the virtual container host without needing to take the containers offline during the process. You do not require a native agent on the ESXi host. The appliance VM does not need to be running for vMotion to occur. Clusters with non-container VMs can also vMotion with fully automated DRS.

## VMware Virtual SAN
The virtual container host maintains filesystem layers inherent in container images by mapping to discrete VMDK files, all of which can be housed in vSphere datastores on Virtual SAN, NFS, or local disks.

<!--
## vSphere Instant Clone
vSphere Integrated Containers Engine allows you to create and run multiple containers rapidly with minimal overhead using vSphere 6 Instant Clone technology, which provisions child VMs forked directly from a parent VM template running a Linux kernel. vSphere Integrated Containers Engine creates the kernel and a few supporting resources to run containers using Photon OS technology.
-->