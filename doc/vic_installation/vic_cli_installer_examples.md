# Examples of Using the vSphere Integrated Containers Command Line Installer #

This document provides examples of the required command line options when installing vSphere Integrated Containers in different environments.

**NOTE**: If you are installing on Mac OS X or Linux and you use the `-passwd` option, and if the password contains special characters, you must wrap the password in single quotation marks ('). If you are installing on Windows, you can wrap the password in double quotation marks (") or omit the quotation marks.

## Install vSphere Integrated Containers Directly on an ESXi Host ##

You can install vSphere Integrated Containers directly on an ESXi host that is not managed by a vCenter Server instance. This example provides the miniumum options required to install vSphere Integrated Containers. The example assumes that you use the `root` account to log in to the ESXi host. The installer prompts you for the password for the `root` account and deploys a vSphere Integrated Containers appliance with the default name `docker-appliance`. You can optionally specify the `user`, `passwd`, and `name` options. If the host contains more than one datastore, you must set the `datastore` option.
<pre>install<i>-win.exe/osx/linux</i> 
-target=<i>esxi_host_address</i> 
-ceip=enable
</pre>

## Install vSphere Integrated Containers in a vCenter Server Cluster ##

If vCenter Server manages multiple clusters, you must specify the `cluster` option. If the specified cluster contains more than one resource pool or datastore, you must specify those options. However, you cannot use the `host` option to select a specific host in a cluster. This example specifies a datacenter, a cluster in that datacenter, and a datastore in that cluster. The example also sets the vCenter Single Sign-On user and password, and gives the vSphere Integrated Containers appliance a custom name.

<pre>install<i>-win.exe/osx/linux</i> 
-target=<i>vcenter_server_address</i> 
-user=Administrator@vsphere.local 
-passwd=<i>vcenter_sso_password</i>
-datacenter=<i>datacenter_name</i>  
-cluster=/<i>datacenter_name</i>/host/<i>cluster_name</i>/ 
-datastore=<i>datastore_name</i>
-name=<i>VIC_appliance_name</i>
-ceip=enable
</pre>
 
## Install vSphere Integrated Containers on vCenter Server with a Single Standalone Host##

If vCenter Server only manages a single host and that host contains a single datastore, you do not need to specify the `host` or `datastore` options. This example also specifies the vCenter Single Sign-On user and password, and gives the vSphere Integrated Containers appliance a custom name.

<pre>install<i>-win.exe/osx/linux</i> 
-target=<i>vcenter_server_address</i> 
-user=Administrator@vsphere.local 
-passwd=<i>vcenter_sso_password</i>
-name=<i>VIC_appliance_name</i>
-ceip=enable
</pre>

## Install vSphere Integrated Containers on vCenter Server with Multiple Standalone Hosts ##

If vCenter Server manages multiple standalone hosts, you must specify the `host` option. If vCenter Server does not manage multiple datacenters or  multiple datastores per host, you do not need to specify those options. This example also specifies the vCenter Single Sign-On user and password, and gives the vSphere Integrated Containers appliance a custom name.

<pre>install<i>-win.exe/osx/linux</i> 
-target=<i>vcenter_server_address</i> 
-user=Administrator@vsphere.local 
-passwd=<i>vcenter_sso_password</i>
-host=<i>host_address</i>
-name=<i>VIC_appliance_name</i>
-ceip=enable
</pre>

## Install vSphere Integrated Containers on vCenter Server with Multiple Datacenters ##

If vCenter Server manages multiple datacenters, you must specify the `datacenter` option. If the specified datacenter contains more than one standalone host, cluster, or datastore, you must specify those options. This example specifies a datacenter and a host in that datacenter that contains a single datastore. The example also sets the vCenter Single Sign-On user and password, and gives the vSphere Integrated Containers appliance a custom name.

<pre>install<i>-win.exe/osx/linux</i> 
-target=<i>vcenter_server_address</i> 
-user=Administrator@vsphere.local 
-passwd=<i>vcenter_sso_password</i>
-datacenter=<i>datacenter_name</i> 
-host=<i>host_address</i>
-name=<i>VIC_appliance_name</i>
-ceip=enable
</pre>

## Install vSphere Integrated Containers on vCenter Server with Multiple Datacenters, Standalone Hosts, and Datastores ##

If vCenter Server manages multiple datacenters, you must specify the `datacenter` option. If the specified datacenter contains more than one host, cluster, or datastore, you must specify those options. This example specifies a datacenter, a host in that datacenter, and a specific datastore on that host. The example also sets the vCenter Single Sign-On user and password, and gives the vSphere Integrated Containers appliance a custom name.

<pre>install<i>-win.exe/osx/linux</i> 
-target=<i>vcenter_server_address</i> 
-user=Administrator@vsphere.local 
-passwd=<i>vcenter_sso_password</i> 
-datacenter=<i>datacenter_name</i>  
-host=<i>host_address</i>
-datastore=<i>datastore_name</i>
-name=<i>VIC_appliance_name</i>
-ceip=enable
</pre>

## Install vSphere Integrated Containers in a Resource Pool on a Standalone Host ##

If a standalone host contains resource pools, and you want to install vSphere Integrated Containers in a resource pool, you must specify the `pool` option. This example specifies a datacenter, a host in that datacenter, and a datastore and resource pool on that host. The example also sets the vCenter Single Sign-On user and password, and gives the vSphere Integrated Containers appliance a custom name.

<pre>install<i>-win.exe/osx/linux</i> 
-target=<i>vcenter_server_address</i> 
-user=Administrator@vsphere.local 
-passwd=<i>vcenter_sso_password</i>
-datacenter=<i>datacenter_name</i>  
-host=<i>host_address</i>
-datastore=<i>datastore_name</i>
-pool=/<i>datacenter</i>/host/<i>host_address</i>/Resources/<i>resource_pool_name</i>/  
-name=<i>VIC_appliance_name</i>
-ceip=enable
</pre>


## Install vSphere Integrated Containers in a Resource Pool in a vCenter Server Cluster ##

If a cluster contains resource pools, and you want to install vSphere Integrated Containers in a resource pool, you must specify the `pool` option. This example specifies a datacenter, a cluster in that datacenter, and a datastore and resource pool in that cluster. The example also sets the vCenter Single Sign-On user and password, and gives the vSphere Integrated Containers appliance a custom name.

<pre>install<i>-win.exe/osx/linux</i> 
-target=<i>vcenter_server_address</i> 
-user=Administrator@vsphere.local 
-passwd=<i>vcenter_sso_password</i>
-datacenter=<i>datacenter_name</i>  
-cluster=/<i>datacenter_name</i>/host/<i>cluster_name</i>/ 
-datastore=<i>datastore_name</i>
-pool=/<i>datacenter</i>/host/<i>cluster_name</i>/Resources/<i>resource_pool_name</i>/ 
-name=<i>VIC_appliance_name</i>
-ceip=enable
</pre>

## Install vSphere Integrated Containers on a Standalone Host in a vCenter Server that also Manages Clusters ##

If a datacenter includes clusters and also includes standalone hosts that are not members of any of the clusters, and if you want to install vSphere Integrated Containers on one of the standalone hosts, you must specify the host address in the `cluster` option and not in the `host` option. This example specifies a datacenter, and uses the `cluster` option to select a standalone host in that datacenter. Setting a host address in the `cluster` option only works if both standalone hosts and clusters are present in the datacenter. The example also sets the vCenter Single Sign-On user and password, and gives the vSphere Integrated Containers appliance a custom name.

<pre>install<i>-win.exe/osx/linux</i> 
-target=<i>vcenter_server_address</i> 
-user=Administrator@vsphere.local 
-passwd=<i>vcenter_sso_password</i> 
-datacenter=<i>datacenter_name</i>  
-cluster=<i>host_address</i> 
-name=<i>VIC_appliance_name</i>
-ceip=enable
</pre>

## Install the vSphere Integrated Containers Appliance with a Static IP Address ##

In environments in which you cannot use DHCP, you can assign a static IP address to the vSphere Integrated Containers appliance. You specify a static IP address by providing an IPv4 CIDR address and the default gateway address. You specify the address of the DNS servers by including one instance of the `-dns` option for each DNS server.

<pre>install<i>-win.exe/osx/linux</i> 
-target=<i>vcenter_server_address</i> 
-user=Administrator@vsphere.local 
-passwd=<i>vcenter_sso_password</i>
-name=<i>VIC_appliance_name</i>
-dns=<i>dns_server_1_address</i>
-dns=<i>dns_server_2_address</i>
-ip=<i>cidr_address</i>%<i>gateway_address</i>
-ceip=enable
</pre>

## Install vSphere Integrated Containers with Customized Appliance Settings ##

You can increase the amount of RAM and the number of CPUs to assign to the vSphere Integrated Containers appliance during installation. This example of a simple installation doubles the default quantities of RAM and CPUs.

<pre>install<i>-win.exe/osx/linux</i> 
-target=<i>vcenter_server_address</i> 
-user=Administrator@vsphere.local 
-passwd=<i>vcenter_sso_password</i>
-name=<i>VIC_appliance_name</i>
-memory=4096
-numCPUs=4
-ceip=enable
</pre>

## Install vSphere Integrated Containers on a vCenter Server Instance that Uses TLS Certificates ##

If vCenter Server verifies extensions by using TLS certificates, you must provide the paths to an X.509 certificate and its key during installation. This example of a simple installation includes the `cert` and `key` options. Specify a path to the certificate and key files that is relative to the location from which you are installing vSphere Integrated Containers.

<pre>install<i>-win.exe/osx/linux</i> 
-target=<i>vcenter_server_address</i> 
-user=Administrator@vsphere.local 
-passwd=<i>vcenter_sso_password</i>
-name=<i>VIC_appliance_name</i>
-cert=<i>path_to_vcenter_server_certificate</i>
-key=<i>path_to_vcenter_server_key</i>
-ceip=enable
</pre>

For information about how to configure TLS authentication in the virtual container host and Docker clients, see [Using TLS Authentication with vSphere Integrated Containers](using_tls_with_vic.md).

## Uninstall vSphere Integrated Containers from an ESXi Host ##
This example uninstalls vSphere Integrated Containers from an ESXi host that is not managed by vCenter Server. If you specified a custom name for the vSphere Integrated Containers appliance during installation, you must set the `name` option. If the host contains more than one datastore, you must set the `datastore` option. This example assumes that you used the `root` account during installation. You can optionally specify the `user` and `passwd` options. The example uses the `yes` option, to answer yes to questions during the uninstallation process. 
<pre>install<i>-win.exe/osx/linux</i> 
-uninstall -yes
-target=<i>esxi_host_address</i> 
-datastore=<i>datastore_name</i>
-name=<i>VIC_appliance_name</i>
</pre>

## Uninstall vSphere Integrated Containers from a vCenter Server Cluster  ##
If vCenter Server manages more than one datacenter, host, cluster, resource pool, or datastore, you must specify those options during uninstallation. If you specified a custom name for the vSphere Integrated Containers appliance during installation, you must set the `name` option. This example assumes that you used the `root` account during installation. The example sets the vCenter Single Sign-On user and password. It uses the `yes` option, to answer yes to questions during the uninstallation process. 

<pre>install<i>-win.exe/osx/linux</i> 
-uninstall -yes
-target=<i>vcenter_server_address</i> 
-user=Administrator@vsphere.local 
-passwd=<i>vcenter_sso_password</i> 
-datacenter=<i>datacenter_name</i>  
-cluster=/<i>datacenter_name</i>/host/<i>cluster_name</i>/ 
-datastore=<i>datastore_name</i>
-pool=/<i>datacenter</i>/host/<i>cluster_name</i>/Resources/<i>resource_pool_name</i>/ 
-name=<i>VIC_appliance_name</i>
</pre>