#  VCH Delete Options #

The command line utility for vSphere Integrated Containers Engine, `vic-machine`, provides a `delete` command that allows you to cleanly remove virtual container hosts (VCHs). 

The `vic-machine delete` command includes one option in addition to the common options described in [Common `vic-machine` Options](common_vic_options.md).

### `--force` ###

Short name: `-f`

Forces `vic-machine delete` to ignore warnings and continue with the deletion of a VCH. Any running container VMs and any volume stores associated with the VCH are deleted. Errors such as an incorrect compute resource still cause the deletion to fail. 

- If you do not specify `--force` and the VCH contains running container VMs, the deletion fails with a warning. 
- If you do not specify `--force` and the VCH has volume stores, the deletion of the VCH succeeds without deleting the volume stores. The list of volume stores appears in the `vic-machine delete` success message for reference and optional manual removal.

If your vSphere environment uses untrusted, self-signed certificates, you can use the `--force` option to delete a VCH without providing the thumbprint of the vCenter Server or ESXi host in the `--thumbprint` option. 

<pre>--force</pre>