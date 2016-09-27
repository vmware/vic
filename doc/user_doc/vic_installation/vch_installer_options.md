# Virtual Container Host Deployment Options

The command line utility for vSphere Integrated Containers Engine, `vic-machine`, provides a `create` command that has options to allow you to customize the deployment of a virtual container host to correspond to your vSphere environment.

- [Mandatory Options](#mandatory)
- [Networking Options](#networking)
- [Compute Resource Options](#compute)
- [Datastore Options](#datastore)
- [Security Options](#security)
- [vApp Deployment Options](#deployment)
- [Advanced Options](#advanced)

<a name="mandatory"></a>
## Mandatory Options ##

The `create` command of the `vic-machine` utility requires you to provide information about where in your vSphere environment to deploy the virtual container host, the user account to use, the bridge network to use, and the location in which to store container image files. The options in this section are either always **mandatory** or mandatory in most circumstances.

### `target` ###

Short name: `-t`

The IPv4 address, fully qualified domain name (FQDN), or URL of the ESXi host or vCenter Server instance on which you are deploying a virtual container host.

To facilitate IP address changes in your infrastructure, provide an FQDN whenever possible, rather than an IP address.

- If the target ESXi host is not managed by vCenter Server, provide the address of the ESXi host.<pre>--target <i>esxi_host_address</i></pre>
- If the target ESXi host is managed by vCenter Server, or if you are deploying to a cluster, provide the address of vCenter Server.<pre>--target <i>vcenter_server_address</i></pre>
- You can include the user name and password in the target URL. <pre>--target <i>vcenter_or_esxi_username</i>:<i>password</i>@<i>vcenter_or_esxi_address</i></pre>

  Wrap the user name or password in single quotes (Linux or Mac OS) or double quotes (Windows) if they include special characters.<pre>'<i>vcenter_or_esxi_usern@me</i>':'<i>p@ssword</i>'@<i>vcenter_or_esxi_address</i></pre>
  
  If you do not include the user name in the target URL, you must specify the `user` option. If you do not specify the `password` option or include the password in the target URL, `vic-machine create` prompts you to enter the password.
- If you are deploying a virtual container host on a vCenter Server instance that includes more than one datacenter, include the datacenter name in the target URL. If you include an invalid datacenter name, `vic-machine create` fails and suggests the available datacenters that you can specify. 

  <pre>--target <i>vcenter_server_address</i>/<i>datacenter_name</i></pre>

  Wrap the datacenter name in single quotes (Linux or Mac OS) or double quotes (Windows) if it includes spaces.
  <pre>--target <i>vcenter_server_address</i>/'<i>datacenter name</i>'</pre>

### `user` ###

Short name: `-u`

The username for the ESXi host or vCenter Server instance on which you are deploying a virtual container host.

If you are deploying a virtual container host on vCenter Server, specify a username for an account that has the Administrator role on that vCenter Server instance. 

<pre>--user <i>esxi_or_vcenter_server_username</i></pre>

Wrap the user name in single quotes (Linux or Mac OS) or double quotes (Windows) if it includes special characters.

<pre>--user '<i>esxi_or_vcenter_server_usern@me</i>'</pre>


You can also specify the username in the URL that you pass to `vic-machine create` in the `target` option, in which case the `user` option is not required.

<a name="image"></a>
### `image-store` ###

Short name: `-i`

The datastore in which to store container image files. The `image-store` option is **mandatory** if there is more than one datastore in your vSphere environment. If there is only one datastore in your vSphere environment, the `image-store` option is not required. 

When you deploy a virtual container host, `vic-machine` creates a folder named `VIC` on the target datastore,  in which to store all of the container images that you pull into a virtual container host. The `vic-machine` utility also places the VM files for the virtual container host in the datastore that you designate as the image store, in a folder that has the same name as the virtual container host. 

If you are deploying the virtual container host to a vCenter Server cluster, the datastore that you designate in the `image-store` option must be shared by at least two ESXi hosts in the cluster. Using non-shared datastores is possible, but limits the use of vSphere features such as DRS and High Availability.

You can designate the same datastore as the image store for multiple virtual container hosts. In this case, only one `VIC` folder is created in the datastore and the container image files are made available to all of the virtual container hosts that use that image store.

vSphere Integrated Containers Engine supports all alphanumeric characters, hyphens, and underscores in datastore paths and datastore names, but no other special characters. 

<pre>--image-store <i>datastore_name</i></pre> 

You can specify a datastore folder to use as the image store in the format <code>datastore/<i>path</i></code>. In this case, the virtual container host uses <code><i>path</i></code> as the image store instead of using the folder with the same name as the virtual container host. If the folder that you specify does not already exist, `vic-machine create` creates it. 

<pre>--image-store <i>datastore_name</i>/<i>path</i></pre> 

Wrap the datastore name or path in single quotes (Linux or Mac OS) or double quotes (Windows) if they include spaces.

<pre>--image-store '<i>datastore name</i>'/'<i>datastore path</i>'</pre> 

If you specify an invalid datastore name, `vic-machine create` fails and suggests valid datastores.

**NOTE**: In the current builds the `container-store` option is not enabled. As a consequence, container VM files are also stored in the datastore that you designate as the image store.

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

Wrap the distributed port group name in single quotes (Linux or Mac OS) or double quotes (Windows) if it includes spaces.

<pre>--bridge-network '<i>distributed port group name</i>'</pre>

<a name="networking"></a>
## Networking Options ##
The `vic-machine create` utility allows you to specify different networks for the different types of traffic between containers, the virtual container host, the external internet, and your vSphere environment.

### `bridge-network` ###

See [bridge-network](#bridge) in the section on mandatory options.

<a name="bridge-range"></a>
### `bridge-network-range` ###

Short name: `--bnr`

The range of IP addresses that additional bridge networks can use when container application developers use `docker network create` to create new bridge networks. If you do not specify the `bridge-network-range` option, the IP range for bridge networks is 172.16.0.0/12.

When you specify the bridge network IP range, you specify the IP range as a CIDR.

<pre>--bridge-network-range 192.168.100.0/24</pre>

<a name="external-network"></a>
### `external-network` ###

Short name: `--en`

The network for containers to use to connect to the Internet. Containers use the external network to pull container images, for example from https://hub.docker.com/, and to publish network services. If you define the external network, you can deploy containers directly on the external interface. 

If not specified, containers use the default VM Network for external traffic. If you specify an invalid network name, `vic-machine create` fails and suggests valid networks.

<pre>--external-network <i>network_name</i></pre>

Wrap the network name in single quotes (Linux or Mac OS) or double quotes (Windows) if it includes spaces.

<pre>--external-network '<i>network name</i>'</pre>

<a name="management-network"></a>
### `management-network` ###

Short name: `--mn`

The network that the virtual container host uses to communicate with vCenter Server and ESXi hosts. Container VMs use this network to communicate with the virtual container host. 

If not specified, the virtual container host uses the external network for management traffic. If you specify an invalid network name, `vic-machine create` fails and suggests valid networks.

<pre>--management-network <i>network_name</i></pre>

Wrap the network name in single quotes (Linux or Mac OS) or double quotes (Windows) if it includes spaces.

<pre>--management-network '<i>network name</i>'</pre>

<a name="client-network"></a>
### `client-network` ###

Short name: `--cln`

The network that the virtual container host uses to generate the Docker API. The Docker API only uses this network.

If not specified, the virtual container host uses the external network for client traffic. If you specify an invalid network name, `vic-machine create` fails and suggests valid networks.

<pre>--client-network <i>network_name</i></pre>

Wrap the network name in single quotes (Linux or Mac OS) or double quotes (Windows) if it includes spaces.

<pre>--client-network '<i>network name</i>'</pre>

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

If the network that you specify in the `container-network` option does not support DHCP, you must also specify the `container-network-gateway` option. 

<pre>--container-network <i>distributed_port_group_name</i>:<i>container_network_name</i></pre>

Wrap the distributed port group name in single quotes (Linux or Mac OS) or double quotes (Windows) if it includes spaces. The descriptive name cannot include spaces.

<pre>--container-network '<i>distributed port group name</i>':<i>container_network_name</i></pre>

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

<a name="compute"></a>
## Compute Resource Options ##

If the vCenter Server instance on which you are deploying a virtual container host only includes a single instance of a standalone host or  cluster, `vic-machine create` automatically detects and uses those resources. If you are deploying to an ESXi host that has no resource pools, `vic-machine create` automatically uses the default resource pool. In these cases, you do not need to specify a compute resource when you run `vic-machine create`.

### `compute-resource` ###

Short name: `-r`

The relative path to the host, cluster, or resource pool in which to deploy the virtual container host. You specify the `compute-resource` option in the following circumstances:

- A vCenter Server instance includes multiple instances of standalone hosts or clusters, or a mixture of standalone hosts and clusters.
- An ESXi host includes multiple resource pools. 
- You want to deploy the virtual container host to a specific resource pool in your environment. 

If you do not specify the `compute-resource` option and multiple possible resources exist, `vic-machine create` fails and suggests valid targets for `compute-resource` in the failure message. 

* To deploy to a specific resource pool on an ESXi host, specify the name of the resource pool: <pre>--compute-resource  <i>resource_pool_name</i></pre>
* To deploy to a vCenter Server instance that has more than one standalone host that are not part of a cluster, specify the IPv4 address or fully qualified domain name (FQDN) of the target host:<pre>--compute-resource <i>host_address</i></pre>
* To deploy to a vCenter Server with more than one cluster, specify the name of the target cluster: <pre>--compute-resource <i>cluster_name</i></pre>
* To deploy to a specific resource pool on a standalone host that is managed by vCenter Server, specify the IPv4 address or FQDN of the target host and name of the resource pool:<pre>--compute-resource <i>host_name</i>/<i>resource_pool_name</i></pre>
* To deploy to a specific resource pool in a cluster, specify the names of the target cluster and the resource pool:<pre>--compute-resource <i>cluster_name</i>/<i>resource_pool_name</i></pre>
* Wrap the resource names in single quotes (Linux or Mac OS) or double quotes (Windows) if they include spaces:<pre>--compute-resource '<i>cluster name</i>'/'<i>resource pool name</i>'</pre>

<a name="datastore"></a>
## Datastore Options ##
The `vic-machine` utility allows you to specify the datastores in which to store container image files, the files for the virtual container host appliance, container VM files, and container volumes. 

- vSphere Integrated Containers Engine fully supports VMware Virtual SAN datastores. 
- vSphere Integrated Containers Engine supports all alphanumeric characters, hyphens, and underscores in datastore paths and datastore names, but no other special characters.
- If you specify different datastores in the different datastore options, and if no single host in a cluster can access all of those datastores, `vic-machine create` fails with an error.<pre>No single host can access all of the requested datastores. 
Installation cannot continue.</pre>
- If you specify different datastores in the different datastore options, and if only one host in a cluster can access all of them, `vic-machine create` succeeds with a warning.<pre>Only one host can access all of the image/container/volume datastores. 
This may be a point of contention/performance degradation and HA/DRS 
may not work as intended.</pre> 

### `image-store` ###

See [image-store](#image) in the section on mandatory options.

### `container-store` ###

Short name: `--cs`

The datastore in which to store container VM files. When you run a container, container VM files are stored in folders at the top level of the designated datastore. If multiple virtual container hosts use the same container store, all of the container VM files appear at the top level of the container store. You cannot currently designate a specific datastore folder for the VM files of the containers that run in a particular virtual container host.

If you are deploying the virtual container host to a vCenter Server cluster, the datastore that you designate in the `container-store` option must be shared by at least two ESXi hosts in the cluster. Using non-shared datastores is possible, but limits the use of vSphere features such as DRS and High Availability.

You can designate the same datastore as the container store for multiple virtual container hosts.

If you do not specify the `container-store` option, vSphere Integrated Containers Engine stores container VM files in the image store.

If you specify an invalid datastore name, `vic-machine create` fails and suggests valid datastores.

<pre>--container-store <i>datastore_name</i></pre>

Wrap the datastore name in single quotes (Linux or Mac OS) or double quotes (Windows) if it includes spaces.

<pre>--container-store '<i>datastore name</i>'</pre>

**NOTE**: In the current builds the `container-store` option is not enabled. Container VM files are stored in the datastore that you designate as the image store. 

<a name="volume-store"></a>
### `volume-store` ###

Short name: `--vs`

The datastore in which to create volumes when container developers use the `docker volume create` command. When you specify the `volume-store` option, you  provide the name of the target datastore and a label for the volume store. You can optionally provide a path to a specific folder in the datastore in which to create the volume store. If you specify an invalid datastore name, `vic-machine create` fails and suggests valid datastores. 

If you are deploying the virtual container host to a vCenter Server cluster, the datastore that you designate in the `volume-store` option must be shared by at least two ESXi hosts in the cluster. Using non-shared datastores is possible, but limits the use of vSphere features such as DRS and High Availability.

You can designate the same datastore as the volume store for multiple virtual container hosts.

The label that you specify is the volume store name that Docker uses. For example, the volume store label appears in the information for a virtual container host when container developers run `docker info`. Container developers also specify the volume store label in the <code>docker volume create --opt VolumeStore=<i>volume_store_label</i></code> option when they create a  volume.

**IMPORTANT** If you do not specify the `volume-store` option, no  volume store is created and container developers cannot use the `docker volume create` command.

- If you only require one volume store, you can set the volume store label to `default`. If you set the volume store label to `default`, container developers do not need to specify the <code>--opt VolumeStore=<i>volume_store_label</i></code> option when they run `docker volume create`. 

  **NOTE**: If container developers intend to create anonymous volumes by using `docker create -v`, you must create a volume store with a label of `default`.

  <pre>--volume-store <i>datastore_name</i>:default</pre>
- If you specify the target datastore and the volume store label, `vic-machine create` creates a folder named `volumes` under the `VIC` folder on the target datastore. Any volumes that container developers create will appear in the `volumes` folder.

  <pre>--volume-store <i>datastore_name</i>:<i>volume_store_label</i></pre>
- If you specify the target datastore, a datastore path, and the volume store label, `vic-machine create` creates a folder named `VIC/volumes` in the location that you specify in the datastore path. If the folders that you specify in the path do not already exist on the datastore, `vic-machine create` creates the appropriate folder structure. Any volumes that container developers create will appear in the <code><i>path</i>/VIC/volumes</code> folder.

  <pre>--volume-store <i>datastore_name</i>/<i>datastore_path</i>:<i>volume_store_label</i></pre>
- Wrap the datastore name and path in single quotes (Linux or Mac OS) or double quotes (Windows) if they include spaces. The volume store label cannot include spaces.

  <pre>--volume-store '<i>datastore name</i>'/'<i>datastore path</i>':<i>volume_store_label</i></pre>
- You can specify the `volume-store` option multiple times, to create multiple volume stores for the virtual container host.

  <pre>--volume-store <i>datastore_name</i>/path:<i>volume_store_label_1</i>
--volume-store <i>datastore_name</i>/<i>path</i>:<i>volume_store_label_2</i>
[...]
--volume-store <i>datastore_name</i>/<i>path</i>:<i>volume_store_label_n</i>
</pre>

<a name="security"></a>
## Security Options ##

 By default, virtual container hosts automatically generate TLS certificates to authenticate connections with Docker clients. You can configure a virtual container host to use no certificate or a Certificate Authority (CA) certificate to authenticate connections with Docker clients.

### `no-tls` ###

Short name: `-k`

If you do not set the `no-tls` option, `vic-machine` by default generates a TLS certificate and key for the virtual container host to  use to authenticate with a Docker client. Set the `no-tls` option if you do not require certificate-based authentication between the virtual container host and the Docker client.

<pre>--no-tls</pre>

If you use the `no-tls` option, you connect Docker clients to the virtual container host via port 2375, instead of via port 2376.

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
--key '<i>path to key file</i>'/<i>key_file_name</i>.pem</pre> 

### `key` ###

Short name: none

The path to the private key file for use with a custom CA certificate. This option is mandatory if your Docker environment uses certificates that are signed by a CA. For information about how to set up a Docker client to use CA certificates, see https://docs.docker.com/engine/security/https/.

Use this option in combination with the `cert` option, that provides the path to an X.509 certificate file. Include the names of the certificate and key files in the paths. 

<pre>--cert <i>path_to_certificate_file</i>/<i>certificate_file_name</i>.pem 
--key <i>path_to_key_file</i>/<i>key_file_name</i>.pem</pre> 

Wrap the folder names in the paths in single quotes (Linux or Mac OS) or double quotes (Windows) if they include spaces.

<pre>--cert '<i>path to certificate file</i>'/<i>certificate_file_name</i>.pem 
--key '<i>path to key file</i>'/<i>key_file_name</i>.pem</pre>

<a name="deployment"></a>
## vApp Deployment Options ##

The `vic-machine` utility provides options to customize the deployment of virtual container hosts.

### `name` ###

Short name: `-n`

A name for the virtual container host appliance. If not specified, `vic-machine` sets the name of the virtual container host to `virtual-container-host`. If a virtual container host of the same name exists on the ESXi host or in the vCenter Server inventory, or if a folder of the same name exists in the target datastore, the deployment of the virtual container host fails.

<pre>--name <i>vch_appliance_name</i></pre>

Wrap the appliance name in single quotes (Linux or Mac OS) or double quotes (Windows) if it includes spaces.

<pre>--name '<i>vch appliance name</i>'</pre>

### `password` ###

Short name: `-p`

The password for the user account on the vCenter Server on which you  are deploying the virtual container host, or the password for the ESXi host if you are deploying directly to an ESXi host. If not specified, `vic-machine` prompts you to enter the password during deployment.

<pre>--password <i>esxi_host_or_vcenter_server_password</i></pre>

Wrap the password in single quotation marks (') on Mac OS and Linux and in double quotation (") marks on Windows if it includes special characters.

<pre>--password '<i>esxi_host_or_vcenter_server_p@ssword</i>'</pre>

### `force` ###

Short name: `-f`

Forces `vic-machine create` to ignore warnings and non-fatal errors and continue with the deployment of a virtual container host. Errors such as an incorrect compute resource still cause the installation to fail. 

<pre>--force</pre>

### `timeout` ###

Short name: none

The timeout period for uploading the vSphere Integrated Containers Engine  appliance and container images to the ESXi host, and for powering on the appliance. Specify a value in the format `XmYs` if the default timeout of 3m0s is insufficient.

<pre>--timeout 5m0s</pre> 

### `appliance-iso` ###

Short name: `--ai`

The ISO image from which the virtual container host appliance boots. Omit this option to boot the appliance from the default ISO that is included with `vic-machine`. Set this option to boot the appliance from a different ISO file, for example to reinstall an existing virtual container host or to update it to a newer version. Include the name of the ISO file in the path.

<pre>--appliance-iso <i>path_to_ISO_file</i>/<i>ISO_file_name</i>.iso</pre>

Wrap the folder names in the path in single quotes (Linux or Mac OS) or double quotes (Windows) if they include spaces.

<pre>--appliance-iso '<i>path to ISO file</i>'/<i>ISO_file_name</i>.iso</pre>

### `bootstrap-iso` ###

Short name: `--bi`

The ISO image from which container VMs boot. Omit this option to boot container VMs from the default Photon OS ISO that is included with `vic-machine`. Set this option to a different ISO file to boot container VMs with an operating system other than Photon OS. Include the name of the ISO file in the path.

<pre>--bootstrap-iso <i>path_to_ISO_file</i>/<i>ISO_file_name</i>.iso</pre>

Wrap the folder names in the path in single quotes (Linux or Mac OS) or double quotes (Windows) if they include spaces.

<pre>--bootstrap-iso '<i>path to ISO file</i>'/<i>ISO_file_name</i>.iso</pre>


### `pool-memory-reservation` ###

Short name: `--pmr`

Reserve a quantity of memory for use by the virtual container host vApp   and container VMs. Specify the memory reservation value in MB. If not specified, `vic-machine create` sets the reservation to 1.

<pre>--pool-memory-reservation 1024</pre>

### `pool-memory-limit` ###

Short name: `--pml`

Limit the amount of memory that is available for use by the virtual container host vApp and container VMs. Specify the memory limit value in MB. If not specified, `vic-machine create` sets the limit to 0 (unlimited).

<pre>--pool-memory-limit 1024</pre>

### `pool-memory-shares` ###

Short name: `--pms`

Set memory shares on the virtual container host vApp. Specify the share value as a level or a number, for example `high`, `normal`, `low`, or `163840`. If not specified, `vic-machine create` sets the share to `normal`.

<pre>--pool-memory-shares low</pre>

### `pool-cpu-reservation` ###

Short name: `--pcr`

Reserve a quantity of CPU capacity for use by the virtual container host vApp and container VMs.  Specify the CPU reservation value in MHz. If not specified, `vic-machine create` sets the reservation to 1.

<pre>--pool-cpu-reservation 1024</pre>

### `pool-cpu-limit` ###

Short name: `--pcl`

Limit the amount of CPU capacity that is available for use by the virtual container host vApp and container VMs. Specify the CPU limit value in MHz. If not specified, `vic-machine create` sets the reservation to 0 (unlimited).

<pre>--pool-cpu-limit 1024</pre>

### `pool-cpu-shares` ###

Short name: `--pcs`

Set CPU shares on the virtual container host vApp. Specify the share value as a level or a number, for example `high`, `normal`, `low`, or `163840`. If not specified, `vic-machine create` sets the share to `normal`.

<pre>--pool-cpu-shares low</pre>

### `debug` ###
Short name: `-v`

Provide verbose logging output, for troubleshooting purposes when running `vic-machine create`. If not specified, the `debug` value is set to 0 and verbose logging is disabled. Provide a value of 1 or greater to increase the verbosity of the logging. Note that setting debug to a value greater than 1 can affect the behavior of `vic-machine create`.

<pre>--debug 1</pre>

<a name="mandatory"></a>
## Advanced Options ##

The options in this section are exposed in the `vic-machine create` help if you run <code>vic-machine create --advanced-options</code>. 

### `appliance-cpu ` ###

Short name: none

The number of virtual CPUs for the virtual container host VM. The default is 1. Set this option to increase the number of CPUs in the virtual container host VM, for example if the virtual container host will handle large volumes of containers, or containers that require a lot of processing power.

**NOTE** Use the `--pool-cpu` options instead of the `--appliance-cpu` option. The `--appliance-cpu` option is mainly intended for use by VMware Support.

<pre>--appliance-cpu <i>number_of_CPUs</i></pre>

### `appliance-memory ` ###

Short name: none

The amount of memory for the virtual container host VM. The default is 2048MB. Set this option to increase the amount of memory in the virtual container host VM, for example if the virtual container host will handle large volumes of containers, or containers that consume a lot of memory.

**NOTE** Use the `--pool-memory` options instead of the `--appliance-memory` option. The `--appliance-memory` option is mainly intended for use by VMware Support.

<pre>--appliance-memory <i>amount_of_memory</i></pre>

### `use-rp` ###

Short name: none

Deploy the virtual container host to a resource pool rather than to a vApp. If you specify this option, `vic-machine create` creates a resource pool with the same name as the virtual container host.

<pre>--use-rp</pre>
