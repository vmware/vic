## Installing Virtual Integrated Containers

The intent is that vSphere Integrated Containers (VIC) should not _require_ an installation step - deploying a [Virtual Container Host](doc/design/arch/vic-container-abstraction.md#virtual-container-host) (VCH) directly without any prior steps should always be possible. At the current time this is the only approach available.

Installation will be required for capabilities such as [self-provisioning](doc/design/validating-proxy.md) and management network isolation via [vmomi proxy](doc/design/vmomi-authenticating-agent.md).

## Deploying a Virtual Container Host

### Requirements

The first three requirements derive from a placeholder install script. These will be addressed in #121.

1. bash - deploying a Virtual Container Host is currently done with an install script written in `bash`
2. govc - [These directions](https://github.com/vmware/govmomi/tree/master/govc#govc) instruct how to install `govc` - release 0.5.0
3. ESX - while this can and will function against vCenter, the placeholder installation script doesn't handle distributed port groups at this time
4. DHCP - the VCH currently requires there be DHCP on the external network (-e flag if not "VM Network")


Replace the `<fields>` in the example with values specific to your environment - this will install to the top-level resource pool of the host (specify -g to generate certificates and configure TLS). Add -f to remove an existing folder or VM with the same name:
```
bin/install.sh -g -t '<user>:<password>@<target-host>' -i <datastore-name> <vch-name>
```
This will, if successful, produce output similar to the following:
```
# Generating certificate/key pair - private key in vch-name-key.pem
# Logging into the target
# Uploading ISOs
[02-04-16 23:16:55] Uploading... OK
[02-04-16 23:16:58] Uploading... OK
# Creating vSwitch
# Creating Portgroup
# Creating the Virtual Container Host appliance
# Adding network interfaces
# Setting component configuration
# Configuring TLS server
# Powering on the Virtual Container Host
# Setting network identities
# Waiting for IP information
#
# SSH to appliance (default=root:password)
# root@x.x.x.x
#
# Log server:
# https://x.x.x.x:2378
#
# Connect to docker:
# docker -H x.x.x.x:2376 --tls --tlscert='vch-name-cert.pem' --tlskey='vch-name-key.pem'
DOCKER_OPTS="--tls --tlscert='vch-name-cert.pem' --tlskey='vch-name-key.pem'"
DOCKER_HOST=x.x.x.x:2376
```



[Issues relating to Virtual Container Host deployment](https://github.com/vmware/vic/labels/component%2Fvic-machine)
