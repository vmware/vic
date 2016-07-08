# Environment Prerequisites for vSphere Integrated Containers Installation

Before you install vSphere Integrated Containers, you must ensure that your infrastructure meets certain requirements.

## Supported Platforms for `vic-machine`

The current builds of the vSphere Integrated Containers installation and management utility, `vic-machine`, have been tested and verified on the following Linux OS, Windows, and Mac OS systems.

|**Platform**|**Supported Versions**|
|---|---|
|Windows|7, 10|
|Mac OS X |10.11 (TBC)|
|Linux|Ubuntu 15.04, others TBD|

## Supported vSphere Configurations

You can install vSphere Integrated Containers in the following vSphere setups:

* Standalone ESXi 6.0 host that is not managed by a vCenter Server instance.
* vCenter Server 6.0, managing one or more standalone ESXi 6.0 hosts.
* vCenter Server 6.0, managing a cluster of ESXi 6.0 hosts, with DRS enabled.

In all cases, your ESXi hosts must have at least 8GB of memory.

Deploying vSphere Integrated Containers to a vCenter Server instance that is running in Enhanced Linked Mode is fully supported.  

## License Requirements
The type of license that vSphere Integrated Containers requires depends on the way in which you deploy the software.

| **Type of Installation** | **vSphere Feature Used** | **Required License** |
| --- | --- | --- |
|ESXi host|Network Serial Port|vSphere Enterprise|
|vCenter Server|Distributed Virtual Switch|vSphere Enterprise Plus|

All of the ESXi hosts in a cluster require an appropriate license. Installation fails if your environment includes one or more ESXi hosts that have inadequate licenses.

## Role and Permissions Requirements
You must use an account with the vSphere Administrator role when you install vSphere Integrated Containers.

<a name="networkreqs"></a>
## Network Requirements
* Use a trusted network for the deployment and use of vSphere Integrated Containers.
* Use a trusted network for connections between Docker clients and the virtual container hosts.
* Each virtual container host requires the following network configuration:
 * One IP address, that can be either static or set by using DHCP.
 * One VLAN, if you use VLAN for network isolation.
 * One IP address for each container that you run with the `docker run --net=host` option.
* In a vCenter Server environment, before you deploy a virtual container host, you must create a distributed virtual switch and a distributed port group. You must add the target ESXi host or hosts to the distributed virtual switch. 
 - For information about how to create a distributed port group, see [Create a vSphere Distributed Switch](https://pubs.vmware.com/vsphere-60/topic/com.vmware.vsphere.networking.doc/GUID-D21B3241-0AC9-437C-80B1-0C8043CC1D7D.html) in the vSphere 6.0 documentation.
 - For information about how to add hosts to a distributed virtual switch, see [Add Hosts to a vSphere Distributed Switch](https://pubs.vmware.com/vsphere-60/topic/com.vmware.vsphere.networking.doc/GUID-E90C1B0D-82CB-4A3D-BE1B-0FDCD6575725.html) in the vSphere 6.0 documentation.