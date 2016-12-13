# Advanced Examples of Deploying a VCH #

This topic provides examples of the options of the `vic-machine create` command to use when deploying virtual container hosts (VCHs) in various vSphere configurations.

- [General Deployment Examples](#general)
  - [Deploy to a vCenter Server Cluster with Multiple Datacenters and Datastores](#cluster)
  - [Deploy to a Specific Standalone Host in vCenter Server](#standalone)
  - [Deploy to a Resource Pool on an ESXi Host](#rp_host)
  - [Deploy to a Resource Pool in a vCenter Server Cluster](#rp_cluster)
  - [Set Limits on Resource Use](#customized)
- [Networking Examples](#networking)
  - [Specify Public, Management, Client, and Container Networks](#networks)
  - [Set a Static IP Address for the VCH Endpoint VM on the Different Networks](#static-ip)
  - [Configure a Non-DHCP Container Network](#ip-range)
  - [Configure a Proxy Server](#proxy)
- [Specify One or More Volume Stores](#volume-stores)
- [Security Examples](#security)
  - [Use Auto-Generated Trusted CA Certificates](#auto_cert)
  - [Use Custom Server Certificates](#custom_cert)
  - [Specify Different User Accounts for VCH Deployment and Operation](#ops-user)
  - [Authorize Access to an Insecure Private Registry Server](#registry)

For simplicity, these examples use the `--force` option to disable the verification of the vCenter Server certificate, so the `--thumbprint` option is not specified. Similarly, all examples that do not relate explicitly to certificate use specify the `--tls-noverify` option.

For detailed descriptions of all of the `vic-machine create` options, see [VCH Deployment Options](vch_installer_options.md).

<a name="general"></a>
## General Deployment Examples ##

The examples in this section demonstrate the deployment of VCHs in different vSphere environments.

<a name="cluster"></a>
### Deploy to a vCenter Server Cluster with Multiple Datacenters and Datastores ###

If vCenter Server has more than one datacenter, you specify the datacenter in the `--target` option.

If vCenter Server manages more than one cluster, you use the `--compute-resource` option to specify the cluster on which to deploy the VCH.

When deploying a VCH to vCenter Server, you must use the `--bridge-network` option to specify an existing port group for container VMs to use to communicate with each other. For information about how to create a distributed virtual switch and port group, see  [vCenter Server Network Requirements](vic_installation_prereqs.md#networkreqs) in *Environment Prerequisites for vSphere Integrated Containers Engine Installation*.

This example deploys a VCH with the following configuration:

- Provides the vCenter Single Sign-On user and password in the `--target` option. Note that the user name is wrapped in quotes, because it contains the `@` character. Use single quotes if you are using `vic-machine` on a Linux or Mac OS system and double quotes on a Windows system. 
- Deploys a VCH named `vch1` to the cluster `cluster1` in datacenter `dc1`. 
- Uses a port group named `vic-bridge` for the bridge network. 
- Designates `datastore1` as the datastore in which to store container images, the files for the VCH appliance, and container VMs. 

<pre>vic-machine-<i>operating_system</i> create
--target 'Administrator@vsphere.local':<i>password</i>@<i>vcenter_server_address</i>/dc1
--compute-resource cluster1
--image-store datastore1
--bridge-network vch1-bridge
--name vch1
--force
--no-tlsverify
</pre>

<a name="standalone"></a> 
### Deploy to a Specific Standalone Host in vCenter Server ###

If vCenter Server manages multiple standalone ESXi hosts that are not part of a cluster, you use the `--compute-resource` option to specify the address of the ESXi host to which to deploy the VCH.

This example deploys a VCH with the following configuration:

- Specifies the user name, password, image store, bridge network, and name for the VCH.
- Deploys the VCH on the ESXi host with the FQDN `esxihost1.organization.company.com` in the datacenter `dc1`. You can also specify an IP address.

<pre>vic-machine-<i>operating_system</i> create
--target 'Administrator@vsphere.local':<i>password</i>@<i>vcenter_server_address</i>/dc1
--image-store datastore1
--bridge-network vch1-bridge
--compute-resource esxihost1.organization.company.com
--name vch1
--force
--no-tlsverify
</pre>

<a name="rp_host"></a>
### Deploy to a Resource Pool on an ESXi Host ###

To deploy a VCH in a specific resource pool on an ESXi host that is not managed by vCenter Server, you specify the resource pool name in the `--compute-resource` option. 

This example deploys a VCH with the following configuration:

- Specifies the user name and password, image store, and a name for the VCH.
- Designates `rp 1` as the resource pool in which to place the VCH. Note that the resource pool name is wrapped in quotes, because it contains a space. Use single quotes if you are using `vic-machine` on a Linux or Mac OS system and double quotes on a Windows system.

<pre>vic-machine-<i>operating_system</i> create
--target root:<i>password</i>@<i>esxi_host_address</i>
--compute-resource 'rp 1'
--image-store datastore1
--name vch1
--force
--no-tlsverify
</pre>

<a name="rp_cluster"></a>
### Deploy to a Resource Pool in a vCenter Server Cluster ###

To deploy a VCH in a resource pool in a vCenter Server cluster, you specify the names of the cluster and resource pool in the `compute-resource` option.

This example deploys a VCH with the following configuration:

- Specifies the user name, password, datacenter, image store, bridge network, and name for the VCH.
- Designates `rp 1` in cluster `cluster 1` as the resource pool in which to place the VCH. Note that the resource pool and cluster names are wrapped in quotes, because they contain spaces. Use single quotes if you are using `vic-machine` on a Linux or Mac OS system and double quotes on a Windows system.

<pre>vic-machine-<i>operating_system</i> create
--target 'Administrator@vsphere.local':<i>password</i>@<i>vcenter_server_address</i>/dc1
--compute-resource 'cluster 1'/'rp 1'
--image-store datastore1
--bridge-network vch1-bridge
--name vch1
--force
--no-tlsverify
</pre>

<a name="customized"></a>
### Set Limits on Resource Use ###

To limit the amount of system resources that the container VMs in a VCH can use, you can set resource limits on the VCH vApp. 

This example deploys a VCH with the following configuration:

- Specifies the user name, password, image store, cluster, bridge network, and name for the VCH.
- Sets resource limits on the VCH by imposing memory and CPU reservations, limits, and shares.

<pre>vic-machine-<i>operating_system</i> create
--target 'Administrator@vsphere.local':<i>password</i>@<i>vcenter_server_address</i>/dc1
--compute-resource cluster1
--image-store datastore1
--bridge-network vch1-bridge
--memory 1024
--memory-reservation 1024
--memory-shares low
--cpu 1024
--cpu-reservation 1024
--cpu-shares low
--name vch1
--force
--no-tlsverify
</pre>

For more information about setting resource use limitations on VCHs, see the [Advanced Deployment Options](vch_installer_options.md#deployment) and [Advanced Resource Management Options](vch_installer_options.md#adv-mgmt) sections in VCH Deployment Options.

<a name="networking"></a>
## Networking Examples ##

The examples in this section demonstrate how to direct traffic to and from VCHs and the other elements in your environment, how to set static IPs, how to configure container VM networks, and how to configure a VCH to use a proxy server.

<a name="networks"></a>
### Specify Public, Management, and Client Networks ###

In addition to the mandatory bridge network, if your vCenter Server environment includes multiple networks, you can direct different types of traffic to different networks. 

- You can direct the traffic between the VCH and the Internet to a specific network by specifying the `--public-network` option. Any container VM traffic that routes through the VCH also uses the public network. If you do not specify the `--public-network` option, the VCH uses the VM Network for public network traffic.
- You can direct traffic between ESXi hosts, vCenter Server, and the VCH to a specific network by specifying the `--management-network` option. If you do not specify the `--management-network` option, the VCH uses the public network for management traffic.
- You can designate a specific network for use by the Docker API by specifying the `--client-network` option. If you do not specify the `--client-network` option, the Docker API uses the public network.

**IMPORTANT**: A VCH supports a maximum of 3 distinct network interfaces. Because the bridge network requires its own port group, at least two of the public, client, and management networks must share a network interface and therefore a port group. Container networks do not go through the VCH, so they are not subject to this limitation. This limitation will be removed in a future release.

This example deploys a VCH with the following configuration:

- Specifies the user name, password, datacenter, cluster, image store, bridge network, and name for the VCH.
- Directs public and management traffic to network 1 and Docker API traffic to network 2. Note that the network names are wrapped in quotes, because they contain spaces. Use single quotes if you are using `vic-machine` on a Linux or Mac OS system and double quotes on a Windows system.

<pre>vic-machine-<i>operating_system</i> create
--target 'Administrator@vsphere.local':<i>password</i>@<i>vcenter_server_address</i>/dc1
--compute-resource cluster1
--image-store datastore1
--bridge-network vch1-bridge
--public-network 'network 1'
--management-network 'network 1'
--client-network 'network 2'
--name vch1
--force
--no-tlsverify
</pre>

For more information about the networking options, see the [Networking Options section](vch_installer_options.md#networking) in VCH Deployment Options.

<a name="static-ip"></a>
### Set a Static IP Address for the VCH Endpoint VM on the Different Networks ###

If you specify networks for any or all of the public, management, and client networks, you can deploy the VCH so that the VCH endpoint VM has a static IP address on one or more of those networks.  

This example deploys a VCH with the following configuration:

- Specifies the user name, password, datacenter, cluster, image store, bridge network, and name for the VCH.
- Directs public and management traffic to network 1 and Docker API traffic to network 2. Note that the network names are wrapped in quotes, because they contain spaces. Use single quotes if you are using `vic-machine` on a Linux or Mac OS system and double quotes on a Windows system.
- Sets a DNS server for use by the public, management, and client networks.
- Sets a static IP address for the VCH endpoint VM on the public and client networks. Because the management network shares a network with the public network, you only need to specify the public network IP address. You cannot specify a management IP address because you are sharing a port group between the management and public network.
- Specifies the gateway for the public network. If you set a static IP address on the public network, you must also specify the gateway address.
- Specifies a gateway for the client network. The `--client-network-gateway` option specifies the routing destination for client network traffic through the VCH endpoint VM, as well as the gateway address. The routing destination informs the VCH that it can reach all of the Docker clients at the network addresses in the ranges that you specify in the routing destinations by sending packets to the specified gateway.
- Because this example specifies a static IP address for the VCH endpoint VM on the client network, `vic-machine create` uses this address as the Common Name with which to create auto-generated trusted certificates. Full TLS authentication is implemented by default, so no TLS options are specified. 

<pre>vic-machine-<i>operating_system</i> create
--target 'Administrator@vsphere.local':<i>password</i>@<i>vcenter_server_address</i>/dc1
--compute-resource cluster1
--image-store datastore1
--bridge-network vch1-bridge
--public-network 'network 1'
--public-network-ip 192.168.1.10
--public-network-gateway 192.168.1.1/24
--management-network 'network 1'
--client-network 'network 2'
--client-network-ip 192.168.3.10
--client-network-gateway 192.168.3.0/24,192.168.128.0/22:192.168.2.1/24
--dns-server <i>dns_server_address</i>
--force
--name vch1
</pre>

For more information about setting static IP addresses, see the [Options for Specifying a Static IP Address for the VCH Endpoint VM](vch_installer_options.md#static-ip) in VCH Deployment Options.

<a name="ip-range"></a>
### Configure a Non-DHCP Network for Container VMs ###

You can designate a specific network for container VMs to use by specifying the `--container-network` option. Containers use this network if the container developer runs `docker run` or `docker create` specifying the `--net` option with one of the specified container networks when they run or create a container. This option requires a port group that must exist before you run `vic-machine create`. You cannot use the same port group that you use for the bridge network. You can provide a descriptive name for the network, for use by Docker. If you do not specify a descriptive name, Docker uses the vSphere network name. For example, the descriptive name appears as an available network in the output of `docker info` and `docker network ls`. 

If the network that you designate as the container network in the `--container-network` option does not support DHCP, you can configure the gateway, DNS server, and a range of IP addresses for container VMs to use. 

This example deploys a VCH with the following configuration:

- Specifies the user name, password, datacenter, cluster, image store, bridge network, and name for the VCH.
- Uses the VM Network for the public, management, and client networks.
- Designates a port group named `vic-containers` for use by container VMs that are run with the `--net` option.
- Gives the container network the name `vic-container-network`, for use by Docker. 
- Specifies the gateway, two DNS servers, and a range of IP addresses on the container network for container VMs to use.

<pre>vic-machine-<i>operating_system</i> create
--target 'Administrator@vsphere.local':<i>password</i>@<i>vcenter_server_address</i>/dc1
--compute-resource cluster1
--image-store datastore1
--bridge-network vch1-bridge
--container-network vic-containers:vic-container-network
--container-network-gateway vic-containers:<i>gateway_ip_address</i>/24
--container-network-dns vic-containers:<i>dns1_ip_address</i>
--container-network-dns vic-containers:<i>dns2_ip_address</i>
--container-network-ip-range vic-containers:192.168.100.0/24
--name vch1
--force
--no-tlsverify
</pre>

For more information about the container network options, see the [`--container-network`](vch_installer_options.md#container-network) and [Options for Configuring a Non-DHCP Network for Container Traffic](vch_installer_options.md#adv-container-net) sections in VCH Deployment Options.

<a name="proxy"></a>
### Configure a Proxy Server ###

If your network access is controlled by a proxy server, you must   configure a VCH to connect to the proxy server when you deploy it, so that it can pull images from external sources.

This example deploys a VCH with the following configuration:

- Specifies the user name, password, image store, cluster, bridge network, and name for the VCH.
- Configures the VCH to access the network via an HTTPS proxy server.

<pre>vic-machine-<i>operating_system</i> create
--target 'Administrator@vsphere.local':<i>password</i>@<i>vcenter_server_address</i>/dc1
--compute-resource cluster1
--image-store datastore1
--bridge-network vch1-bridge
--https-proxy https://<i>proxy_server_address</i>:<i>port</i>
--name vch1
--force
--no-tlsverify
</pre>


<a name="volume-stores"></a>
## Specify Volume Stores ##

If container application developers will use the `docker volume create` command to create containers that use volumes, you must create volume stores when you deploy VCHs. You specify volume stores in the `--volume-store` option. You can specify `--volume-store` multiple times to create multiple volume stores. 

When you create a volume store, you specify the name of the datastore to use and an optional path to a folder on that datastore. You also specify a descriptive name for that volume store for use by Docker.

This example deploys a VCH with the following configuration:

- Specifies the user name, password, datacenter, cluster, bridge network, and name for the VCH.
- Specifies the `volumes` folder on `datastore 1` as the default volume store. Creating a volume store named `default` allows container application developers to create anonymous or named volumes by using `docker create -v`. 
- Specifies a second volume store named `volume_store_2` in the `volumes` folder on `datastore 2`. 
- Note that the datastore names are wrapped in quotes, because they contain spaces. Use single quotes if you are using `vic-machine` on a Linux or Mac OS system and double quotes on a Windows system.

<pre>vic-machine-<i>operating_system</i> create
--target 'Administrator@vsphere.local':<i>password</i>@<i>vcenter_server_address</i>/dc1
--compute-resource cluster1
--bridge-network vch1-bridge
--image-store 'datastore 1'
--volume-store 'datastore 1'/volumes:default</i>
--volume-store 'datastore 2'/volumes:volume_store_2</i>
--name vch1
--force
--no-tlsverify
</pre>

For more information about volume stores, see the [volume-store section](vch_installer_options.md#volume-store) in VCH Deployment Options. 

<a name="security"></a>
## Security Examples ##

The examples in this section demonstrate how to configure a VCH to use Certificate Authority (CA) certificates to enable `TLSVERIFY` in your Docker environment, and to allow access to insecure registries of Docker images.

<a name="auto_cert"></a>
###  Use Auto-Generated Trusted CA Certificates ###

You can deploy a VCH that implements two-way authentication with trusted auto-generated TLS certificates that are signed by a CA. 

To automatically generate a server certificate that can pass client verification, you must specify the Common Name (CN) for the certificate by using the [`--tls-cname`](vch_installer_options.md#tls-cname) option. The CN should be the FQDN or IP address of the server, or a domain with a wildcard. The CN value must match the name or address that clients will use to connect to the server. You can use the `--organization` option to add basic descriptive information to the server certificate. This information is visible to clients if they inspect the server certificate.

If you specify an existing CA file with which to validate clients, you must also provide an existing server certificate that is compatible with the `--tls-cname` value or the IP address of the client interface.

This example deploys a VCH with the following configuration:

- Specifies the user, password, datacenter, image store, cluster, bridge network, and name for the VCH.
- Provides a wildcard domain `*.example.org` as the FQDN for the VCH, for use as the Common Name in the certificate. This assumes that there is a DHCP server offering IP addresses on VM Network, and that those addresses have corresponding DNS entries such as `dhcp-a-b-c.example.com`.
- Specifies a folder in which to store the auto-generated certificates.

<pre>vic-machine-<i>operating_system</i> create
--target 'Administrator@vsphere.local':<i>password</i>@<i>vcenter_server_address</i>/dc1
--compute-resource cluster1
--image-store datastore1
--bridge-network vch1-bridge
--tls-cname *.example.org
--cert-path <i>path_to_cert_folder</i>
--force
--name vch1
</pre>

The Docker API for this VCH will be accessible at `https://dhcp-a-b-c.example.com:2376`.

For more information about using auto-generated CA certificates, see the [Security Options section](vch_installer_options.md#security) in VCH Deployment Options.

<a name="custom_cert"></a>
### Use Custom Server Certificates ###

You can create a VCH that uses a custom server certificate, for example  a server certificate that has been signed by Verisign or another public root. You use the `--cert` and `--key` options to provide the paths to a custom X.509 certificate and its key when you deploy a VCH. The paths to the certificate and key files must be relative to the location from which you are running `vic-machine create`.

This example deploys a VCH with the following configuration:

- Specifies the user name, password, image store, cluster, bridge network, and name for the VCH.
- Provides the paths relative to the current location of the `*.pem` files for the custom server certificate and key files.

<pre>vic-machine-<i>operating_system</i> create
--target 'Administrator@vsphere.local':<i>password</i>@<i>vcenter_server_address</i>/dc1
--compute-resource cluster1
--image-store datastore1
--bridge-network vch1-bridge
--cert ../some/relative/path/<i>certificate_file</i>.pem
--key ../some/relative/path/<i>key_file</i>.pem
--name vch1
--force
</pre>

For more information about using custom server certificates, see the [Advanced Security Options section](vch_installer_options.md#adv-security) in VCH Deployment Options.

<a name="ops-user"></a>
### Specify Different User Accounts for VCH Deployment and Operation ###

When you deploy a VCH, you can use different vSphere user accounts for deployment and for operation. This allows you to run VCHs with lower levels of privileges than are required for deployment.

This example deploys a VCH with the following configuration:

- Specifies the image store, cluster, bridge network, and name for the VCH.
- Specifies <i>vsphere_admin</i> in the `--target` option, to identify the user account with which to deploy the VCH.
- Specifies <i>vsphere_user</i> and its password in the `--ops-user` and `--ops-password` options, to identify the user account with which the VCH runs.

<pre>vic-machine-<i>operating_system</i> create
--target <i>vsphere_admin</i>:<i>vsphere_admin_password</i>@<i>vcenter_server_address</i>/dc1
--compute-resource cluster1
--image-store datastore1
--bridge-network vch1-bridge
--name vch1
--ops-user <i>vsphere_user</i>
--ops-password <i>vsphere_user_password</i>
--force
--no-tlsverify
</pre>

For more information about using custom server certificates, see [--ops-user](vch_installer_options.md#ops-user) in VCH Deployment Options.

<a name="registry"></a>
### Authorize Access to an Insecure Private Registry Server ###

To authorize connections from a VCH to an insecure private registry server, set the `insecure-registry` option. You can specify `insecure-registry` multiple times to allow connections from the VCH to multiple insecure private registry servers.

This example deploys a VCH with the following configuration:

- Specifies the user name, password, image store, cluster, bridge network, and name for the VCH.
- Authorizes the VCH to pull Docker images from the insecure private registry servers located at the URLs <i>registry_URL_1</i> and <i>registry_URL_2</i>.
- The registry server at <i>registry_URL_2</i> listens for connections on port 5000. 

<pre>vic-machine-<i>operating_system</i> create
--target 'Administrator@vsphere.local':<i>password</i>@<i>vcenter_server_address</i>/dc1
--compute-resource cluster1
--image-store datastore1
--bridge-network vch1-bridge
--insecure-registry <i>registry_URL_1</i>
--insecure-registry <i>registry_URL_2:5000</i>
--name vch1
--force
--no-tlsverify
</pre>

For more information about configuring VCHs to connect to insecure private registry servers, see the section on the [`insecure-registry` option](vch_installer_options.md#insecure-registry) in VCH Deployment Options.