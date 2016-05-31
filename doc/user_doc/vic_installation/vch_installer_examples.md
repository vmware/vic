# Examples of Deploying a Virtual Container Host #

This document provides examples of the required command line options when deploying virtual container hosts in different vSphere configurations.

## Deploy a Virtual Container Host Directly on an ESXi Host with no Resource Pools##

You can install vSphere Integrated Containers directly on an ESXi host that is not managed by a vCenter Server instance. This example provides the miniumum options required to install vSphere Integrated Containers. The installer prompts you for the password for the ESXi host and deploys a virtual container host appliance with the default name `docker-appliance`. When an ESXi host has no resource pools, you must provide the root resource pool in the `compute-resource` option, exactly as shown in the example.
  
<pre>vic-machine
-target <i>esxi_host_IPv4_address_or_FQDN</i>
-user root
-compute-resource /ha-datacenter/host/localhost.eng.vmware.com/Resources/
-image-store <i>datastore_name</i>
</pre>

## Deploy a Virtual Container Host in a vCenter Server Cluster ##

If vCenter Server manages clusters, you use the `compute-resource` option to specify the name of the cluster to which to deploy the virtual container host. In clusters that have more than one host, you must use the `-bridge-network` option to specify an existing network for container VMs to use to communicate with each other. The example sets the vCenter Single Sign-On user and password, and gives the vSphere Integrated Containers appliance a custom name, `vch1`.
<pre>vic-machine
-target <i>vcenter_server_IPv4_address_or_FQDN</i>
-user Administrator@vsphere.local
-passwd <i>vcenter_sso_password</i>
-compute-resource /<i>datacenter_name</i>/host/<i>cluster_name</i>/
-image-store <i>datastore_name</i>
-bridge-network <i>network_name</i>
-name vch1
</pre>


## Deploy a Virtual Container Host and Specify External, Management, and Container Networks ##

If your vSphere environment includes multiple networks, or if you need to send all traffic between vCenter Server and the virtual container host over a specific network, you must specify the `externalNetwork` option. If you do not specify the `externalNetwork` option, the vSphere Integrated Containers uses the default VM Network for external traffic. 

If you want to send all container-related traffic over a specific network, specify the `containerNetwork` option. The value of the `containerNetwork` option is the name of a private port group that you  must create on the ESXi host before you install vSphere Integrated Containers. For information about creating a private port group, see [Create a Private Port Group for Virtual Container Hosts](create_a_private_port_group_for_vch.md). 

<pre>vic-machine
-target <i>vcenter_server_IPv4_address_or_FQDN</i>
-user Administrator@vsphere.local
-passwd <i>vcenter_sso_password</i>
-compute-resource /<i>datacenter_name</i>/host/<i>cluster_name</i>/
-image-store <i>datastore_name</i>
-bridge-network <i>network_name</i>
-name vch1
</pre>

 
## Deploy a Virtual Container Host on a Standalone Host in vCenter Server ##

If vCenter Server only manages a single host and that host contains a single datastore, you do not need to specify the `host` or `datastore` options. This example also specifies the vCenter Single Sign-On user and password, and gives the vSphere Integrated Containers appliance a custom name.

<pre>vic-machine
-target <i>vcenter_server_IPv4_address_or_FQDN</i>
-user
-compute-resource
-image-store <i>datastore_name</i>
</pre>

## Deploy a Virtual Container Host in a Resource Pool on an ESXi Host ##

If a standalone host contains resource pools, and you want to install vSphere Integrated Containers in a resource pool, you must specify the `pool` option. This example specifies a datacenter, a host in that datacenter, and a datastore and resource pool on that host. The example also sets the vCenter Single Sign-On user and password, and gives the vSphere Integrated Containers appliance a custom name.

<pre>vic-machine
-target <i>esxi_host_IPv4_address_or_FQDN</i>
-user root
-compute-resource /ha-datacenter/host/localhost.eng.vmware.com/Resources/resource_pool_name/
-image-store <i>datastore_name</i>
</pre>

## Deploy a Virtual Container Host in a Resource Pool in a vCenter Server Cluster ##

If a cluster contains resource pools, and you want to install vSphere Integrated Containers in a resource pool, you must specify the `pool` option. This example specifies a datacenter, a cluster in that datacenter, and a datastore and resource pool in that cluster. The example also sets the vCenter Single Sign-On user and password, and gives the vSphere Integrated Containers appliance a custom name.

<pre>vic-machine
-target <i>vcenter_server_IPv4_address_or_FQDN</i>
-user
-compute-resource
-image-store <i>datastore_name</i>
</pre>

## Deploy a Virtual Container Host for Use with a Docker Client that Uses CA Certificates ##

If vCenter Server verifies extensions by using TLS certificates, you must provide the paths to an X.509 certificate and its key during installation. This example of a simple installation includes the `cert` and `key` options. Specify a path to the certificate and key files that is relative to the location from which you are installing vSphere Integrated Containers.

<pre>vic-machine
-target <i>vcenter_server_IPv4_address_or_FQDN</i>
-user
-compute-resource
-image-store <i>datastore_name</i>
</pre>
