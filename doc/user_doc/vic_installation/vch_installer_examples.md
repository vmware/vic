# Examples of Deploying a Virtual Container Host #

This topic provides examples of the options of the `vic-machine` `create` command to use when deploying virtual container hosts in different vSphere configurations.

- [Deploy a Virtual Container Host Directly on an ESXi Host with no Resource Pools](#esxi)
- [Deploy a Virtual Container Host in a vCenter Server Cluster](#cluster)
- [Deploy a Virtual Container Host to a Cluster and Specify External, Management, and Container Networks](#networks)
- [Deploy a Virtual Container Host and Specify a Different Datastore for Container VM Files](#datastores)
- [Deploy a Virtual Container Host on a Standalone Host in vCenter Server](#standalone)
- [Deploy a Virtual Container Host in a Resource Pool on an ESXi Host](#rp_host)
- [Deploy a Virtual Container Host in a Resource Pool in a vCenter Server Cluster](#rp_cluster)
- [Deploy a Virtual Container Host for Use with a Docker Client that Requires CA Certificates](#ca_cert)
- [Deploy a Virtual Container Host with no TLS Authentication](#no_tls)
- [Deploy a Virtual Container Host with Customized VM Settings](#customized)

<a name="esxi"></a>
## Deploy a Virtual Container Host Directly on an ESXi Host with no Resource Pools##

You can install vSphere Integrated Containers directly on an ESXi host that is not managed by a vCenter Server instance. This example provides the miniumum options required to install vSphere Integrated Containers. The installer prompts you for the password for the ESXi host and deploys a virtual container host appliance with the default name `virtual-container-host`. You specify the datastore in which to store the virtual container host files and container image files in the mandatory `image-store` option.

<pre>vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> create
--target <i>esxi_host_IPv4_address_or_FQDN</i>
--user root
--image-store <i>datastore_name</i>
</pre>

<a name="cluster"></a>
## Deploy a Virtual Container Host in a vCenter Server Cluster ##

If vCenter Server manages more than one cluster, you use the `compute-resource` option to specify the address of the cluster to which to deploy the virtual container host. 

You must use the `bridge-network` option to specify an existing distributed port group for container VMs to use to communicate with each other and with the virtual container host.  

In addition to the mandatory options for deployment to a cluster, this example sets the vCenter Single Sign-On user and password and gives the virtual container host a custom name, `vch1`.

<pre>vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> create
--target <i>vcenter_server_IPv4_address_or_FQDN</i>
--user Administrator@vsphere.local
--password <i>vcenter_sso_password</i>
--compute-resource <i>cluster_name</i>
--image-store <i>datastore_name</i>
--bridge-network <i>distributed_port_group_name</i>
--name vch1
</pre>

<a name="networks"></a>
## Deploy a Virtual Container Host to a Cluster and Specify External, Management, and Container Networks ##

If your vSphere environment includes multiple networks, you can direct different types of traffic to different networks. 

- You direct the traffic between the virtual container host, container VMs, and the internet to a specific network by specifying the `external-network` option. If you do not specify the `external-network` option, the virtual container host uses the default VM Network for external traffic.
- You direct traffic between ESXi hosts, vCenter Server, and the virtual container host to a specific network by specifying the `management-network` option. If you do not specify the `management-network` option, the virtual container host uses the bridge network for management traffic.
- You direct the traffic between container VMs and between container VMs and the virtual container host to a specific network by specifying the `bridge-network` option. This option is mandatory in environments with more than one ESXi host. 

To specify networking options, you must provide the name of a distributed port group that exists in your vCenter Server environment before you deploy a virtual container host. For information about how to create a distributed virtual switch and port group, see *Network Requirements* in [Environment Prerequisites for vSphere Integrated Containers Installation](vic_installation_prereqs.md#networkreqs).

In addition to the mandatory options for deployment to a cluster, this example sets the vCenter Single Sign-On user and password, specifies different networks for the different types of traffic, and gives the virtual container host a custom name, `vch1`.

<pre>vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> create
--target <i>vcenter_server_IPv4_address_or_FQDN</i>
--user Administrator@vsphere.local
--password <i>vcenter_sso_password</i>
--compute-resource <i>cluster_name</i>
--image-store <i>datastore_name</i>
--bridge-network <i>network_1_name</i>
--management-network <i>network_2_name</i>
--external-network <i>network_3_name</i>
--name vch1
</pre>

<a name="datastores"></a>
## Deploy a Virtual Container Host and Specify a Different Datastore for Container VM Files ##

If your vSphere environment includes multiple datastores, you can specify different datastores for the container image files and the container VM files. If you specify the `container-store` option, the VM files for container VMs are placed in folders at the top level of the designated datastore. The virtual container host files and the container image files are placed in the datastore that you specify in the `image-store` option. If you do not specify the `container-store` option, the virtual container host uses the `image-store` datastore for container VM files.

In addition to the mandatory options for deployment to a cluster, this example sets the vCenter Single Sign-On user and password, specifies a different datastore for container VM files, and gives the virtual container host a custom name, `vch1`.

<pre>vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> create
--target <i>vcenter_server_IPv4_address_or_FQDN</i>
--user Administrator@vsphere.local
--password <i>vcenter_sso_password</i>
--compute-resource <i>cluster_name</i>
--image-store <i>datastore_1_name</i>
--container-datastore <i>datastore_2_name</i>
--bridge-network <i>network_1_name</i>
--name vch1
</pre>

It is also possible to specify a specific folder in which to store your images, instead of just a datastore. This can be done by providing a path following the `image-store` option:

<pre>vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> create
--target <i>vcenter_server_IPv4_address_or_FQDN</i>
--user Administrator@vsphere.local
--password <i>vcenter_sso_password</i>
--compute-resource <i>cluster_name</i>
--image-store <i>datastore_1_name/path/to/some/specific/location</i>
--container-datastore <i>datastore_2_name</i>
--bridge-network <i>network_1_name</i>
--name vch1
</pre>

<a name="standalone"></a> 
## Deploy a Virtual Container Host on a Standalone Host in vCenter Server ##

If vCenter Server manages multiple standalone ESXi hosts that are not part of a cluster, you use the `compute-resource` option to specify the address of the ESXi host to which to deploy the virtual container host.   If the vCenter Server manages more than one standalone ESXi host, you must specify the `bridge-network` option.

In addition to the mandatory options, this example specifies the vCenter Single Sign-On user and password, and gives the virtual container host a custom name, `vch1`.

<pre>vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> create
--target <i>vcenter_server_IPv4_address_or_FQDN</i> 
--user Administrator@vsphere.local
--password <i>vcenter_sso_password</i>
--compute-resource <i>host_IPv4_address_or_FQDN</i>
--image-store <i>datastore_name</i>
--bridge-network <i>network_name</i>
--name vch1
</pre>

<a name="rp_host"></a>
## Deploy a Virtual Container Host in a Resource Pool on an ESXi Host ##

To deploy a virtual container host in a resource pool on an ESXi host that is not managed by vCenter Server, you specify the resource pool address in the `compute-resource` option. 

This example uses the minimum required options.

<pre>vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> create
--target <i>esxi_host_IPv4_address_or_FQDN</i>
--user root
--compute-resource <i>resource_pool_name</i>
--image-store <i>datastore_name</i>
</pre>

<a name="rp_cluster"></a>
## Deploy a Virtual Container Host in a Resource Pool in a vCenter Server Cluster ##

To deploy a virtual container host in a resource pool in a vCenter Server Cluster, you specify the resource pool address in the `compute-resource` option.

In addition to the mandatory options, this example sets the vCenter Single Sign-On user and password, and gives the virtual container host a custom name, `vch1`.

<pre>vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> create
--target <i>vcenter_server_IPv4_address_or_FQDN</i>
--user Administrator@vsphere.local
--password <i>vcenter_sso_password</i>
--compute-resource <i>cluster_name</i>/<i>resource_pool_name</i>
--image-store <i>datastore_name</i>
--bridge-network <i>network_name</i>
--name vch1
</pre>

<a name="ca_cert"></a>
## Deploy a Virtual Container Host for Use with a Docker Client that Requires CA Certificates ##

If your development environment uses custom CA certificates to authenticate connections between Docker clients and virtual container hosts, use the `cert` and `key` options to provide the paths to an X.509 certificate and its key when you deploy a virtual container host. The paths to the certificate and key files must be relative to the location from which you are deploying a virtual container host.

In addition to the mandatory options for deployment to a cluster, this example sets the vCenter Single Sign-On user and password, and gives the virtual container host a custom name, `vch1`.

<pre>vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> create
--target <i>vcenter_server_IPv4_address_or_FQDN</i>
--user Administrator@vsphere.local
--password <i>vcenter_sso_password</i>
--compute-resource <i>cluster_name</i>
--image-store <i>datastore_name</i>
--bridge-network <i>network_name</i>
--cert <i>path_to_certificate_file</i>
--key <i>path_to_certificate_file</i>
--name vch1
</pre>

<a name="no_tls"></a>
## Deploy a Virtual Container Host with no TLS Authentication ##

If your development environment does not require authentication of the connections between Docker clients and virtual container hosts, use the `no-tls` option to deactivate TLS in the virtual container host. If you do not specify `no-tls` when you deploy a virtual container host, vic-machine automatically generates a certificate and key file that you must provide to the Docker client when you connect it to this virtual container host.

In addition to the mandatory options for deployment to a cluster, this example sets the vCenter Single Sign-On user and password, and gives the virtual container host a custom name, `vch1`.

<pre>vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> create
--target <i>vcenter_server_IPv4_address_or_FQDN</i>
--user Administrator@vsphere.local
--password <i>vcenter_sso_password</i>
--compute-resource <i>cluster_name</i>
--image-store <i>datastore_name</i>
--bridge-network <i>network_name</i>
--name vch1
--no-tls
</pre>

<a name="customized"></a>
## Deploy a Virtual Container Host with Customized VM Settings ##

If a virtual container host is likely to handle heavy container workloads, you can increase the number of CPUs and the amount of memory in the virtual container host VM by setting the `appliance-cpu` and `appliance-memory` options. For example, set these options if the virtual container host will handle large volumes of containers, containers that consume a lot of memory, or containers that require a lot of processing power. If these options are not set, `vic-machine` deploys a virtual container host with 1 CPU and 2048 MB of memory.

In addition to the mandatory options for deployment to a cluster, this example sets the vCenter Single Sign-On user and password and gives the virtual container host a custom name, `vch1`.

<pre>vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> create
--target <i>vcenter_server_IPv4_address_or_FQDN</i>
--user Administrator@vsphere.local
--password <i>vcenter_sso_password</i>
--compute-resource <i>cluster_name</i>
--image-store <i>datastore_name</i>
--bridge-network <i>network_name</i>
--name vch1
--appliance-cpu <i>number_of_CPUs</i>
--appliance-memory <i>amount_of_memory</i>
</pre>
