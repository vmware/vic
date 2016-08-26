#  Virtual Container Host List Options #

The command line utility for vSphere Integrated Containers, `vic-machine`, provides an `ls` command that lists the virtual container hosts that are running on an ESXi host or vCenter Server instance. The options that `vic-machine ls` provides allow you to list the virtual container hosts that are running in specific areas of your vSphere environment.

### `target` ###

Short name: `-t`

The IPv4 address, fully qualified domain name (FQDN), or URL of the ESXi host or vCenter Server instance on which you deployed the virtual container hosts. This option is mandatory.

- If the target ESXi host is not managed by vCenter Server, provide the address of the host.<pre>--target <i>esxi_host_address</i></pre>
- If the target ESXi host is managed by vCenter Server, or if you deployed the virtual container hosts to a cluster, provide the address of vCenter Server.<pre>--target <i>vcenter_server_address</i></pre>
- You can include the user name and password in the target URL. <pre>--target <i>vcenter_or_esxi_username</i>:<i>password</i>@<i>vcenter_or_esxi_address</i></pre>

  Wrap the user name or password in single quotes (Linux or Mac OS) or double quotes (Windows) if they include special characters.<pre>'<i>vcenter_or_esxi_usern@me</i>':'<i>p@ssword</i>'@<i>vcenter_or_esxi_address</i></pre>
  
  If you do not include the user name in the target URL, you must specify the `user` option. If you do not specify the `password` option or include the password in the target URL, `vic-machine ls` prompts you to enter the password.
- If you deployed the virtual container hosts on a vCenter Server instance that includes more than one datacenter, you can include the datacenter name in the target URL. If you include a datacenter name, `vic-machine ls` lists all of the virtual container hosts that are running in that datacenter. If you do not include a datacenter name, `vic-machine ls` lists all of the virtual container hosts that are running on that vCenter Server instance, for all datacenters.<pre>--target <i>vcenter_server_address</i></pre>
- <pre>--target <i>vcenter_server_address</i>/<i>datacenter_name</i></pre>
- If you do not specify the `password` option or include the password in the target URL, `vic-machine ls` prompts you to enter the password.

### `user` ###

Short name: `-u`

The username for the ESXi host or vCenter Server instance on which you deployed the virtual container hosts. This option is mandatory if you do not specify the username in the `target` option.

<pre>--user <i>esxi_or_vcenter_server_username</i></pre>

Wrap the user name in single quotes (Linux or Mac OS) or double quotes (Windows) if it includes special characters.

<pre>--user '<i>esxi_or_vcenter_server_usern@me</i>'</pre>

### `password` ###

Short name: `-p`

The password for the user account on the vCenter Server on which you  deployed the virtual container hosts, or the password for the ESXi host if you deployed directly to an ESXi host. If not specified, `vic-machine ls` prompts you to enter the password.

<pre>--password '<i>esxi_host_or_vcenter_server_p@ssword</i>'</pre>

Wrap the password in single quotation marks (') on Mac OS and Linux and in double quotation (") marks on Windows if it includes special characters.

<pre>--password '<i>esxi_host_or_vcenter_server_p@ssword</i>'</pre>

### `compute-resource` ###

Short name: `-r`

The relative path to a resource pool on an ESXi host, or to a host, cluster, or resource pool in vCenter Server, in which you deployed virtual container hosts. 

If you specify the `compute-resource` option when the target is a vCenter Server instance that has more than one datacenter, you must include the datacenter name in the `target` option.

If you do not specify the `compute-resource` option, `vic-machine ls` lists all of the virtual container hosts that are running in all of the resource pools on the target ESXi host, or in all of the hosts, clusters, and resource pools on the target vCenter Server instance.

* To list the virtual container hosts in a specific resource pool on an ESXi host that is not managed by vCenter Server, specify the name of the resource pool: <pre>--compute-resource  <i>resource_pool_name</i></pre>
* To list the virtual container hosts on a specific standalone host  that is managed by vCenter Server, specify the IPv4 address or fully qualified domain name (FQDN) of the host:<pre>--compute-resource <i>host_address</i></pre>
* To list the virtual container hosts in a specific vCenter Server cluster, specify the name of the target cluster: <pre>--compute-resource <i>cluster_name</i></pre>
* To list the virtual container hosts in a specific resource pool on a standalone host that is managed by vCenter Server, specify the IPv4 address or FQDN of the target host and name of the resource pool:<pre>--compute-resource <i>host_name</i>/<i>resource_pool_name</i></pre>
* To list the virtual container hosts in a specific resource pool in a cluster, specify the names of the target cluster and the resource pool:<pre>--compute-resource <i>cluster_name</i>/<i>resource_pool_name</i></pre>
* Wrap the resource names in single quotes (Linux or Mac OS) or double quotes (Windows) if they include spaces:<pre>--compute-resource '<i>cluster name</i>'/'<i>resource pool name</i>'</pre>

### `debug` ###
Short name: `-v`

Provide verbose logging output, for troubleshooting purposes when running `vic-machine ls`. If not specified, the `debug` value is set to 0 and verbose logging is disabled. Provide a value of 1 or greater to increase the verbosity of the logging.

<pre>--debug 1</pre>