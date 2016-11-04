# Obtain Information About a Virtual Container Host #

You can obtain information about a virtual container host by using the `vic-machine inspect` command.

The `vic-machine inspect` command does not include any options in addition to the common options described in [Common `vic-machine` Options](common_vic_options.md).

**Prerequisites**

You have deployed a virtual container host.

**Procedure**

1. On the system on which you run `vic-machine`, navigate to the directory that contains the `vic-machine` utility.
2. Run the `vic-machine inspect` command. 

   The following example includes the options required to obtain information about a named instance of a virtual container host from a simple  vCenter Server environment. 

  - You must specify the username and optionally the password, either in the `--target` option or separately in the `--user` and `--password` options. 
  - If the virtual container host has a name other than the default name, `virtual-container-host`, you must specify the `--name` or `--id` option. 
  - If multiple compute resources exist in the datacenter, you must specify the `--compute-resource` or `--id` option.
  - If your vSphere environment uses untrusted, self-signed certificates, you must also specify the thumbprint of the vCenter Server instance or ESXi host in the `--thumbprint` option. To obtain the thumbprint of the vCenter Server or ESXi host certificate, run `vic-machine` without the specifying the `--thumbprint` option. The inspection of the virtual container host fails, but the resulting error message includes the required certificate thumbprint. You can copy the thumbprint from the error message and run `vic-machine` again, including the `--thumbprint` option.

   <pre>$ vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> inspect
--target <i>vcenter_server_username</i>:<i>password</i>@<i>vcenter_server_address</i>
--thumbprint <i>certificate_thumbprint</i>
--name <i>vch_name</i></pre>

**Result**

The `vic-machine inspect` command displays information about the virtual container host:

- The virtual container host ID:
  
  <pre>VCH ID: VirtualMachine:vm-101</pre> The vSphere Managed Object Reference, or moref, of the virtual container host. You can use virtual container host ID when you run the `vic-machine delete`, `upgrade` or `debug` commands. Using a virtual container host ID reduces the number of options that you need to specify when you run those commands.
- The version of the `vic-machine` utility and the version of the virtual container host that you are inspecting.
  <pre>Installer version: <i>vic_machine_version</i>-<i>vic_machine_build</i>-<i>tag</i>
VCH version: <i>vch_version</i>-<i>vch_build</i>-<i>tag</i>
VCH upgrade status: 
Installer has same version as VCH
No upgrade available with this installer version</pre>
  If `vic-machine inspect` reports a difference between the version or build number of `vic-machine` and the version or build number of the virtual container host, the upgrade status is `Upgrade available`. 
- The address of the VCH Admin portal for the virtual container host.
  
  <pre>vic-admin portal:
https://<i>vch_address</i>:2378</pre>

- The address at which the virtual container host publishes ports.
<pre><i>vch_address</i></pre>
- The Docker environment variables that container developers can use when connecting to this virtual container host.
  - Virtual container host with full TLS authentication with trusted Certificate Authority certificates:
  <pre>DOCKER_TLS_VERIFY=1 
DOCKER_CERT_PATH=<i>path_to_certificates</i>
DOCKER_HOST=<i>vch_address</i>:2376</pre>
  - Virtual container host with TLS authentication with untrusted self-signed certificates:
  <pre>
DOCKER_HOST=<i>vch_address</i>:2376</pre>
  - Virtual container host with no TLS authentication:
  <pre>DOCKER_HOST=<i>vch_address</i>:2375</pre>
- The Docker command to use to connect to the Docker endpoint.
  - Virtual container host with full TLS authentication with trusted Certificate Authority certificates:
  <pre>docker -H <i>vch_address</i>:2376 --tlsverify info</pre>
  - Virtual container host with TLS authentication with untrusted self-signed certificates or no TLS authentication:
  <pre>docker -H <i>vch_address</i>:2376 --tls info</pre>
  - Virtual container host with no TLS authentication:
  <pre>docker -H <i>vch_address</i>:2375 info</pre>