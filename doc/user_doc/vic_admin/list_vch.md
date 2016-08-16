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
--target <i>esxi_host <i>esxi_host_address</i>
--user root
--password <i>esxi_host_passwor</i>
</pre>

   <pre>$ vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> ls
--target <i>vcenter_server_username</i>:<i>password</i>@<i>vcenter_server_address</i>
</pre>


**Result**

The `vic-machine ls` command lists the virtual container hosts that are running on the ESXi host or vCenter Server instance that you specified.

- Virtual container hosts running on an ESXi host:

  <pre>ID      PATH                                         NAME
1       /ha-datacenter/host/<i>host_name</i>/Resources      <i>vch_1</i>
2       /ha-datacenter/host/<i>host_name</i>/Resources      <i>vch_2</i>
[...]   [...]                                        [...]
<i>n</i>       /ha-datacenter/host/<i>host_name</i>/Resources      <i>vch_n</i></pre>

- Virtual container hosts running on a standalone host that is managed by vCenter Server:

  <pre>ID           PATH                                           NAME
vm-<i>id_1</i>      /<i>datacenter</i>/host/<i>host_address</i>/Resources        <i>vch_1</i>
vm-<i>id_2</i>      /<i>datacenter</i>/host/<i>host_address</i>/Resources        <i>vch_2</i>
[...]        [...]                                          [...]
vm-<i>id_n</i>      /<i>datacenter</i>/host/<i>host_address</i>/Resources       <i> vch_n</i>
</pre>

- Virtual container hosts running in a vCenter Server cluster:

  <pre>ID           PATH                                           NAME
vm-<i>id_1</i>      /<i>datacenter</i>/host/<i>cluster_name</i>/Resources        <i>vch_1</i>
vm-<i>id_2</i>      /<i>datacenter</i>/host/<i>cluster_name</i>/Resources        <i>vch_2</i>
[...]        [...]                                          [...]
vm-<i>id_n</i>      /<i>datacenter</i>/host/<i>cluster_name</i>/Resources       <i> vch_n</i>
</pre>

If virtual container hosts are deployed in resource pools on hosts or clusters, the resource pool names appear after `Resources` in the `PATH`.
