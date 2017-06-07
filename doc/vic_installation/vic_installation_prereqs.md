# Environment Prerequisites for vSphere Integrated Containers Installation

Before you install vSphere Integrated Containers, you must ensure that your vSphere infrastructure meets certain requirements.


## Supported vSphere Configurations
 The required vSphere configurations are different depending on whether you use the command line installer or the OVA installer.

You can install vSphere Integrated Containers in the following vSphere setups:

* Command Line Installer (Default):
 * Standalone ESXi 6.0 host that is not managed by a vCenter Server instance.
 * vCenter Server 6.0, managing a single ESXi 6.0 or 5.5 host.
 * vCenter Server 6.0, managing a cluster of ESXi 6.0 or 5.5 hosts, with DRS enabled.
* OVA Installer (Experimental):
 * vCenter Server 6.0, managing a single ESXi 6.0 or 5.5 host.
 * vCenter Server 6.0, managing a cluster of ESXi 6.0 or 5.5 hosts, with DRS enabled.

In all cases, your ESXi hosts must have at least 8GB of memory. 

## License Requirements
The type of license that vSphere Integrated Containers requires depends on the way in which you deploy the software.

| **Type of Installation** | **vSphere Feature Used** | **Required License** |
| --- | --- | --- |
|Standalone ESXi host or vCenter Server managing a single ESXi host|Network Serial Port|vSphere Enterprise|
|vCenter Server managing a cluster|Distributed Virtual Switch|vSphere Enterprise Plus|


All of the ESXi hosts in a cluster require an appropriate license. Installation fails if your environment includes one or more ESXi hosts that have inadequate licenses.

## Role and Permissions Requirements
You must use an account with the vSphere Administrator role when you install vSphere Integrated Containers.

## Network Requirements
* Use a trusted network for the deployment and use of vSphere Integrated Containers.
* Use a trused network for connections between Docker clients and the virtual container hosts.
* Use a Gigabit connection between the machine on which you run the command line installer and the vCenter Server or ESXi hosts on which you are installing vSphere Integrated Containers.
* Use a Gigabit connection between vCenter Server and the ESXi hosts, and between the ESXi hosts and the virtual container hosts.
* Open outgoing port 2377 on all ESXi hosts, for communication between the hosts and virtual container hosts.
 * For an example of how to open port 2377 on an ESXi host, see [Open an Outgoing Port on ESXi Hosts](open_an_outgoing_port_on_esxi_hosts.md).
 * In test environments, you can disable the firewalls on the ESXi hosts instead of opening port 2377. To disable the firewall, log into the ESXi hosts as ```root```, and run the following command: 
   
   ```$ esxcli network firewall set --enabled false```  

* Each virtual container host requires the following network configuration:
 * One IP address, that can be either static or set by using DHCP.
 * One VLAN, if you use VLAN for network isolation.
 * One IP address for each container that you run with the `docker run --net=host` option.
* Deployment to a vCenter Server cluster requires a private port group for each virtual container host. 
 * If you install vSphere Integrated Containers by using the command line installer, you can either create a custom private port group in advance, or the installer can create a default port group. 
 * If you install vSphere Integrated Containers by using the OVA deployment, you must create a private port group in advance. 
 
 For instructions about how to create a private port group, see [Create a Private Port Group for Virtual Container Hosts](create_a_private_port_group_for_vch.md).


