## Installing Virtual Integrated Containers

The intent is that vSphere Integrated Containers (VIC) should not _require_ an installation step - deploying a [Virtual Container Host](../design/arch/vic-container-abstraction.md#virtual-container-host) (VCH) directly without any prior steps should always be possible. At the current time this is the only approach available.

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
vic-machine-linux create --target target-host[/datacenter] --image-store <datastore name> --name <vch-name> --user root --password <password> --compute-resource <resource pool path>
```
This will, if successful, produce output similar to the following when deploying VIC onto an ESXi:
```
INFO[2016-06-29T14:41:08-07:00] ### Installing VCH ####
INFO[2016-06-29T14:41:08-07:00] Generating certificate/key pair - private key in ./vic-001-key.pem
INFO[2016-06-29T14:41:09-07:00] Validating supplied configuration
INFO[2016-06-29T14:41:09-07:00] Firewall status: DISABLED on /ha-datacenter/host/localhost.localdomain/localhost.localdomain
INFO[2016-06-29T14:41:09-07:00] Firewall configuration OK on hosts:
INFO[2016-06-29T14:41:09-07:00]   /ha-datacenter/host/localhost.localdomain/localhost.localdomain
WARN[2016-06-29T14:41:09-07:00] Evaluation license detected. VIC may not function if evaluation expires or insufficient license is later assigned.
INFO[2016-06-29T14:41:09-07:00] License check OK
INFO[2016-06-29T14:41:09-07:00] DRS check SKIPPED - target is standalone host
INFO[2016-06-29T14:41:09-07:00] Creating Resource Pool vic-001
INFO[2016-06-29T14:41:09-07:00] Creating VirtualSwitch
INFO[2016-06-29T14:41:09-07:00] Creating Portgroup
INFO[2016-06-29T14:41:09-07:00] Creating appliance on target
INFO[2016-06-29T14:41:09-07:00] Network role external is sharing NIC with management
INFO[2016-06-29T14:41:09-07:00] Network role client is sharing NIC with management
INFO[2016-06-29T14:41:09-07:00] Uploading images for container
INFO[2016-06-29T14:41:09-07:00]   /home/user/go/src/github.com/vmware/vic/bin/appliance.iso
INFO[2016-06-29T14:41:09-07:00]   /home/user/go/src/github.com/vmware/vic/bin/bootstrap.iso
INFO[2016-06-29T14:41:15-07:00] Waiting for IP information
INFO[2016-06-29T14:41:35-07:00] Waiting for major appliance components to launch
INFO[2016-06-29T14:41:36-07:00] Initialization of appliance successful
INFO[2016-06-29T14:41:36-07:00]
INFO[2016-06-29T14:41:36-07:00] SSH to appliance (default=root:password)
INFO[2016-06-29T14:41:36-07:00] ssh root@192.168.218.129
INFO[2016-06-29T14:41:36-07:00]
INFO[2016-06-29T14:41:36-07:00] Log server:
INFO[2016-06-29T14:41:36-07:00] https://192.168.218.129:2378
INFO[2016-06-29T14:41:36-07:00]
INFO[2016-06-29T14:41:36-07:00] DOCKER_HOST=192.168.218.129:2376
INFO[2016-06-29T14:41:36-07:00]
INFO[2016-06-29T14:41:36-07:00] Connect to docker:
INFO[2016-06-29T14:41:36-07:00] docker -H 192.168.218.129:2376 --tls info
INFO[2016-06-29T14:41:36-07:00] Installer completed successfully
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


## Inspecting a Virtual Container Host

Specify the same resource pool and VCH name used to create a VCH, vic-machine inspect can show the VCH information.

```
vic-machine-linux inspect --target target-host[/datacenter] --user root --password <password> --compute-resource <resource pool path> --name <vch-name>
INFO[2016-06-29T16:03:17-05:00] ### Inspecting VCH ####
INFO[2016-06-29T16:03:20-05:00]
INFO[2016-06-29T16:03:20-05:00] VCH: <resource pool path>/<vch-name>
INFO[2016-06-29T16:03:21-05:00]
INFO[2016-06-29T16:03:21-05:00] SSH to appliance (default=root:password)
INFO[2016-06-29T16:03:21-05:00] ssh root@x.x.x.x
INFO[2016-06-29T16:03:21-05:00]
INFO[2016-06-29T16:03:21-05:00] Log server:
INFO[2016-06-29T16:03:21-05:00] https://x.x.x.x:2378
INFO[2016-06-29T16:03:21-05:00]
INFO[2016-06-29T16:03:21-05:00] DOCKER_HOST=x.x.x.x:2376
INFO[2016-06-29T16:03:21-05:00] DOCKER_OPTS="-H x.x.x.x:2376"
INFO[2016-06-29T16:03:21-05:00]
INFO[2016-06-29T16:03:21-05:00] Connect to docker:
INFO[2016-06-29T16:03:21-05:00] docker -H x.x.x.x:2376 info
INFO[2016-06-29T16:03:21-05:00] Completed successfully
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
