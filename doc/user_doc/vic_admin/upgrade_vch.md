# Upgrade a VCH #

You upgrade virtual container hosts (VCHs) by downloading a new version of vSphere Integrated Containers Engine and running the `vic-machine upgrade` command.

You can use `vic-machine upgrade` to upgrade VCHs from version 0.8 and above. You can run `vic-machine upgrade` on VCHs that are either running or powered off. When you upgrade a running VCH, the VCH goes temporarily offline, but container workloads continue as normal during the upgrade process. Upgrading a VCH does not affect any mapped container networks that you defined by setting the `vic-machine create --container-network` option. The following operations are not available during upgrade:

- You cannot access container logs
- You cannot attach to a container
- NAT based port forwarding is unavailable

**IMPORTANT**: Upgrading a VCH does not upgrade any existing container VMs that the VCH manages. For container VMs to boot from the latest version of `bootstrap.iso`, container developers must recreate them.

For descriptions of the options that `vic-machine upgrade` includes in addition to the [Common `vic-machine` Options](common_vic_options.md) , see [VCH Upgrade Options](upgrade_vch_options.md).

**Prerequisites**

- You deployed one or more VCHs with an older version of the `vic-machine create` command.
- You downloaded a new version of the vSphere Integrated Containers Engine bundle.
- Run the `vic-machine ls` command by using the new version of `vic-machine` to see the upgrade status of all of the VCHs that are running on a vCenter Server instance or ESXi host. For information about running `vic-machine ls`, see [List VCHs and Obtain Their IDs](list_vch.md).
- Optionally note the IDs of the VCHs.

**Procedure**

1. On the system on which you run `vic-machine`, navigate to the directory that contains the new version of the `vic-machine` utility.
2. Run the `vic-machine upgrade` command. 

     The following example includes the options required to upgrade a VCH in a simple vCenter Server environment. 

  - You must specify the username and optionally the password, either in the `target` option or separately in the `--user` and `--password` options. 
  - If the VCH has a name other than the default name, `virtual-container-host`, you must specify the `--name` or `--id` option. 
  - If multiple compute resources exist in the datacenter, you must specify the `--compute-resource` or `--id` option. 
  - If your vSphere environment uses untrusted, self-signed certificates, you must also specify the thumbprint of the vCenter Server instance or ESXi host in the `--thumbprint` option. To obtain the thumbprint of the vCenter Server or ESXi host certificate, run `vic-machine` without the specifying the `--thumbprint` or `--force` options. The upgrade of the VCH fails, but the resulting error message includes the required certificate thumbprint. You can copy the thumbprint from the error message and run `vic-machine` again, including the `--thumbprint` option.

     <pre>$ vic-machine-<i>operating_system</i> upgrade
--target <i>vcenter_server_username</i>:<i>password</i>@<i>vcenter_server_address</i>
--thumbprint <i>certificate_thumbprint</i>
--name <i>vch_name</i></pre>

3. If the upgrade operation fails with error messages, run `vic-machine upgrade` again with the `--force` option.

    If your vSphere environment uses untrusted, self-signed certificates, running `vic-machine upgrade` with the `--force` option allows you to omit the `--thumbprint` option.

     <pre>$ vic-machine-<i>operating_system</i> upgrade
--target <i>vcenter_server_username</i>:<i>password</i>@<i>vcenter_server_address</i>
--name <i>cluster_name</i></i>
--force</pre>

**Result**

During the upgrade process, `vic-machine upgrade` performs the following operations:

- Validates whether the configuration of the existing VCH is compatible with the new version. If not, the upgrade fails. 
- Uploads the new versions of the `appliance.iso` and `bootstrap.iso` files to the VCH. There is no timeout for this stage of the upgrade process, so that the ISO files can upload over slow connections.
- Creates a snapshot of the VCH endpoint VM, to use in case the upgrade fails and has to roll back.
- Boots the VCH by using the new version of the `appliance.iso` file.
- Deletes the snapshot of the VCH endpoint VM once the upgrade has succeeded.
- After you upgrade a VCH, any new container VMs will boot from the new version of the `bootstrap.iso` file.
- If the upgrade times out while waiting for the VCH service to start, the upgrade fails and rolls back to the previous version.

**What to Do Next**

[Upgrade the HTML5 vSphere Client Plug-In](upgrade_h5_plugin.md).
