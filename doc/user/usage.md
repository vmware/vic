## Installing Virtual Integrated Containers

The intent is that vSphere Integrated Containers (VIC) should not _require_ an installation step - deploying a [Virtual Container Host](doc/design/arch/vic-container-abstraction.md#virtual-container-host) (VCH) directly without any prior steps should always be possible. At the current time this is the only approach available.

Installation will be required for capabilities such as [self-provisioning](doc/design/validating-proxy.md) and management network isolation via [vmomi proxy](doc/design/vmomi-authenticating-agent.md).

## Deploying a Virtual Container Host

### Requirements

- ESXi/vCenter - the target virtualization environment.
   - ESXi - Enterprise license
   - vCenter - Enterprise plus license, only very simple configurations have been tested. 
- DHCP - the VCH currently requires there be DHCP on all networks other than the bridge. The external network can be set via -external-network and will default to "VM Network" and the client and management roles will share external network NIC if not explicitly configured.
- Bridge network - when installed in a vCenter environment vic-machine does not automatically create a bridge network. An existing vSwitch or Distributed Portgroup should be specified via the -bridge-network flag, should not be the same as the external network, and should not have DHCP.

Replace the `<fields>` in the example with values specific to your environment - this will install VCH to the specified resource pool of ESXi or vCenter, and the container VMs will be created under that resource pool.

- -target is URL of the destination vSphere environment in the form `https://user:password@ip-or-fqdn/datacenter-name`. Any or all of the protocol, user, and password may also be omitted as can the datacenter if targetting ESXi or only one datacenter is configured.
- -compute-resource is the resource pool where VCH will be deployed to. For vCenter this should start with the cluster, followed by the resource pool path, e.g. `mycluster/a/resource/pool/path` (vCenter), or `/a/resource/pool` (ESXi) .
- -generate-cert flag is to generate certificates and configure TLS. 
- -force flag is to remove an existing datastore folder or VM with the same name.

```
vic-machine-linux --target target-host[/datacenter] --image-store <datastore name> --name <vch-name> --user root --password <password> --compute-resource <resource pool path> --generate-cert
```
This will, if successful, produce output similar to the following:
```
INFO[2016-06-19T05:16:09Z] ### Installing VCH ####
INFO[2016-06-19T05:16:09Z] Generating certificate/key pair - private key in ./vch-name-key.pem
INFO[2016-06-19T05:16:10Z] Validating supplied configuration
INFO[2016-06-19T05:16:10Z] Creating appliance on target
INFO[2016-06-19T05:16:10Z] Network role client is sharing NIC with external
INFO[2016-06-19T05:16:10Z] Network role management is sharing NIC with external
INFO[2016-06-19T05:16:11Z] Uploading images for container
INFO[2016-06-19T05:16:11Z]      bootstrap.iso
INFO[2016-06-19T05:16:11Z]      appliance.iso
INFO[2016-06-19T05:16:15Z] Waiting for IP information
INFO[2016-06-19T05:16:28Z] Waiting for major appliance components to launch
INFO[2016-06-19T05:16:28Z] Initialization of appliance successful
INFO[2016-06-19T05:16:28Z]
INFO[2016-06-19T05:16:28Z] SSH to appliance (default=root:password)
INFO[2016-06-19T05:16:28Z] ssh root@x.x.x.x
INFO[2016-06-19T05:16:28Z]
INFO[2016-06-19T05:16:28Z] Log server:
INFO[2016-06-19T05:16:28Z] https://x.x.x.x:2378
INFO[2016-06-19T05:16:28Z]
INFO[2016-06-19T05:16:28Z] DOCKER_HOST=x.x.x.x:2376
INFO[2016-06-19T05:16:28Z]
INFO[2016-06-19T05:16:28Z] Connect to docker:
INFO[2016-06-19T05:16:28Z] docker -H x.x.x.x:2376 --tls info
INFO[2016-06-19T05:16:28Z] Installer completed successfully          
```


## Deleting a Virtual Container Host

Specify the same resource pool and VCH name used to create a VCH, then the VCH will removed, toghether with the created containers, images, and volumes, if --force is provided. Here is an example command and output - replace the `<fields>` in the example with values specific to your environment.

```
vic-machine-linux delete --target target-host[/datacenter] --user root --password <password> --compute-resource <resource pool path> --name <vch-name>
INFO[2016-06-27T00:09:25Z] ### Removing VCH ####
INFO[2016-06-27T00:09:26Z] Removing VMs
INFO[2016-06-27T00:09:26Z] Removing images
INFO[2016-06-27T00:09:26Z] Removing volumes
INFO[2016-06-27T00:09:26Z] Removing appliance VM network devices
INFO[2016-06-27T00:09:27Z] Removing Portgroup XXX
INFO[2016-06-27T00:09:27Z] Removing VirtualSwitch XXX
INFO[2016-06-27T00:09:27Z] Removing Resource Pool XXX
INFO[2016-06-27T00:09:27Z] Completed successfully
```


[Issues relating to Virtual Container Host deployment](https://github.com/vmware/vic/labels/component%2Fvic-machine)
