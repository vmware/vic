# Examples of Deploying a Virtual Container Host #

This topic provides examples of the options of the `vic-machine` `create` command to use when deploying virtual container hosts in different vSphere configurations.

- [Deploy a Virtual Container Host Directly on an ESXi Host with no Resource Pools and a Single Datastore](#esxi)
- [Deploy a Virtual Container Host in a vCenter Server Cluster](#cluster)
- [Deploy a Virtual Container Host to a Cluster and Specify External, Management, Client, and Container Networks](#networks)
- [Deploy a Virtual Container Host and Configure a Non-DHCP Container Network](#ip-range)
- [Deploy a Virtual Container Host and Specify an Image Store Folder and a Different Datastore for Container VM Files](#datastores)
- [Deploy a Virtual Container Host and Specify One or More Volume Stores](#volume-stores)
- [Deploy a Virtual Container Host on a Standalone Host in vCenter Server](#standalone)
- [Deploy a Virtual Container Host in a Resource Pool on an ESXi Host](#rp_host)
- [Deploy a Virtual Container Host in a Resource Pool in a vCenter Server Cluster](#rp_cluster)
- [Deploy a Virtual Container Host for Use with a Docker Client that Requires CA Certificates](#ca_cert)
- [Deploy a Virtual Container Host with Limits on Resource Use](#customized)
- [Deploy a Virtual Container Host and Authorize Access to an Insecure Registry](#registry)

For detailed descriptions of all of the `vic-machine create` options, see [Virtual Container Host Deployment Options](vch_installer_options.md)

<a name="esxi"></a>
## Deploy a Virtual Container Host Directly on an ESXi Host with no Resource Pools and a Single Datastore##

You can deploy a virtual container host directly on an ESXi host that is not managed by a vCenter Server instance. This example provides the minimum options required to deploy a virtual container host. The `vic-machine create` command prompts you for the password for the ESXi host and deploys a virtual container host with the default name `virtual-container-host`. If there is only one datastore on the host and there are no resource pools, you do not need to specify the `image-store` or `compute-resource` options. When deploying to an ESXi host, `vic-machine create` creates a standard virtual switch and a distributed port group, so you do not need to specify any network options if you do not have specific network requirements.

<pre>vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> create
--target <i>esxi_host_address</i>
--user root
</pre>

<a name="cluster"></a>
## Deploy a Virtual Container Host in a vCenter Server Cluster ##

If vCenter Server has more than one datacenter, you specify the datacenter in the `target` option.

If vCenter Server manages more than one cluster, you use the `compute-resource` option to specify the cluster on which to deploy the virtual container host.

When deploying a virtual container host to vCenter Server, you must use the `bridge-network` option to specify an existing distributed port group for container VMs to use to communicate with each other. For information about how to create a distributed virtual switch and port group, see *Network Requirements* in [Environment Prerequisites for vSphere Integrated Containers Engine Installation](vic_installation_prereqs.md#networkreqs).

This example deploys a virtual container host with the following configuration:

- Provides the vCenter Single Sign-On user and password in the `target` option. Note that the user name is wrapped in quotes, because it contains the `@` character. Use single quotes if you are using `vic-machine` on a Linux or Mac OS system and double quotes on a Windows system. 
- Deploys a virtual container host named `vch1` to the cluster `cluster1` in datacenter `dc1`. 
- Uses a distributed port group named `vic-bridge` for the bridge network. 
- Designates `datastore1` as the datastore in which to store container images and the files for the virtual container host appliance. 

<pre>vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> create
--target 'Administrator@vsphere.local':<i>password</i>@<i>vcenter_server_address</i>/dc1
--compute-resource cluster1
--image-store datastore1
--bridge-network vic-bridge
--name vch1
</pre>

<a name="networks"></a>
## Deploy a Virtual Container Host to a Cluster and Specify External, Management, Client, and Container Networks ##

In addition to the mandatory bridge network, if your vCenter Server environment includes multiple networks, you can direct different types of traffic to different networks. 

- You can direct the traffic between the virtual container host, container VMs, and the internet to a specific network by specifying the `external-network` option. If you do not specify the `external-network` option, the virtual container host uses the default VM Network for external traffic.
- You can direct traffic between ESXi hosts, vCenter Server, and the virtual container host to a specific network by specifying the `management-network` option. If you do not specify the `management-network` option, the virtual container host uses the external network for management traffic.
- You can designate a specific network for use by the Docker API by specifying the `client-network` option. If you do not specify the `client-network` option, the Docker API uses the external network.
- You can designate a specific network for container VMs to use by specifying the `container-network` option. Containers use this network if the container developer runs `docker run` or `docker create` with the `--net` option when they run or create a container. This option requires a distributed port group that must exist before you run `vic-machine create`. You cannot use the same distributed port group that you use for the bridge network. You can provide a descriptive name for the network, for use by Docker. If you do not specify a descriptive name, Docker uses the vSphere network name. For example, the descriptive name appears as an available network in the output of `docker info`. 

This example deploys a virtual container host with the following configuration:

- Specifies the user name, password, datacenter, cluster, image store, bridge network, and name for the virtual container host.
- Directs external, management, and Docker API traffic to network 1, network 2, and network 3 respectively. Note that the network names are wrapped in quotes, because they contain spaces. Use single quotes if you are using `vic-machine` on a Linux or Mac OS system and double quotes on a Windows system.
- Designates a distributed port group named `vic-containers` for use by container VMs that are run with the `--net` option.
- Gives the container network the name `vic-container-network`, for use by Docker.  

<pre>vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> create
--target 'Administrator@vsphere.local':<i>password</i>@<i>vcenter_server_address</i>/dc1
--compute-resource cluster1
--image-store datastore1
--bridge-network vic-bridge
--external-network 'network 1'
--management-network 'network 2'
--client-network 'network 3'
--container-network vic-containers:vic-container-network
--name vch1
</pre>

For more information about the networking options, see the [Networking Options section in Virtual Container Host Deployment Options](vch_installer_options.md#networking).

<a name="ip-range"></a>
## Deploy a Virtual Container Host and Configure a Non-DHCP Container Network ##

If the network that you designate as the container network in the `container-network` option does not support DHCP, you can configure the gateway, DNS server, and a range of IP addresses for container VMs to use. 

This example deploys a virtual container host with the following configuration:

- Specifies the user name, password, datacenter, cluster, image store, bridge network, and name for the virtual container host.
- Uses the default VM Network for the external, managment, and client networks.
- Designates a distributed port group named `vic-containers` for use by container VMs that are run with the `--net` option.
- Gives the container network the name `vic-container-network`, for use by Docker. 
- Specifies the gateway, two DNS servers, and a range of IP addresses on the container network for container VMs to use.

<pre>vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> create
--target 'Administrator@vsphere.local':<i>password</i>@<i>vcenter_server_address</i>/dc1
--compute-resource cluster1
--image-store datastore1
--bridge-network vic-bridge
--container-network vic-containers:vic-container-network
--container-network-gateway vic-containers:<i>gateway_ip_address</i>/255.255.255.0
--container-network-dns vic-containers:<i>dns1_ip_address</i>
--container-network-dns vic-containers:<i>dns2_ip_address</i>
--container-network-ip-range vic-containers:192.168.100.0/24
--name vch1
</pre>

For more information about the container network options, see the [container network section in Virtual Container Host Deployment Options](vch_installer_options.md#container-network).


<a name="datastores"></a>
## Deploy a Virtual Container Host and Specify an Image Store Folder and a Different Datastore for Container VM Files ##

If your vSphere environment includes multiple datastores, you can specify different datastores for the container image files and the container VM files. If you specify the `container-store` option, the VM files for container VMs are placed in folders at the top level of the designated datastore. The virtual container host files and the container image files are placed in the datastore that you specify in the `image-store` option. If you do not specify the `container-store` option, the virtual container host uses the `image-store` datastore for container VM files.

You specify a specific folder in which to store your container images by providing a path in the `image-store` option.

**NOTE**: In the current builds of vSphere Integrated Containers Engine, the `container-store` option is not enabled. Container VM files are stored in the datastore that you designate as the image store.

This example deploys a virtual container host with the following configuration:

- Specifies the user name, password, datacenter, cluster, bridge network, and name for the virtual container host.
- Designates the folder `vic images` on `datastore 1` as the location in which to store container image files. If the `vic images` folder does not already exist, `vic-machine create` creates it.
- Designates `datastore 2` as the datastore in which to store container VM files.
- Note that the datastore and folder names are wrapped in quotes, because they contain spaces. Use single quotes if you are using `vic-machine` on a Linux or Mac OS system and double quotes on a Windows system.

<pre>vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> create
--target 'Administrator@vsphere.local':<i>password</i>@<i>vcenter_server_address</i>/dc1
--compute-resource cluster1
--bridge-network vic-bridge
--image-store 'datastore 1/vic images'
--container-store 'datastore 2'
--name vch1
</pre>

For more information about the datastore options, see the [Datastore Options section in Virtual Container Host Deployment Options](vch_installer_options.md#datastore). 

<a name="volume-stores"></a>
## Deploy a Virtual Container Host and Specify One or More Volume Stores ##

If container application developers will use the `docker volume create` command to create containers that use volumes, you must create volume stores when you deploy virtual container hosts. You specify volume stores in the `volume-store` option. You can specify `volume-store` multiple times to create multiple volume stores. 

When you create a volume store, you specify the name of the datastore to use and an optional path to a folder on that datastore. You also specify a descriptive name for that volume store for use by Docker.

This example deploys a virtual container host with the following configuration:

- Specifies the user name, password, datacenter, cluster, bridge network, and name for the virtual container host.
- Specifies the `volumes` folder on `datastore 1` as the default volume store. Creating a volume store named `default` allows container application developers to create anonymous volumes by using `docker create -v`. 
- Specifies a second volume store named `volume_store_2` in the `volumes` folder on `datastore 2`. 
- Note that the datastore names are wrapped in quotes, because they contain spaces. Use single quotes if you are using `vic-machine` on a Linux or Mac OS system and double quotes on a Windows system.

<pre>vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> create
--target 'Administrator@vsphere.local':<i>password</i>@<i>vcenter_server_address</i>/dc1
--compute-resource cluster1
--bridge-network vic-bridge
--image-store 'datastore 1'
--volume-store 'datastore 1'/volumes:default</i>
--volume-store 'datastore 2'/volumes:volume_store_2</i>
--name vch1
</pre>

For more information about volume stores, see the [volume-store section in Virtual Container Host Deployment Options](vch_installer_options.md#volume-store). 

<a name="standalone"></a> 
## Deploy a Virtual Container Host on a Standalone Host in vCenter Server ##

If vCenter Server manages multiple standalone ESXi hosts that are not part of a cluster, you use the `compute-resource` option to specify the address of the ESXi host to which to deploy the virtual container host.

This example deploys a virtual container host with the following configuration:

- Specifies the user name, password, image store, bridge network, and name for the virtual container host.
- Deploys the virtual container host on the ESXi host with the FQDN `esxihost1.organization.company.com` in the datacenter `dc1`. You can also specify an IP address.

<pre>vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> create
--target 'Administrator@vsphere.local':<i>password</i>@<i>vcenter_server_address</i>/dc1
--image-store datastore1
--bridge-network vic-bridge
--compute-resource esxihost1.organization.company.com
--name vch1
</pre>

<a name="rp_host"></a>
## Deploy a Virtual Container Host in a Resource Pool on an ESXi Host ##

To deploy a virtual container host in a specific resource pool on an ESXi host that is not managed by vCenter Server, you specify the resource pool name in the `compute-resource` option. 

This example deploys a virtual container host with the following configuration:

- Specifies the user name and password and a name for the virtual container host.
- Designates `rp 1` as the resource pool in which to place the virtual container host. Note that the resource pool name is wrapped in quotes, because it contains a space. Use single quotes if you are using `vic-machine` on a Linux or Mac OS system and double quotes on a Windows system.

<pre>vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> create
--target root:<i>password</i>@<i>esxi_host_address</i>
--compute-resource 'rp 1'
--name vch1
</pre>

<a name="rp_cluster"></a>
## Deploy a Virtual Container Host in a Resource Pool in a vCenter Server Cluster ##

To deploy a virtual container host in a resource pool in a vCenter Server cluster, you specify the names of the cluster and resource pool in the `compute-resource` option.

This example deploys a virtual container host with the following configuration:

- Specifies the user name, password, datacenter, image store, bridge network, and name for the virtual container host.
- Designates `rp 1` in cluster `cluster 1` as the resource pool in which to place the virtual container host. Note that the resource pool and cluster names are wrapped in quotes, because they contain spaces. Use single quotes if you are using `vic-machine` on a Linux or Mac OS system and double quotes on a Windows system.

<pre>vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> create
--target 'Administrator@vsphere.local':<i>password</i>@<i>vcenter_server_address</i>/dc1
--compute-resource 'cluster 1'/'rp 1'
--image-store datastore1
--bridge-network vic-bridge
--name vch1
</pre>

<a name="ca_cert"></a>
## Deploy a Virtual Container Host for Use with a Docker Client that Requires CA Certificates ##

If your development environment uses custom Certificate Authority certificates to authenticate connections between Docker clients and virtual container hosts, use the `cert` and `key` options to provide the paths to an X.509 certificate and its key when you deploy a virtual container host. The paths to the certificate and key files must be relative to the location from which you are running `vic-machine create`.

This example deploys a virtual container host with the following configuration:

- Specifies the user name, password, image store, cluster, bridge network, and name for the virtual container host.
- Provides the paths relative to the current location of the `*.pem` files for the custom CA certificate and key files.

<pre>vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> create
--target 'Administrator@vsphere.local':<i>password</i>@<i>vcenter_server_address</i>/dc1
--compute-resource cluster1
--image-store datastore1
--bridge-network vic-bridge
--cert ../some/relative/path/<i>certificate_file</i>.pem
--key ../some/relative/path/<i>key_file</i>.pem
--name vch1
</pre>

For more information about using CA certificates, see the [Security  Options section in Virtual Container Host Deployment Options](vch_installer_options.md#security).

<a name="customized"></a>
## Deploy a Virtual Container Host with Limits on Resource Use ##

To limit the amount of system resources that the container VMs in a virtual container host can use, you can set resource limits on the virtual container host vApp. 

This example deploys a virtual container host with the following configuration:

- Specifies the user name, password, image store, cluster, bridge network, and name for the virtual container host.
- Sets resource limits on the virtual container host by imposing memory and CPU reservations, limits, and shares.

<pre>vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> create
--target 'Administrator@vsphere.local':<i>password</i>@<i>vcenter_server_address</i>/dc1
--compute-resource cluster1
--image-store datastore1
--bridge-network vic-bridge
--pool-memory-reservation 1024
--pool-memory-limit 1024
--pool-memory-shares low
--pool-cpu-reservation 1024
--pool-cpu-limit 1024
--pool-cpu-shares low
--name vch1
</pre>

For more information about setting resource use limitations on virtual container hosts, see the [vApp Deployment Options section in Virtual Container Host Deployment Options](vch_installer_options.md#deployment).

<a name="registry"></a>
## Deploy a Virtual Container Host and Authorize Access to an Insecure Registry ##

An insecure registry server is a registry server for Docker images that is secured by self-signed certificates rather than by TLS. To authorize connections from a virtual container host to an insecure registry server, set the `docker-insecure-registry` option. You can specify `docker-insecure-registry` multiple times to allow connections from the virtual container host to multiple insecure registry servers.

This example deploys a virtual container host with the following configuration:

- Specifies the user name, password, image store, cluster, bridge network, and name for the virtual container host.
- Authorizes the virtual container host to pull Docker images from the insecure registry servers located at the URLs <i>registry_URL_1</i> and <i>registry_URL_2</i>.

<pre>vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> create
--target 'Administrator@vsphere.local':<i>password</i>@<i>vcenter_server_address</i>/dc1
--compute-resource cluster1
--image-store datastore1
--bridge-network vic-bridge
--docker-insecure-registry <i>registry_URL_1</i>
--docker-insecure-registry <i>registry_URL_2</i>
--name vch1
</pre>

For more information about setting resource use limitations on virtual container hosts, see the section on the `docker-insecure-registry` option in [Virtual Container Host Deployment Options](vch_installer_options.md#registry).

