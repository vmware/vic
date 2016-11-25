# VCH Deployment Fails with Unknown or Non-Specified Argument Error or Incorrect User Name Error #
When you use the command line installer to deploy a virtual container host (VCH), the deployment fails with an error about unknown CLI arguments, unspecified mandatory options, or an invalid user name and password.

## Problem ##
Deployment fails during the validation of the configuration that you provided, even if you did specify the options cited as missing or incorrect. For example: 

<pre>Image datastore path must be specified; use format datastore/path</pre>
<pre>Unknown argument: <i>argument</i>
vic-machine failed: invalid CLI arguments</pre>
<pre>vic-machine failed: Failed to log in to <i>vcenter_server_or_esxi_host_address</i>: 
ServerFaultCode: Cannot complete login due to an incorrect user name or password</pre>

## Cause ##
String values that you provided for certain options contain spaces, or the user name and password contain special characters.

## Solution ##
Wrap any arguments that contain spaces or special characters in single quotation marks (') on Mac OS and Linux and in double quotation (") marks on Windows. 

Option arguments that might require quotation marks include the following:

- User names and passwords in `target`, or in `user` and `password`
- Datacenter names in `target`
- VCH names in `name`
- Datastore names and paths in `image-store`, `container-store`, and `volume-store`
- Network and distributed port group names in all networking options.
- Cluster and resource pool names in `compute-resource`
- Folder names in the paths for `cert`, `key`, `appliance-iso`, and `bootstrap-iso`

For information about when to use quotation marks for different options, see the descriptions of those options in [VCH Deployment Options](vch_installer_options.md).