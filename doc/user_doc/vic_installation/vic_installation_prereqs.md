# Environment Prerequisites for vSphere Integrated Containers Engine Installation #

Before you install vSphere Integrated Containers Engine, you must ensure that your infrastructure meets certain requirements.

## Supported Platforms for `vic-machine` ##

The vSphere Integrated Containers Engine installation and management utility, `vic-machine`, has been tested and verified on the following  64-bit OS, Windows, Mac OS, and Photon OS systems.

|**Platform**|**Supported Versions**|
|---|---|
|Windows|7, 10|
|Mac OS X |10.11 (El Capitan)|
|Linux|Ubuntu 16.04 LTS|

Other recent 64-bit OS versions should work but are untested.

## Supported vSphere Configurations ##

You can install vSphere Integrated Containers Engine in the following vSphere setups:

* Standalone ESXi 6.0 or 6.5 host that is not managed by a vCenter Server instance.
* vCenter Server 6.0 or 6.5, managing one or more standalone ESXi 6.0 or 6.5 hosts.
* vCenter Server 6.0 or 6.5, managing a cluster of ESXi 6.0 or 6.5 hosts, with VMware vSphere Distributed Resource Scheduler&trade; (DRS) enabled.

Caveats and limitations:

- Deploying vSphere Integrated Containers Engine to vSphere 5.5 environments works but is unsupported.
- Deploying vSphere Integrated Containers Engine to a nested ESXi host, namely ESXi running in a virtual machine, is not supported in production environments. Deploying vSphere Integrated Containers Engine to a nested ESXi host is acceptable for testing purposes only.
- Deploying vSphere Integrated Containers Engine to a vCenter Server instance that is running in Enhanced Linked Mode is fully supported.  

## ESXi Host Requirements ##

To be valid targets for virtual container hosts and container VMs, standalone ESXi hosts and all ESXi hosts in vCenter Server clusters must meet the following criteria:

- In vCenter Server clusters, at least two ESXi hosts must be attached to shared storage for use as container stores, image stores, and volume stores. Using non-shared datastores is possible, but limits the use of vSphere features such as DRS and High Availability. The use of VMware vSAN&trade; datastores is fully supported.
- The firewall on all ESXi hosts must be configured to allow connections on the back channel and to allow outbound connections on port 2377. For instruction about how to open port 2377 on ESXi hosts, see [VCH Deployment Fails with Firewall Validation Error](ts_firewall_error.md).
- All ESXi hosts must be attached to the distributed virtual switch for the bridge network in vCenter Server. For more information about distributed virtual switches, see [Network Requirements](#networkreqs) below.
- All ESXi hosts must be attached to any mapped vSphere networks.

During deployment of virtual container hosts, the `vic-machine` utilty checks that the target ESXi hosts meet the requirements, and issues warnings if they do not.

## License Requirements ##
The type of license that vSphere Integrated Containers Engine requires depends on the way in which you deploy the software.

| **Type of Installation** | **vSphere Feature Used** | **Required License** |
| --- | --- | --- |
|ESXi host|Network Serial Port|vSphere Enterprise|
|vCenter Server|Distributed Virtual Switch|vSphere Enterprise Plus|

All of the ESXi hosts in a cluster require an appropriate license. Installation fails if your environment includes one or more ESXi hosts that have inadequate licenses.

## Role and Permissions Requirements
You must use an account with the vSphere Administrator role when you install vSphere Integrated Containers Engine.

<a name="networkreqs"></a>
## Network Requirements
* Use a trusted network for the deployment and use of vSphere Integrated Containers Engine.
* Use a trusted network for connections between Docker clients and the virtual container hosts.
* Each virtual container host requires the following network configuration:
 * An IP address that can be either static or obtained by DHCP.
 * A network for use as the public network. You can share this network between multiple virtual container hosts.
* In vCenter Server environment, before you deploy a virtual container host, you must create a bridge network for use by container VMs.
 - Create a distributed virtual switch with a distributed port group for each virtual container host, for use as the bridge network. You can create multiple port groups on the same distributed virtual switch, but each virtual container host requires its own port group.  For information about how to create a distributed virtual switch and a distributed port group, see [Create a vSphere Distributed Switch](https://pubs.vmware.com/vsphere-60/topic/com.vmware.vsphere.networking.doc/GUID-D21B3241-0AC9-437C-80B1-0C8043CC1D7D.html) in the vSphere  documentation. 
 - Add the target ESXi host or hosts to the distributed virtual switch. For information about how to add hosts to a distributed virtual switch, see [Add Hosts to a vSphere Distributed Switch](https://pubs.vmware.com/vsphere-60/topic/com.vmware.vsphere.networking.doc/GUID-E90C1B0D-82CB-4A3D-BE1B-0FDCD6575725.html) in the vSphere  documentation.
 - If you are not using private VLANs, assign a VLAN ID to the port group, to ensure that the bridge network is isolated. For information about how to assign a VLAN ID to a port group, see [VMware KB 1003825](https://kb.vmware.com/kb/1003825). For more information about private VLAN, see [VMware KB 1010691](https://kb.vmware.com/kb/1010691).

