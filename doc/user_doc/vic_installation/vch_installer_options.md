# Virtual Container Host Deployment Options

The command line utility for vSphere Integrated Containers, `vic-machine`, provides options that allow you to customize the deployment of a virtual container host to correspond to your vSphere environment.

- [Mandatory Options](#mandatory)
- [Networking Options](#networking)
- [Datastore Options](#datastore)
- [Security Options](#security)
- [Deployment Options](#deployment)

<a name="mandatory"></a>
## Mandatory Options ##

The `vic-machine` utility requires you to provide information about where in your vSphere environment to deploy the virtual container host, the user account to use, and the location in which to store files. The options in this section are **mandatory** in all cases.

### `target` ###
The IPv4 address or fully qualified domain name (FQDN) of the ESXi host or vCenter Server instance on which you are installing vSphere Integrated containers.

- If an ESXi host is managed by vCenter Server, you must provide the address of vCenter Server rather than the address of the host.
- If you are installing vSphere Integrated Containers on a vCenter Server instance, you must specify the `user` option. 
- If you are installing vSphere Integrated Containers directly on an ESXi host and you do not specify the `user` option, `vic-machine` uses the `root` account for installation.
- If you do not specify the `passwd` option, `vic-machine` prompts you to enter the password.

To facilitate IP address changes in your infrastructure, provide an FQDN whenever possible, rather than an IP address.

<pre>-target <i>esxi_host_or_vcenter_server_address</i></pre>

### `user` ###
The username for the ESXi host or vCenter Server instance on which you are installing vSphere Integrated containers.

If you are deploying vSphere Integrated Containers on vCenter Server, specify a username for an account that has the Administrator role on that vCenter Server.

<pre>-user <i>esxi_or_vcenter_server_username</i></pre>

### `compute-resource` ###

The path to the host, cluster, or resource pool in which to deploy the virtual container host. You must specify the paths in the `govc` format, including the forward slashes at the beginning and end of the paths. 

* To deploy to an ESXi host that is not managed by vCenter Server, specify the name of an existing resource pool: <pre>-compute-resource 
/ha-datacenter/host/localhost.eng.vmware.com/Resources/<i>resource_pool_name</i>/</pre>
* To deploy to an ESXi host that is not managed by vCenter Server, when the host has no existing resource pools, specify the root resource pool:<pre>-compute-resource /ha-datacenter/host/localhost.eng.vmware.com/Resources/</pre>
* To deploy to a vCenter Server instance that has one or more hosts but no clusters, specify the IPv4 address or fully qualified domain name (FQDN) of the target host:<pre>-compute-resource /<i>datacenter_name</i>/host/<i>host_address</i>/</pre>
* To deploy to a vCenter Server with one or more clusters, specify the name of the target cluster: <pre>-compute-resource /<i>datacenter_name</i>/host/<i>cluster_name</i>/</pre>
* To deploy to a specific resource pool on a standalone host that is managed by vCenter Server, specify the IPv4 address or FQDN of the target host and name of the resource pool:<pre>-compute-resource /<i>datacenter_name</i>/host/<i>host_name</i>/Resources/<i>resource_pool_name</i>/</pre>
* To deploy to a specific resource pool in a cluster, specify the names of the target cluster and the resource pool:<pre>-compute-resource /<i>datacenter_name</i>/host/<i>cluster_name</i>/Resources/<i>resource_pool_name</i>/</pre>

<a name="image"></a>
### `image-store` ###

The datastore in which to store container image files. When you deploy a virtual container host, `vic-machine` creates a folder named `VIC` on the target datastore,  in which to store all of the container images that you pull into a virtual container host. The `vic-machine` utility also places the VM files for the virtual container host in the datastore that you designate as the image store, in a folder that has the same name as the virtual container host.

You can designate the same datastore as the image store for multiple virtual container hosts. In this case, only one `VIC` folder is created in the datastore and the container image files are made available to all of the virtual container hosts that use that image store. 

**NOTE**: In the current builds the `container-store` option is not enabled. As a consequence, container VM files are also stored in the datastore that you designate as the image store.

<pre>-image-store <i>datastore_name</i></pre> 

<a name="networking"></a>
## Networking Options ##
The `vic-machine` utility allows you to specify different networks for the different types of traffic between containers, the virtual container host, the external internet, and your vSphere environment.

### `bridge-network` ###

The network that container VMs use to communicate with each other. You can assign the same network to multiple virtual container hosts. 

The `bridge-network` option is **optional** when your environment only has one host. For example:

* ESXi host with no vCenter Server 
* vCenter Server with no cluster and one ESXi host
* vCenter Server with a cluster and one ESXi host

In a single host environment, if you do not specify `bridge-network`, `vic-machine` creates a standard virtual switch and a port group that each have the same name as the virtual container host. You can optionally specify this option to assign an existing port group for use as the bridge network for container VMs. You can also optionally specify this option to create a new virtual switch and port group that have a different name to the virtual container host.

The `bridge-network` option is **mandatory** in environments with more than one host. For example:

* vCenter Server with no cluster and more than one standalone ESXi host 
* vCenter Server with a cluster with more than one ESXi host

In a multiple-host environment, you must create a private port group before you run `vic-machine`, and pass the port group name to the `bridge-network` option. For information about how to create a port group, see [Create a Private Port Group for Virtual Container Hosts](create_a_private_port_group_for_vch.md).

<pre>-bridge-network <i>network_name</i></pre>

### `external-network` ###

The network for containers to use to connect to the Internet. Containers use the external network to pull container images, for example from https://hub.docker.com/, and to publish network services. If you define the external network, you can deploy containers directly on the external interface. 

If not specified, containers use the default VM Network for external traffic.

<pre>-external-network <i>network_name</i></pre>

### `management-network` ###

The network that the virtual container host uses to communicate with vCenter Server and ESXi hosts. Container VMs use this network to communicate with the virtual container host. 

If not specified, the virtual container host uses the bridge network for management traffic.

<pre>-management-network <i>network_name</i></pre>

<a name="datastore"></a>
## Datastore Options ##
The `vic-machine` utility allows you to specify the datastores in which to store container VM files and container image files. 

### `image-store` ###

See [image-store](#image) in the section on mandatory options.

### `container-store` ###

The datastore in which to store container VM files. When you run a container, container VM files are stored in folders at the top level of the designated datastore. If multiple virtual container hosts use the same container store, all of the container VM files appear at the top level of the container store. You cannot currently designate a specific datastore folder for the VM files of the containers that run in a particular virtual container host.

If you do not specify the `container-store` option, vSphere Integrated Containers stores container VM files in the same datastore that you specify in the mandatory `image-store` option.

**NOTE**: In the current builds the `container-store` option is not enabled. Container VM files are stored in the datastore that you designate as the image store.

<pre>-container-store <i>datastore_name</i></pre> 

<a name="security"></a>
## Security Options ##

You can configure a virtual container host to use an automatically generated certificate, no certificate, or a Certificate Authority (CA) certificate to authenticate connections with Docker clients.

### `generate-cert` ###

If you do not set the `generate-cert` option, `vic-machine` by default generates a TLS certificate and key for the virtual container host to  use to authenticate with a Docker client. Set the `generate-cert` option to `false` if you do not require certificate-based authentication between the virtual container host and the Docker client. If you use the `cert` and `key` options to upload a custom CA certificate, `vic-machine` does not auto-generate a certificate, without requiring you to set `generate-cert` to `false`.

<pre>-generate-cert=false</pre>

### `cert` ###

The path to an X.509 certificate for the Docker API to use to authenticate the virtual container host with a Docker client.

- This option is mandatory if your Docker environment uses TLS certificates that are signed by a CA. For information about how to set up a Docker client to use CA certificates, see https://docs.docker.com/engine/security/https/.
- Use this option in combination with the `key` option, that provides the path to the private key file for the CA certificate.

If you use the `cert` and `key` options, `vic-machine` does not automatically generate certificates. Omit this option if your Docker environment does not use certificates that are signed by a CA. 

<pre>-cert <i>path_to_vcenter_server_certificate</i> -key <i>path_to_vcenter_server_key</i></pre> 

### `key` ###
The path to the private key file for use with a custom CA certificate. This option is mandatory if your Docker environment uses certificates that are signed by a CA. For information about how to set up a Docker client to use CA certificates, see https://docs.docker.com/engine/security/https/.

Use this option in combination with the `cert` option, that provides the path to an X.509 certificate file. 

<pre>-cert <i>path_to_vcenter_server_certificate</i> -key <i>path_to_vcenter_server_key</i></pre>


<a name="deployment"></a>
## Deployment Options ##

The `vic-machine` utility provides options to customize the deployment of virtual container hosts.

### `name` ###
A name for the virtual container host appliance. If not specified, `vic-machine` sets the name of the virtual container host to `docker-appliance`. If a virtual container host of the same name exists on the ESXi host or in the vCenter Server inventory, or if a folder of the same name exists in the target datastore, the deployment of the virtual container host fails. Use the `force` option to replace an existing virtual container host with a new one with the same name.
<pre>-name <i>vch_appliance_name</i></pre>

### `passwd` ###
The password for the user account on the vCenter Server on which you  are deploying the virtual container host, or the password for the ESXi host if you are deploying directly to an ESXi host. If not specified, `vic-machine` prompts you to enter the password during deployment.

**NOTE**: If your password contains special characters, you must wrap the password in single quotation marks (').

<pre>-passwd '<i>esxi_host_or_vcenter_server_p@ssword</i>'</pre>

### `force` ###
Forces `vic-machine` to ignore warnings and non-fatal errors and continue with the deployment of a virtual container host. Errors such as an incorrect compute resource still cause the installation to fail. Using the `force` option deletes any existing virtual container host appliances or datastore folders that have the same name as the one you are specifying in the current deployment. 

<pre>-force</pre>

### `timeout` ###
The timeout period for uploading the vSphere Integrated Containers  appliance and container images to the ESXi host, and for powering on the appliance. Specify a value in the format `XmYs` if the default timeout of 3m0s is insufficient.

<pre>-timeout 5m0s</pre> 

### `appliance-iso` ###

The ISO image from which the virtual container host appliance boots. Omit this option to boot the appliance from the default ISO that is included with `vic-machine`. Set this option to boot the appliance from a different ISO file, for example to reinstall an existing virtual container host or to update it to a newer version.

<pre>-appliance-iso <i>path_to_ISO_file</i></pre>

### `bootstrap-iso` ###

The ISO image from which container VMs boot. Omit this option to boot container VMs from the default Photon OS ISO that is included with `vic-machine`. Set this option to a different ISO file to boot container VMs with an operating system other than Photon OS.

<pre>-bootstrap-iso <i>path_to_ISO_file</i></pre>

### `appliance-cpu ` ###
The number of virtual CPUs for the virtual container host appliance VM. The default is 1. Set this option to increase the number of CPUs in the virtual container host VM, for example if the virtual container host will handle large volumes of containers, or containers that require a lot of processing power.

<pre>-appliance-cpu <i>number_of_CPUs</i></pre>

### `appliance-memory ` ###
The amount of memory for the virtual container host appliance VM. The default is 2048MB. Set this option to increase the amount of memory in the virtual container host VM, for example if the virtual container host will handle large volumes of containers, or containers that consume a lot of memory.

<pre>-appliance-memory <i>amount_of_memory</i></pre>