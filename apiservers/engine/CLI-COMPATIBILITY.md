# Docker CLI Compatibility
VIC attempts to be compatible with the Docker client (CLI); however, there are some differences. The CLI differences are documented below, alphabetically. Note, the operations listed as supported but with differences does not mean incompatible. It means most of these operations support either a smaller subset of the parameters or the request/response parameters have different meanings that apply to running in an ESX or vCenter environment.

The list is organized into 3 separate tables. The first table lists CLI operations that are supported in the first version of VIC. The second table lists operations that VIC may support in the future.

### VIC Support

|Docker CLI Operation | support viable | technical priority |
|--- | ---| ---|
|attach | YES| P0 - feel good factor and interactive debug |
|cp | YES| P1 - live copy |
|create | *Yes - differences*| P0 - core function: use |
|events | *Yes - differences*| P1 - necessary for automation |
|history | YES| - |
|images | YES| P0 - core function: admin |
|login | YES| - |
|logout | YES| - |
|logs | YES| P0 - core function: use |
|network connect| YES| P2 - necessary only for dynamic topology changes |
|network create| *Yes - differences*| P0 - core function: use/admin |
|network disconnect| YES| P2 - symmetrical with connect |
|network ls| YES| P0 - core function: use |
|network rm| YES| - |
|port | YES| P1 - core function: use - necessity reduced if not using NAT model |
|ps | YES| P0 - core function: use/admin |
|pull | YES| P0 - core function: use |
|rm | YES| P0 - core function: use/admin |
|rmi | YES| P1 - core function: admin - maybe less critical in production |
|run | *Yes - differences*| P0 - core function: use |
|rename | YES| - |
|restart | YES| P0 - core function: use |
|start | YES| P0 - core function: use |
|stop | YES| P0 - core function: use |
|tag| YES| - |
|version | YES| P0 - core function: admin |
|volume create| *Yes - differences*| P0 - core function: use/admin |
|volume inspect| YES| P1 - core function: admin - maybe less critical in production | 
|volume ls| YES| P0 - core function: use |
|volume rm| YES| P1 - core function: use/admin - depends on dynamism of volume use | 
|exec | YES | P1 - interactive debug |
|pause | YES| - |
|stats | YES| P1 - core function: diagnostics |
|top | YES| - |
|unpause | YES| - |
|wait | YES | P1 - necessary for automation |

The technical priority ranking is as follows (assessed with regard to production deployment):
P0 - delay GA if not available
P1 - significant reduction in capability - consider slipping
P2 - confined reduction in capability
-  - not significant in initial use cases

### Future Support

|Docker CLI Operation |
|--- |
|build |
|commit |
|diff |
|export |
|import |
|load |
|push |
|save |
|search |

The following sections dive further into these differences.

## create

```
Usage: docker create [OPTIONS] IMAGE [COMMAND] [ARG...]

Create a new container
```

| params | Docker | VIC |
|---|---|---|
|-a, --attach=[]|Attach to STDIN, STDOUT or STDERR|YES|
|--add-host=[] |Add a custom host-to-IP mapping (host:ip)|YES|
|-c, --cpu-shares=0 | CPU shares (relative weight)|*diff*|
|--cap-add=[]|Add Linux capabilities|NO, VIC containers do not run in cgroups|
|--cap-drop=[]|Drop Linux capabilities|NO, VIC containers do not run in cgroups|
|--cgroup-parent=""|Optional parent cgroup for the container|NO, VIC containers do not run in cgroups|
|--cidfile=""|Write the container ID to the file|YES|
|--cpuset-cpus=""|CPUs in which to allow execution (0-3, 0,1)|NO|
|--device=[]|Add a host device to the container|*diff*|
|--dns=[]|Set custom DNS servers|YES|
|--dns-search=[]|Set custom DNS search domains|YES|
|-e, --env=[]|Set environment variables|YES|
|--entrypoint=""|Overwrite the default ENTRYPOINT of the image|YES|
|--env-file=[]|Read in a file of environment variables|YES|
|--expose=[]|Expose a port or a range of ports|YES|
|-h, --hostname=""|Container host name|YES|
|-i, --interactive=false|Keep STDIN open even if not attached|*maybe*|
|--ipc=""|IPC namespace to use|NO, VIC containers run in separate VMs|
|-l, --label=[]|Set metadata on the container (e.g., --label=com.example.key=value)|YES|
|--label-file=[]|Read in a line delimited file of labels|YES|
|--link=[]|Add link to another container|*maybe*|
|--log-driver=""|Logging driver for container|YES|
|--lxc-conf=[]|Add custom lxc options|NO, VIC containers do not run in cgroups|
|-m, --memory=""|Memory limit|*diff*|
|--mac-address=""|Container MAC address (e.g. 92:d0:c6:0a:29:33)|YES|
|--name=""|Assign a name to the container|YES|
|--net="bridge"|Set the Network mode for the container|YES|
|-P, --publish-all=false|Publish all exposed ports to random ports|*maybe*|
|-p, --publish=[]|Publish a container's port(s) to the host|YES|
|--privileged=false|Give extended privileges to this container|*diff*|
|--read-only=false|Mount the container's root filesystem as read only|*diff*|
|--restart="no"|Restart policy (no, on-failure[:max-retry], always)|*maybe*|
|--security-opt=[]|Security options|*maybe*|
|-t, --tty=false|Allocate a pseudo-TTY|*diff*|
|-u, --user=""|Username or UID|*diff*|
|-v, --volume=[]|Bind mount a volume|*diff*|
|--volumes-from=[]|Mount volumes from the specified container(s)|*maybe*|
|-w, --workdir=""|Working directory inside the container|YES|


