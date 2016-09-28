# List Virtual Container Hosts and Obtain Their IDs #

You can obtain a list of the virtual container hosts that are running in vCenter Server or on an ESXi host by using the `vic-machine ls` command. 

The `vic-machine ls` command lists virtual container hosts with their IDs. You can use virtual container host IDs when you run the `vic-machine inspect` and `vic-machine delete` commands. Using virtual container host IDs reduces the number of options that you need to specify when you use `vic-machine inspect` and `vic-machine delete`.

**Prerequisites**

You have deployed at least one virtual container host.

**Procedure**

1. On the system on which you run `vic-machine`, navigate to the directory that contains the `vic-machine` utility.
2. Run the `vic-machine ls` command. 

   To obtain a list of all virtual container hosts that are running on an ESXi host or vCenter Server instance, you must provide the address of the target ESXi host or vCenter Server. You must specify the username and optionally the password, either in the `target` option or separately in the `user` and `password` options. 

  <pre>$ vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> ls
--target <i>esxi_host_address</i>
--user root
--password <i>esxi_host_password</i>
</pre>

   <pre>$ vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> ls
--target <i>vcenter_server_username</i>:<i>password</i>@<i>vcenter_server_address</i>
</pre>


**Result**

The `vic-machine ls` command lists the virtual container hosts that are running on the ESXi host or vCenter Server instance that you specified.

<pre>ID         PATH     NAME    VERSION                     UPGRADE STATUS
vm-101     <i>path</i>     <i>vch_1</i>   <i>vch_version</i>-<i>vch_build</i>-<i>tag</i>   Up to date
vm-102     <i>path</i>     <i>vch_2</i>   <i>vch_version</i>-<i>vch_build</i>-<i>tag</i>   Up to date
[...]
vm-<i>n</i>       <i>path</i>     <i>vch_n</i>   <i>vch_version</i>-<i>vch_build</i>-<i>tag</i>   Up to date
</pre>

- You can use the virtual container host ID when running `vic-machine inspect` or `delete`. Using a virtual container host ID reduces the number of options that you need to specify when you use `vic-machine inspect` or `delete`.
- The `PATH` value depends on where the virtual container host is deployed:

  - ESXi host that is not managed by vCenter Server:
  <pre>/ha-datacenter/host/<i>host_name</i>/Resources</pre>
  - Standalone host that is managed by vCenter Server: 
  <pre>/<i>datacenter</i>/host/<i>host_address</i>/Resources</pre>
  - vCenter Server cluster:
  <pre>/<i>datacenter</i>/host/<i>cluster_name</i>/Resources</pre>
  - If virtual container hosts are deployed in resource pools on hosts or clusters, the resource pool names appear after `Resources` in the path.
- The `VERSION` value includes the version of `vic-machine`  that was used to create the virtual container host, the build number of this version, and a hashtag to identify the build.
- The `UPGRADE STATUS` reflects whether the version of `vic-machine` that you are using is the same as the version of the virtual container host. If the version or build number of the virtual container host does not match that of `vic-machine`, `UPGRADE STATUS` is <code>Upgradeable to <i>vch_version</i>-<i>vch_build</i>-<i>tag</i></code>.

  **NOTE**: In the current builds, virtual container host upgrade is not yet implemented.
