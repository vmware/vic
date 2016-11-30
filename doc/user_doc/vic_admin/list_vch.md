# List VCHs and Obtain Their IDs #

You can obtain a list of the virtual container hosts (VCHs) that are running in vCenter Server or on an ESXi host by using the `vic-machine ls` command. 

The `vic-machine ls` command lists VCHs with their IDs, names, and versions.
<!--
 The `vic-machine ls` command informs you whether upgrades are available for the VCHs.
-->
The `vic-machine ls` command does not include any options in addition to the common options described in [Common `vic-machine` Options](common_vic_options.md).


**Prerequisites**

You have deployed at least one VCH.

**Procedure**

1. On the system on which you run `vic-machine`, navigate to the directory that contains the `vic-machine` utility.
2. Run the `vic-machine ls` command. 

   - To obtain a list of all VCHs that are running on an ESXi host or vCenter Server instance, you must provide the address of the target ESXi host or vCenter Server. 
   - You must specify the username and optionally the password, either in the `--target` option or separately in the `--user` and `--password` options. 
   - If your vSphere environment uses untrusted, self-signed certificates, you must also specify the thumbprint of the vCenter Server instance or ESXi host in the `--thumbprint` option. To obtain the thumbprint of the vCenter Server or ESXi host certificate, run `vic-machine` without the specifying the `--thumbprint` option. The listing of the VCHs fails, but the resulting error message includes the required certificate thumbprint. You can copy the thumbprint from the error message and run `vic-machine` again, including the `--thumbprint` option.

  <pre>$ vic-machine-<i>operating_system</i> ls
--target <i>esxi_host_address</i>
--user root
--password <i>esxi_host_password</i>
--thumbprint <i>certificate_thumbprint</i>
</pre>

   <pre>$ vic-machine-<i>operating_system</i> ls
--target <i>vcenter_server_username</i>:<i>password</i>@<i>vcenter_server_address</i>
--thumbprint <i>certificate_thumbprint</i>
</pre>


**Result**

The `vic-machine ls` command lists the VCHs that are running on the ESXi host or vCenter Server instance that you specified.
<!--
<pre>ID         PATH     NAME    VERSION                     UPGRADE STATUS
vm-101     <i>path</i>     <i>vch_1</i>   <i>vch_version</i>-<i>vch_build</i>-<i>tag</i>   Up to date
vm-102     <i>path</i>     <i>vch_2</i>   <i>vch_version</i>-<i>vch_build</i>-<i>tag</i>   Up to date
[...]
vm-<i>n</i>       <i>path</i>     <i>vch_n</i>   <i>vch_version</i>-<i>vch_build</i>-<i>tag</i>   Up to date
</pre>
-->

<pre>ID         PATH     NAME    VERSION 
vm-101     <i>path</i>     <i>vch_1</i>   <i>vch_version</i>-<i>vch_build</i>-<i>tag</i>
vm-102     <i>path</i>     <i>vch_2</i>   <i>vch_version</i>-<i>vch_build</i>-<i>tag</i>
[...]
vm-<i>n</i>       <i>path</i>     <i>vch_n</i>   <i>vch_version</i>-<i>vch_build</i>-<i>tag</i>
</pre>

- The IDs are the vSphere Managed Object References, or morefs, for the VCH endpoint VMs. You can use VCH IDs when you run the  `vic-machine inspect`, `debug`, and `delete` commands. Using VCH IDs reduces the number of options that you need to specify when you run those commands.
- The `PATH` value depends on where the VCH is deployed:

  - ESXi host that is not managed by vCenter Server:
  <pre>/ha-datacenter/host/<i>host_name</i>/Resources</pre>
  - Standalone host that is managed by vCenter Server: 
  <pre>/<i>datacenter</i>/host/<i>host_address</i>/Resources</pre>
  - vCenter Server cluster:
  <pre>/<i>datacenter</i>/host/<i>cluster_name</i>/Resources</pre>
  - If VCHs are deployed in resource pools on hosts or clusters, the resource pool names appear after `Resources` in the path.
- The `VERSION` value includes the version of `vic-machine`  that was used to create the VCH, the build number of this version, and a hashtag to identify the build.

<!--
- The `UPGRADE STATUS` reflects whether the version of `vic-machine` that you are using is the same as the version of the VCH. If the version or build number of the VCH does not match that of `vic-machine`, `UPGRADE STATUS` is <code>Upgradeable to <i>vch_version</i>-<i>vch_build</i>-<i>tag</i></code>.
-->