# Authorize SSH Access to the VCH Endpoint VM #

You can use the `vic-machine debug` command to enable shell access to a VCH endpoint VM by setting a root password on the VM. Setting a root password enables access to the VCH endpoint VM via the VM console. You can also use `debug` to authorize SSH access to the VCH endpoint VM. You can optionally upload a key file for public key authentication when accessing the endpoint VM by using SSH. 

**IMPORTANT**: Any changes that you make to a VCH by using `vic-machine debug` are non-persistent and are discarded if the VCH endpoint VM reboots.

**Prerequisites**

You deployed a VCH.

**Procedure**

1. On the system on which you run `vic-machine`, navigate to the directory that contains the `vic-machine` utility.
2. Run the `vic-machine debug` command. 

    - Specify the vSphere target and its credentials, either in the `--target` option or separately in the `--user` and `--password` options. 
      
       The credentials that you provide must have the following privilege on the endpoint VM:<pre>Virtual machine.Guest Operations.Guest Operation Program Execution</pre>
   - Specify the ID or name of the VCH to debug.
   - Potentially provide the thumbprint of the vCenter Server or ESXi host certificate.
   - Specify the `--rootpw` option. Wrap the password in single quotes (Linux or Mac OS) or double quotes (Windows) if it includes shell characters such as `$`, `!` or `%`.
   - Authorize SSH access by specifying `--enable-ssh`.
   - Optionally, specify the `--authorized-key` option to upload a public key file to `/root/.ssh/authorized_keys` folder in the endpoint VM. Include the name of the `*.pub` file in the path.

    <pre>
     $ vic-machine-<i>operating_system</i> debug
     --target <i>vcenter_server_or_esxi_host_address</i>
     --user <i>vcenter_server_or_esxi_host_username</i>
     --password <i>vcenter_server_or_esxi_host_password</i>
     --id <i>vch_id</i>
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