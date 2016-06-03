# Environment Prerequisites for vSphere Integrated Containers Installation

Before you install vSphere Integrated Containers, you must ensure that your infrastructure meets certain requirements.

## Supported Installation Environments

The current builds of vSphere Integrated Containers run on Linux OS, Windows, and Mac OS systems.

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
|Standalone ESXi host or vCenter Server managing standalone ESXi hosts|Network Serial Port|vSphere Enterprise|
|vCenter Server managing a cluster|Distributed Virtual Switch|vSphere Enterprise Plus|

All of the ESXi hosts in a cluster require an appropriate license. Installation fails if your environment includes one or more ESXi hosts that have inadequate licenses.

## Role and Permissions Requirements
You must use an account with the vSphere Administrator role when you install vSphere Integrated Containers.

## Network Requirements
* Use a trusted network for the deployment and use of vSphere Integrated Containers.
* Use a trused network for connections between Docker clients and the virtual container hosts.
* Each virtual container host requires the following network configuration:
 * One IP address, that can be either static or set by using DHCP.
 * One VLAN, if you use VLAN for network isolation.
 * One IP address for each container that you run with the `docker run --net=host` option.