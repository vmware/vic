# Virtual Container Host Deployment Options

The command line utility for vSphere Integrated Containers Engine, `vic-machine`, provides a `create` command with options that allow you to customize the deployment of virtual container hosts to correspond to your vSphere environment.

- [vSphere Target Options](#vsphere)
- [Security Options](#security)
- [Datastore Options](#datastore)
- [Networking Options](#networking)
- [Appliance Deployment Options](#deployment)

To allow you to fine-tune the deployment of virtual container hosts, `vic-machine create` provides [Advanced Options](#advanced).

- [Advanced Security Options](#adv-security)
- [Options for Specifying a Static IP Address for the Virtual Container Host Endpoint VM](#static-ip)
- [Options for Configuring a Non-DHCP Network for Container Traffic](#adv-container-net)
- [Advanced Resource Management Options](#adv-mgmt)
- [Other Advanced Options](#adv-other)

<a name="vsphere"></a>
## vSphere Target Options ##

The `create` command of the `vic-machine` utility requires you to provide information about where in your vSphere environment to deploy the virtual container host and the vCenter Server or ESXi user account to use.

### `target` ###

Short name: `-t`

The IPv4 address, fully qualified domain name (FQDN), or URL of the ESXi host or vCenter Server instance on which you are deploying a virtual container host. This option is always **mandatory**.

To facilitate IP address changes in your infrastructure, provide an FQDN whenever possible, rather than an IP address.

- If the target ESXi host is not managed by vCenter Server, provide the address of the ESXi host.<pre>--target <i>esxi_host_address</i></pre>
- If the target ESXi host is managed by vCenter Server, or if you are deploying to a cluster, provide the address of vCenter Server.<pre>--target <i>vcenter_server_address</i></pre>
- You can include the user name and password in the target URL. <pre>--target <i>vcenter_or_esxi_username</i>:<i>password</i>@<i>vcenter_or_esxi_address</i></pre>

  Wrap the user name or password in single quotes (Linux or Mac OS) or double quotes (Windows) if they include special characters.<pre>'<i>vcenter_or_esxi_usern@me</i>':'<i>p@ssword</i>'@<i>vcenter_or_esxi_address</i></pre>
  
  If you do not include the user name in the target URL, you must specify the `user` option. If you do not specify the `password` option or include the password in the target URL, `vic-machine create` prompts you to enter the password.
- If you are deploying a virtual container host on a vCenter Server instance that includes more than one datacenter, include the datacenter name in the target URL. If you include an invalid datacenter name, `vic-machine create` fails and suggests the available datacenters that you can specify. 

  <pre>--target <i>vcenter_server_address</i>/<i>datacenter_name</i></pre>

  Wrap the datacenter name in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if it includes spaces.
  <pre>--target <i>vcenter_server_address</i>/'<i>datacenter name</i>'</pre>

### `user` ###

Short name: `-u`

The username for the ESXi host or vCenter Server instance on which you are deploying a virtual container host.

If you are deploying a virtual container host on vCenter Server, specify a username for an account that has the Administrator role on that vCenter Server instance. 

<pre>--user <i>esxi_or_vcenter_server_username</i></pre>

Wrap the user name in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if it includes special characters.

<pre>--user '<i>esxi_or_vcenter_server_usern@me</i>'</pre>

You can also specify the username in the URL that you pass to `vic-machine create` in the `target` option, in which case the `user` option is not required.

### `password` ###

Short name: `-p`

The password for the user account on the vCenter Server on which you  are deploying the virtual container host, or the password for the ESXi host if you are deploying directly to an ESXi host. If not specified, `vic-machine` prompts you to enter the password during deployment.

<pre>--password <i>esxi_host_or_vcenter_server_password</i></pre>

Wrap the password in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if it includes special characters.

<pre>--password '<i>esxi_host_or_vcenter_server_p@ssword</i>'</pre>

You can also specify the username and password in the URL that you pass to `vic-machine create` in the `target` option, in which case the `password` option is not required.

### `compute-resource` ###

Short name: `-r`

The relative path to the host, cluster, or resource pool in which to deploy the virtual container host. 

If the vCenter Server instance on which you are deploying a virtual container host only includes a single instance of a standalone host or  cluster, `vic-machine create` automatically detects and uses those resources. If you are deploying to an ESXi host that has no resource pools, `vic-machine create` automatically uses the default resource pool. In these cases, you do not need to specify a compute resource when you run `vic-machine create`.

You specify the `compute-resource` option in the following circumstances:

- A vCenter Server instance includes multiple instances of standalone hosts or clusters, or a mixture of standalone hosts and clusters.
- An ESXi host includes multiple resource pools. 
- You want to deploy the virtual container host to a specific resource pool in your environment. 

If you do not specify the `compute-resource` option and multiple possible resources exist, `vic-machine create` fails and suggests valid targets for `compute-resource` in the failure message. 

* To deploy to a specific resource pool on an ESXi host, specify the name of the resource pool: <pre>--compute-resource  <i>resource_pool_name</i></pre>
* To deploy to a vCenter Server instance that has more than one standalone host that are not part of a cluster, specify the IPv4 address or fully qualified domain name (FQDN) of the target host:<pre>--compute-resource <i>host_address</i></pre>
* To deploy to a vCenter Server with more than one cluster, specify the name of the target cluster: <pre>--compute-resource <i>cluster_name</i></pre>
* To deploy to a specific resource pool on a standalone host that is managed by vCenter Server, specify the IPv4 address or FQDN of the target host and name of the resource pool:<pre>--compute-resource <i>host_name</i>/<i>resource_pool_name</i></pre>
* To deploy to a specific resource pool in a cluster, specify the names of the target cluster and the resource pool:<pre>--compute-resource <i>cluster_name</i>/<i>resource_pool_name</i></pre>
* Wrap the resource names in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if they include spaces:<pre>--compute-resource '<i>cluster name</i>'/'<i>resource pool name</i>'</pre>

<a name="thumbprint"></a>
### `thumbprint` ###

Short name: None

The thumbprint of the vCenter Server or ESXi host certificate. This option is **mandatory** unless you run `vic-machine create` with the `--force` option. 

**IMPORTANT** Running `vic-machine create` with the `--force` option rather than providing the certificate thumbprint is not recommended, because it permits man-in-the-middle attacks to go undetected.

To obtain the thumbprint of the vCenter Server or ESXi host certificate, run `vic-machine create` without the specifying the `--thumbprint` or `--force` options. The deployment of the virtual container host fails, but the resulting error message includes the required certificate thumbprint. 

<pre>Failed to verify certificate for target=<i>vcenter_or_esxi_host</i> (thumbprint=<i>thumbprint</i>)
</pre>

You can copy the thumbprint from the error message and run vic-machine create again, including the `thumbprint` option.

<pre>--thumbprint <i>certificate_thumbprint</i></pre>

<a name="security"></a>
## Security Options ##

For information about the security requirements for virtual container hosts, see [vSphere Integrated Containers Engine Security Overview](security.md).

### `tls-cname` ###

Short name: None

The Common Name to use in generated CA certificate when requiring client certificate authentication.

### `no-tlsverify` ###

Short name: `--kv`

Disables authentication via client certificates

<a name="registry"></a>
### `insecure-registry` ###

Short name: `--dir`

If your Docker environment stores Docker images in an insecure private registry server, you must configure virtual container hosts to connect to this private registry server when you deploy them. An insecure  private registry server is a private registry server that is secured by self-signed certificates rather than by TLS. You authorize connections from a virtual container host to an insecure private registry server by setting the URL of a registry server in the `insecure-registry` option. If the registry server listens on a specific port, add the port number to the URL.

You can specify `insecure-registry` multiple times to allow connections from the virtual container host to multiple insecure  private registry servers.

<pre>--insecure-registry <i>registry_URL_1</i>
--insecure-registry <i>registry_URL_2</i>:<i>port_number</i>
</pre>

**NOTE**: The current builds of vSphere Integrated Containers do not yet support private registry servers that you secure by using TLS certificates.

<a name="datastore"></a>
## Datastore Options ##
The `vic-machine` utility allows you to specify the datastore in which to store container image files, container VM files, and the files for the virtual container host appliance. You can also specify datastores in which to create container volumes. 

- vSphere Integrated Containers Engine fully supports VMware vSAN datastores. 
- vSphere Integrated Containers Engine supports all alphanumeric characters, hyphens, and underscores in datastore paths and datastore names, but no other special characters.
- If you specify different datastores in the different datastore options, and if no single host in a cluster can access all of those datastores, `vic-machine create` fails with an error.<pre>No single host can access all of the requested datastores. 
Installation cannot continue.</pre>
- If you specify different datastores in the different datastore options, and if only one host in a cluster can access all of them, `vic-machine create` succeeds with a warning.<pre>Only one host can access all of the image/container/volume datastores. 
This may be a point of contention/performance degradation and HA/DRS 
may not work as intended.</pre> 

<a name="image"></a>
### `image-store` ###

Short name: `-i`

The datastore in which to store container image files, container VM files, and the files for the virtual container host appliance. The `image-store` option is **mandatory** if there is more than one datastore in your vSphere environment. If there is only one datastore in your vSphere environment, the `image-store` option is not required. 

When you deploy a virtual container host, `vic-machine` creates a folder named `VIC` on the target datastore,  in which to store all of the container images that you pull into a virtual container host. The `vic-machine` utility also places the VM files for the virtual container host appliance in the datastore that you designate as the image store, in a folder that has the same name as the virtual container host. The files for a container VM are stored in the image store in a folder that has the same name as the container.

If you are deploying the virtual container host to a vCenter Server cluster, the datastore that you designate in the `image-store` option must be shared by at least two ESXi hosts in the cluster. Using non-shared datastores is possible, but limits the use of vSphere features such as vSphere vMotion&reg; and VMware vSphere Distributed Resource Scheduler&trade; (DRS).

You can designate the same datastore as the image store for multiple virtual container hosts. In this case, only one `VIC` folder is created in the datastore and container image files are made available to all of the virtual container hosts that use that image store.

vSphere Integrated Containers Engine supports all alphanumeric characters, hyphens, and underscores in datastore paths and datastore names, but no other special characters. 

<pre>--image-store <i>datastore_name</i></pre> 

You can specify a datastore folder to use as the image store in the format <code>datastore/<i>path</i></code>. In this case, the virtual container host uses <code><i>path</i></code> as the image store instead of using the folder with the same name as the virtual container host. If the folder that you specify does not already exist, `vic-machine create` creates it. 

<pre>--image-store <i>datastore_name</i>/<i>path</i></pre> 

Wrap the datastore name and path in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if they include spaces.

<pre>--image-store '<i>datastore name</i>'/'<i>datastore path</i>'</pre> 

If you specify an invalid datastore name, `vic-machine create` fails and suggests valid datastores.

<a name="volume-store"></a>
### `volume-store` ###

Short name: `--vs`

The datastore in which to create volumes when container developers use the `docker volume create` or `docker create -v` commands. When you specify the `volume-store` option, you  provide the name of the target datastore and a label for the volume store. You can optionally provide a path to a specific folder in the datastore in which to create the volume store. If you specify an invalid datastore name, `vic-machine create` fails and suggests valid datastores. 

If you are deploying the virtual container host to a vCenter Server cluster, the datastore that you designate in the `volume-store` option should be shared by at least two ESXi hosts in the cluster. Using non-shared datastores is possible and `vic-machine create` succeeds, but it issues a warning that this configuration limits the use of vSphere features such as vSphere vMotion and DRS.

You can designate the same datastore as the volume store for multiple virtual container hosts.

The label that you specify is the volume store name that Docker uses. For example, the volume store label appears in the information for a virtual container host when container developers run `docker info`. Container developers specify the volume store label in the <code>docker volume create --opt VolumeStore=<i>volume_store_label</i></code> option when they create a volume.

**IMPORTANT** If you do not specify the `volume-store` option, no  volume store is created and container developers cannot use the `docker volume create` or `docker create -v` commands.

- If you only require one volume store, you can set the volume store label to `default`. If you set the volume store label to `default`, container developers do not need to specify the <code>--opt VolumeStore=<i>volume_store_label</i></code> option when they run `docker volume create`. 

  **NOTE**: If container developers intend to use `docker create -v` to create containers that are attached to anonymous or named volumes, you must create a volume store with a label of `default`.

  <pre>--volume-store <i>datastore_name</i>:default</pre>
 
- If you specify the target datastore and the volume store label, `vic-machine create` creates a folder named `volumes` under the `VIC` folder on the target datastore. Any volumes that container developers create will appear in the `volumes` folder.

  <pre>--volume-store <i>datastore_name</i>:<i>volume_store_label</i></pre>
- If you specify the target datastore, a datastore path, and the volume store label, `vic-machine create` creates a folder named `VIC/volumes` in the location that you specify in the datastore path. If the folders that you specify in the path do not already exist on the datastore, `vic-machine create` creates the appropriate folder structure. Any volumes that container developers create will appear in the <code><i>path</i>/VIC/volumes</code> folder.

  <pre>--volume-store <i>datastore_name</i>/<i>datastore_path</i>:<i>volume_store_label</i></pre>
- Wrap the datastore name and path in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if they include spaces. The volume store label cannot include spaces.

  <pre>--volume-store '<i>datastore name</i>'/'<i>datastore path</i>':<i>volume_store_label</i></pre>
- You can specify the `volume-store` option multiple times, to create multiple volume stores for the virtual container host.

  <pre>--volume-store <i>datastore_name</i>/path:<i>volume_store_label_1</i>
--volume-store <i>datastore_name</i>/<i>path</i>:<i>volume_store_label_2</i>
[...]
--volume-store <i>datastore_name</i>/<i>path</i>:<i>volume_store_label_n</i>
</pre>

<a name="networking"></a>
## Networking Options ##
The `vic-machine create` utility allows you to specify different networks for the different types of traffic between containers, the virtual container host, the external internet, and your vSphere environment. For information about the different networks that virtual container hosts use, see [vSphere Integrated Containers Engine Networking Overview](networks.html).

By default, `vic-machine create` obtains IP addresses for virtual container host endpoint VMs by using DHCP. For information about how to specify a static IP address for the virtual container host endpoint VM on the client, external, and management networks, see [Specify a Static IP Address for the Virtual Container Host Endpoint VM](#static-ip) in Advanced Options.

<a name="bridge"></a>
### `bridge-network` ###

Short name: `-b`

A distributed port group that container VMs use to communicate with each other. 

The `bridge-network` option is **mandatory** if you are deploying a virtual container host to vCenter Server.

In a vCenter Server environment, before you run `vic-machine create`, you must create a distributed virtual switch and a distributed port group. You must add the target ESXi host or hosts to the distributed virtual switch, and assign a VLAN ID to the port group, to ensure that the bridge network is isolated. For information about how to create a distributed virtual switch and port group, see *Network Requirements* in [Environment Prerequisites for vSphere Integrated Containers Engine Installation](vic_installation_prereqs.md#networkreqs).

You pass the name of the distributed port group to the `bridge-network` option. Each virtual container host requires its own distributed port group. 

**IMPORTANT** 
- Do not assign the same `bridge-network` distributed port group to multiple virtual container hosts. Sharing a distributed port group between virtual container hosts might result in multiple container VMs being assigned the same IP address. 
- Do not use the `bridge-network` distributed port group as the target for any of the other `vic-machine create` networking options.

If you specify an invalid network name, `vic-machine create` fails and suggests valid networks.

The `bridge-network` option is **optional** when you are deploying a virtual container host to an ESXi host with no vCenter Server. In this case, if you do not specify `bridge-network`, `vic-machine` creates a  virtual switch and a port group that each have the same name as the virtual container host. You can optionally specify this option to assign an existing port group for use as the bridge network for container VMs. You can also optionally specify this option to create a new virtual switch and port group that have a different name to the virtual container host.

<pre>--bridge-network <i>distributed_port_group_name</i></pre>

Wrap the distributed port group name in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if it includes spaces.

<pre>--bridge-network '<i>distributed port group name</i>'</pre>

For information about how to specify a range of IP addresses for additional bridge networks, see [`bridge-network-range`](#bridge-range) in Advanced Networking Options.

<a name="client-network"></a>
### `client-network` ###

Short name: `--cln`

The network that the virtual container host uses to generate the Docker API. The Docker API only uses this network.

If not specified, the virtual container host uses the external network for client traffic. If you specify an invalid network name, `vic-machine create` fails and suggests valid networks.

<pre>--client-network <i>network_name</i></pre>

Wrap the network name in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if it includes spaces.

<pre>--client-network '<i>network name</i>'</pre>

<a name="external-network"></a>
### `external-network` ###

Short name: `--en`

The network for containers to use to connect to the Internet. Virtual container hosts use the external network to pull container images, for example from https://hub.docker.com/. Container VMs use the external network to publish network services. If you define the external network, you can deploy containers directly on the external interface. 

If not specified, containers use the default VM Network for external traffic. If you specify an invalid network name, `vic-machine create` fails and suggests valid networks.

<pre>--external-network <i>network_name</i></pre>

Wrap the network name in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if it includes spaces.

<pre>--external-network '<i>network name</i>'</pre>

<a name="management-network"></a>
### `management-network` ###

Short name: `--mn`

The network that the virtual container host uses to communicate with vCenter Server and ESXi hosts. Container VMs use this network to communicate with the virtual container host. 

If not specified, the virtual container host uses the external network for management traffic. If you specify an invalid network name, `vic-machine create` fails and suggests valid networks.

<pre>--management-network <i>network_name</i></pre>

Wrap the network name in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if it includes spaces.

<pre>--management-network '<i>network name</i>'</pre>

<a name="container-network"></a>
### `container-network` ###

Short name: `--cn`

A network for container VMs to use for external communication when container developers  run `docker run` or `docker create` with the `--net` option. 

To specify a container network, you provide the name of a distributed port group for the container VMs to use, and an optional descriptive name for the container network for use by Docker.  If you do not specify a descriptive name, Docker uses the vSphere network name. If you specify an invalid network name, `vic-machine create` fails and suggests valid networks.

- The distributed port group must exist before you run `vic-machine create`. 
- You cannot use the same distributed port group as you use for the bridge network. 
- You can create the distributed port group on the same distributed virtual switch as the distributed port group that you use for the bridge network.
- The descriptive name appears under `Networks` when you run `docker info` on the deployed virtual container host.
- Container developers use the descriptive name in the `--net` option when they run `docker run` or `docker create`.

If you do not specify the `container-network` option, or if container developers run `docker run` or `docker create` without specifying `--net`, container VMs use the bridge network. 

<pre>--container-network <i>distributed_port_group_name</i>:<i>container_network_name</i></pre>

Wrap the distributed port group name in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if it includes spaces. The descriptive name cannot include spaces.

<pre>--container-network '<i>distributed port group name</i>':<i>container_network_name</i></pre>

If the network that you specify in the `container-network` option does not support DHCP, see [Options for Configuring a Non-DHCP Network for Container Traffic](#adv-container-net) in Advanced Options. 

<a name="deployment"></a>
## Appliance Deployment Options ##

The `vic-machine` utility provides options to customize the virtual container host appliance.

### `name` ###

Short name: `-n`

A name for the virtual container host appliance. If not specified, `vic-machine` sets the name of the virtual container host to `virtual-container-host`. If a virtual container host of the same name exists on the ESXi host or in the vCenter Server inventory, or if a folder of the same name exists in the target datastore, the deployment of the virtual container host fails.

<pre>--name <i>vch_appliance_name</i></pre>

Wrap the appliance name in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if it includes spaces.

<pre>--name '<i>vch appliance name</i>'</pre>

### `memory` ###

Short name: `--mem`

Limit the amount of memory that is available for use by the virtual container host appliance and container VMs. Specify the memory limit value in MB. If not specified, `vic-machine create` sets the limit to 0 (unlimited).

<pre>--memory 1024</pre>

### `cpu` ###

Short name: None

Limit the amount of CPU capacity that is available for use by the virtual container host appliance and container VMs. Specify the CPU limit value in MHz. If not specified, `vic-machine create` sets the limit to 0 (unlimited).

<pre>--cpu 1024</pre>

### `force` ###

Short name: `-f`

Forces `vic-machine create` to ignore warnings and non-fatal errors and continue with the deployment of a virtual container host. Errors such as an incorrect compute resource still cause the installation to fail.

You can use the `--force` option to deploy a virtual container host without providing the thumbprint of the vCenter Server or ESXi host in the `thumbprint` option. 

**IMPORTANT** Running `vic-machine create` with the `--force` option rather than providing the certificate thumbprint is not recommended, because it permits man-in-the-middle attacks to go undetected.

<pre>--force</pre>

### `timeout` ###

Short name: none

The timeout period for uploading the vSphere Integrated Containers Engine  appliance and container images to the ESXi host, and for powering on the appliance. Specify a value in the format `XmYs` if the default timeout of 3m0s is insufficient.

<pre>--timeout 5m0s</pre> 

<a name="advanced"></a>
# Advanced Options #

The options in this section are exposed in the `vic-machine create` help if you run <code>vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> create --extended-help</code>, or <code>vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> create -x</code>. 

<a name="adv-security"></a>
## Advanced Security Options ##

### `--tls-ca` ###

Short name: `--ca`

xx

### `--certificate-key-size` ###

Short name: `--ksz`

xx

### `--organization` ###

Short name: None

xx

### `no-tls` ###

Short name: `-k`

Set the `no-tls` option if you do not require certificate-based  TLS authentication between the virtual container host and the Docker client. If you use the `no-tls` option, you connect Docker clients to the virtual container host via port 2375, instead of via port 2376.

<pre>--no-tls</pre>

### `cert` ###

Short name: none

The path to an X.509 certificate for the Docker API to use to authenticate the virtual container host with a Docker client.

- This option is mandatory if your Docker environment uses TLS certificates that are signed by a CA. For information about how to set up a Docker client to use CA certificates, see https://docs.docker.com/engine/security/https/.
- Use this option in combination with the `key` option, that provides the path to the private key file for the CA certificate.

If you use the `cert` and `key` options, `vic-machine` does not automatically generate certificates. Omit this option if your Docker environment does not use certificates that are signed by a CA. Include the names of the certificate and key files in the paths.

<pre>--cert <i>path_to_certificate_file</i>/<i>certificate_file_name</i>.pem 
--key <i>path_to_key_file</i>/<i>key_file_name</i>.pem</pre> 

Wrap the folder names in the paths in single quotes (Linux or Mac OS) or double quotes (Windows) if they include spaces.

<pre>--cert '<i>path to certificate file</i>'/<i>certificate_file_name</i>.pem 
--key '<i>path to key file</i>'/<i>key_file_name</i>.pem
</pre> 

### `key` ###

Short name: none

The path to the private key file for use with a custom CA certificate. This option is mandatory if your Docker environment uses certificates that are signed by a CA. For information about how to set up a Docker client to use CA certificates, see https://docs.docker.com/engine/security/https/.

Use this option in combination with the `cert` option, that provides the path to an X.509 certificate file. Include the names of the certificate and key files in the paths. 

<pre>--cert <i>path_to_certificate_file</i>/<i>certificate_file_name</i>.pem 
--key <i>path_to_key_file</i>/<i>key_file_name</i>.pem</pre> 

Wrap the folder names in the paths in single quotes (Linux or Mac OS) or double quotes (Windows) if they include spaces.

<pre>--cert '<i>path to certificate file</i>'/<i>certificate_file_name</i>.pem 
--key '<i>path to key file</i>'/<i>key_file_name</i>.pem
</pre>

<a name="static-ip"></a>
## Options for Specifying a Static IP Address for the Virtual Container Host Endpoint VM ##

You can specify a static IP address for the virtual container host endpoint VM on each of the client, external, and management networks. DHCP is used for the endpoint VM for any network on which you do not specify a static IP address.

If you specify static IP addresses, you can only specify one static IP address on a given port group. If more than one of the client, external, or management networks shares a port group, you can only specify an IP address for one of those networks. The same address is then used for all of the networks that share that port group.

Assigning the same subnet to multiple port groups can cause routing problems.  If `vic-machine create` detects that you have assigned the same subnet to multiple port groups, it issues a warning.

To specify a static IP address for the endpoint VM, you provide an IP address, and a gateway address. You can also optionally specify one or more DNS server addresses.

### `dns-server` ###

Short name: None

A DNS server to use if you specify static IP addresses for the virtual container host on the client, external, and management networks. You can specify `dns-server` multiple times, to configure multiple DNS servers.  

If you specify `dns-server` but you do not specify a static IP address for one or more of the client, external, and management networks, `vic-machine create` ignores the `dns-server` setting for that network and uses the DNS servers that are provided by DHCP. 

If you use a mixture of static and DHCP addresses for the virtual container host on the different networks, the virtual container host uses the DNS servers that you specify in `dns-server` and those that DHCP provides.

If you specify static IP address for the virtual container host on any of the client, external, and management networks and you do not specify `dns-server`, the DNS server defaults to 8.8.8.8 and 8.8.4.4. 

<pre>--dns=172.16.10.10
--dns=172.16.10.11
</pre>

### `client-network-ip`, `external-network-ip`, `management-network-ip` ###

Short name: None

A static IP address for the virtual container host on the client, external, or management network. If you specify an IP address by using `client/external/management-network-ip`, you must also specify a  corresponding gateway address by using `client/external/management-network-gateway`. If you specify neither a gateway nor an IP address for a given network, `vic-machine create` uses DHCP to obtain an IP address for the virtual container host endpoint VM on that network.

You can specify IP addresses in CIDR format.

<pre>--external-network-ip 192.168.X.N/24
--management-network-ip 192.168.Y.N/24
--client-network-ip 192.168.Z.N/24
</pre>

You can also specify IP addresses as resolvable FQDNs. If you specify an FQDN, `vic-machine create` uses the netmask from the gateway.

<pre>--external-network-ip=vch27-team-a.internal.domain.com
--management-network-ip=vch27-team-b.internal.domain.com
--client-network-ip=vch27-team-c.internal.domain.com
</pre>

### `client-network-gateway`, `external-network-gateway`, `management-network-gateway` ###

Short name: None

The gateway to use if you specify a static IP address for the virtual container host on the client, external, or management network. If you specify a gateway address by using `client/external/management-network-gateway`, you must also specify a corresponding IP address by using `client/external/management-network-ip`. If you specify neither a gateway nor an IP address for a given network, `vic-machine create` uses DHCP to obtain an IP address for the virtual container host endpoint VM on that network.

You specify gateway addresses in CIDR format.

<pre>
--external-network-gateway 192.168.X.1/24
--management-network-gateway 192.168.Y.1/24
--client-network-gateway 192.168.Z.1/24
</pre>


<a name="adv-container-net"></a>
## Options for Configuring a Non-DHCP Network for Container Traffic ##
If the network that you specify in the `container-network` option does not support DHCP, you must specify the `container-network-gateway` option. You can optionally specify one or more DNS servers and a range of IP addresses for container VMs on the container network. 

For information about the container network, see the section on the [`container-network` option](#container-network).

### `container-network-gateway` ###

Short name: `--cng`

The gateway for the subnet of the container network. This option is required if the network that you specify in the `container-network` option does not support DHCP. Specify the gateway in the format <code><i>container_network</i>:<i>subnet</i></code>. If you specify this option, it is recommended that you also specify the  `container-network-dns` option.

When you specify the container network gateway, you must use the distributed port group that you specify in the `container-network` option. If you specify `container-network-gateway` but you do not specify `container-network`, or if you specify a different distributed port group to the one that you specify in `container-network`, `vic-machine create` fails with an error.

<pre>--container-network-gateway <i>distributed_port_group_name</i>:<i>gateway_ip_address</i>/<i>subnet_mask</i></pre>

Wrap the distributed port group name in single quotes (Linux or Mac OS) or double quotes (Windows) if it includes spaces.

<pre>--container-network-gateway '<i>distributed port group name</i>':<i>gateway_ip_address</i>/<i>subnet_mask</i></pre>

### `container-network-dns` ###

Short name: `--cnd`

The address of the DNS server for the container network. This option is recommended if the network that you specify in the `container-network` option does not support DHCP. 

When you specify the container network DNS server, you must use the distributed port group that you specify in the `container-network` option. You can specify `container-network-dns` multiple times, to configure multiple DNS servers. If you specify `container-network-dns` but you do not specify `container-network`, or if you specify a different distributed port group to the one that you specify in `container-network`, `vic-machine create` fails with an error.

<pre>--container-network-dns <i>distributed_port_group_name</i>:8.8.8.8</pre>

Wrap the distributed port group name in single quotes (Linux or Mac OS) or double quotes (Windows) if it includes spaces.

<pre>--container-network-dns '<i>distributed port group name</i>':8.8.8.8</pre>

### `container-network-ip-range` ###

Short name: `--cnr`

The range of IP addresses that container VMs can use if the network that you specify in the `container-network` option does not support DHCP. If you do not specify this option, the IP range for container VMs is the entire subnet that you specify in `container-network-gateway`.

When you specify the container network IP range, you must use the distributed port group that you specify in the `container-network `option. If you specify `container-network-ip-range` but you do not specify `container-network`, or if you specify a different distributed port group to the one that you specify in `container-network`, `vic-machine create` fails with an error.

<pre>--container-network-ip-range <i>distributed_port_group_name</i>:192.168.100.2-192.168.100.254</pre>

You can also specify the IP range as a CIDR.

<pre>--container-network-ip-range <i>distributed_port_group_name</i>:192.168.100.0/24</pre>

Wrap the distributed port group name in single quotes (Linux or Mac OS) or double quotes (Windows) if it includes spaces.

<pre>--container-network-ip-range '<i>distributed port group name</i>':192.168.100.0/24</pre>


<a name="adv-mgmt"></a>
## Advanced Resource Management Options ##

### `memory-reservation` ###

Short name: `--memr`

Reserve a quantity of memory for use by the virtual container host appliance and container VMs. Specify the memory reservation value in MB. If not specified, `vic-machine create` sets the reservation to 1.

<pre>--memory-reservation 1024</pre>

### `memory-shares` ###

Short name: `--mems`

Set memory shares on the virtual container host appliance. Specify the share value as a level or a number, for example `high`, `normal`, `low`, or `163840`. If not specified, `vic-machine create` sets the share to `normal`.

<pre>--memory-shares low</pre>

### `cpu-reservation` ###

Short name: `--cpur`

Reserve a quantity of CPU capacity for use by the virtual container host appliance and container VMs.  Specify the CPU reservation value in MHz. If not specified, `vic-machine create` sets the reservation to 1.

<pre>--cpu-reservation 1024</pre>

### `cpu-shares` ###

Short name: `--cpus`

Set CPU shares on the virtual container host appliance. Specify the share value as a level or a number, for example `high`, `normal`, `low`, or `163840`. If not specified, `vic-machine create` sets the share to `normal`.

<pre>--cpu-shares low</pre>

### `appliance-cpu ` ###

Short name: none

The number of virtual CPUs for the virtual container host endpoint VM. The default is 1. Set this option to increase the number of CPUs in the virtual container host VM, for example if the virtual container host will handle large volumes of containers, or containers that require a lot of processing power.

**NOTE** Use the `--cpu` option instead of the `--appliance-cpu` option. The `--appliance-cpu` option is mainly intended for use by VMware Support.

<pre>--appliance-cpu <i>number_of_CPUs</i></pre>

### `appliance-memory ` ###

Short name: none

The amount of memory for the virtual container host endpoint VM. The default is 2048MB. Set this option to increase the amount of memory in the virtual container host VM, for example if the virtual container host will handle large volumes of containers, or containers that consume a lot of memory.

**NOTE** Use the `--memory` option instead of the `--appliance-memory` option. The `--appliance-memory` option is mainly intended for use by VMware Support.

<pre>--appliance-memory <i>amount_of_memory</i></pre>

<a name="adv-other"></a>
## Other Advanced Options ##

<a name="bridge-range"></a>
### `bridge-network-range` ###

Short name: `--bnr`

The range of IP addresses that additional bridge networks can use when container application developers use `docker network create` to create new bridge networks. If you do not specify the `bridge-network-range` option, the IP range for bridge networks is 172.16.0.0/12.

When you specify the bridge network IP range, you specify the IP range as a CIDR.

<pre>--bridge-network-range 192.168.100.0/24</pre>


### `base-image-size` ###

Short name: None

The size of the base image from which to create other images. You should not normally need to use this setting. Specify the size in `GB` or `MB`. The default is 8GB. 

<pre>--base-image-size 4GB</pre>

### `container-store` ###

Short name: `--cs`

The `container-store` option is not enabled. Container VM files are stored in the datastore that you designate as the image store. 

### `appliance-iso` ###

Short name: `--ai`

The ISO image from which the virtual container host appliance boots. Omit this option to boot the appliance from the default ISO that is included with `vic-machine`. Set this option to boot the appliance from a different ISO file, for example to reinstall an existing virtual container host or to update it to a newer version. Include the name of the ISO file in the path.

<pre>--appliance-iso <i>path_to_ISO_file</i>/<i>ISO_file_name</i>.iso</pre>

Wrap the folder names in the path in single quotes (Linux or Mac OS) or double quotes (Windows) if they include spaces.

<pre>--appliance-iso '<i>path to ISO file</i>'/<i>ISO_file_name</i>.iso</pre>

### `bootstrap-iso` ###

Short name: `--bi`

The ISO image from which container VMs boot. Omit this option to boot container VMs from the default Photon OS kernel ISO that is included with `vic-machine`. Set this option to a different ISO file to boot container VMs with an operating system other than the Photon OS kernel. Include the name of the ISO file in the path.

<pre>--bootstrap-iso <i>path_to_ISO_file</i>/<i>ISO_file_name</i>.iso</pre>

Wrap the folder names in the path in single quotes (Linux or Mac OS) or double quotes (Windows) if they include spaces.

<pre>--bootstrap-iso '<i>path to ISO file</i>'/<i>ISO_file_name</i>.iso</pre>

### `use-rp` ###

Short name: none

Deploy the virtual container host appliance to a resource pool on vCenter Server rather than to a vApp. If you specify this option, `vic-machine create` creates a resource pool with the same name as the virtual container host.

<pre>--use-rp</pre>


### `debug` ###
Short name: `-v`

Provide verbose logging output, for troubleshooting purposes when running `vic-machine create`. If not specified, the `debug` value is set to 0 and verbose logging is disabled. Provide a value of 1 or greater to increase the verbosity of the logging. Note that setting debug to a value greater than 1 can affect the behavior of `vic-machine create`.

<pre>--debug 1</pre>