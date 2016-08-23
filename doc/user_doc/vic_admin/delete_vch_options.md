#  Virtual Container Host Delete Options #

The command line utility for vSphere Integrated Containers, `vic-machine`, provides a `delete` command that allows you to cleanly remove virtual container hosts. The options that `vic-machine delete` requires depend on the location in your vSphere environment in which you deployed the virtual container host.

### `target` ###

Short name: `-t`

The IPv4 address, fully qualified domain name (FQDN), or URL of the ESXi host or vCenter Server instance on which you deployed the virtual container host. This option is mandatory.

- If the target ESXi host is not managed by vCenter Server, provide the address of the host.<pre>--target <i>esxi_host_address</i></pre>
- If the target ESXi host is managed by vCenter Server, or if you deployed the virtual container host to a cluster, provide the address of vCenter Server.<pre>--target <i>vcenter_server_address</i></pre>
- You can optionally include the user name and password of the ESXi host or vCenter Server in the target URL. Wrap the user name or password in single quotes (Linux or Mac OS) or double quotes (Windows) if they include special characters.<pre>--target <i>username</i>:<i>password</i>@<i>esxi_or_vcenter_server_address</i></pre>
- If you deployed the virtual container host on a vCenter Server instance that includes more than one datacenter, include the datacenter name in the target URL. If you include an invalid datacenter name, `vic-machine delete` fails and suggests the available datacenters that you can specify.<pre>--target <i>vcenter_server_address</i>/<i>datacenter_name</i></pre>
- If you do not specify the `password` option or include the password in the target URL, `vic-machine inspect` prompts you to enter the password.

### `user` ###

Short name: `-u`

The username for the ESXi host or vCenter Server instance on which you deployed the virtual container host. This option is mandatory if you do not specify the username in the `target` option.

<pre>--user <i>esxi_or_vcenter_server_username</i></pre>

### `password` ###

Short name: `-p`

The password for the user account on the vCenter Server on which you  deployed the virtual container host, or the password for the ESXi host if you deployed directly to an ESXi host. If not specified, `vic-machine inspect` prompts you to enter the password.

**NOTE**: If your password contains special characters, you must wrap the password in single quotation marks (') on Mac OS and Linux and in double quotation (") marks on Windows.

<pre>--password '<i>esxi_host_or_vcenter_server_p@ssword</i>'</pre>

### `compute-resource` ###

Short name: `-r`

The relative path to the host, cluster, or resource pool in which you deployed the virtual container host. Specify `--compute-resource` with exactly the same value that you used when you ran `vic-machine create`. You specify the `compute-resource` option in the following circumstances:

- vCenter Server includes multiple instances of standalone hosts or clusters, or a mixture of standalone hosts and clusters.
- The ESXi host includes multiple resource pools. 
- You deployed the virtual container host in a specific resource pool in your environment. 

If you specify the `id` option, you do not need to specify the `compute-resource` option.

If you do not specify the `compute-resource` or `id` options and multiple possible resources exist, `vic-machine delete` fails and suggests valid targets for `compute-resource` in the failure message. 

* If the virtual container host is in a specific resource pool on an ESXi host, specify the name of the resource pool: <pre>--compute-resource  <i>resource_pool_name</i></pre>
* If the virtual container host is on a vCenter Server instance that has more than one standalone host but no clusters, specify the IPv4 address or fully qualified domain name (FQDN) of the target host:<pre>--compute-resource <i>host_address</i></pre>
* If the virtual container host is on a vCenter Server with more than one cluster, specify the name of the target cluster: <pre>--compute-resource <i>cluster_name</i></pre>
* If the virtual container host is in a specific resource pool on a standalone host that is managed by vCenter Server, specify the IPv4 address or FQDN of the target host and name of the resource pool:<pre>--compute-resource <i>host_name</i>/<i>resource_pool_name</i></pre>
* If the virtual container host is in a specific resource pool in a cluster, specify the names of the target cluster and the resource pool:<pre>--compute-resource <i>cluster_name</i>/<i>resource_pool_name</i></pre>

### `name` ###

Short name: `-n`

The name of the virtual container host appliance to delete. This option is mandatory if the virtual container host to delete has a name other than the default name, `virtual-container-host`, or if you do not use the `id` option. Specify `--name` with exactly the same value that you used when you ran `vic-machine create`.

<pre>--name <i>vch_appliance_name</i></pre>

### `id` ###

Short name: None

The ID of the virtual container host to delete, for example `vm-100`.  You obtain the ID of a virtual container host by running `vic-machine ls`. If you specify the `id` option, you do not need to specify the `name` or `compute-resource` options.

<pre>--id <i>vch_id</i></pre>

### `force` ###

Short name: `-f`

Forces `vic-machine delete` to ignore warnings and continue with the deletion of a virtual container host. Any running container VMs and any volume stores associated with the virtual container host are deleted. Errors such as an incorrect compute resource still cause the deletion to fail. 

- If you do not specify `force` and the virtual container host contains running container VMs, the deletion fails with a warning. 
- If you do not specify `force` and the virtual container host has volume stores, the deletion of the virtual container host succeeds without deleting the volume stores. The list of volume stores appears in the `vic-machine delete` success message for reference and optional manual removal.

<pre>--force</pre>

### `timeout` ###

Short name: none

The timeout period for deleting the virtual container host. Specify a value in the format `XmYs` if the default timeout of 3m0s is insufficient.

<pre>--timeout 5m0s</pre> 

### `debug` ###
Short name: `-v`

Provide verbose logging output, for troubleshooting purposes when running `vic-machine delete`. If not specified, the `debug` value is set to 0 and verbose logging is disabled. Provide a value of 1 or greater to increase the verbosity of the logging.

<pre>--debug 1</pre>