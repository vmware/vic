# Authorize SSH Access to the VCH Endpoint VM #

By default, SSH access to the virtual container host (VCH) endpoint VM is disabled. The command line utility for vSphere Integrated Containers Engine, `vic-machine`, provides a `debug` command that allows you to enable SSH access to the VCH endpoint VM. The `debug` command also allows you to set a password for the root user account on the endpoint VM. You can also use `debug` to upload a key file for public key authentication when accessing the endpoint VM. 

**IMPORTANT**: If you set a password for the VCH endpoint VM, this password does not persist if you reboot the VM. You must run `vic-machine debug` to reset the password each time you reboot the VCH endpoint VM.

For descriptions of the options that `vic-machine debug` includes in addition to the [Common `vic-machine` Options](common_vic_options.md) , see [VCH Debug Options](debug_vch_options.md).

**Prerequisites**

You have deployed at least one VCH.

**Procedure**

1. On the system on which you run `vic-machine`, navigate to the directory that contains the `vic-machine` utility.
2. Run the `vic-machine debug` command. 

   - You must specify the username and optionally the password, either in the `--target` option or separately in the `--user` and `--password` options. 
   - If your vSphere environment uses untrusted, self-signed certificates, you must also specify the thumbprint of the vCenter Server instance or ESXi host in the `--thumbprint` option. To obtain the thumbprint of the vCenter Server or ESXi host certificate, run `vic-machine` without the specifying the `--thumbprint` option. The operation fails, but the resulting error message includes the required certificate thumbprint. You can copy the thumbprint from the error message and run `vic-machine` again, including the `--thumbprint` option.
   - Specify the `--enable-ssh` and `--rootpw` options. Wrap the password in single quotes (Linux or Mac OS) or double quotes (Windows) if it includes special characters.
   - Optionally, specify the `--authorized-key` option to upload a public key file to `/root/.ssh/authorized_keys` folder in the endpoint VM. Include the name of the `*.pub` file in the path.

  <pre>$ vic-machine-<i>operating_system</i> debug
--target <i>esxi_host_address</i>
--user root
--password <i>esxi_host_password</i>
--thumbprint <i>certificate_thumbprint</i>
--enable-ssh
--rootpw '<i>new_p@ssword</i>'
--authorized-key <i>path_to_public_key_file</i>/<i>key_file</i>.pub
</pre>
  
**Result**

The output of the `vic-machine debug` command includes confirmation that SSH access is enabled:

<pre>### Configuring VCH for debug ####
[...]
SSH to appliance:
ssh root@<i>vch_address</i>
[...]
Completed successfully</pre>   