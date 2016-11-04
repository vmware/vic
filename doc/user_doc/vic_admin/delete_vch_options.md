#  Virtual Container Host Delete Options #

The command line utility for vSphere Integrated Containers Engine, `vic-machine`, provides a `delete` command that allows you to cleanly remove virtual container hosts. 

The `vic-machine delete` command includes one option in addition to the common options described in [Common `vic-machine` Options](common_vic_options.md).

### `--force` ###

Short name: `-f`

Forces `vic-machine delete` to ignore warnings and continue with the deletion of a virtual container host. Any running container VMs and any volume stores associated with the virtual container host are deleted. Errors such as an incorrect compute resource still cause the deletion to fail. 

- If you do not specify `--force` and the virtual container host contains running container VMs, the deletion fails with a warning. 
- If you do not specify `--force` and the virtual container host has volume stores, the deletion of the virtual container host succeeds without deleting the volume stores. The list of volume stores appears in the `vic-machine delete` success message for reference and optional manual removal.

If your vSphere environment uses untrusted, self-signed certificates, you can use the `--force` option to delete a virtual container host without providing the thumbprint of the vCenter Server or ESXi host in the `--thumbprint` option. 

<pre>--force</pre>