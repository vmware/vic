# Debugging the VCH #

By default, all shell access to the virtual container host (VCH) endpoint VM is disabled. Login shells for all users are set to `/bin/false`. The command line utility for vSphere Integrated Containers Engine, `vic-machine`, provides a `debug` command that allows you to enable shell access to the virtual container host (VCH) endpoint VM, either by using the VM console or via SSH.

Do not confuse the `vic-machine debug` command with the `vic-machine create --debug` option. The `vic-machine debug` command allows you to log into and debug a VCH endpoint VM that you have already deployed. The `vic-machine create --debug` option deploys a new VCH that has increased levels of logging and other modifications, to allow you to debug the environment in which you deploy VCHs. For information about the `vic-machine create --debug` option, see the section on `--debug` in [VCH Deployment Options](../vic_installation/vch_installer_options.html#debug) in *vSphere Integrated Containers Engine Installation*.

* [Enable Shell Access to the VCH Endpoint VM](vch_shell_access.md)
* [Authorize SSH Access to the VCH Endpoint VM](vch_ssh_access.md) 
* [VCH Debug Options](debug_vch_options.md)