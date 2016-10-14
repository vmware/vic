## Installing Virtual Integrated Containers Engine

The intent is that vSphere Integrated Containers Engine (VIC Engine) should not _require_ an installation step - deploying a [Virtual Container Host](../design/arch/vic-container-abstraction.md#virtual-container-host) (VCH) directly without any prior steps should always be possible. At the current time this is the only approach available.

## Deploying a Virtual Container Host

### Requirements

- ESXi/vCenter - the target virtualization environment.
   - ESXi - Enterprise license
   - vCenter - Enterprise plus license, only very simple configurations have been tested.
- DHCP - the VCH currently requires there be DHCP on all networks other than the bridge. The external network can be set via -external-network and will default to "VM Network" and the client and management roles will share external network NIC if not explicitly configured.
- Bridge network - when installed in a vCenter environment vic-machine does not automatically create a bridge network. An existing vSwitch or Distributed Portgroup should be specified via the -bridge-network flag, should not be the same as the external network, and should not have DHCP.

Replace the `<fields>` in the example with values specific to your environment - this will install VCH to the specified resource pool of ESXi or vCenter, and the container VMs will be created under that resource pool.

- --target is the URL of the destination vSphere environment in the form `https://user:password@ip-or-fqdn/datacenter-name`. Protocol, user, and password are OPTIONAL. Datacenter is OPTIONAL if targetting ESXi or only one datacenter is configured.
- --compute-resource is the resource pool where VCH will be deployed to. For vCenter this should start with the cluster, followed by the resource pool path, e.g. `mycluster/a/resource/pool/path` (vCenter), or `/a/resource/pool` (ESXi) .
- --force flag is to remove an existing datastore folder or VM with the same name.

```
vic-machine-linux create --target target-host[/datacenter] --image-store <datastore name> --name <vch-name> --user root --password <password> --compute-resource <resource pool path> --no-tlsverify
```
This will, if successful, produce output similar to the following when deploying VIC Engine onto an ESXi:
```
INFO[2016-10-08T23:37:34Z] Generating self-signed certificate/key pair - private key in ./XXX/key.pem
WARN[2016-10-08T23:37:34Z] Configuring without TLS verify - client authentication disabled
INFO[2016-10-08T23:37:34Z] ### Installing VCH ####
INFO[2016-10-08T23:37:34Z] Validating supplied configuration
INFO[2016-10-08T23:37:34Z] Firewall status: DISABLED on "/ha-datacenter/host/esx-a.localdomain/esx-a.localdomain"
WARN[2016-10-08T23:37:34Z] Firewall configuration will be incorrect if firewall is reenabled on hosts:
WARN[2016-10-08T23:37:34Z]   "/ha-datacenter/host/esx-a.localdomain/esx-a.localdomain"
WARN[2016-10-08T23:37:34Z] Firewall must permit 2377/tcp outbound if firewall is reenabled
INFO[2016-10-08T23:37:34Z] License check OK
INFO[2016-10-08T23:37:34Z] DRS check SKIPPED - target is standalone host
INFO[2016-10-08T23:37:34Z] Creating Resource Pool "XXX"
INFO[2016-10-08T23:37:35Z] Creating directory [datastore1] volumes
INFO[2016-10-08T23:37:35Z] Datastore path is [datastore1] volumes
INFO[2016-10-08T23:37:35Z] Creating appliance on target
INFO[2016-10-08T23:37:35Z] Network role "management" is sharing NIC with "client"
INFO[2016-10-08T23:37:35Z] Network role "external" is sharing NIC with "client"
INFO[2016-10-08T23:37:35Z] Uploading images for container
INFO[2016-10-08T23:37:35Z]      "bootstrap.iso"
INFO[2016-10-08T23:37:35Z]      "appliance.iso"
INFO[2016-10-08T23:37:39Z] Waiting for IP information
INFO[2016-10-08T23:37:51Z] Waiting for major appliance components to launch
INFO[2016-10-08T23:38:01Z] Initialization of appliance successful
INFO[2016-10-08T23:38:01Z]
INFO[2016-10-08T23:38:01Z] vic-admin portal:
INFO[2016-10-08T23:38:01Z] https://x.x.x.x:2378
INFO[2016-10-08T23:38:01Z]
INFO[2016-10-08T23:38:01Z] Docker environment variables:
INFO[2016-10-08T23:38:01Z]   DOCKER_HOST=x.x.x.x:2376
INFO[2016-10-08T23:38:01Z]
INFO[2016-10-08T23:38:01Z] Environment saved in XXX/XXX.env
INFO[2016-10-08T23:38:01Z]
INFO[2016-10-08T23:38:01Z] Connect to docker:
INFO[2016-10-08T23:38:01Z] docker -H x.x.x.x:2376 --tls info
INFO[2016-10-08T23:38:01Z] Installer completed successfully

```


## Deleting a Virtual Container Host

Specify the same resource pool and VCH name used to create a VCH, then the VCH will removed, together with the created containers, images, and volumes, if --force is provided. Here is an example command and output - replace the `<fields>` in the example with values specific to your environment.

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


## Inspecting a Virtual Container Host

Specify the same resource pool and VCH name used to create a VCH, vic-machine inspect can show the VCH information.

```
vic-machine-linux inspect --target target-host[/datacenter] --user root --password <password> --compute-resource <resource pool path> --name <vch-name>
INFO[2016-10-08T23:40:28Z] ### Inspecting VCH ####
INFO[2016-10-08T23:40:29Z]
INFO[2016-10-08T23:40:29Z] VCH ID: VirtualMachine:286
INFO[2016-10-08T23:40:29Z]
INFO[2016-10-08T23:40:29Z] Installer version: v0.6.0-0-0fac2c0
INFO[2016-10-08T23:40:29Z] VCH version: v0.6.0-0-0fac2c0
INFO[2016-10-08T23:40:29Z]
INFO[2016-10-08T23:40:29Z] VCH upgrade status:
INFO[2016-10-08T23:40:29Z] Installer has same version as VCH
INFO[2016-10-08T23:40:29Z] No upgrade available with this installer version
INFO[2016-10-08T23:40:29Z]
INFO[2016-10-08T23:40:29Z] vic-admin portal:
INFO[2016-10-08T23:40:29Z] https://x.x.x.x:2378
INFO[2016-10-08T23:40:29Z]
INFO[2016-10-08T23:40:29Z] Docker environment variables:
INFO[2016-10-08T23:40:29Z]   DOCKER_HOST=x.x.x.x:2376
INFO[2016-10-08T23:40:29Z]
INFO[2016-10-08T23:40:29Z]
INFO[2016-10-08T23:40:29Z] Connect to docker:
INFO[2016-10-08T23:40:29Z] docker -H x.x.x.x:2376 --tls info
INFO[2016-10-08T23:40:29Z] Completed successfully
```


## Enabling SSH to Virtual Container Host appliance
Specify the same resource pool and VCH name used to create a VCH, vic-machine debug will enable SSH on the appliance VM and then display the VCH information, now with SSH entry.

```
vic-machine-linux debug --target target-host[/datacenter] --user root --password <password> --compute-resource <resource pool path> --name <vch-name> --enable-ssh --rootpw <other password> --authorized-key <keyfile>
INFO[2016-10-08T23:41:16Z] ### Configuring VCH for debug ####
INFO[2016-10-08T23:41:16Z]
INFO[2016-10-08T23:41:16Z] VCH ID: VirtualMachine:286
INFO[2016-10-08T23:41:16Z]
INFO[2016-10-08T23:41:16Z] Installer version: v0.6.0-0-0fac2c0
INFO[2016-10-08T23:41:16Z] VCH version: v0.6.0-0-0fac2c0
INFO[2016-10-08T23:41:16Z]
INFO[2016-10-08T23:41:16Z] SSH to appliance
INFO[2016-10-08T23:41:16Z] ssh root@x.x.x.x
INFO[2016-10-08T23:41:16Z]
INFO[2016-10-08T23:41:16Z] vic-admin portal:
INFO[2016-10-08T23:41:16Z] https://x.x.x.x:2378
INFO[2016-10-08T23:41:16Z]
INFO[2016-10-08T23:41:16Z] Docker environment variables:
INFO[2016-10-08T23:41:16Z]   DOCKER_HOST=x.x.x.x:2376
INFO[2016-10-08T23:41:16Z]
INFO[2016-10-08T23:41:16Z]
INFO[2016-10-08T23:41:16Z] Connect to docker:
INFO[2016-10-08T23:41:16Z] docker -H x.x.x.x:2376 --tls info
INFO[2016-10-08T23:41:16Z] Completed successfully
```

## List Virtual Container Hosts

vic-machine ls can list all VCHs in your VC/ESXi, or list all VCHs under the provided resource pool by compute-resource parameter.
```
vic-machine-linux ls --target target-host --user root --password <password>
INFO[2016-08-08T16:21:57-05:00] ### Listing VCHs ####

ID                           PATH                                                   NAME
VirtualMachine:vm-189        /dc1/host/cluster1/Resources/test1/test1-2        test1-2-1
VirtualMachine:vm-189        /dc2/host/cluster2/Resources/test2/test2-2        test2-2-1


vic-machine-linux ls --target target-host/dc1 --user root --password <password> --compute-resource cluster1/test1
INFO[2016-08-08T16:25:50-02:00] ### Listing VCHs ####

ID                           PATH                                                   NAME
VirtualMachine:vm-189        /dc1/host/cluster1/Resources/test1/test1-2        test1-2-1
```


## Configuring Volumes in a Virtual Container Host

Volumes are implemented as VMDKs and mounted as block devices on a containerVM. This means that they cannot be used concurrently by multiple running, containers. Attempting to start a container that has a volume attached that is in use by a running container will result in an error.

The location in which volumes are created can be specified at creation time via the `--volume-store` argument. This can be supplied multiple times to configure multiple datastores or paths:
```
vic-machine-linux create --volume-store=datastore1/some/path:default --volume-store=ssdDatastore/other/path:fast ...
```

The volume store to use is specified via driver optons when creating volumes (capacity is in MB):
```
docker volume create --name=reports --opts VolumeStore=fast --opt Capacity=1024
```

Providing a volume store named `default` allows the driver options to be omitted in the example above and enables anoymous volumes, e.g.:
```
docker run -v /var/lib/data -it busybox 
```
  

## Exposing vSphere networks within a Virtual Container Host

vSphere networks can be directly mapped into the VCH for use by containers. This allows a container to expose services to the wider world without using port-forwarding (which is not yet implemented):

```
vic-machine-linux create --container-network=vpshere-network:descriptive-name
```

A container can then be attached directly to this network and a bridge network via:
```
docker create --net=descriptive-name haproxy
docker network connect bridge <container-id>
```

Currently the container does **not** have a firewall configured [#692](https://github.com/vmware/vic/issues/692) in this circumstance.

[Issues relating to Virtual Container Host deployment](https://github.com/vmware/vic/labels/component%2Fvic-machine)
