# vSphere Integrated Containers Command Line Installer Options

The command line installer for vSphere Integrated Containers, `vic-machine`, provides options that allow you to customize the deployment of a virtual container host to correspond to your vSphere environment.

## `appliance-iso` ##

The ISO image from which the virtual container host appliance boots. Omit this option to boot the appliance from the default ISO that is included with the `vic-machine` installer. Set this option to boot the appliance from a different ISO file, for example to update an existing virtual container host to a newer version.

<pre>-appliance-iso <i>path_to_ISO_file</i></pre>

## `bootstrap-iso` ##

The ISO image from which container VMs boot. Omit this option to boot container VMs from the default Photon OS ISO that is included with the `vic-machine` installer. Set this option to a different ISO file to boot container VMs with an operating system other than Photon OS.

<pre>-bootstrap-iso <i>path_to_ISO_file</i></pre>

## `bridge-network` ##

The network that container VMs use to communicate with each other. If not specified, the installer creates a virtual switch and a port group that each has the same name as the virtual container host that you are deploying. Specify this option to assign an existing virtual switch or port group for use as the bridge network for container VMs. You can also specify this option to create a new virtual switch and port group with a different name to the virtual container host. You can assign the same virtual switch and port group to multiple virtual container hosts.

<pre>-bridge-network <i>network_name</i></pre>

## `cert` ##

The path to an X.509 certificate for vSphere to use to authenticate with the virtual container host.

- This option is mandatory if your vSphere environment uses certificates that are signed by a Certificate Authority (CA).
- Use this option in combination with the `key` option, that provides the path to the private key file for the vCenter Server certificate.

Omit this option if vSphere does not use certificates that are signed by a CA.

<pre>-cert <i>path_to_vcenter_server_certificate</i> -key <i>path_to_vcenter_server_key</i></pre>

For information about how to use CA signed certificates in vSphere 6.0, see [VMware KB 2111219](http://kb.vmware.com/kb/2111219). 

For information about how to configure TLS authentication between the virtual container host and Docker clients, see [Using TLS Authentication with vSphere Integrated Containers](using_tls_with_vic.md).

## `compute-resource` ##

This option is **mandatory**.

- If you are deploying a virtual container host to a standalone ESXi host, you must specify an existing resource pool.<pre> -compute-resource /Resources/<i>resource_pool_name</i>/</pre>
- If you are deploying a virtual container host to a vCenter Server instance, you must specify an existing host, cluster, or resource pool.<pre>-compute-resource /<i>datacenter_name</i>/host/<i>cluster_name</i>/</pre>

## `container-store` ##

Datastore in which to store container VM files.

## `external-network` ##

to connect to the internet. Containers can use this external network to publish network services. After defining the external network, you can deploy containers directly on the external interface.

## `force` ##
Forces the installer to ignore warnings and non-fatal errors and continue with the deployment. Errors such as an incorrect compute resource name still cause the installation to fail. Deletes existing virtual container host appliances or datastore folders that have the same name as the current deployment. 

<pre>-force</pre>

## `generate-cert` ##

Generate a TLS certificate for the virtual container host. The default is to generate a certificate, even if this option is not specified. Set the option to false if you do not require certificate-based authentication.

<pre>-generate-cert false</pre>

## `image-store` ##

Datastore in which to store container image files. This option is **mandatory**.

## `key` ##
The path to the private key file for the vCenter Server certificate. This option is mandatory if your vSphere environment uses certificates that are signed by a Certificate Authority (CA). Use this option in combination with the `cert` option, that provides the path to an X.509 certificate file for vSphere to use to authenticate with the vSphere Integrated Containers extension. 

<pre>-cert <i>path_to_vcenter_server_certificate</i> -key <i>path_to_vcenter_server_key</i></pre>

For information about how to use CA signed certificates in vSphere 6.0, see [VMware KB 2111219](http://kb.vmware.com/kb/2111219). 

For information about how to configure TLS authentication between the virtual container host and Docker clients, see [Using TLS Authentication with vSphere Integrated Containers](using_tls_with_vic.md).

## `management-network` ##

to communicate with vCenter and ESXi hosts. This network also serves as a tether within the containers to communicate with the vSphere Container Host.

## `name` ##
A name for the virtual container host appliance. If not specified, the installer sets the name to `docker-appliance`.
<pre>-name <i>VIC_appliance_name</i></pre>

## `passwd` ##
The password for the vCenter Server user account that you are using to deploy the virtual container host, or the password for the ESXi host. If not specified, the installer prompts you to enter the password during installation.

**NOTE**: If your password contains special characters, you must wrap the password in single quotation marks (').

<pre>-passwd '<i>esxi_host_or_vcenter_server_p@ssword</i>'</pre>

## `target` ##
The address of the ESXi host or vCenter Server instance on which you are installing vSphere Integrated containers. This option is **mandatory**.

- If an ESXi host is managed by a vCenter Server instance, you must provide the address of vCenter Server rather than of the host.
- If you are installing vSphere Integrated Containers on a vCenter Server instance, you must specify the `user` option.
- If you are installing vSphere Integrated Containers directly on an ESXi host and you do not specify the `user` option, the installer uses the `root` account for installation.
- If you do not specify the `passwd` option, the installer prompts you to enter the password.

To facilitate IP address changes in your infrastructure, provide a fully qualified domain name (FQDN) whenever possible, rather than an IP address.

<pre>-target <i>esxi_host_or_vcenter_server_address</i></pre>

## `timeout` ##
The timeout period for uploading the vSphere Integrated Containers  appliance and container images to the ESXi host, and for powering on the appliance. Specify a value in the format `XmYs` if the default timeout of 3m0s is insufficient.

<pre>-timeout 5m0s</pre>

## `user` ##
The username for the ESXi host or vCenter Server instance on which you are installing vSphere Integrated containers. This option is **mandatory**. 

- If you are installing vSphere Integrated Containers directly on an ESXi host and you do not specify this option, the installer uses the `root` account for installation.
- If you are deploying vSphere Integrated Containers on vCenter Server, specify a username for an account that has the Administrator role on that vCenter Server.

<pre>-user <i>esxi_or_vcenter_server_username</i></pre>