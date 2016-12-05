# VCH Debug Options #

The command line utility for vSphere Integrated Containers Engine, `vic-machine`, provides a `debug` command that allows you to enable SSH access to the virtual container host (VCH) endpoint VM, set a password for the root user account, and upload a key file for automatic public key authentication. 

If you authorize SSH access to the VCH endpoint VM, you can edit system configuration files that you cannot edit by running `vic-machine` commands.

**NOTE**: Modifications that you make to the configuration of the VCH endpoint VM do not persist if you reboot the VM.

The `vic-machine debug` command includes the following options in addition to the common options described in [Common `vic-machine` Options](common_vic_options.md).

### `--enable-ssh` ###

Short name: `--ssh`

Enable an SSH server in the VCH endpoint VM. The `sshd` service runs until the VCH endpoint VM reboots. The `--enable-ssh` takes no arguments.

<pre>--enable-ssh</pre>

### `--rootpw` ###

Short name: `--pw`

Set a new password for the root user account on the VCH endpoint VM.

**IMPORTANT**: If you set a password for the VCH endpoint VM, this password does not persist if you reboot the VM. You must run vic-machine debug to reset the password each time that the VCH endpoint VM reboots.

Wrap the password in single quotes (Linux or Mac OS) or double quotes (Windows) if it includes shell characters such as `$`, `!` or `%`.

<pre>--rootpw '<i>new_p@ssword</i>'</pre>

### `--authorized-key` ###

Short name: `--key`

Upload a public key file to `/root/.ssh/authorized_keys` to enable SSH key authentication for the `root` user. Include the name of the `*.pub` file in the path.

<pre>--authorized-key <i>path_to_public_key_file</i>/<i>key_file</i>.pub</pre>