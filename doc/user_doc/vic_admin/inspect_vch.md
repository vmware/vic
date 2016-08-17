# Obtain Information About a Virtual Container Host #

You can obtain information about a virtual container host by using the `vic-machine inspect` command.

**Prerequisites**

You have deployed a virtual container host.

**Procedure**

1. On the system on which you run `vic-machine`, navigate to the directory that contains the `vic-machine` utility.
2. Run the `vic-machine inspect` command. 

   The following example includes the options required to obtain information about a named instance of a virtual container host from a simple  vCenter Server environment.

   <pre>$ vic-machine<i>-darwin</i><i>-linux</i><i>-windows</i> inspect
--target <i>vcenter_server_username</i>:<i>password</i>@<i>vcenter_server_address</i>
--name <i>vch_name</i></pre>

**Result**

The `vic-machine inspect` command displays the connection information about the virtual container host:

<pre>vic-admin portal:
https://<i>vch_address</i>:2378
DOCKER_HOST=<i>vch_address</i>:2376
Connect to docker:
docker -H <i>vch_address</i>:2376 --tls info
Completed successfully</pre>