# vSphere Integrated Containers Installation Prerequisites

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

In all cases, your ESXi hosts must have at least 8GB of RAM.

## License Requirements
vSphere Integrated Containers uses virtual serial ports to connect the virtual container hosts to the ESXi hosts. Virtual serial ports are available with the following vSphere licenses:

* vSphere Enterprise
* vSphere Enterprise Plus 

All of the ESXi hosts in your environment require an appropriate license. Installation fails if your environment includes one or more ESXi hosts that have inadequate licenses.

## Role and Permissions Requirements
You must use an account with the vSphere Administrator role when you install vSphere Integrated Containers. **(TO BE CONFIRMED)**

## Network Requirements
* Use a trusted network for the deployment and use of vSphere Integrated Containers.
* Use a trused network for connections between Docker clients and the virtual container hosts.
* Use a Gigabit connection between the machine on which you run the command line installer and the vCenter Server or ESXi hosts on which you are installing vSphere Integrated Containers. **(TO BE CONFIRMED)**
* Use a Gigabit connection between vCenter Server and the ESXi hosts, and between the ESXi hosts and the virtual container hosts. **(TO BE CONFIRMED)**
* Open port 2377 on all ESXi hosts, for communication between the hosts and virtual container hosts.
 * For instructions about how to open a port on an ESXi host, see [VMware KB 2008226]( http://kb.vmware.com/selfservice/microsites/search.do?language=en_US&cmd=displayKC&externalId=2008226).
 * In test environments, you can disable the firewalls on the ESXi hosts instead of opening port 2377. To disable the firewall, log into the ESXi hosts as ```root```, and run the following command: 
   
   ```$ esxcli network firewall set --enabled false.```  

* Create a private port group for the virtual container hosts. For instructions about how to create a private port group, see **LINK TO NEW TOPIC**.


