#  VCH Upgrade Options #

The command line utility for vSphere Integrated Containers Engine, `vic-machine`, provides an `upgrade` command that allows you to upgrade virtual container hosts (VCHs) to a newer version. 

The `vic-machine upgrade` command includes the following options in addition to the common options described in [Common `vic-machine` Options](common_vic_options.md).

### `--appliance-iso` ###

Short name: `--ai`

The path to the new version of the ISO image from which to upgrade the VCH appliance. Set this option if you have moved the `appliance.iso` file to a folder that is not the folder that contains the `vic-machine` binary or is not the folder from which you are running `vic-machine`. Include the name of the ISO file in the path.

**NOTE**: Do not use the `--appliance-iso` option to point `vic-machine` to an `--appliance-iso` file that is of a different version to the version of `vic-machine` that you are running.

<pre>--appliance-iso <i>path_to_ISO_file</i>/<i>ISO_file_name</i>.iso</pre>

Wrap the folder names in the path in single quotes (Linux or Mac OS) or double quotes (Windows) if they include spaces.

<pre>--appliance-iso '<i>path to ISO file</i>'/appliance.iso</pre>

### `--bootstrap-iso` ###

Short name: `--bi`

The path to the new version of the ISO image from which to upgrade the container VMs that the VCH manages. Set this option if you have moved the `bootstrap.iso` file to a folder that is not the folder that contains the `vic-machine` binary or is not the folder from which you are running `vic-machine`. Include the name of the ISO file in the path.

**NOTE**: Do not use the `--bootstrap-iso` option to point `vic-machine` to a `--bootstrap-iso` file that is of a different version to the version of `vic-machine` that you are running.

<pre>--bootstrap-iso <i>path_to_ISO_file</i>/bootstrap.iso</pre>

Wrap the folder names in the path in single quotes (Linux or Mac OS) or double quotes (Windows) if they include spaces.

<pre>--bootstrap-iso '<i>path to ISO file</i>'/<i>ISO_file_name</i>.iso</pre>

### `--force` ###

Short name: `-f`

Forces `vic-machine upgrade` to ignore warnings and continue with the upgrade of a VCH. Errors such as an incorrect compute resource still cause the upgrade to fail. 

If your vSphere environment uses untrusted, self-signed certificates, you can use the `--force` option to upgrade a VCH without providing the thumbprint of the vCenter Server or ESXi host in the `thumbprint` option. 

<pre>--force</pre>
