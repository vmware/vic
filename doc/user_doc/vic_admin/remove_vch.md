# Delete a VCH #

You delete virtual container hosts (VCHs) by using the `vic-machine delete` command.

For descriptions of the options that `vic-machine delete` includes in addition to the [Common `vic-machine` Options](common_vic_options.md) , see [VCH Delete Options](delete_vch_options.md).

When you delete a VCH that uses TLS authentication with trusted Certificate Authority (CA) certificates, `vic-machine delete` does not delete the certificates or the certificate folder, even if you specify the `--force` option. Because `vic-machine delete` does not delete the certificates, you can delete VCHs and create new ones that reuse the same certificates. This is useful if you have already distributed the client certificates for VCHs that you need to recreate.

**Prerequisites**

You have deployed a VCH that you no longer require.

**Procedure**

1. On the system on which you run `vic-machine`, navigate to the directory that contains the `vic-machine` utility.
2. Run the `vic-machine delete` command. 

   The following example includes the options required to remove a VCH from a simple vCenter Server environment. 

  - You must specify the username and optionally the password, either in the `--target` option or separately in the `--user` and `--password` options. 
  - If the VCH has a name other than the default name, `virtual-container-host`, you must specify the `--name` or `--id` option. 
  - If multiple compute resources exist in the datacenter, you must specify the `--compute-resource` or `--id` option.
  - If your vSphere environment uses untrusted, self-signed certificates, you must also specify the thumbprint of the vCenter Server instance or ESXi host in the `--thumbprint` option. To obtain the thumbprint of the vCenter Server or ESXi host certificate, run `vic-machine` without the specifying the `--thumbprint` or `--force` options. The deletion of the VCH fails, but the resulting error message includes the required certificate thumbprint. You can copy the thumbprint from the error message and run `vic-machine` again, including the `--thumbprint` option.

   <pre>$ vic-machine-<i>operating_system</i> delete
--target <i>vcenter_server_username</i>:<i>password</i>@<i>vcenter_server_address</i>
--thumbprint <i>certificate_thumbprint</i>
--name <i>vch_name</i></pre>

3. If the delete operation fails with a message about container VMs that are powered on, run `docker stop` on the containers and run `vic-machine delete`. Alternatively, run `vic-machine delete` with the `--force` option.

  **CAUTION** Running `vic-machine delete` with the `--force` option removes all running container VMs that the VCH manages, as well as any associated volumes and volume stores. It is not recommended to use the `--force` option to remove running containers.

  If your vSphere environment uses untrusted, self-signed certificates, running `vic-machine delete` with the `--force` option allows you to omit the `--thumbprint` option.

   <pre>$ vic-machine-<i>operating_system</i> delete
--target <i>vcenter_server_username</i>:<i>password</i>@<i>vcenter_server_address</i>
--name <i>vch_name</i></i>
--force</pre>