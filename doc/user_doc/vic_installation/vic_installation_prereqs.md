# Environment Prerequisites for VCH Deployment #

Before you install vSphere Integrated Containers Engine, you must ensure that your vSphere infrastructure meets the requirements.

- [Supported Platforms for `vic-machine`](#platforms)
- [Supported vSphere Configurations](#configs)
- [License Requirements](#license)
- [ESXi Host Firewall Requirements](#firewall)
- [ESXi Host Storage Requirements for vCenter Server Clusters](#storage)
- [General Network Requirements](#general_networkreqs)
- [vCenter Server Network Requirements](#networkreqs)

<a name="platforms"></a>
## Supported Platforms for `vic-machine` ##

The vSphere Integrated Containers Engine installation and management utility, `vic-machine`, has been tested and verified on the following  64-bit Windows, Mac OS, and Linux OS systems.

|**Platform**|**Supported Versions**|
|---|---|
|Windows|7, 10|
|Mac OS X |10.11 (El Capitan)|
|Linux|Ubuntu 16.04 LTS|

Other recent 64-bit OS versions should work but are untested.

<a name="configs"></a>
## Supported vSphere Configurations ##

You can install vSphere Integrated Containers Engine in the following vSphere setups:

* vCenter Server 6.0 or 6.5, managing a cluster of ESXi 6.0 or 6.5 hosts, with VMware vSphere Distributed Resource Scheduler&trade; (DRS) enabled.
* vCenter Server 6.0 or 6.5, managing one or more standalone ESXi 6.0 or 6.5 hosts.
* Standalone ESXi 6.0 or 6.5 host that is not managed by a vCenter Server instance.

Caveats and limitations:

- VMware does not support the use of nested ESXi hosts, namely running ESXi in virtual machines. Deploying vSphere Integrated Containers Engine to a nested ESXi host is acceptable for testing purposes only.
- Deploying vSphere Integrated Containers Engine to a vCenter Server instance that is running in Enhanced Linked Mode is fully supported. 

<a name="license"></a>
## License Requirements ##
vSphere Integrated Containers Engine requires a vSphere Enterprise Plus license.

All of the ESXi hosts in a cluster require an appropriate license. Installation fails if your environment includes one or more ESXi hosts that have inadequate licenses. 

<a name="firewall"></a>
## ESXi Host Firewall Requirements ##

To be valid targets for virtual container hosts (VCHs) and container VMs, ESXi hosts must have the following firewall configuration:
- Allow outbound TCP traffic to port 2377 on the endpoint VM, for use by the interactive container shell.
- Allow inbound HTTPS/TCP traffic on port 443, for uploading to and downloading from datastores.

These requirements apply to standalone ESXi hosts and to ESXi hosts in vCenter Server clusters.

For information about how to open ports on ESXi hosts, see [VCH Deployment Fails with Firewall Validation Error](ts_firewall_error.md).

<a name="storage"></a>
## ESXi Host Storage Requirements for vCenter Server Clusters ##

ESXi hosts in vCenter Server clusters must meet the following storage requirements in order to be usable by a VCH:
- Be attached to the datastores that you will use for image stores and volume stores. 
- Have access to shared storage to allow VCHs to use more than one host in the cluster.

For information about image stores and volumes stores, see the [Datastore Options](vch_installer_options.md#datastore) section of *VCH Deployment Options*.

<a name="general_networkreqs"></a>
## General Network Requirements ##

The following network requirements apply to deployment of VCHs to standalone ESXi hosts and to vCenter Server:

- Use a trusted network for the deployment and use of vSphere Integrated Containers Engine.
- Use a trusted network for the management network.
- Connections between Docker clients and the VCH are encrypted via TLS unless you explicitly disable TLS. The client network does not need to be trusted.
- Each VCH requires an IPv4 address on each of the networks that it is connected to. The bridge network is handled internally, but other interfaces must have a static IP configured on them, or be able to acquire one via DHCP.
- Each VCH requires access to at least one network, for use as the public network. You can share this network between multiple VCHs. The public network does not need to be trusted.
- Do not share the bridge network interface with with any other network, unless you ensure that the bridge IP ranges do not conflict with other VCHs or VMs on that network. It is highly recommended that a bridge network be solely for use by only one VCH.

<a name="networkreqs"></a>
## vCenter Server Network Requirements ##
The following network requirements apply to the deployment of VCHs to vCenter Server: 
 
- Create a distributed virtual switch with a port group for each VCH, for use as the bridge network. You can create multiple port groups on the same distributed virtual switch, but each VCH requires its own port group for the bridge network.
- Optionally create port groups for use as mapped container networks.  
- All hosts in a cluster must be attached to the port groups that you will use for the VCH bridge network and for any mapped container networks.
- If you are not using private VLANs, assign a VLAN ID to the port groups, to ensure that the bridge network and mapped container networks are isolated. 

For information about bridge networks and container networks, see the [`--bridge-network`](vch_installer_options.md#bridge) and [`--container-network`](vch_installer_options.md#container-network) options in *VCH Deployment Options*. 

For information about how to create a distributed virtual switch and a port group, see [Create a vSphere Distributed Switch](https://pubs.vmware.com/vsphere-65/topic/com.vmware.vsphere.networking.doc/GUID-D21B3241-0AC9-437C-80B1-0C8043CC1D7D.html) in the vSphere  documentation. 

For information about how to add hosts to a distributed virtual switch, see [Add Hosts to a vSphere Distributed Switch](https://pubs.vmware.com/vsphere-65/topic/com.vmware.vsphere.networking.doc/GUID-E90C1B0D-82CB-4A3D-BE1B-0FDCD6575725.html) in the vSphere  documentation.

For information about how to assign a VLAN ID to a port group, see [VMware KB 1003825](https://kb.vmware.com/kb/1003825). For more information about private VLAN, see [VMware KB 1010691](https://kb.vmware.com/kb/1010691).


<!--
## Role and Permissions Requirements
You must use an account with the vSphere Administrator role when you install vSphere Integrated Containers Engine.
-->