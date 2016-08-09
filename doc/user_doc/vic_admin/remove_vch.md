# Remove a Virtual Container Host #

You remove virtual container hosts by using the `vic-machine delete` command.

**Prerequisites**

You have deployed a virtual container host that you no longer require.

**Procedure**

1. On the system on which you run `vic-machine`, navigate to the directory that contains the `vic-machine` utility.
2. Run the `vic-machine delete` command. 

   The following example includes the options required to remove a named instance of a virtual container host from a simple vCenter Server environment.

   <pre>$ vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> delete
--target <i>vcenter_server_username</i>:<i>password</i>@<i>vcenter_server_address</i>
--name <i>vch_name</i></pre>

3. If the delete operation fails with a message about container VMs that are powered on, run `vic-machine delete` again with the `--force` option.

  **CAUTION** Running `vic-machine delete` with the `--force` option removes all running container VMs that the virtual container host manages, as well as any associated volumes and volume stores.

   <pre>$ vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> delete
--target <i>vcenter_server_username</i>:<i>password</i>@<i>vcenter_server_address</i>
--name <i>cluster_name</i></i>
--force</pre>