## events

Events **may** may also report ESX and vCenter events.


## info

The information provided by *info* appears similar to those reported by Docker but many of them are specific to vSphere.

## inspect

Some parameters **returned** by *inspect* are specific to vSphere.  Others are not relevant.

### Response:

The response for inspect is a JSON document composed of several inner documents.  The following are the primary subdocuments that differs.  For the rest of the response, please visit,
```
https://docs.docker.com/engine/reference/api/docker_remote_api_v1.22/#inspect-a-container.
```

**AppArmorProfile**

A VIC container isolation unit is a VM. Processes running in a VIC container do not use AppArmor profiles.

**HostConfig**

|Param|Docker|VIC|
|---|---|---|
|"Binds"|volume binds|*diff*|
|"BlkioWeight|Block IO weight (relative weight) accepts a weight value between 10 and 1000.|NO|
|"BlkioWeightDevice"|Block IO weight (relative device weight) in the form of: "BlkioWeightDevice": [{"Path": "device_path", "Weight": weight}]|NO|
|"BlkioDeviceReadBps"|Limit read rate (bytes per second) from a device in the form of: "BlkioDeviceReadBps": [{"Path": "device_path", "Rate": rate}], for example: "BlkioDeviceReadBps": [{"Path": "/dev/sda", "Rate": "1024"}]"|NO|
|"BlkioDeviceWriteBps"|Limit write rate (bytes per second) to a device in the form of: "BlkioDeviceWriteBps": [{"Path": "device_path", "Rate": rate}], for example: "BlkioDeviceWriteBps": [{"Path": "/dev/sda", "Rate": "1024"}]"|NO|
|"BlkioDeviceReadIOps"|Limit read rate (IO per second) from a device in the form of: "BlkioDeviceReadIOps": [{"Path": "device_path", "Rate": rate}], for example: "BlkioDeviceReadIOps": [{"Path": "/dev/sda", "Rate": "1000"}]|NO|
|"BlkioDeviceWiiteIOps"|Limit write rate (IO per second) to a device in the form of: "BlkioDeviceWriteIOps": [{"Path": "device_path", "Rate": rate}], for example: "BlkioDeviceWriteIOps": [{"Path": "/dev/sda", "Rate": "1000"}]|NO|
|"CapAdd"|A list of kernel capabilities to add to the container.|NO|
|"CapDrop"|A list of kernel capabilities to drop from the container.|NO|
|"ContainerIDFile"|-|NO|
|"CpusetCpus"|CPUs in which to allow execution (e.g., 0-3, 0,1).|NO|
|"CpusetMems"|Memory nodes (MEMs) in which to allow execution (0-3, 0,1). Only effective on NUMA systems.|NO|
|"CpuShares"|CPU shares (relative weight).|NO|
|"CpuPeriod"|The length of a CPU period in microseconds.|NO|
|"Devices"|A list of devices to add to the container specified as a JSON object in the form { "PathOnHost": "/dev/deviceName", "PathInContainer": "/dev/deviceName", "CgroupPermissions": "mrw"}|NO|
|"Dns"|A list of DNS servers for the container to use.|YES|
|"DnsOptions"|A list of DNS options|NO|
|"DnsSearch"|A list of DNS search domains|YES|
|"ExtraHosts"|A list of hostnames/IP mappings to add to the container’s /etc/hosts file. Specified in the form ["hostname:IP"]|NO|
|"IpcMode"|-|NO|
|"Links"|A list of links for the container. Each link entry should be in the form of container_name:alias|NO|
|"LxcConf"|-|NO|
|"Memory"|Memory limit in bytes|NO|
|"MemorySwap"|Total memory limit (memory + swap); set -1 to enable unlimited swap. You must use this with memory and make the swap value larger than memory.|NO|
|"MemoryReservation"|Memory soft limit in bytes|NO|
|"KernelMemory"|Kernel memory limit in bytes.|NO|
|"OomKillDisable"|Boolean value, whether to disable OOM Killer for the container or not.|NO|
|"OomScoreAdj"|An integer value containing the score given to the container in order to tune OOM killer preferences.|NO|
|"NetworkMode"|Sets the networking mode for the container. Supported standard values are: bridge, host, none, and container:name&#124;id. Any other value is taken as a custom network’s name to which this container should connect to.|*diff*|
|"PortBindings"|A map of exposed container ports and the host port they should map to. A JSON object in the form { <port>/<protocol>: [{ "HostPort": "<port>" }] } Take note that port is specified as a string and not an integer value.|NO|
|"Privileged"|Gives the container full access to the host. Specified as a boolean value.|NO, vms are the isolation unit in VIC and the commands inside of the container has access to the host's vm|
|"ReadonlyRootfs"|Mount the container’s root filesystem as read only. Specified as a boolean value.|NO|
|"PublishAllPorts"|Allocates a random host port for all of a container’s exposed ports. Specified as a boolean value.|NO|
|"RestartPolicy"|The behavior to apply when the container exits.|NO|
|"LogConfig"|Log configuration for the container, specified as a JSON object in the form { "Type": "<driver_name>", "Config": {"key1": "val1"}}. Available types: json-file, syslog, journald, gelf, awslogs, splunk, none. json-file logging driver.|YES|
|"SecurityOpt"|A list of string values to customize labels for MLS systems, such as SELinux.|NO|
|"VolumesFrom"|A list of volumes to inherit from another container. Specified in the form <container name>[:<ro|rw>]|NO|
|"Ulimits"|A list of ulimits to set in the container, specified as { "Name": <name>, "Soft": <soft limit>, "Hard": <hard limit> }, for example: Ulimits: { "Name": "nofile", "Soft": 1024, "Hard": 2048 }|NO|
|"VolumeDriver"|Driver that this container users to mount volumes.|NO, VIC has its own volume driver|
|"ShmSize"|Size of /dev/shm in bytes. The size must be greater than 0. If omitted the system uses 64MB|NO|

