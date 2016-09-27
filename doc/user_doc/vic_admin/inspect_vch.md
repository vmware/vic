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

The `vic-machine inspect` command displays information about the virtual container host:

- The virtual container host ID:
  
  <pre>VCH ID: VirtualMachine:vm-101</pre> You can use virtual container host ID when you run the `vic-machine delete` command. Using a virtual container host ID reduces the number of options that you need to specify when you use `vic-machine delete`.
- The version of the `vic-machine` utility and the version of the virtual container host that you are inspecting.
  <pre>Installer version: <i>vic_machine_version</i>-<i>vic_machine_build</i>-<i>tag</i>
VCH version: <i>vch_version</i>-<i>vch_build</i>-<i>tag</i>
VCH upgrade status: 
Installer has same version as VCH
No upgrade available with this installer version</pre>
  If `vic-machine inspect` reports a difference between the version or build number of `vic-machine` and the version or build number of the virtual container host, the upgrade status is `Upgrade available`. 

  **NOTE**: In the current builds, virtual container host upgrade is not yet implemented.
- The address of the VIC Admin portal for the virtual container host.
  
  <pre>vic-admin portal:
https://<i>vch_address</i>:2378</pre>

- The address of the Docker endpoint.

  <pre>DOCKER_HOST=<i>vch_address</i>:2376</pre>
- The Docker command to use to connect to the Docker endpoint.
  <pre>Connect to docker:
docker -H <i>vch_address</i>:2376 --tls info</pre>