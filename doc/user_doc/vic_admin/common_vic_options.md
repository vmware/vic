# Common `vic-machine` Options #

This section describes the options that are common to all `vic-machine` commands. The common options that `vic-machine` requires relate to the vSphere environment in which you deployed the virtual container host (VCH), and to the VCH itself.  

### `--target` ###

Short name: `-t`

The IPv4 address, fully qualified domain name (FQDN), or URL of the ESXi host or vCenter Server instance on which you deployed the VCH. This option is always **mandatory**.

- If the target ESXi host is not managed by vCenter Server, provide the address of the host.<pre>--target <i>esxi_host_address</i></pre>
- If the target ESXi host is managed by vCenter Server, or if you deployed the VCH to a cluster, provide the address of vCenter Server.<pre>--target <i>vcenter_server_address</i></pre>
- You can include the user name and password in the target URL. <pre>--target <i>vcenter_or_esxi_username</i>:<i>password</i>@<i>vcenter_or_esxi_address</i></pre>

  Wrap the user name or password in single quotes (Linux or Mac OS) or double quotes (Windows) if they include special characters.<pre>'<i>vcenter_or_esxi_usern@me</i>':'<i>p@ssword</i>'@<i>vcenter_or_esxi_address</i></pre>
  
  If you do not include the user name in the target URL, you must specify the `user` option. If you do not specify the `password` option or include the password in the target URL, `vic-machine` prompts you to enter the password.
- If you deployed the VCH on a vCenter Server instance that includes more than one datacenter, include the datacenter name in the target URL. If you include an invalid datacenter name, `vic-machine` fails and suggests the available datacenters that you can specify.<pre>--target <i>vcenter_server_address</i>/<i>datacenter_name</i></pre>


### `--user` ###

Short name: `-u`

The username for the ESXi host or vCenter Server instance on which you deployed the VCH. This option is mandatory if you do not specify the username in the `target` option.

<pre>--user <i>esxi_or_vcenter_server_username</i></pre>

Wrap the user name in single quotes (Linux or Mac OS) or double quotes (Windows) if it includes special characters.

<pre>--user '<i>esxi_or_vcenter_server_usern@me</i>'</pre>


### `--password` ###

Short name: `-p`

The password for the user account on the vCenter Server on which you  deployed the VCH, or the password for the ESXi host if you deployed directly to an ESXi host. If not specified, `vic-machine` prompts you to enter the password.

<pre>--password <i>esxi_host_or_vcenter_server_password</i></pre>

Wrap the password in single quotation marks (') on Mac OS and Linux and in double quotation (") marks on Windows if it includes special characters.

<pre>--password '<i>esxi_host_or_vcenter_server_p@ssword</i>'</pre>

### `--thumbprint` ###

Short name: None

The thumbprint of the vCenter Server or ESXi host certificate. Specify this option if your vSphere environment uses untrusted, self-signed certificates. Alternatively, specifying the `--force` option allows you to omit the `--thumbprint` option. If your vSphere environment uses trusted certificates that are signed by a known Certificate Authority (CA), you do not need to specify the `--thumbprint` option.

To obtain the thumbprint of the vCenter Server or ESXi host certificate, run `vic-machine` without the specifying the `--thumbprint` or `--force` options. The operation fails, but the resulting error message includes the required certificate thumbprint. You can copy the thumbprint from the error message and run `vic-machine` again, including the `thumbprint` option. If you obtain the thumbprint by other means, use upper-case letters and colon delimitation rather than space delimitation when you specify `--thumbprint`.

<pre>--thumbprint <i>certificate_thumbprint</i></pre>

### `--compute-resource` ###

Short name: `-r`

The relative path to the host, cluster, or resource pool in which you deployed the VCH. Specify `--compute-resource` with exactly the same value that you used when you ran `vic-machine create`. You specify the `compute-resource` option in the following circumstances:

- vCenter Server includes multiple instances of standalone hosts or clusters, or a mixture of standalone hosts and clusters.
- You deployed the VCH in a specific resource pool in your environment. 

If you specify the `id` option, you do not need to specify the `compute-resource` option.

If you do not specify the `compute-resource` or `id` options and multiple possible resources exist, `vic-machine` fails and suggests valid targets for `compute-resource` in the failure message. 

* If the VCH is in a specific resource pool on an ESXi host, specify the name of the resource pool: <pre>--compute-resource  <i>resource_pool_name</i></pre>
* If the VCH is on a vCenter Server instance that has more than one standalone host but no clusters, specify the IPv4 address or fully qualified domain name (FQDN) of the target host:<pre>--compute-resource <i>host_address</i></pre>
* If the VCH is on a vCenter Server with more than one cluster, specify the name of the target cluster: <pre>--compute-resource <i>cluster_name</i></pre>
* If the VCH is in a specific resource pool on a standalone host that is managed by vCenter Server, specify the IPv4 address or FQDN of the target host and name of the resource pool:<pre>--compute-resource <i>host_name</i>/<i>resource_pool_name</i></pre>
* If the VCH is in a specific resource pool in a cluster, specify the names of the target cluster and the resource pool:<pre>--compute-resource <i>cluster_name</i>/<i>resource_pool_name</i></pre>
* Wrap the resource names in single quotes (Linux or Mac OS) or double quotes (Windows) if they include spaces:<pre>--compute-resource '<i>cluster name</i>'/'<i>resource pool name</i>'</pre>

### `--name` ###

Short name: `-n`

The name of the VCH. This option is mandatory if the VCH has a name other than the default name, `virtual-container-host`, or if you do not use the `id` option. Specify `--name` with exactly the same value that you used when you ran `vic-machine create`. This option is not used by `vic-machine ls`.

<pre>--name <i>vch_appliance_name</i></pre>

Wrap the appliance name in single quotes (Linux or Mac OS) or double quotes (Windows) if it includes spaces.

<pre>--name '<i>vch appliance name</i>'</pre>

### `--id` ###

Short name: None

The vSphere Managed Object Reference, or moref, of the VCH, for example `vm-100`.  You obtain the ID of a VCH by running `vic-machine ls`. If you specify the `id` option, you do not need to specify the `--name` or `--compute-resource` options. This option is not used by `vic-machine create` or `vic-machine version`.

<pre>--id <i>vch_id</i></pre>

### `--timeout` ###

Short name: none

The timeout period for performing operations on the VCH. Specify a value in the format `XmYs` if the default timeout is insufficient.

<pre>--timeout 5m0s</pre> 
