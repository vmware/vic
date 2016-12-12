# VCH Deployment Fails with Missing Common Name Error Even When TLS Options Are Specified Correctly #
If you deploy a virtual container host (VCH) and you have specified one of the `vic-machine create --tls-cname`, `--no-tlsverify`, or `--no-tls` options, or you set a static IP address on the client network, the deployment fails with an error about the certificate Common Name being missing.

## Problem ##
Deployment fails during the validation of the configuration that you provided, even if you did specify a TLS option or you set a static IP address on the client network. For example: 

<pre>$ vic-machine-windows create 
--target 'Administrator@vsphere.local:<i>password</i>'@<i>vcenter_server</i> 
--bridge-network vic bridge --no-tls
### Installing VCH ####
[...]
Common Name must be provided when generating certificates for client
authentication:
[...]
Create cannot continue: unable to generate certificates
--------------------
vic-machine-windows.exe failed: provide Common Name for server certificate</pre>

If you include a TLS option at the beginning of the `vic-machine create` command rather than the end, you see the following error: 
<pre>$ vic-machine-windows create 
--target 'Administrator@vsphere.local:<i>password</i>'@<i>vcenter_server</i> 
 --no-tls
--bridge-network vic bridge
### Installing VCH ####
[...]
Unknown argument: bridge
--------------------
vic-machine-windows.exe failed: invalid CLI arguments
</pre>


## Cause ##
String values that you provided for certain options contain spaces or special characters that you did not escape with quotations marks. The `vic-machine create` input validator validates the arguments that you provide only as far as the argument that includes the space or special character. If you specify the TLS option before the argument with the space or special character, `vic-machine create` throws the correct error message. However, if you specify the TLS option after the argument that includes the space or special character, the `vic-machine create` validator stops before it reaches the TLS option, and so throws the error about the missing Common Name.

## Solution ##
Wrap any arguments that contain spaces or special characters in single quotation marks (') on Mac OS and Linux and in double quotation marks (")  on Windows. 

Option arguments that might require quotation marks include the following:

- User names and passwords in `--target`, or in `--user` and `--password`
- Datacenter names in `--target`
- VCH names in `--name`
- Datastore names and paths in `--image-store` and `--volume-store`
- Network and port group names in all networking options.
- Cluster and resource pool names in `--compute-resource`
- Folder names in the paths for `--cert-path`, `--cert`, `--key`, `--appliance-iso`, and `--bootstrap-iso`

For information about when to use quotation marks for different options, see the descriptions of those options in [VCH Deployment Options](vch_installer_options.md).