# VCH Deployment Options

The command line utility for vSphere Integrated Containers Engine, `vic-machine`, provides a `create` command with options that allow you to customize the deployment of virtual container hosts (VCHs) to correspond to your vSphere environment.

- [vSphere Target Options](#vsphere)
- [Security Options](#security)
- [Private Registry Options](#registry)
- [Datastore Options](#datastore)
- [Networking Options](#networking)
- [Additional Deployment Options](#deployment)

To allow you to fine-tune the deployment of VCHs, `vic-machine create` provides [Advanced Options](#advanced).

- [Advanced Security Options](#adv-security)
- [Options for Specifying a Static IP Address for the VCH Endpoint VM](#static-ip)
- [Options for Configuring a Non-DHCP Network for Container Traffic](#adv-container-net)
- [Options to Configure VCHs to Use Proxy Servers](#proxy)
- [Advanced Resource Management Options](#adv-mgmt)
- [Other Advanced Options](#adv-other)

<a name="vsphere"></a>
## vSphere Target Options ##

The `create` command of the `vic-machine` utility requires you to provide information about where in your vSphere environment to deploy the VCH and the vCenter Server or ESXi user account to use.

### `--target` ###

Short name: `-t`

The IPv4 address, fully qualified domain name (FQDN), or URL of the ESXi host or vCenter Server instance on which you are deploying a VCH. This option is always **mandatory**.

To facilitate IP address changes in your infrastructure, provide an FQDN whenever possible, rather than an IP address.

- If the target ESXi host is not managed by vCenter Server, provide the address of the ESXi host.<pre>--target <i>esxi_host_address</i></pre>
- If the target ESXi host is managed by vCenter Server, or if you are deploying to a cluster, provide the address of vCenter Server.<pre>--target <i>vcenter_server_address</i></pre>
- You can include the user name and password in the target URL. <pre>--target <i>vcenter_or_esxi_username</i>:<i>password</i>@<i>vcenter_or_esxi_address</i></pre>

  Wrap the user name or password in single quotes (Linux or Mac OS) or double quotes (Windows) if they include special characters.<pre>'<i>vcenter_or_esxi_usern@me</i>':'<i>p@ssword</i>'@<i>vcenter_or_esxi_address</i></pre>
  
  If you do not include the user name in the target URL, you must specify the `user` option. If you do not specify the `password` option or include the password in the target URL, `vic-machine create` prompts you to enter the password.
- If you are deploying a VCH on a vCenter Server instance that includes more than one datacenter, include the datacenter name in the target URL. If you include an invalid datacenter name, `vic-machine create` fails and suggests the available datacenters that you can specify. 

  <pre>--target <i>vcenter_server_address</i>/<i>datacenter_name</i></pre>

  Wrap the datacenter name in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if it includes spaces.
  <pre>--target <i>vcenter_server_address</i>/'<i>datacenter name</i>'</pre>

### `--user` ###

Short name: `-u`

The username for the ESXi host or vCenter Server instance on which you are deploying a VCH.

If you are deploying a VCH on vCenter Server, specify a username for an account that has the Administrator role on that vCenter Server instance. 

<pre>--user <i>esxi_or_vcenter_server_username</i></pre>

Wrap the user name in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if it includes special characters.

<pre>--user '<i>esxi_or_vcenter_server_usern@me</i>'</pre>

You can also specify the username in the URL that you pass to `vic-machine create` in the `target` option, in which case the `user` option is not required.

### `--password` ###

Short name: `-p`

The password for the user account on the vCenter Server on which you  are deploying the VCH, or the password for the ESXi host if you are deploying directly to an ESXi host. If not specified, `vic-machine` prompts you to enter the password during deployment.

<pre>--password <i>esxi_host_or_vcenter_server_password</i></pre>

Wrap the password in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if it includes special characters.

<pre>--password '<i>esxi_host_or_vcenter_server_p@ssword</i>'</pre>

You can also specify the username and password in the URL that you pass to `vic-machine create` in the `target` option, in which case the `password` option is not required.

### `--compute-resource` ###

Short name: `-r`

The relative path to the host, cluster, or resource pool in which to deploy the VCH. 

If the vCenter Server instance on which you are deploying a VCH only includes a single instance of a standalone host or a cluster, `vic-machine create` automatically detects and uses those resources. In this case, you do not need to specify a compute resource when you run `vic-machine create`. If you are deploying to an ESXi host and you do not specify `--compute-resource`, `vic-machine create` automatically uses the default resource pool. 

You specify the `--compute-resource` option in the following circumstances:

- A vCenter Server instance includes multiple instances of standalone hosts or clusters, or a mixture of standalone hosts and clusters.
- You want to deploy the VCH to a specific resource pool in your environment. 

If you do not specify the `--compute-resource` option and multiple possible resources exist, or if you specify an invalid resource name, `vic-machine create` fails and suggests valid targets for `--compute-resource` in the failure message. 

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

To obtain the thumbprint of the vCenter Server or ESXi host certificate, run `vic-machine create` without the specifying the `--thumbprint` or `--force` options. The deployment of the VCH fails, but the resulting error message includes the required certificate thumbprint. You can copy the thumbprint from the error message and run vic-machine create again, including the `--thumbprint` option. If you obtain the thumbprint by other means, use upper-case letters and colon delimitation rather than space delimitation when you specify `--thumbprint`.

<pre>--thumbprint <i>certificate_thumbprint</i></pre>

<a name="security"></a>
## Security Options ##

As a convenience, `vic-machine create` provides the option of generating a client certificate, server certificate, and certificate authority (CA) as appropriate when you deploy a VCH. The generated certificates are functional, but they do not allow for fine control over aspects such as expiration, intermediate certificate authorities, and so on. To exercise fine control over the certificates used, obtain or generate custom certificates yourself before deploying a VCH. Use the [`--key`](#key), [`--cert`](#cert), and [`--tls-ca`](#tls-ca) options to pass the custom certificates to `vic-machine create`.

### Restrict access to the Docker API ###
vSphere Integrated Containers Engine authenticates Docker API clients by using client certificates. This configuration is commonly referred to as `tlsverify` in documentation about containers and Docker. A client certificate is accepted if it is signed by a CA that you provide by specifying one or more instances of the `--tls-ca` option. In the case of the certificates that `vic-machine create` generates, `vic-machine create` creates a CA and uses it to create and sign a single client certificate.

When using the Docker client, the client validates the server either by using CAs that are present in the root certificate bundle of the client system, or that are provided explicitly by using the `--tlscacert` option when running Docker commands. As a part of this validation, the server certificate must explicitly state at least one of the following, and must match the name or address that the client uses to access the server:

- The FQDN used to communicate with the server
- The IP address used to communicate with the server
- A wildcard domain that matches all of the FQDNs in a specific subdomain. For an example of a domain wildcard, see [https://en.wikipedia.org/wiki/Wildcard_certificate#Example](https://en.wikipedia.org/wiki/Wildcard_certificate#Example).

<a name="tls-cname"></a>
#### `--tls-cname` ####

Short name: None

The FQDN or IP address to embed in the generated server certificate. Specify an FQDN, IP address, or a domain wildcard. If you provide a  custom server certificate by using the `--cert` option, you can use `--tls-cname` as a sanity check to ensure that the certificate is valid for the deployment.

If you do not specify `--tls-cname` but you do set a static address for the VCH on the client network interface, `vic-machine create` uses that  address for the Common Name, with the same results as if you had specified `--tls-cname=x.x.x.x`. For information about setting a static IP address on the client network, see [Options for Specifying a Static IP Address for the VCH Endpoint VM](#static-ip).

When you specify the `--tls-cname` option, `vic-machine create` performs the following actions during the deployment of the VCH:

- Checks for an existing certificate in either a folder that has the same name as the VCH that you are deploying, or in a location that you specify in the [`--cert-path`](#cert-path) option. If a valid certificate exists that includes the same Common Name attribute as the one that you specify in `--tls-cname`, `vic-machine create` reuses it. Reusing certificates allows you to delete and recreate VCHs for which you have already distributed the certificates to container developers.
- If certificates are present in the certificate folder that include a different Common Name attribute to the one that you specify in `--tls-cname`, `vic-machine create` fails.  
- If a certificate folder does not exist, `vic-machine create` creates a folder with the same name as the VCH, or creates a folder in the location that you specify in the `--cert-path` option. 
- If valid certificates do not already exist, `vic-machine create` creates the following trusted CA, server, and client certificate/key pairs in the certificate folder:
  - `ca.pem`
  - `ca-key.pem`
  - `cert.pem`
  - `key.pem`
  - `server-cert.pem`
  - `server-key.pem`
- Creates a browser-friendly PFX client certificate, `cert.pfx`, to use to authenticate connections to the VCH Admin portal for the VCH.

**NOTE**: The folder and file permissions for the generated certificate and key are readable only by the user who created them.

Running `vic-machine create` with the `--tls-cname` option also creates an environment file named <code><i>vch_name</i>.env</code>, that contains Docker environment variables that container developers can use to configure their Docker client environment:

- Activates TLS client verification.<pre>DOCKER_TLS_VERIFY=1</pre>
- The path to the client certificates.<pre>DOCKER_CERT_PATH=<i>path_to_certs</i></pre>
- The address of the VCH.<pre>DOCKER_HOST=<i>vch_address</i>:2376</pre>

You must provide copies of the `cert.pem` and `key.pem` client certificate files and the environment file to container developers so that they can connect Docker clients to the VCH. If you deploy the VCH with the `--tls-cname` option, container developers must configure the client appropriately with one of the following options:
- By using the following `tlsverify`, `tlscert`, and `tlskey` Docker options, adding `tlscacert` if a custom CA was used to sign the server certificate.
- By setting `DOCKER_CERT_PATH=/path/to/client/cert.pem` and `DOCKER_TLS_VERIFY=1`.

<pre>--tls-cname vch-name.example.org</pre>
<pre>--tls-cname *.example.org</pre>

<a name="tls-ca"></a>
#### `--tls-ca` ####

Short name: `--ca`

You can specify `--tls-ca` multiple times, to point `vic-machine create` to a file that contains the public portion of a CA. `vic-machine create` uses these CAs to validate client certificates that are offered as credentials for Docker API access. This does not need to be the same CA that you use to sign the server certificate.

<pre>--tls-ca <i>path_to_ca_file</i></pre>

**NOTE**: The `--tls-ca` option appears in the extended help that you see by running <code>vic-machine-<i>os</i> create --extended-help</code> or <code>vic-machine-<i>os</i> create -x</code>.

### Unrestricted access to the Docker API ###
To deploy a VCH that does not restrict access to the Docker API,  use the `--no-tlsverify` option. 

<a name="no-tlsverify"></a>
#### `--no-tlsverify` ####

Short name: `--kv`

The `--no-tlsverify` option prevents the use of CAs for client authentication. You still require a server certificate if you use `--no-tlsverify`. You can supply a custom server certificate by using the  [`--cert`](#cert) and [`--key`](#key)  options. If you do not use `--cert` and `--key` to supply a custom server certificate, `vic-machine create` generates a self-signed server certificate.

When you specify the `--no-tlsverify` option, `vic-machine create` performs the following actions during the deployment of the VCH.

- Generates a self-signed server certificate if you do not specify `--cert` and `--key`.
- Creates a folder with the same name as the VCH in the location in which you run `vic-machine create`.
- Creates an environment file named <code><i>vch_name</i>.env</code> in that folder, that contains the `DOCKER_HOST=vch_address` environment variable, that you can provide to container developers to use to set up their Docker client environment.

If you deploy a VCH with the `--no-tlsverify` option, container developers run Docker commands with the `--tls` option, and the `DOCKER_TLS_VERIFY` environment variable must not be set. Note that setting `DOCKER_TLS_VERIFY` to 0 or `false` has no effect. 

The `--no-tlsverify` option takes no arguments. 

<pre>--no-tlsverify</pre>

<a name="registry"></a>
## Private Registry Options ##

If container developers need to access Docker images that are stored in private registry servers, you must configure VCHs to allow them to connect to the private registry servers when you deploy the VCHs. VCHs can connect to both secure and insecure private registry servers.

<a name="registry-ca"></a>
### `--registry-ca` ###

Short name: `--rc`

The path to a CA certificate that can validate the registry's server certificate. You can specify `--registry-ca` multiple times to specify multiple CA certificates for different registries. The use of these certificates is independent of the client security options that you specify. For example, it is possible to disable TLS for client authentication by using `--no-tls`, and to use `--registry-ca` tospecify CA certificates to validate a private registry.

<pre>--registry-ca <i>path_to_ca_cert_1</i>
--registry-ca <i>path_to_ca_cert_2</i>
</pre>

**NOTE**: The `--registry-ca` option appears in the extended help that you see by running <code>vic-machine-<i>os</i> create --extended-help</code> or <code>vic-machine-<i>os</i> create -x</code>.

<a name="insecure-registry"></a>
### `--insecure-registry` ###

Short name: `--dir`

An insecure private registry server is a private registry server for Docker images that does not provide TLS. The VCH cannot confirm the identity of the remote system that it is pulling images from and the communication is not encrypted. Insecure private registries are not recommended in production environments.  

If you authorize a VCH to connect to an insecure private registry server, the VCH attempts to access the registry server via HTTP if access via HTTPS fails. VCHs always use HTTPS when connecting to registry servers for which you have not authorized insecure access.

To permit pulling images from an insecure private registry, use the `--insecure-registry` option. You can specify `--insecure-registry` multiple times if multiple insecure registries are permitted. If the registry server listens on a specific port, add the port number to the URL

<pre>--insecure-registry <i>registry_URL_1</i>
--insecure-registry <i>registry_URL_2</i>:<i>port_number</i>
</pre>

<a name="datastore"></a>
## Datastore Options ##
The `vic-machine` utility allows you to specify the datastore in which to store container image files, container VM files, and the files for the VCH. You can also specify datastores in which to create container volumes. 

- vSphere Integrated Containers Engine fully supports VMware vSAN datastores. 
- vSphere Integrated Containers Engine supports all alphanumeric characters, hyphens, and underscores in datastore paths and datastore names, but no other special characters.
- If you specify different datastores in the different datastore options, and if no single host in a cluster can access all of those datastores, `vic-machine create` fails with an error.
  <pre>No single host can access all of the requested datastores. 
  Installation cannot continue.</pre>
- If you specify different datastores in the different datastore options, and if only one host in a cluster can access all of them, `vic-machine create` succeeds with a warning.
  <pre>Only one host can access all of the image/container/volume datastores. 
  This may be a point of contention/performance degradation and HA/DRS 
  may not work as intended.</pre> 
- VCHs do not support datastore name changes. If a datastore changes name after you have deployed a VCH that uses that datastore, that VCH will no longer function.

<a name="image"></a>
### `--image-store` ###

Short name: `-i`

The datastore in which to store container image files, container VM files, and the files for the VCH. The `--image-store` option is **mandatory** if there is more than one datastore in your vSphere environment. If there is only one datastore in your vSphere environment, the `--image-store` option is not required. 

If you do not specify the `--image-store` option and multiple possible datastores exist, or if you specify an invalid datastore name, `vic-machine create` fails and suggests valid datastores in the failure message. 

If you are deploying the VCH to a vCenter Server cluster, the datastore that you designate in the `image-store` option must be shared by at least two ESXi hosts in the cluster. Using non-shared datastores is possible, but limits the use of vSphere features such as vSphere vMotion&reg; and VMware vSphere Distributed Resource Scheduler&trade; (DRS).

To specify a whole datastore as the image store, specify the datastore name in the `--image-store` option:

<pre>--image-store <i>datastore_name</i></pre>

If you designate a whole datastore as the image store, `vic-machine` creates the following set of folders in the target datastore: 

-  <code><i>datastore_name</i>/VIC/<i>vch_uuid</i>/images</code>, in which to store all of the container images that you pull into the VCH.
- <code><i>datastore_name</i>/<i>vch_name</i></code>, that contains the VM files for the VCH.
- <code><i>datastore_name</i>/<i>vch_name</i>/kvstores</code>, a key-value store folder for the VCH.

You can specify a datastore folder to use as the image store by specifying a path in the `--image-store` option</code>: 

<pre>--image-store <i>datastore_name</i>/<i>path</i></pre> 

If the folder that you specify in `/path` does not already exist, `vic-machine create` creates it. Wrap the datastore name and path in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if they include spaces:  <pre>--image-store '<i>datastore name</i>'/'<i>datastore path</i>'</pre>  

If you designate a datastore folder as the image store, `vic-machine` creates the following set of folders in the target datastore:

- <code><i>datastore_name</i>/<i>path</i>/VIC/<i>vcu_uuid</i>/images</code>, in which to store all of the container images that you pull into the VCH. 
- <code><i>datastore_name</i>/<i>vch_name</i></code>, that contains the VM files for the VCH. This is the same as if you specified a datastore as the image store.
- <code><i>datastore_name</i>/<i>vch_name</i>/kvstores</code>, a key-value store folder for the VCH. This is the same as if you specified a datastore as the image store.

By specifying the path to a datastore folder in the `--image-store` option, you can designate the same datastore folder as the image store for multiple VCHs. In this way, `vic-machine create` creates only one `VIC` folder in the datastore, at the path that you specify. The `VIC` folder contains one <code><i>vch_uuid</i>/images</code> folder for each VCH that you deploy. By creating one <code><i>vch_uuid</i>/images</code> folder for each VCH, vSphere Integrated Containers Engine limits the potential for conflicts of image use between VCHs, even if you share the same image store folder between multiple hosts.

When container developers create containers, vSphere Integrated Containers Engine stores the files for container VMs at the top level of the image store, in folders that have the same name as the containers.

<a name="volume-store"></a>
### `--volume-store` ###

Short name: `--vs`

The datastore in which to create volumes when container developers use the `docker volume create` or `docker create -v` commands. When you specify the `volume-store` option, you  provide the name of the target datastore and a label for the volume store. You can optionally provide a path to a specific folder in the datastore in which to create the volume store. If the folders that you specify in the path do not already exist on the datastore, `vic-machine create` creates the appropriate folder structure. 

The `vic-machine create` creates command creates the `volumes` folder independently from the folders for VCH files so that you can share volumes between VCHs. If you delete a VCH, any volumes that the VCH managed will remain available in the volume store unless you specify the `--force` option when you delete the VCH. You can then assign an existing volume store that already contains data to a newly created VCH. 

**IMPORTANT**: If multiple VCHs will use the same datastore for their volume stores, specify a different datastore folder for each VCH. Do not designate the same datastore folder as the volume store for multiple VCHs.

If you are deploying the VCH to a vCenter Server cluster, the datastore that you designate in the `volume-store` option should be shared by at least two ESXi hosts in the cluster. Using non-shared datastores is possible and `vic-machine create` succeeds, but it issues a warning that this configuration limits the use of vSphere features such as vSphere vMotion and DRS.

The label that you specify is the volume store name that Docker uses. For example, the volume store label appears in the information for a VCH when container developers run `docker info`. Container developers specify the volume store label in the <code>docker volume create --opt VolumeStore=<i>volume_store_label</i></code> option when they create a volume.

If you specify an invalid datastore name, `vic-machine create` fails and suggests valid datastores. 

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
- You can specify the `volume-store` option multiple times, to create multiple volume stores for the VCH.

  <pre>--volume-store <i>datastore_name</i>/path:<i>volume_store_label_1</i>
--volume-store <i>datastore_name</i>/<i>path</i>:<i>volume_store_label_2</i>
[...]
--volume-store <i>datastore_name</i>/<i>path</i>:<i>volume_store_label_n</i>
</pre>

<a name="networking"></a>
## Networking Options ##
The `vic-machine create` utility allows you to specify different networks for the different types of traffic between containers, the VCH, the external internet, and your vSphere environment. For information about the different networks that VCHs use, see [Networks Used by vSphere Integrated Containers Engine](networks.md).

**IMPORTANT**: A VCH supports a maximum of 3 distinct network interfaces. Because the bridge network requires its own port group, at least two of the public, client, and management networks must share a network interface and therefore a port group. Container networks do not go through the VCH, so they are not subject to this limitation. This limitation will be removed in a future release.

By default, `vic-machine create` obtains IP addresses for VCH endpoint VMs by using DHCP. For information about how to specify a static IP address for the VCH endpoint VM on the client, public, and management networks, see [Specify a Static IP Address for the VCH Endpoint VM](#static-ip) in Advanced Options.

If your network access is controlled by a proxy server, see [Options to Configure VCHs to Use Proxy Servers](#proxy) in Advanced Options. 

When you specify different network interfaces for the different types of traffic, `vic-machine create` checks that the firewalls on the ESXi hosts allow connections to port 2377 from those networks. If access to port 2377 on one or more ESXi hosts is subject to IP address restrictions, and if those restrictions block access to the network interfaces that you specify, `vic-machine create` fails with a firewall configuration error:
<pre>Firewall configuration incorrect due to allowed IP restrictions on hosts: 
"/ha-datacenter/host/localhost.localdomain/localhost.localdomain" 
Firewall must permit dst 2377/tcp outbound to the VCH management interface
</pre>

<a name="bridge"></a>
### `--bridge-network` ###

Short name: `-b`

A port group that container VMs use to communicate with each other. 

The `bridge-network` option is **mandatory** if you are deploying a VCH to vCenter Server.

In a vCenter Server environment, before you run `vic-machine create`, you must create a distributed virtual switch and a port group. You must add the target ESXi host or hosts to the distributed virtual switch, and assign a VLAN ID to the port group, to ensure that the bridge network is isolated. For information about how to create a distributed virtual switch and port group, see [vCenter Server Network Requirements](vic_installation_prereqs.md#networkreqs) in *Environment Prerequisites for vSphere Integrated Containers Engine Installation*.

You pass the name of the port group to the `bridge-network` option. Each VCH requires its own port group. 

**IMPORTANT** 
- Do not assign the same `bridge-network` port group to multiple VCHs. Sharing a port group between VCHs might result in multiple container VMs being assigned the same IP address. 
- Do not use the `bridge-network` port group as the target for any of the other `vic-machine create` networking options.

If you specify an invalid port group name, `vic-machine create` fails and suggests valid port groups.

The `bridge-network` option is **optional** when you are deploying a VCH to an ESXi host with no vCenter Server. In this case, if you do not specify `bridge-network`, `vic-machine` creates a  virtual switch and a port group that each have the same name as the VCH. You can optionally specify this option to assign an existing port group for use as the bridge network for container VMs. You can also optionally specify this option to create a new virtual switch and port group that have a different name to the VCH.

<pre>--bridge-network <i>distributed_port_group_name</i></pre>

Wrap the port group name in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if it includes spaces.

<pre>--bridge-network '<i>port group name</i>'</pre>

For information about how to specify a range of IP addresses for additional bridge networks, see [`bridge-network-range`](#bridge-range) in Advanced Networking Options.

<a name="client-network"></a>
### `--client-network` ###

Short name: `--cln`

A port group on which the VCH will make the Docker API available to Docker clients. Docker clients use this network to issue Docker API requests to the VCH.

If not specified, the VCH uses the public network for client traffic. If you specify an invalid port group name, `vic-machine create` fails and suggests valid port groups.

<pre>--client-network <i>port_group_name</i></pre>

Wrap the port group name in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if it includes spaces.

<pre>--client-network '<i>port group name</i>'</pre>

<a name="public-network"></a>
### `--public-network` ###

Short name: `--pn`

A port group for containers to use to connect to the Internet. VCHs use the public network to pull container images, for example from https://hub.docker.com/. Containers that use use port mapping expose network services on the public interface. 

**NOTE**: vSphere Integrated Containers Engine adds a new capability to Docker that allows you to directly map containers to a network by using the `--container-network` option. This is the recommended way to deploy container services.

If not specified, containers use the VM Network for public network traffic. If you specify an invalid port group name, `vic-machine create` fails and suggests valid port groups.

<pre>--public-network <i>port_group</i></pre>

Wrap the network name in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if it includes spaces.

<pre>--public-network '<i>port group name</i>'</pre>

<a name="management-network"></a>
### `--management-network` ###

Short name: `--mn`

A port group that the VCH uses to communicate with vCenter Server and ESXi hosts. Container VMs use this network to communicate with the VCH. 

**IMPORTANT**: Because the management network provides access to your vSphere environment, and because container VMs use this network to communicate with the VCH, always use a secure network for the management network.

When you create a VCH, `vic-machine create` checks that the firewall on ESXi hosts allows connections to port 2377 from the management network of the VCH. If access to port 2377 on ESXi hosts is subject to IP address restrictions, and if those restrictions block access to the management network interface, `vic-machine create` fails with a firewall configuration error:
<pre>Firewall configuration incorrect due to allowed IP restrictions on hosts: 
"/ha-datacenter/host/localhost.localdomain/localhost.localdomain" 
Firewall must permit dst 2377/tcp outbound to the VCH management interface
</pre>

**NOTE**: If the management network uses DHCP, `vic-machine` checks the firewall status of the management network before the VCH receives an IP address. It is therefore not possible to fully assess whether the firewall permits the IP address of the VCH. In this case, `vic-machine create` issues a warning. 

<pre>Unable to fully verify firewall configuration due to DHCP use on management network 
VCH management interface IP assigned by DHCP must be permitted by allowed IP settings 
Firewall allowed IP configuration may prevent required connection on hosts: 
"/ha-datacenter/host/localhost.localdomain/localhost.localdomain" 
Firewall must permit dst 2377/tcp outbound to the VCH management interface
</pre>

If not specified, the VCH uses the public network for management traffic. If you specify an invalid port group name, `vic-machine create` fails and suggests valid port groups.

<pre>--management-network <i>port_group_name</i></pre>

Wrap the network name in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if it includes spaces.

<pre>--management-network '<i>port group name</i>'</pre>

<a name="container-network"></a>
### `--container-network` ###

Short name: `--cn`

A port group for container VMs to use for external communication when container developers  run `docker run` or `docker create` with the `--net` option. 

You can optionally specify one or more container networks. Container networks allow containers to directly attach to a network without having to route through the VCH via network address translation (NAT). Container networks that you add by using the `--container-network` option appear when you run the `docker network ls` command. These networks are available for use by containers. Containers that use these networks are directly attached to the container network, and do not go through the VCH or share the public IP of the VCH.

**IMPORTANT**: For security reasons, whenever possible, use separate port groups for the container network and the management network.

To specify a container network, you provide the name of a port group for the container VMs to use, and an optional descriptive name for the container network for use by Docker.  If you do not specify a descriptive name, Docker uses the vSphere network name. If you specify an invalid port group name, `vic-machine create` fails and suggests valid port groups.

- You can specify a vSphere network as the container network.
- The port group must exist before you run `vic-machine create`. 
- You cannot use the same port group as you use for the bridge network. 
- You can create the port group on the same distributed virtual switch as the port group that you use for the bridge network.
- If the port group that you specify in the `container-network` option does not support DHCP, see [Options for Configuring a Non-DHCP Network for Container Traffic](#adv-container-net) in Advanced Options. 
- The descriptive name appears under `Networks` when you run `docker info` or `docker network ls` on the deployed VCH.
- Container developers use the descriptive name in the `--net` option when they run `docker run` or `docker create`.

You can specify `--container-network` multiple times to add multiple vSphere networks to Docker.

If you do not specify `--container-network`, or if you deploy containers that do not use a container network, the containers' network services are still be available via port mapping through the VCH, by using NAT through the public interface of the VCH.

<pre>--container-network <i>port_group_name</i>:<i>container_port _group_name</i></pre>

Wrap the port group name in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if it includes spaces. The descriptive name cannot include spaces.

<pre>--container-network '<i>port group name</i>':<i>container port group name</i></pre>

<a name="deployment"></a>
## Additional Deployment Options ##

The `vic-machine` utility provides options to customize the VCH.

### `--name` ###

Short name: `-n`

A name for the VCH. If not specified, `vic-machine` sets the name of the VCH to `virtual-container-host`. If a VCH of the same name exists on the ESXi host or in the vCenter Server inventory, or if a folder of the same name exists in the target datastore, `vic-machine create` creates a folder named <code><i>vch_name</i>_1</code>
 
<pre>--name <i>vch_name</i></pre>

Wrap the name in single quotes (') on Mac OS and Linux and in double quotes (") on Windows if it includes spaces.

<pre>--name '<i>vch name</i>'</pre>

### `--memory` ###

Short name: `--mem`

Limit the amount of memory that is available for use by the VCH vApp in vCenter Server, or for the VCH resource pool on an ESXi host. This limit also applies to the container VMs that run in the VCH vApp or resource pool. Specify the memory limit value in MB. If not specified, `vic-machine create` sets the limit to 0 (unlimited).

<pre>--memory 1024</pre>

### `--cpu` ###

Short name: None

Limit the amount of CPU capacity that is available for use by the VCH vApp in vCenter Server, or for the VCH resource pool on an ESXi host. This limit also applies to the container VMs that run in the VCH vApp or resource pool. Specify the CPU limit value in MHz. If not specified, `vic-machine create` sets the limit to 0 (unlimited).

<pre>--cpu 1024</pre>

### `--force` ###

Short name: `-f`

Forces `vic-machine create` to ignore warnings and non-fatal errors and continue with the deployment of a VCH. Errors such as an incorrect compute resource still cause the installation to fail.

If your vSphere environment uses untrusted, self-signed certificates, you can use the `--force` option to deploy a VCH without providing the thumbprint of the vCenter Server or ESXi host in the `thumbprint` option. 

**IMPORTANT** Running `vic-machine create` with the `--force` option rather than providing the certificate thumbprint is not recommended, because it permits man-in-the-middle attacks to go undetected.

<pre>--force</pre>

### `--timeout` ###

Short name: none

The timeout period for uploading the vSphere Integrated Containers Engine files and ISOs to the ESXi host, and for powering on the VCH. Specify a value in the format `XmYs` if the default timeout of 3m0s is insufficient.

<pre>--timeout 5m0s</pre> 

<a name="advanced"></a>
# Advanced Options #

The options in this section are exposed in the `vic-machine create` help if you run <code>vic-machine-<i>operating_system</i> create --extended-help</code>, or <code>vic-machine-<i>operating_system</i> create -x</code>. 

<a name="adv-security"></a>
## Advanced Security Options ##

The advanced security options allow you to customize the authentication of connections from Docker clients to VCHs.

- Add optional information to auto-generated trusted TLS certificates by specifying the `--certificate-key-size`, and `--organization` options.
- Use custom server certificates by using the `--cert` and `--key` options.
- Disable TLS authentication completely by using the `--no-tls` option.

### `--certificate-key-size` ###

Short name: `--ksz`

The size of the key for `vic-machine create` to use when it creates auto-generated trusted certificates. If not specified, `vic-machine create` creates keys with default size of 2048 bits. It is not recommended to use key sizes of less than 2048 bits.

<pre>--certificate-key-size 3072</pre>

### `--organization` ###

Short name: None

A list of identifiers to record in certificates generated by vic-machine. If not specified,`vic-machine create` uses the name of the VCH as the organization value.

**NOTE**: The `client-ip-address` is used for `CommonName` but not  for  `Organisation`.

<pre>--organization <i>organization_name</i></pre>

<a name="cert"></a>
### `--cert` ###

Short name: none

The path to a custom X.509 server certificate. This certificate identifies the VCH endpoint VM both to Docker clients and to browsers  that connect to the VCH Admin portal.

- This certificate should have the following certificate usages:
  - `KeyEncipherment`
  - `DigitalSignature`
  - `KeyAgreement`
  - `ServerAuth`
- This option is mandatory if you use custom TLS certificates, rather than auto-generated certificates.
- Use this option in combination with the `--key` option, that provides the path to the private key file for the custom certificate.
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

The path to the private key file to use with a custom server certificate. This option is mandatory if you specify the `--cert` option, that provides the path to a custom X.509 certificate file. Include the names of the certificate and key files in the paths. 

<pre>--cert <i>path_to_certificate_file</i>/<i>certificate_file_name</i>.pem 
--key <i>path_to_key_file</i>/<i>key_file_name</i>.pem
</pre> 

Wrap the folder names in the paths in single quotes (Linux or Mac OS) or double quotes (Windows) if they include spaces.

<pre>--cert '<i>path to certificate file</i>'/<i>certificate_file_name</i>.pem 
--key '<i>path to key file</i>'/<i>key_file_name</i>.pem
</pre>

<a name="cert-path"></a>
### `--cert-path` ###

Short name: none

By default `--cert-path` is a folder in the current directory, that takes its name from the VCH name that you specify in the `--name` option. `vic-machine create` checks in `--cert-path` for existing certificates with the standard names and uses those certificates if  they are present:
* `server-cert.pem` 
* `server-key.pem`
* `ca.pem`

If `vic-machine create` does not find existing certificates with the standard names in `--cert-path`, or if you do not specify certificates directly by using the `--cert`, `--key`, and `--tls-ca` options, `vic-machine create` generates certificates. Generated certificates are saved in the `--cert-path` folder with the standard names listed. `vic-machine create` additionally generates other certificates:
* `cert.pem` and `key.pem` for client certificates, if required.
* `ca-key.pem`, the private key for the certificate authority. 

<pre>--cert-path '<i>path_to_certificate_folder</i>'
</pre>

<a name="no-tls"></a>
### `--no-tls` ###

Short name: `-k`

Disables TLS authentication of connections between the Docker client and  the VCH. 

Set the `no-tls` option if you do not require TLS authentication between the VCH and the Docker client. Any Docker client can connect to the VCH if you disable TLS authentication. 

If you use the `no-tls` option, container developers connect Docker clients to the VCH via port 2375, instead of via port 2376.

<pre>--no-tls</pre>


<a name="static-ip"></a>
## Options for Specifying a Static IP Address for the VCH Endpoint VM ##

You can specify a static IP address for the VCH endpoint VM on each of the client, public, and management networks. DHCP is used for the endpoint VM for any network on which you do not specify a static IP address.

To specify a static IP address for the endpoint VM on the client, public, or management network, you provide an IP address in the `client/public/management-network-ip` option. If you set a static IP address, you must also provide a gateway address. You can optionally specify one or more DNS server addresses.

### `--dns-server` ###

Short name: None

A DNS server for the VCH endpoint VM to use on the client, public, or management networks. You can specify `dns-server` multiple times, to configure multiple DNS servers.  

- If you specify `dns-server`, `vic-machine create` always uses the `--dns-server` setting for all three of the client, public, and management networks.
- If you do not specify `dns-server` and you specify a static IP address for the endpoint VM on all three of the client, public, and management networks, `vic-machine create` uses the Google public DNS service. 
- If you do not specify `dns-server` and you use a mixture of static IP addresses and DHCP for the client, public, and management networks, `vic-machine create` uses the DNS servers that DHCP provides.
- If you do not specify `dns-server` and you use DHCP for all of the client, public, and management networks, `vic-machine create` uses the DNS servers that DHCP provides.

<pre>--dns-server=172.16.10.10
--dns-server=172.16.10.11
</pre>

### `--client-network-ip`, `--public-network-ip`, `--management-network-ip` ###

Short name: None

A static IP address for the VCH endpoint VM on the public, client, or management network. 

You specify a static IP address for the endpoint VM on the public, client, or management networks by using the `--public/client/management-network-ip` options. If you set a static IP address for the endpoint VM on any of the networks, you must specify a corresponding gateway address by using the `--public/client/management-network-gateway` option. 

- You can only specify one static IP address on a given port group. If more than one of the client, public, or management networks 
- s a port group, you can only specify a static IP address on one of those networks. All of the networks that share that port group use the IP address that you specify. 
- If either of the client or management networks shares a port group with the public network, you can only specify a static IP address on the public network.
- If either or both of the client or management networks do not use the same port group as the public network, you can specify a static IP address for the endpoint VM on those networks by using `--client-network-ip` or `--management-network-ip`, or both. In this case, you must specify a corresponding gateway address by using `client/management-network-gateway`. 
- If the client and management networks both use the same port group, and the public network does not use that port group, you can set a static IP address for the endpoint VM on either or both of the client and management networks.
- If you assign a static IP address to the VCH endpoint VM on the client network by setting the `--client-network-ip` option, and you do not specify one of the TLS options, `vic-machine create` uses this address as the Common Name with which to auto-generate trusted CA certificates. If you do not specify `--tls-cname`, `--no-tls` or `--no-tlsverify`, two-way TLS authentication with trusted certificates is implemented by default when you deploy the VCH with a static IP address on the client network. If you assign a static IP address to the endpoint VM on the client network, `vic-machine create` creates the same certificate and environment variable files as described in the [`--tls-cname` option](#tls-cname).
 
  **IMPORTANT**: If the client network shares a port group with the public network you cannot set a static IP address for the endpoint VM on the client network. To assign a static IP address to the endpoint VM you must set a static IP address on the public network by using the `--public-network-ip` option. In this case, `vic-machine create` uses the public network IP address as the Common Name with which to auto-generate trusted CA certificates, in the same way as it would for the client network.

- If you do not specify an IP address for the endpoint VM on a given network, `vic-machine create` uses DHCP to obtain an IP address for the endpoint VM on that network.
-  When you specify an address, `vic-machine create` uses the netmask from the gateway.

You can specify addresses as IPv4 addresses. Do not use CIDR notation.

<pre>--public-network-ip 192.168.X.N
--management-network-ip 192.168.Y.N
--client-network-ip 192.168.Z.N
</pre>

You can also specify addresses as resolvable FQDNs.

<pre>--public-network-ip=vch27-team-a.internal.domain.com
--management-network-ip=vch27-team-b.internal.domain.com
--client-network-ip=vch27-team-c.internal.domain.com
</pre>

### `--client-network-gateway`, `--public-network-gateway`, `--management-network-gateway` ###

Short name: None

The gateway to use if you use `--public/client/management-network-ip` to specify a static IP address for the VCH endpoint VM on the public, client, or management networks. If you specify a static IP address on any network, you must specify a gateway by using the `--public/client/management-network-gateway` options. 

You specify the public network gateway address in CIDR format.

<pre>--public-network-gateway 192.168.X.1/24</pre>

**IMPORTANT**: Assigning the same subnet to multiple port groups can cause routing problems.  If `vic-machine create` detects that you have assigned the same subnet to multiple port groups, it issues a warning.

The default route for the VCH endpoint VM is always on the public network. As a consequence, if you specify a static IP address on either of the management or client networks, you must specify the routing destination for those networks in the `--management-network-gateway` and `--client-network-gateway` options. You specify the routing destination or destinations in a comma-separated list, with the address of the gateway separated from the routing destinations by a colon (:). You specify the gateway addresses in CIDR format:

<pre>--management-network-gateway <i>routing_destination_1</i>/<i>subnet_mask</i>,
<i>routing_destination_2</i>/<i>subnet_mask</i>:
<i>gateway_address</i>/<i>subnet_mask</i></pre>
<pre>--client-network-gateway <i>routing_destination_1</i>/<i>subnet_mask</i>,
<i>routing_destination_2</i>/<i>subnet_mask</i>:
<i>gateway_address</i>/<i>subnet_mask</i>
</pre>

In the following example, `--management-network-gateway` informs the VCH that it can reach all of the vSphere management endoints that are in the ranges 192.168.3.0-255 and 192.168.128.0-192.168.131.255 by sending packets to the gateway at 192.168.2.1. Ensure that the address ranges that you specify include all of the systems that will connect to this VCH instance. 

<pre>--management-network-gateway 192.168.3.0/24,192.168.128.0/22:192.168.2.1/24
</pre>


<a name="adv-container-net"></a>
## Options for Configuring a Non-DHCP Network for Container Traffic ##
If the network that you specify in the `container-network` option does not support DHCP, you must specify the `container-network-gateway` option. You can optionally specify one or more DNS servers and a range of IP addresses for container VMs on the container network. 

For information about the container network, see the section on the [`container-network` option](#container-network).

### `--container-network-gateway` ###

Short name: `--cng`

The gateway for the subnet of the container network. This option is required if the network that you specify in the `--container-network` option does not support DHCP. Specify the gateway in the format <code><i>container_network</i>:<i>subnet</i></code>. If you specify this option, it is recommended that you also specify the  `--container-network-dns` option.

When you specify the container network gateway, you must use the port group that you specify in the `--container-network` option. If you specify `--container-network-gateway` but you do not specify `--container-network`, or if you specify a different port group to the one that you specify in `--container-network`, `vic-machine create` fails with an error.

<pre>--container-network-gateway <i>port_group_name</i>:<i>gateway_ip_address</i>/<i>subnet_mask</i></pre>

Wrap the port group name in single quotes (Linux or Mac OS) or double quotes (Windows) if it includes spaces.

<pre>--container-network-gateway '<i>port group name</i>':<i>gateway_ip_address</i>/<i>subnet_mask</i></pre>

### `--container-network-dns` ###

Short name: `--cnd`

The address of the DNS server for the container network. This option is recommended if the network that you specify in the `--container-network` option does not support DHCP. 

When you specify the container network DNS server, you must use the  port group that you specify in the `--container-network` option. You can specify `--container-network-dns` multiple times, to configure multiple DNS servers. If you specify `--container-network-dns` but you do not specify `--container-network`, or if you specify a different port group to the one that you specify in `--container-network`, `vic-machine create` fails with an error.

<pre>--container-network-dns <i>port_group_name</i>:8.8.8.8</pre>

Wrap the port group name in single quotes (Linux or Mac OS) or double quotes (Windows) if it includes spaces.

<pre>--container-network-dns '<i>port group name</i>':8.8.8.8</pre>

### `--container-network-ip-range` ###

Short name: `--cnr`

The range of IP addresses that container VMs can use if the network that you specify in the `container-network` option does not support DHCP. If you specify `--container-network-ip-range`, VCHs manage the addresses for containers within that range. The range that you specify must not be used by other computers or VMs on the network. If you specify `container-network-gateway` but do not specify `--container-network-ip-range`, the IP range for container VMs is the entire subnet that you specify in `--container-network-gateway`. 

When you specify the container network IP range, you must use the port group that you specify in the `--container-network `option. If you specify `--container-network-ip-range` but you do not specify `--container-network`, or if you specify a different port group to the one that you specify in `--container-network`, `vic-machine create` fails with an error.

<pre>--container-network-ip-range <i>port_group_name</i>:192.168.100.2-192.168.100.254</pre>

You can also specify the IP range as a CIDR.

<pre>--container-network-ip-range <i>port_group_name</i>:192.168.100.0/24</pre>

Wrap the port group name in single quotes (Linux or Mac OS) or double quotes (Windows) if it includes spaces.

<pre>--container-network-ip-range '<i>port group name</i>':192.168.100.0/24</pre>

<a name="proxy"></a>
## Options to Configure VCHs to Use Proxy Servers ##

If access to the Internet or to your private image registries requires the use of a proxy server, you must configure a VCH to connect to the proxy server when you deploy it. The proxy is used only when pulling images, and not for any other purpose.

**IMPORTANT**: Configuring a VCH to use a proxy server does not configure proxy support on the containers that this VCH runs. Container developers must configure proxy servers on containers when they create them. 

### `--https-proxy` ###

Short name: `--sproxy`

The address of the HTTPS proxy server through which the VCH accesses image registries when using HTTPS. Specify the address of the proxy server as either an FQDN or an IP address.

<pre>--https-proxy https://<i>proxy_server_address</i>:<i>port</i></pre>

### `--http-proxy` ###

Short name: `--hproxy`

The address of the HTTP proxy server through which the VCH accesses image registries when using HTTP. Specify the address of the proxy server as either an FQDN or an IP address.

<pre>--http-proxy http://<i>proxy_server_address</i>:<i>port</i></pre>

<a name="adv-mgmt"></a>
## Advanced Resource Management Options ##

You can set limits on the memory and CPU shares and reservations on the VCH. For information about memory and CPU shares and reservations, see [Allocate Memory Resources](https://pubs.vmware.com/vsphere-65/topic/com.vmware.vsphere.vm_admin.doc/GUID-49D7217C-DB6C-41A6-86B3-7AFEB8BF575F.html), and [Allocate CPU Resources](https://pubs.vmware.com/vsphere-65/topic/com.vmware.vsphere.vm_admin.doc/GUID-6C9023B2-3A8F-48EB-8A36-44E3D14958F6.html) in the vSphere documentation.

### `--memory-reservation` ###

Short name: `--memr`

Reserve a quantity of memory for use by the VCH vApp in vCenter Server, or for the VCH resource pool on an ESXi host. This limit also applies to the container VMs that run in the VCH vApp or resource pool. Specify the memory reservation value in MB. If not specified, `vic-machine create` sets the reservation to 1.

<pre>--memory-reservation 1024</pre>

### `--memory-shares` ###

Short name: `--mems`

Set memory shares on the VCH vApp in vCenter Server, or on the VCH resource pool on an ESXi host.  This limit also applies to the container VMs that run in the VCH vApp or resource pool. Specify the share value as a level or a number, for example `high`, `normal`, `low`, or `163840`. If not specified, `vic-machine create` sets the share to `normal`.

<pre>--memory-shares low</pre>

### `--cpu-reservation` ###

Short name: `--cpur`

Reserve a quantity of CPU capacity for use by the VCH vApp in vCenter Server, or for the VCH resource pool on an ESXi host. This limit also applies to the container VMs that run in the VCH vApp or resource pool.  Specify the CPU reservation value in MHz. If not specified, `vic-machine create` sets the reservation to 1.

<pre>--cpu-reservation 1024</pre>

### `--cpu-shares` ###

Short name: `--cpus`

Set CPU shares on the VCH vApp in vCenter Server, or on the VCH resource pool on an ESXi host.  This limit also applies to the container VMs that run in the VCH vApp or resource pool. Specify the share value as a level or a number, for example `high`, `normal`, `low`, or `163840`. If not specified, `vic-machine create` sets the share to `normal`.

<pre>--cpu-shares low</pre>

### `--appliance-cpu ` ###

Short name: none

The number of virtual CPUs for the VCH endpoint VM. The default is 1. Set this option to increase the number of CPUs in the VCH VM, for example if the VCH will handle large volumes of containers, or containers that require a lot of processing power.

**NOTE** Use the `--cpu` option instead of the `--appliance-cpu` option to increase the overall CPU capacity of the VCH vApp, rather than on the VCH endpoint VM. The `--appliance-cpu` option is mainly intended for use by VMware Support.

<pre>--appliance-cpu <i>number_of_CPUs</i></pre>

### `--appliance-memory ` ###

Short name: none

The amount of memory for the VCH endpoint VM. The default is 2048MB. Set this option to increase the amount of memory in the VCH endpoint VM if the VCH will pull large container images.

**NOTE** With the exception of VCHs that pull large container images, always use the `--memory` option instead of the `--appliance-memory` option to increase the overall amount of memory for the VCH vApp, rather than on the VCH endpoint VM. Use docker create -m to set the amounts of memory on container VMs. The `--appliance-memory` option is mainly intended for use by VMware Support.

<pre>--appliance-memory <i>amount_of_memory</i></pre>

<a name="adv-other"></a>
## Other Advanced Options ##

<a name="bridge-range"></a>
### `--bridge-network-range` ###

Short name: `--bnr`

The range of IP addresses that additional bridge networks can use when container application developers use `docker network create` to create new bridge networks. If you do not specify the `bridge-network-range` option, the IP range for bridge networks is 172.16.0.0/12.

When you specify the bridge network IP range, you specify the IP range as a CIDR. The smallest subnet that you can specify is /16.

<pre>--bridge-network-range 192.168.100.0/16</pre>


### `--base-image-size` ###

Short name: None

The size of the base image from which to create other images. You should not normally need to use this option. Specify the size in `GB` or `MB`. The default size is 8GB. Images are thin-provisioned, so they do not usually consume 8GB of space.  

<pre>--base-image-size 4GB</pre>

### `--container-store` ###

Short name: `--cs`

The `container-store` option is not enabled. Container VM files are stored in the datastore that you designate as the image store. 

### `--appliance-iso` ###

Short name: `--ai`

The path to the ISO image from which the VCH appliance boots. Set this option if you have moved the `appliance.iso` file to a folder that is not the folder that contains the `vic-machine` binary or is not the folder from which you are running `vic-machine`. Include the name of the ISO file in the path.

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

Deploy the VCH appliance to a resource pool on vCenter Server rather than to a vApp. If you specify this option, `vic-machine create` creates a resource pool with the same name as the VCH.

<pre>--use-rp</pre>

<a name="debug"></a>
### `--debug` ###
Short name: `-v`

Deploy the VCH with a more verbose level of logging, for troubleshooting purposes. Specifying the `--debug` option increases the verbosity of the logging for all aspects of VCH operation, not just deployment. For example, by setting `--debug`, you increase the verbosity of the logging for VCH initialization, VCH services, container VM initialization, and so on. If not specified, the `debug` value is set to 0 and verbose logging is disabled. Provide a value of 1 or greater to increase the verbosity of the logging. Note that setting debug to a value greater than 1 can affect the behavior of `vic-machine create`. For example, setting `--debug` to 3 suppresses the restart of failed components.

Additionally, deploying a VCH with a `debug` value of 3 or higher enables access to the VCH endpoint VM console by default, with a root password of `password`. This functionality enables you to perform targeted interactive diagnostics in environments in which a VCH endpoint VM failure occurs consistently and in a fashion that prevents `vic-machine debug` from functioning. 

**IMPORTANT**: There is no provision for persistently changing the default root password. This configuration should not be used other than for debugging in a secured environment

<pre>--debug 1</pre>