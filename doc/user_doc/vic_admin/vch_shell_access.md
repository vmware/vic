# Enable shell access to the VCH Endpoint VM #

You can use the `vic-machine debug` command to enable shell access to a VCH endpoint VM by setting a root password on the VM. Setting a root password enables access to the VCH endpoint VM via the VM console only. If you require SSH access to the VCH endoint VM, rather than just shell access, see [Authorize SSH Access to the VCH Endpoint VM](vch_ssh_access.md).  

**IMPORTANT**: Any changes that you make to a VCH by using `vic-machine debug` are non-persistent and are discarded if the VCH endpoint VM reboots.

For descriptions of the options that `vic-machine debug` includes in addition to the [Common `vic-machine` Options](common_vic_options.md) , see [VCH Debug Options](debug_vch_options.md).



**Prerequisites**

You deployed a VCH.

**Procedure**

1. On the system on which you run `vic-machine`, navigate to the directory that contains the `vic-machine` utility.
2. Run the `vic-machine debug` command. 
    - Specify the vSphere target and its credentials, either in the `--target` option or separately in the `--user` and `--password` options. 
      
       The credentials that you provide must have the following privilege on the endpoint VM:<pre>Virtual machine.Guest Operations.Guest Operation Program Execution</pre>
    - Specify the ID or name of the VCH to debug.
    - Potentially provide the thumbprint of the vCenter Server or ESXi host certificate.
    - Provide a password for the root user on the VCH endpoint VM by specifying the `--rootpw` option. Wrap the password in single quotes (Linux or Mac OS) or double quotes (Windows) if it includes shell characters such as `$`, `!` or `%`.

    <pre>$ vic-machine-<i>operating_system</i> debug
     --target <i>vcenter_server_or_esxi_host_address</i>
     --user <i>vcenter_server_or_esxi_host_username</i>
     --password <i>vcenter_server_or_esxi_host_password</i>
     --id <i>vch_id</i>
     --thumbprint <i>certificate_thumbprint</i>
     --rootpw '<i>new_p@ssword</i>' 
</pre>

**Result**

The output of the `vic-machine debug` command includes confirmation that SSH access is enabled:

<pre>### Configuring VCH for debug ####
[...]
SSH to appliance:
ssh root@<i>vch_address</i>
[...]
Completed successfully</pre>   
