# vSphere Integrated Containers Command Line Installer Options

The command line installer for vSphere Integrated Containers provides options that allow you to customize the installation to correspond to your vSphere environment.

## `ceip` ##

Enable or disable participation in the VMware Customer Experience Improvement Program. This option is mandatory.

<pre>-ceip=<i>enable</i>/<i>disable</i></pre>

## `cert` ##

The path to an X.509 certificate for vSphere to use to authenticate with the vSphere Integrated Containers extension.

- This option is mandatory if your vSphere environment uses certificates that are signed by a Certificate Authority (CA).
- Use this option in combination with the `key` option, that provides the path to the private key file for the vCenter Server certificate.

Omit this option if vSphere does not use certificates that are signed by a CA.

<pre>-cert=<i>path_to_vcenter_server_certificate</i> -key=<i>path_to_vcenter_server_key</i></pre>

For information about how to use CA signed certificates in vSphere 6.0, see [VMware KB 2111219](http://kb.vmware.com/kb/2111219). 

For information about how to configure TLS authentication between the virtual container host and Docker clients, see [Using TLS Authentication with vSphere Integrated Containers](using_tls_with_vic.md).

## `cidr` ##
IPv4 CIDR notation to designate the network routing for containers to use. If not specified, vSphere Integrated Containers uses the CIDR address 172.17.0.1/16.

<pre>-cidr=192.168.100.0/22</pre>

## `cluster` ##

The path to the cluster on which to install vSphere Integrated Containers. This option is mandatory if the target datacenter contains more than one cluster.

- If your datacenter includes clusters and also includes standalone hosts that are not members of any of the clusters, and if you want to install vSphere Integrated Containers on one of the standalone hosts, you must specify the host address in the `cluster` option.
- You cannot use the `host` option in combination with the `cluster` option, to specify a particular host in a cluster.

Omit this option if vCenter Server only manages one cluster. Specify the path by using the vSphere `govc` CLI format, including the leading and trailing forward slashes.

<pre>-cluster=/<i>datacenter_name</i>/host/<i>cluster_name</i>/</pre>

## `containerNetwork` ##

The name of a private port group and network for containers to use. Set this option if you want containers to use a specific private port group and network. You must create the private port group and network before you run the installer with this option. For information about creating a private group group, see [Create a Private Port Group for Virtual Container Hosts](create_a_private_port_group_for_vch.md).

If not specified, the installer creates a private port group named `container_port_group` and a network with the same name as the vSphere Integrated Containers appliance.

<pre>-containerNetwork=<i>network_name</i></pre>

## `datacenter` ##

The name of the datacenter in which to install vSphere Integrated Containers. This option is mandatory if vCenter Server manages more than one datacenter.

Omit this option if vCenter Server only manages one datacenter.

<pre>-datacenter=<i>datacenter_name</i></pre>


## `datastore` ##
The name of the datastore in which to store the files of vSphere Integrated Containers appliance. vSphere Integrated Containers uses this datastore to store container images and the files of container virtual machines. This option is mandatory if you are installing vSphere Integrated Containers on an ESXi host or cluster that contains more than one datastore.

- vSphere Integrated Containers supports VMware Virtual SAN datastores.
- vSphere Integrated Containers does not support datastore clusters.

Omit this option if the ESXi host or vCenter Server instance only has one datastore.

<pre>-datastore=<i>datastore_name</i></pre>

## `dns` ##
The address of a DNS server, to allow you to assign static IP addresses by using the `ip` option. You can specify the `dns` option multiple times, to identify multiple DNS servers. If not specified, the installer assigns IP addresses by using DHCP.
<pre>-dns=<i>dns_server_1</i> -dns=<i>dns_server_2</i></pre>

## `dockerOpts` ##
Set Docker options to configure the connection between the vSphere Integrated Containers appliance and your Docker client instance. For example, set the path to the configuration files of the Docker client, set the Docker logging levels, use TLS certificate authentication, or use an insecure registry.

Set standard Docker options in double quotes.   
<pre>-dockerOpts="--log-level=info --debug"</pre>

## `externalNetwork` ##
The name of the network to use for traffic between the vSphere Integrated Containers appliance and vSphere. If not specified, the installer uses the default `VM Network`.

<pre>-externalNetwork=<i>external_network_name</i></pre>

## `force` ##
Forces the installer to ignore warnings and non-fatal errors and continue with installation or uninstallation. Errors such as an incorrect host, cluster, datacenter or datastore name still cause the installation or uninstallation to fail.

<pre>-force</pre>

## `host` ##
The address of the ESXi host on which to install vSphere Integrated Containers. Specify this option if you are installing vSphere Integrated Containers on a vCenter Server instance that manages more than one ESXi host and the hosts are not included in a cluster.

- If your datacenter includes clusters and also includes standalone hosts that are not members of any of the clusters, and if you want to install vSphere Integrated Containers on one of the standalone hosts, you must specify the host address in the `cluster` option.
- You cannot use the `host` option in combination with the `cluster` option, to specify a particular host in a cluster.

To facilitate IP address changes in your infrastructure, provide a fully qualified domain name (FQDN) whenever possible, rather than an IP address.

Omit this option if vCenter Server only manages one ESXi host.

<pre>-host=<i>host_address</i></pre>

## `ip` ##
An IPv4 CIDR address for the installer to use to obtain a static IP address for the vSphere Integrated Containers appliance. Requires you to specify the `dns` option. If not specified, the installer assigns an IP address by using DHCP.

<pre>-dns=<i>dns_server_1</i> -dns=<i>dns_server_2</i> -ip=<i>cidr_address</i>%<i>gateway_address</i></pre>

## `key` ##
The path to the private key file for the vCenter Server certificate. This option is mandatory if your vSphere environment uses certificates that are signed by a Certificate Authority (CA). Use this option in combination with the `cert` option, that provides the path to an X.509 certificate file for vSphere to use to authenticate with the vSphere Integrated Containers extension. 

<pre>-cert=<i>path_to_vcenter_server_certificate</i> -key=<i>path_to_vcenter_server_key</i></pre>

For information about how to use CA signed certificates in vSphere 6.0, see [VMware KB 2111219](http://kb.vmware.com/kb/2111219). 

For information about how to configure TLS authentication between the virtual container host and Docker clients, see [Using TLS Authentication with vSphere Integrated Containers](using_tls_with_vic.md).

## `logfile` ##
The location in which to save the vSphere Integrated Containers installation log file. If not specified, the installer saves the log file in the directory in which you run it, with the name `install.log`. 
<pre>-logfile=<i>path</i>/<i>log_file_name</i>.log</pre>

## `memoryMB` ##
The amount of RAM to assign to the virtual container host. Specify this option if you intend to run large numbers of containers in this virtual container host. If not specified, the installer assigns 2048 MB of RAM to the virtual container host.
<pre>-memoryMB=4096</pre>

## `name` ##
A name for the vSphere Integrated Containers appliance. If not specified, the installer sets the name to `docker-appliance`.
<pre>-name=<i>VIC_appliance_name</i></pre>

## `numCPUs` ##
The number of CPUs to assign to the virtual container host. Specify this option if you intend to run large numbers of containers in this virtual container host. If not specified, the installer creates the appliance with 2 CPUs.

<pre>-numCPUs=4</pre>

## `passwd` ##
The password for the vCenter Server user account that you are using to install vSphere Integrated Containers, or the password for the ESXi host. If not specified, the installer prompts you to enter the password during installation.

**NOTE**: If you are installing on Mac OS X or Linux and your password contains special characters, you must wrap the password in single quotation marks ('). If you are installing on Windows, you can wrap the password in double quotation marks (") or omit the quotation marks.

Windows:
<pre>-passwd="<i>esxi_host_or_vcenter_server_p@ssword</i>"</pre>
Mac OS X or Linux:
<pre>-passwd='<i>esxi_host_or_vcenter_server_p@ssword</i>'</pre>

## `pool` ##
The path to a resource pool in which to place the vSphere Integrated Containers appliance. vSphere Integrated Containers places containers in this resource pool. 

If not specified, the installer places the appliance and containers at the root level of the cluster. Specify the path by using the vSphere `govc` CLI format, including the leading and trailing forward slashes.

<pre>-pool=/<i>datacenter</i>/host/<i>cluster_name</i>/Resources/<i>resource_pool_name</i>/</pre>

## `target` ##
The address of the ESXi host or vCenter Server instance on which you are installing vSphere Integrated containers. This option is mandatory.

- If an ESXi host is managed by a vCenter Server instance, you must provide the address of vCenter Server rather than of the host.
- If you are installing vSphere Integrated Containers on a vCenter Server instance, you must specify the `user` option.
- If you are installing vSphere Integrated Containers directly on an ESXi host and you do not specify the `user` option, the installer uses the `root` account for installation.
- If you do not specify the `passwd` option, the installer prompts you to enter the password.

To facilitate IP address changes in your infrastructure, provide a fully qualified domain name (FQDN) whenever possible, rather than an IP address.

<pre>-target=<i>esxi_host_or_vcenter_server_address</i></pre>

## `timeout` ##
The timeout period for uploading the vSphere Integrated Containers  appliance and container images to the ESXi host, and for powering on the appliance. Specify a value in the format `XmYs` if the default timeout of 3m0s is insufficient.

<pre>-timeout=5m0s</pre>

## `uninstall` ##
Uninstalls vSphere Integrated Containers. Removes the vSphere Integrated Containers vApp and virtual machines from the vCenter Server inventory. Deletes all of the vAPP and VM files from storage. Does not delete any private port groups or networks that the installer created during installation.

- Requires the `target` option.
- If you installed vSphere Integrated Containers on a vCenter Server instance, you must specify the `user` option.
- If you installed vSphere Integrated Containers directly on an ESXi host and you do not specify this option, the installer uses the `root` account.
- If you do not specify the `passwd` option, the installer prompts you to enter the password.
- If you specified a name for the vSphere Integrated Containers appliance during installation, you must specify the `name` option. 
- If your vSphere infrastructure includes more than one datacenter, host, cluster, or datastore, you must specify the relevant options.
- If you specified a resource pool for the vSphere Integrated Containers appliance during installation, you must specify the `pool` option. 
- Specify the `yes` option to answer yes to all questions during the uninstallation process.

**NOTE**: If you do not specify the `yes` option, the installer prompts you to confirm that you want to uninstall vSphere Integrated Containers. Enter the word `yes` to confirm. If you enter `y`, the uninstall operation quits.

<pre>install<i>-win.exe/osx/linux</i> 
-uninstall -yes 
-target=<i>vcenter_server_address</i> 
-user=<i>vcenter_server_username</i> 
-passwd=<i>vcenter_server_password</i> 
-datacenter=<i>datacenter_name</i>  
-cluster=/<i>datacenter_name</i>/host/<i>cluster_name</i>/ 
-pool=/<i>datacenter</i>/host/<i>cluster_name</i>/Resources/<i>resource_pool_name</i>/ 
-name=<i>VIC_appliance_name</i></pre>

## `user` ##
The username for the ESXi host or vCenter Server instance on which you are installing vSphere Integrated containers.

- If you are installing vSphere Integrated Containers directly on an ESXi host and you do not specify this option, the installer uses the root account for installation.
- This option is **mandatory** if you are installing vSphere Integrated Containers on a vCenter Server instance. Specify a username for an account that has the Administrator role on vCenter Server.

<pre>-user=<i>vcenter_server_username</i></pre>

## `verify` ##
Verifies that the installation succeeded by pulling the BusyBox container from Docker Hub. The installer starts, stops, and removes the container before confirming that installation was successful.
<pre>-verify</pre>

## `yes` ##
Use in combination with the `uninstall` option to automatically answer yes to all questions during uninstallation.

<pre>-uninstall -yes</pre>
