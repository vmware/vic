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
- [Options to Configure Virtual Container Hosts to Use Proxy Servers](#proxy)
- [Advanced Resource Management Options](#adv-mgmt)
- [Other Advanced Options](#adv-other)

<a name="vsphere"></a>
## vSphere Target Options ##

The `create` command of the `vic-machine` utility requires you to provide information about where in your vSphere environment to deploy the virtual container host and the vCenter Server or ESXi user account to use.

### `--target` ###

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

### `--user` ###

Short name: `-u`

The username for the ESXi host or vCenter Server instance on which you are deploying a virtual container host.

If you are deploying a virtual container host on vCenter Server, specify a username for an account that has the Administrator role on that vCenter Server instance. 

<pre>--user <i>esxi_or_vcenter_server_username</i></pre>

Wrap the user name in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if it includes special characters.

<pre>--user '<i>esxi_or_vcenter_server_usern@me</i>'</pre>

You can also specify the username in the URL that you pass to `vic-machine create` in the `target` option, in which case the `user` option is not required.

### `--password` ###

Short name: `-p`

The password for the user account on the vCenter Server on which you  are deploying the virtual container host, or the password for the ESXi host if you are deploying directly to an ESXi host. If not specified, `vic-machine` prompts you to enter the password during deployment.

<pre>--password <i>esxi_host_or_vcenter_server_password</i></pre>

Wrap the password in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if it includes special characters.

<pre>--password '<i>esxi_host_or_vcenter_server_p@ssword</i>'</pre>

You can also specify the username and password in the URL that you pass to `vic-machine create` in the `target` option, in which case the `password` option is not required.

### `--compute-resource` ###

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
### `--thumbprint` ###

Short name: None

The thumbprint of the vCenter Server or ESXi host certificate. Specify this option if your vSphere environment uses untrusted, self-signed certificates. If your vSphere environment uses trusted certificates that are signed by a known Certificate Authority (CA), you do not need to specify the `--thumbprint` option.

**NOTE** If your vSphere environment uses untrusted, self-signed certificates, you can run `vic-machine create` without the `--thumbprint` option by using the `--force` option. However, running `vic-machine create` with the `--force` option rather than providing the certificate thumbprint is not recommended, because it permits man-in-the-middle attacks to go undetected.

To obtain the thumbprint of the vCenter Server or ESXi host certificate, run `vic-machine create` without the specifying the `--thumbprint` or `--force` options. The deployment of the virtual container host fails, but the resulting error message includes the required certificate thumbprint. You can copy the thumbprint from the error message and run vic-machine create again, including the `thumbprint` option.

<pre>--thumbprint <i>certificate_thumbprint</i></pre>

<a name="security"></a>
## Security Options ##

When you deploy a virtual container host, you must specify the type of authentication to use when Docker clients connect to that virtual container host. 
<!--
- Two-way authentication with trusted auto-generated TLS certificates that are signed by a Certificate Authority (CA). Specify the [`tls-cname`](#tls-cname) option when you deploy the virtual container host.
- Server-side authentication with auto-generated, untrusted TLS certificates that are not signed by a CA, with no client-side verification. Specify the [`no-tlsverify`](#no-tlsverify) option when you deploy the virtual container host.
- Authentication with trusted custom TLS certificates that are signed by a CA.  Specify the [`cert`](#cert) and [`key`](#key) advanced options when you deploy the virtual container host.
- No TLS authentication. Any Docker client can connect to the virtual container host. Specify the [`no-tls`](#no-tls) advanced option when you deploy the virtual container host.

For more information about the possible security configurations for virtual container hosts, see [Securing Virtual Container Host Connections](security.md).

**IMPORTANT**: If you assign a static IP address to a virtual container host on the client network and you do not specify any authentication options, `vic-machine` behaves in the same way as if you set the `--tls-cname` option. If you do not set a static IP address on the virtual container host, it is **mandatory** to specify an authentication option when you deploy a virtual container host. For information about setting a static IP address on a virtual container host, see [Options for Specifying a Static IP Address for the Virtual Container Host Endpoint VM](#static-ip) in Advanced Options.
-->
The security options also allow you to configure virtual container hosts to connect to insecure registries and download container images by setting the `--insecure-registry` option.

<a name="tls-cname"></a>
### `--tls-cname` ###

Short name: None

The Common Name to use in an auto-generated CA certificate if you require two-way, trusted TLS certificate authentication when connecting Docker clients to the virtual container host.

The `--tls-cname` option is the minimum option that you must specify when using auto-generated trusted TLS certificates. For information about further options that you can specify when using auto-generated trusted certificates, see  the descriptions of the `--tls-ca`, `--certificate-key-size`, and `--organization` options in [Advanced Security Options](#adv-security).

If you specify a static IP address for the virtual container host on the client network by setting the `--client-network-ip` and `--client-network-gateway` options, `vic-machine create` uses this address as the Common Name when it creates auto-generated trusted certificates. In this case, you do not need to specify `--tls-cname` or any other authentication options. For information about setting a static IP address on a virtual container host, see [Options for Specifying a Static IP Address for the Virtual Container Host Endpoint VM](#static-ip) in Advanced Options.

When you specify the `--tls-cname` option, and potentially other options for auto-generating trusted TLS certificates, `vic-machine create` performs the following actions during the deployment of the virtual container host.

- Creates a folder with the same name as the virtual container host in the location in which you run `vic-machine create`.
- Creates trusted CA, server, and client certificate/key pairs in the newly created folder:
  - `ca.pem`
  - `ca-key.pem`
  - `cert.pem`
  - `key.pem`
  - `server-cert.pem`
  - `server-key.pem`
- Creates a browser-friendly PFX client certificate, `cert.pfx`, to use to authenticate connections to the VCH Admin portal for the virtual container host.

Running `vic-machine create` with the `--tls-cname` option also creates an environment file named <code><i>vch_name</i>.env</code>, that contains Docker environment variables that container developers can use to configure their Docker client environment:

- Activates TLS client verification.<pre>DOCKER_TLS_VERIFY=1</pre>
- The path to the client certificates.<pre>DOCKER_CERT_PATH=<i>path_to_certs</i></pre>
- The address of the virtual container host.<pre>DOCKER_HOST=<i>vch_address</i>:2376</pre>

You must provide copies of the certificate files and the environment file to container developers so that they can connect Docker clients to the virtual container host. 

If you use trusted certificates, container developers run Docker commands with the `--tlsverify`, `--tlscacert`, `--tlscert`, and `--tlskey` options.

When you specify the `--tls-cname` option, you must provide an FQDN for the virtual container host or the name of the domain to which the virtual container host will belong. The system on which you run `vic-machine create` and the remote vCenter Server system must agree on the vCenter Server system's FQDN or domain. As a consequence, to use the `--tls-cname` option, you must have a DNS service running on the client network that the virtual container host uses. You cannot specify an IP address in the `--tls-cname` option. If you do not have a DNS service on the client network, you can still implement full TLS authentication with trusted certificates by either specifying a static IP address or by using the `--cert` and `--key` options to upload custom certificates.  

<pre>--tls-cname vch-name.example.org</pre>
<pre>--tls-cname *.example.org</pre>

<a name="no-tlsverify"></a>
### `--no-tlsverify` ###

Short name: `--kv`

Authentication of the virtual container host with auto-generated TLS certificates that are not signed by a CA, with no client-side verification. The `vic-machine create` command still generates certificates, but these are untrusted, self-signed certificates. 

If you configure the virtual container host for untrusted TLS certificate authentication, clients are not verified. Consequently, container developers do not require copies of the certificate and key files.

When you specify the `--no-tlsverify` option, `vic-machine create` performs the following actions during the deployment of the virtual container host.

- Creates a folder with the same name as the virtual container host in the location in which you run `vic-machine create`.
- Creates an environment file named <code><i>vch_name</i>.env</code>, that contains the `DOCKER_HOST=vch_address` environment variable, that you can provide to container developers to use to set up their Docker client environment.

If you use untrusted certificates, container developers run Docker commands with the `--tls` option. The `--no-tlsverify` option takes no arguments. 

<pre>--no-tlsverify</pre>

<a name="registry"></a>
### `--insecure-registry` ###

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
### `--image-store` ###

Short name: `-i`

The datastore in which to store container image files, container VM files, and the files for the virtual container host appliance. The `--image-store` option is **mandatory** if there is more than one datastore in your vSphere environment. If there is only one datastore in your vSphere environment, the `--image-store` option is not required. 

If you are deploying the virtual container host to a vCenter Server cluster, the datastore that you designate in the `image-store` option must be shared by at least two ESXi hosts in the cluster. Using non-shared datastores is possible, but limits the use of vSphere features such as vSphere vMotion&reg; and VMware vSphere Distributed Resource Scheduler&trade; (DRS).

When you deploy a virtual container host, `vic-machine` creates a set of folders in the target datastore: 

- A folder with the same name as the virtual container host, at the top level of the datastore. This folder contains the VM files for the virtual container host appliance.
- A folder named `VIC` inside the virtual container host folder. The `VIC` folder contains a folder that uses the UUID of the virtual container host endpoint VM as its name. The <code>VIC/<i>vch_uuid</i></code> folder contains a subfolder named `images`, in which to store all of the container images that you pull into the virtual container host. 

You can specify a datastore folder to use as the image store in the format <code><i>datastore_name</i>/<i>path</i></code>. If the path to the folder that you specify does not already exist, `vic-machine create` creates it. In this case, `vic-machine` still creates the folder for the files of the virtual container host appliance at the top level of the datastore. However, `vic-machine create` creates the `VIC` folder inside the <code><i>datastore_name</i>/<i>path</i></code> folder, rather than in the same folder as the virtual container host files. 

By specifying the path to a datastore folder in the `--image-store` option, you can designate the same datastore folder as the image store for multiple virtual container hosts. In this way, `vic-machine create` creates only one `VIC` folder in the datastore, at the path that you specify. The `VIC` folder contains one <code><i>vch_uuid</i>/images</code> folder for each virtual container host that you deploy. By creating one <code><i>vch_uuid</i>/images</code> folder for each virtual container host, vSphere Integrated Containers Engine limits the potential for conflicts of image use between virtual container hosts, even if you share the same image store folder between multiple hosts.

**NOTE**: In the current builds of vSphere Integrated Containers Engine, sharing an image store folder between multiple virtual container hosts can lead to inconsistent behavior. Designate a different folder for the image store for each virtual container host, or omit the datastore folder from the `--image-store` option.

When container developers create containers, vSphere Integrated Containers Engine stores the files for container VMs at the top level of the image store, in folders that have the same name as the containers.

vSphere Integrated Containers Engine supports all alphanumeric characters, hyphens, and underscores in datastore paths and datastore names, but no other special characters. 

- Specify a datastore as the image store:<pre>--image-store <i>datastore_name</i></pre> 
- Specify a datastore folder as the image store:<pre>--image-store <i>datastore_name</i>/<i>path</i></pre> 
- Wrap the datastore name and path in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if they include spaces:  <pre>--image-store '<i>datastore name</i>'/'<i>datastore path</i>'</pre> 

If you specify an invalid datastore name, `vic-machine create` fails and suggests valid datastores.

<a name="volume-store"></a>
### `--volume-store` ###

Short name: `--vs`

The datastore in which to create volumes when container developers use the `docker volume create` or `docker create -v` commands. When you specify the `volume-store` option, you  provide the name of the target datastore and a label for the volume store. You can optionally provide a path to a specific folder in the datastore in which to create the volume store. If the folders that you specify in the path do not already exist on the datastore, `vic-machine create` creates the appropriate folder structure. If you specify an invalid datastore name, `vic-machine create` fails and suggests valid datastores. 

**IMPORTANT**: If multiple virtual container hosts will use the same datastore for their volume stores, specify a different datastore folder for each virtual container host. Do not designate the same datastore folder as the volume store for multiple virtual container hosts.

If you are deploying the virtual container host to a vCenter Server cluster, the datastore that you designate in the `volume-store` option should be shared by at least two ESXi hosts in the cluster. Using non-shared datastores is possible and `vic-machine create` succeeds, but it issues a warning that this configuration limits the use of vSphere features such as vSphere vMotion and DRS.

The label that you specify is the volume store name that Docker uses. For example, the volume store label appears in the information for a virtual container host when container developers run `docker info`. Container developers specify the volume store label in the <code>docker volume create --opt VolumeStore=<i>volume_store_label</i></code> option when they create a volume.

**IMPORTANT** If you do not specify the `volume-store` option, no  volume store is created and container developers cannot use the `docker volume create` or `docker create -v` commands.

- If you only require one volume store, you can set the volume store label to `default`. If you set the volume store label to `default`, container developers do not need to specify the <code>--opt VolumeStore=<i>volume_store_label</i></code> option when they run `docker volume create`. 

  **NOTE**: If container developers intend to use `docker create -v` to create containers that are attached to anonymous or named volumes, you must create a volume store with a label of `default`.

  <pre>--volume-store <i>datastore_name</i>:default</pre>
 
- If you specify the target datastore and the volume store label, `vic-machine create` creates a folder named `VIC/volumes` at the top level of the target datastore. Any volumes that container developers create will appear in the `VIC/volumes` folder.

  <pre>--volume-store <i>datastore_name</i>:<i>volume_store_label</i></pre>
- If you specify the target datastore, a datastore path, and the volume store label, `vic-machine create` creates a folder named `volumes` in the location that you specify in the datastore path. Any volumes that container developers create will appear in the <code><i>path</i>/volumes</code> folder.

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
The `vic-machine create` utility allows you to specify different networks for the different types of traffic between containers, the virtual container host, the external internet, and your vSphere environment. For information about the different networks that virtual container hosts use, see [Networks Used by vSphere Integrated Containers Engine](networks.md).

**IMPORTANT**: A virtual container host supports a maximum of 3 distinct networks. Because the bridge and container networks require  their own distributed port groups, at least two of the external, client, and management networks must share a network.

By default, `vic-machine create` obtains IP addresses for virtual container host endpoint VMs by using DHCP. For information about how to specify a static IP address for the virtual container host endpoint VM on the client, external, and management networks, see [Specify a Static IP Address for the Virtual Container Host Endpoint VM](#static-ip) in Advanced Options.

If your network access is controlled by a proxy server, see [Options to Configure Virtual Container Hosts to Use Proxy Servers](#proxy) in Advanced Options. 

<a name="bridge"></a>
### `--bridge-network` ###

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
### `--client-network` ###

Short name: `--cln`

The network that the virtual container host uses to generate the Docker API. The Docker API only uses this network.

If not specified, the virtual container host uses the external network for client traffic. If you specify an invalid network name, `vic-machine create` fails and suggests valid networks.

<pre>--client-network <i>network_name</i></pre>

Wrap the network name in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if it includes spaces.

<pre>--client-network '<i>network name</i>'</pre>

<a name="external-network"></a>
### `--external-network` ###

Short name: `--en`

The network for containers to use to connect to the Internet. Virtual container hosts use the external network to pull container images, for example from https://hub.docker.com/. Container VMs use the external network to publish network services. If you define the external network, you can deploy containers directly on the external interface. 

If not specified, containers use the default VM Network for external traffic. If you specify an invalid network name, `vic-machine create` fails and suggests valid networks.

<pre>--external-network <i>network_name</i></pre>

Wrap the network name in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if it includes spaces.

<pre>--external-network '<i>network name</i>'</pre>

<a name="management-network"></a>
### `--management-network` ###

Short name: `--mn`

The network that the virtual container host uses to communicate with vCenter Server and ESXi hosts. Container VMs use this network to communicate with the virtual container host. 

If not specified, the virtual container host uses the external network for management traffic. If you specify an invalid network name, `vic-machine create` fails and suggests valid networks.

<pre>--management-network <i>network_name</i></pre>

Wrap the network name in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if it includes spaces.

<pre>--management-network '<i>network name</i>'</pre>

<a name="container-network"></a>
### `--container-network` ###

Short name: `--cn`

A network for container VMs to use for external communication when container developers  run `docker run` or `docker create` with the `--net` option. 

To specify a container network, you provide the name of a distributed port group for the container VMs to use, and an optional descriptive name for the container network for use by Docker.  If you do not specify a descriptive name, Docker uses the vSphere network name. If you specify an invalid network name, `vic-machine create` fails and suggests valid networks.

- You can specify a vSphere network as the container network.
- The distributed port group must exist before you run `vic-machine create`. 
- You cannot use the same distributed port group as you use for the bridge network. 
- You can create the distributed port group on the same distributed virtual switch as the distributed port group that you use for the bridge network.
- If the network that you specify in the `container-network` option does not support DHCP, see [Options for Configuring a Non-DHCP Network for Container Traffic](#adv-container-net) in Advanced Options. 
- The descriptive name appears under `Networks` when you run `docker info` on the deployed virtual container host.
- Container developers use the descriptive name in the `--net` option when they run `docker run` or `docker create`.

If you do not specify the `container-network` option, or if container developers run `docker run` or `docker create` without specifying `--net`, container VMs use the bridge network. 

<pre>--container-network <i>distributed_port_group_name</i>:<i>container_network_name</i></pre>

Wrap the distributed port group name in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if it includes spaces. The descriptive name cannot include spaces.

<pre>--container-network '<i>distributed port group name</i>':<i>container_network_name</i></pre>

<a name="deployment"></a>
## Appliance Deployment Options ##

The `vic-machine` utility provides options to customize the virtual container host appliance.

### `--name` ###

Short name: `-n`

A name for the virtual container host appliance. If not specified, `vic-machine` sets the name of the virtual container host to `virtual-container-host`. If a virtual container host of the same name exists on the ESXi host or in the vCenter Server inventory, or if a folder of the same name exists in the target datastore, the deployment of the virtual container host fails.

<pre>--name <i>vch_appliance_name</i></pre>

Wrap the appliance name in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if it includes spaces.

<pre>--name '<i>vch appliance name</i>'</pre>

### `--memory` ###

Short name: `--mem`

Limit the amount of memory that is available for use by the virtual container host appliance and container VMs. Specify the memory limit value in MB. If not specified, `vic-machine create` sets the limit to 0 (unlimited).

<pre>--memory 1024</pre>

### `--cpu` ###

Short name: None

Limit the amount of CPU capacity that is available for use by the virtual container host appliance and container VMs. Specify the CPU limit value in MHz. If not specified, `vic-machine create` sets the limit to 0 (unlimited).

<pre>--cpu 1024</pre>

### `--force` ###

Short name: `-f`

Forces `vic-machine create` to ignore warnings and non-fatal errors and continue with the deployment of a virtual container host. Errors such as an incorrect compute resource still cause the installation to fail.

If your vSphere environment uses untrusted, self-signed certificates, you can use the `--force` option to deploy a virtual container host without providing the thumbprint of the vCenter Server or ESXi host in the `thumbprint` option. 

**IMPORTANT** Running `vic-machine create` with the `--force` option rather than providing the certificate thumbprint is not recommended, because it permits man-in-the-middle attacks to go undetected.

<pre>--force</pre>

### `--timeout` ###

Short name: none

The timeout period for uploading the vSphere Integrated Containers Engine  appliance and container images to the ESXi host, and for powering on the appliance. Specify a value in the format `XmYs` if the default timeout of 3m0s is insufficient.

<pre>--timeout 5m0s</pre> 

<a name="advanced"></a>
# Advanced Options #

The options in this section are exposed in the `vic-machine create` help if you run <code>vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> create --extended-help</code>, or <code>vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> create -x</code>. 

<a name="adv-security"></a>
## Advanced Security Options ##

The advanced security options allow you to customize the authentication of connections from Docker clients to virtual container hosts.

- Add optional information to auto-generated trusted TLS certificates by specifying the `--tls-ca`, `--certificate-key-size`, and `--organization` options.
- Use custom trusted TLS certificates by using the `--cert` and `--key` options.
- Disable TLS authentication completely by using the `--no-tls` option.

### `--tls-ca` ###

Short name: `--ca`

Certificate Authority (CA) files to use to verify Docker client certificates. Specify the `--tls-ca` option if your certificates are validated by a CA that is not commonly recognized. Specify the `--tls-ca` option multiple times to specify multiple CA files. 

<pre>--tls-ca <i>path_to_ca_file</i></pre>

### `--certificate-key-size` ###

Short name: `--ksz`

The size of the key for `vic-machine create` to use when it creates auto-generated trusted certificates. If not specified, `vic-machine create` creates keys with default size of 2048 bits. It is not recommended to use key sizes of less than 2048 bits.

<pre>--certificate-key-size 3072</pre>

### `--organization` ###

Short name: None

A list of identifiers to record in auto-generated trusted certificates. If not specified,`vic-machine create` uses the name of the virtual container host as the organization value. It also uses the IP address or FQND of the virtual container host as the organization if you set a static IP address by using the `--client-network-ip` and `--client-network-gateway` options.

<pre>--organization <i>organization_name</i></pre>

<a name="cert"></a>
### `--cert` ###

Short name: none

The path to a custom X.509 certificate that has been signed by a CA, for the Docker API to use to authenticate the virtual container host with a Docker client.

- This option is mandatory if you use custom TLS certificates, rather than auto-generated certificates, to authenticate connections between Docker clients and the virtual container hosts.
- Use this option in combination with the `key` option, that provides the path to the private key file for the custom certificate.
- Include the names of the certificate and key files in the paths.
- If you use trusted custom certificates, container developers run Docker commands with the `--tlsverify`, `--tlscacert`, `--tlscert`, and `--tlskey` options.

<pre>--cert <i>path_to_certificate_file</i>/<i>certificate_file_name</i>.pem 
--key <i>path_to_key_file</i>/<i>key_file_name</i>.pem
</pre> 

Wrap the folder names in the paths in single quotes (Linux or Mac OS) or double quotes (Windows) if they include spaces.

<pre>--cert '<i>path to certificate file</i>'/<i>certificate_file_name</i>.pem 
--key '<i>path to key file</i>'/<i>key_file_name</i>.pem
</pre> 

<a name="key"></a>
### `--key` ###

Short name: none

The path to the private key file to use with a custom CA certificate. This option is mandatory if you specify the `cert` option, that provides the path to a custom X.509 certificate file. Include the names of the certificate and key files in the paths. 

<pre>--cert <i>path_to_certificate_file</i>/<i>certificate_file_name</i>.pem 
--key <i>path_to_key_file</i>/<i>key_file_name</i>.pem
</pre> 

Wrap the folder names in the paths in single quotes (Linux or Mac OS) or double quotes (Windows) if they include spaces.

<pre>--cert '<i>path to certificate file</i>'/<i>certificate_file_name</i>.pem 
--key '<i>path to key file</i>'/<i>key_file_name</i>.pem
</pre>

<a name="no-tls"></a>
### `--no-tls` ###

Short name: `-k`

Disables TLS authentication of connections between the Docker client and  the virtual container host. 

Set the `no-tls` option if you do not require TLS authentication between the virtual container host and the Docker client. Any Docker client can connect to the virtual container host if you disable TLS authentication. 

If you use the `no-tls` option, container developers connect Docker clients to the virtual container host via port 2375, instead of via port 2376.

<pre>--no-tls</pre>


<a name="static-ip"></a>
## Options for Specifying a Static IP Address for the Virtual Container Host Endpoint VM ##

You can specify a static IP address for the virtual container host endpoint VM on each of the client, external, and management networks. DHCP is used for the endpoint VM for any network on which you do not specify a static IP address.

If you specify static IP addresses, you can only specify one static IP address on a given port group. If more than one of the client, external, or management networks shares a port group, you can only specify an IP address for one of those networks. The same address is then used for all of the networks that share that port group.

Assigning the same subnet to multiple port groups can cause routing problems.  If `vic-machine create` detects that you have assigned the same subnet to multiple port groups, it issues a warning.

To specify a static IP address for the endpoint VM, you provide an IP address, and a gateway address. You can also optionally specify one or more DNS server addresses.

**IMPORTANT**: If you assign a static IP address to a virtual container host on the client network by setting the `--client-network-ip` and `--client-network-gateway` options, `vic-machine create` uses this address to auto-generate trusted CA certificates. In this case, two-way TLS authentication with trusted certificates is implemented by default, and you do not need to perform any additional TLS configuration when you deploy the virtual container host. If you assign a static IP to a virtual container host on the client network, `vic-machine create` creates the same certificate and environment variable files as described in the [`--tls-cname` option](#tls-cname).

### `--dns-server` ###

Short name: None

A DNS server to use if you specify static IP addresses for the virtual container host on the client, external, and management networks. You can specify `dns-server` multiple times, to configure multiple DNS servers.  

If you specify `dns-server` but you do not specify a static IP address for one or more of the client, external, and management networks, `vic-machine create` ignores the `dns-server` setting for that network and uses the DNS servers that are provided by DHCP. 

If you use a mixture of static and DHCP addresses for the virtual container host on the different networks, the virtual container host uses the DNS servers that you specify in `dns-server` and those that DHCP provides.

If you specify static IP address for the virtual container host on any of the client, external, and management networks and you do not specify `dns-server`, the DNS server defaults to 8.8.8.8 and 8.8.4.4. 

<pre>--dns-server=172.16.10.10
--dns-server=172.16.10.11
</pre>

### `--client-network-ip`, `--external-network-ip`, `--management-network-ip` ###

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

### `--client-network-gateway`, `--external-network-gateway`, `--management-network-gateway` ###

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

### `--container-network-gateway` ###

Short name: `--cng`

The gateway for the subnet of the container network. This option is required if the network that you specify in the `container-network` option does not support DHCP. Specify the gateway in the format <code><i>container_network</i>:<i>subnet</i></code>. If you specify this option, it is recommended that you also specify the  `container-network-dns` option.

When you specify the container network gateway, you must use the distributed port group that you specify in the `container-network` option. If you specify `container-network-gateway` but you do not specify `container-network`, or if you specify a different distributed port group to the one that you specify in `container-network`, `vic-machine create` fails with an error.

<pre>--container-network-gateway <i>distributed_port_group_name</i>:<i>gateway_ip_address</i>/<i>subnet_mask</i></pre>

Wrap the distributed port group name in single quotes (Linux or Mac OS) or double quotes (Windows) if it includes spaces.

<pre>--container-network-gateway '<i>distributed port group name</i>':<i>gateway_ip_address</i>/<i>subnet_mask</i></pre>

### `--container-network-dns` ###

Short name: `--cnd`

The address of the DNS server for the container network. This option is recommended if the network that you specify in the `container-network` option does not support DHCP. 

When you specify the container network DNS server, you must use the distributed port group that you specify in the `container-network` option. You can specify `container-network-dns` multiple times, to configure multiple DNS servers. If you specify `container-network-dns` but you do not specify `container-network`, or if you specify a different distributed port group to the one that you specify in `container-network`, `vic-machine create` fails with an error.

<pre>--container-network-dns <i>distributed_port_group_name</i>:8.8.8.8</pre>

Wrap the distributed port group name in single quotes (Linux or Mac OS) or double quotes (Windows) if it includes spaces.

<pre>--container-network-dns '<i>distributed port group name</i>':8.8.8.8</pre>

### `--container-network-ip-range` ###

Short name: `--cnr`

The range of IP addresses that container VMs can use if the network that you specify in the `container-network` option does not support DHCP. If you specify `--container-network-ip-range`, virtual container hosts manage the addresses for containers within that range. The range that you specify must not be used by other computers or VMs on the network. If you specify `container-network-gateway` but do not specify `--container-network-ip-range`, the IP range for container VMs is the entire subnet that you specify in `container-network-gateway`. 

When you specify the container network IP range, you must use the distributed port group that you specify in the `container-network `option. If you specify `container-network-ip-range` but you do not specify `container-network`, or if you specify a different distributed port group to the one that you specify in `container-network`, `vic-machine create` fails with an error.

<pre>--container-network-ip-range <i>distributed_port_group_name</i>:192.168.100.2-192.168.100.254</pre>

You can also specify the IP range as a CIDR.

<pre>--container-network-ip-range <i>distributed_port_group_name</i>:192.168.100.0/24</pre>

Wrap the distributed port group name in single quotes (Linux or Mac OS) or double quotes (Windows) if it includes spaces.

<pre>--container-network-ip-range '<i>distributed port group name</i>':192.168.100.0/24</pre>

<a name="proxy"></a>
## Options to Configure Virtual Container Hosts to Use Proxy Servers ##

If your network access is controlled by a proxy server, you must   configure a virtual container host to connect to the proxy server when you deploy it.

**IMPORTANT**: Configuring a virtual container host to use a proxy server does not configure proxy support on the containers that this virtual container host runs. Container developers must configure proxy servers on containers when they create them. 

### `--http-proxy` ###

Short name: `--hproxy`

The address of the HTTP proxy server through which the virtual container host accesses the network. Specify the address of the proxy server as either an FQDN or an IP address.

<pre>--http-proxy http://<i>proxy_server_address</i>:<i>port</i></pre>

### `--https-proxy` ###

Short name: `--sproxy`

The address of the HTTPS proxy server through which the virtual container host accesses the network. Specify the address of the proxy server as either an FQDN or an IP address.

<pre>--https-proxy https://<i>proxy_server_address</i>:<i>port</i></pre>

<a name="adv-mgmt"></a>
## Advanced Resource Management Options ##

### `--memory-reservation` ###

Short name: `--memr`

Reserve a quantity of memory for use by the virtual container host appliance and container VMs. Specify the memory reservation value in MB. If not specified, `vic-machine create` sets the reservation to 1.

<pre>--memory-reservation 1024</pre>

### `--memory-shares` ###

Short name: `--mems`

Set memory shares on the virtual container host appliance. Specify the share value as a level or a number, for example `high`, `normal`, `low`, or `163840`. If not specified, `vic-machine create` sets the share to `normal`.

<pre>--memory-shares low</pre>

### `--cpu-reservation` ###

Short name: `--cpur`

Reserve a quantity of CPU capacity for use by the virtual container host appliance and container VMs.  Specify the CPU reservation value in MHz. If not specified, `vic-machine create` sets the reservation to 1.

<pre>--cpu-reservation 1024</pre>

### `--cpu-shares` ###

Short name: `--cpus`

Set CPU shares on the virtual container host appliance. Specify the share value as a level or a number, for example `high`, `normal`, `low`, or `163840`. If not specified, `vic-machine create` sets the share to `normal`.

<pre>--cpu-shares low</pre>

### `--appliance-cpu ` ###

Short name: none

The number of virtual CPUs for the virtual container host endpoint VM. The default is 1. Set this option to increase the number of CPUs in the virtual container host VM, for example if the virtual container host will handle large volumes of containers, or containers that require a lot of processing power.

**NOTE** Use the `--cpu` option instead of the `--appliance-cpu` option. The `--appliance-cpu` option is mainly intended for use by VMware Support.

<pre>--appliance-cpu <i>number_of_CPUs</i></pre>

### `--appliance-memory ` ###

Short name: none

The amount of memory for the virtual container host endpoint VM. The default is 2048MB. Set this option to increase the amount of memory in the virtual container host VM, for example if the virtual container host will handle large volumes of containers, or containers that consume a lot of memory.

**NOTE** Use the `--memory` option instead of the `--appliance-memory` option. The `--appliance-memory` option is mainly intended for use by VMware Support.

<pre>--appliance-memory <i>amount_of_memory</i></pre>

<a name="adv-other"></a>
## Other Advanced Options ##

<a name="bridge-range"></a>
### `--bridge-network-range` ###

Short name: `--bnr`

The range of IP addresses that additional bridge networks can use when container application developers use `docker network create` to create new bridge networks. If you do not specify the `bridge-network-range` option, the IP range for bridge networks is 172.16.0.0/12.

When you specify the bridge network IP range, you specify the IP range as a CIDR.

<pre>--bridge-network-range 192.168.100.0/24</pre>


### `--base-image-size` ###

Short name: None

The size of the base image from which to create other images. You should not normally need to use this option. Specify the size in `GB` or `MB`. The default is 8GB. 

<pre>--base-image-size 4GB</pre>

### `--container-store` ###

Short name: `--cs`

The `container-store` option is not enabled. Container VM files are stored in the datastore that you designate as the image store. 

### `--appliance-iso` ###

Short name: `--ai`

The path to the ISO image from which the virtual container host appliance boots. Set this option if you have moved the `appliance.iso` file to a folder that is not the folder that contains the `vic-machine` binary or is not the folder from which you are running `vic-machine`. Include the name of the ISO file in the path.

**NOTE**: Do not use the `--appliance-iso` option to point `vic-machine` to an `--appliance-iso` file that is of a different version to the version of `vic-machine` that you are running.

<pre>--appliance-iso <i>path_to_ISO_file</i>/appliance.iso</pre>

Wrap the folder names in the path in single quotes (Linux or Mac OS) or double quotes (Windows) if they include spaces.

<pre>--appliance-iso '<i>path to ISO file</i>'/appliance.iso</pre>

### `--bootstrap-iso` ###

Short name: `--bi`

The path to the ISO image from which to boot container VMs. Set this option if you have moved the `bootstrap.iso` file to a folder that is not the folder that contains the `vic-machine` binary or is not the folder from which you are running `vic-machine`. Include the name of the ISO file in the path.

**NOTE**: Do not use the `--bootstrap-iso` option to point `vic-machine` to a `--bootstrap-iso` file that is of a different version to the version of `vic-machine` that you are running.

<pre>--bootstrap-iso <i>path_to_ISO_file</i>/bootstrap.iso</pre>

Wrap the folder names in the path in single quotes (Linux or Mac OS) or double quotes (Windows) if they include spaces.

<pre>--bootstrap-iso '<i>path to ISO file</i>'/bootstrap.iso</pre>

### `--use-rp` ###

Short name: none

Deploy the virtual container host appliance to a resource pool on vCenter Server rather than to a vApp. If you specify this option, `vic-machine create` creates a resource pool with the same name as the virtual container host.

<pre>--use-rp</pre>


### `--debug` ###
Short name: `-v`

Provide verbose logging output, for troubleshooting purposes when running `vic-machine create`. If not specified, the `debug` value is set to 0 and verbose logging is disabled. Provide a value of 1 or greater to increase the verbosity of the logging. Note that setting debug to a value greater than 1 can affect the behavior of `vic-machine create`.

<pre>--debug 1</pre>