**misc params**

|Param|Docker|VIC|
|---|:---|---|
|HostnamePath|e.g. "/var/lib/docker/containers/ba033ac4401106a3b513bc9d639eee123ad78ca3616b921167cd74b20e25ed39/hostname"|This has no meaning in VIC.|
|HostsPath|e.g. "/var/lib/docker/containers/ba033ac4401106a3b513bc9d639eee123ad78ca3616b921167cd74b20e25ed39/hosts"|This has no meaning in VIC|
|LogPath|e.g. "/var/lib/docker/containers/1eb5fabf5a03807136561b3c00adcd2992b535d624d5e18b6cdc6a6844d9767b/1eb5fabf5a03807136561b3c00adcd2992b535d624d5e18b6cdc6a6844d9767b-json.log"|This has no meaning in VIC|
|"Id"|Image ID, Sha256 content addressable hash|YES|
|"Image"|-|NO?|
|"MountLabel"|-|NO|
|"Name"|Docker's human readable name|YES|
|"Path"|-|NO|
|"ProcessLabel"|-|NO|
|"ResolvConfPath"|-|NO|
|"RestartCount"|-|NO|

**NetworkSettings**

|Param|Docker|VIC|
|---|---|---|
|"Bridge"|""|?|
|"SandboxID"|""|?|
|"HairpinMode"|""|?|
|"LinkLocalIPv6Address"|""|?|
|"LinkLocalIPv6PrefixLen"|""|?|
|"Ports"|""|?|
|"SandboxKey"|""|?|
|"SecondaryIPAddresses"|""|?|
|"SecondaryIPv6Addresses"|""|?|
|"EndpointID"|""|?|
|"Gateway"|""|?|
|"GlobalIPv6Address"|""|?|
|"GlobalIPv6PrefixLen"|""|?|
|"IPAddress"|""|?|
|"IPPrefixLen"|""|?|
|"IPv6Gateway"|""|?|
|"MacAddress"|""|?|
|"Networks"|-|VIC will have some form of this data|

**State**

|Param|Docker|VIC|
|---|---|---|
|"Error"|-|NO|
|"ExitCode"|-|YES|
|"FinishedAt"|-|YES|
|"OOMKilled"|-|?|
|"Dead"|-|YES|
|"Paused"|-|YES|
|"Pid"|-|NO|
|"Restarting"|-|?|
|"Running"|-|?|
|"StartedAt"|-|YES|
|"Status"|-|YES|

**Mounts**

An array of mount points.  VIC will provide some form of this.  Here is an example of what Docker returns.

```
    "Mounts": [
        {
            "Name": "fac362...80535",
            "Source": "/data",
            "Destination": "/data",
            "Driver": "local",
            "Mode": "ro,Z",
            "RW": false,
            "Propagation": ""
        }
    ]
```

## network create

The network creation workflow will differ from standard Docker.  The -driver option is ignored.  All networks in VIC are of type *bridged*.  Policies on these networks are administered with vSphere.  This makes exposing containers to the external world much simpler.  In the future, support for VMWare's NSX will be incorporated.  For more information on VIC's networking model, please visit, https://github.com/vmware/vic/blob/master/doc/design/networking/MVPnetworking.md

## volume create

Much like network creation, volume creation also leverages the vSphere infrastructure.  Volumes are just VMDKs created within a VCH.  This gives VIC containers a lot of advantages.  A container can run on any ESX server that is part of the VCH and it's data volume can be on another ESX server.  The container and the data volume can be moved from one ESX server to another within the VCH, and they would be able to connect seamlessly.

Shared volumes between containers will eventually be supported.  Stay tune.

VIC can also leverage a vSan infrastructure.

## run

Docker run is a composite operation for pull, create, start.  As such, the parameter support is equivalent to that of docker create.
