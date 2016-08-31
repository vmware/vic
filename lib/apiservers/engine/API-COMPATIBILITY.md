# Docker Remote API Compatibility
VIC attempts to be compatible with the Docker remote API; however, there are some differences.  The API differences are documented below, grouped together by functional areas.

Here are the [Docker v1.22 docs](https://docs.docker.com/engine/reference/api/docker_remote_api_v1.22/) for the readers convenience. Only the calls that will have differences in implementation are documented below.

# Support Matrix

## Container API Support

|API|support?|
|---|---|
|List Containers|YES|
|Create a Container|YES|
|inspect a container| YES|
|List processes running inside a container|YES|
|Get Container Logs|YES|
|Inspect changes on a containers filesystem|Future|
|Export a container|Future|
|Get Container stats based on resource usage|Future|
|Resize a container TTY|YES|
|Start a Container|YES|
|Stop a container|YES|
|Restart a Container|YES|
|kill a container|YES|
|Update a Container|*Maybe*|
|Rename a Container|YES|
|Pause a Container|*maybe*|
|Unpause a Container|*maybe*|
|Attach a Container|YES|
|Attach a Container(websocket)|*maybe*|
|wait a container|YES|
|Remove a Container|YES|
|Copy files or folders from a Container|YES|
|Retrieve information about files and folders in a container|YES|
|Get and archive of a filesystem resource in a container|YES|
|Extract an archive of files or folders to a driectory in a container|YES|

## Image API Support

|API|support?|
|---|---|
|List Images|YES|
|Build image from a Dockerfile|Future|
|Create an Image|YES|
|Inspect and Image |YES|
|Get the History of an Image|YES|
|Push an image on the registry|Future|
|Tag an image into a repository|YES|
|Remove and Image|YES|
|Search Images|YES|

## Misc API Support

|API|support?|
|---|---|
|Check auth Configuration|YES|
|Display System-wide Information|YES|
|Show the docker version information|YES|
|Ping the Docker server|YES|
|Create a New Image from a containers changes|Future|
|Monitor Docker Events|YES|
|Get a tarball containing all images in a repository|Future|
|Get a tarball containing all images|Future|
|Load a Tarball with a set of images and tags into docker|YES|
|Exec Create |*maybe*|
|Exec Start |*maybe*|
|Exec Resize|*maybe*|
|Exec Inspect|*maybe*|


# Container APIs

## Create a container

```
POST /containers/create
```

### Json Fields

**misc**

|field|description|supported?|
|---|---|---|
|AttachStdin|Attach to Stdin|YES|
|AttachStdout|Attach to STDOUT|YES|
|AttachStderr|Attach to STDERR|YES|
|Cmd|Array of strings which specify commands to run|YES|
|Domainname|string to use for the containers domain|YES|
|Env|Set environment variables|YES|
|Entrypoint=""|Overwrite the default ENTRYPOINT of the image|YES|
|Hostname|Container host name|YES|
|Image|Image name to be used for the container|*maybe*|
|OpenStdin|Keep STDIN open even if not attached|YES|
|StdinOnce|close stdin once single attached client disconnects|YES|
|NetworkDisabled|boolean value that controls the presence of a network for the container.|YES|
|ExposedPorts|Object mapping to an empty object |YES|
|StopSignal|Signal to stop the container. can be a string or unsigned integer| *maybe*|
|Tty|Allocate a pseudo-TTY|*diff*|
|User|Username or UID|*diff*|
|WorkingDir|Working directory inside the container|YES|


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



**Host Config field**

|field|description|supported|
|---|---|---|
|field|description|Support|
|Binds|volume binds|*diff*|
|BlkioWeight|Block IO weight (relative weight) accepts a weight value between 10 and 1000.|NO|
|BlkioWeightDevice|Block IO weight (relative device weight) in the form of: "BlkioWeightDevice": [{"Path": "device_path", "Weight": weight}]|NO|
|"BlkioDeviceReadBps"|Limit read rate (bytes per second) from a device in the form of: "BlkioDeviceReadBps": [{"Path": "device_path", "Rate": rate}], for example: "BlkioDeviceReadBps": [{"Path": "/dev/sda", "Rate": "1024"}]"|NO|
|BlkioDeviceWriteBps|Limit write rate (bytes per second) to a device in the form of: "BlkioDeviceWriteBps": [{"Path": "device_path", "Rate": rate}], for example: "BlkioDeviceWriteBps": [{"Path": "/dev/sda", "Rate": "1024"}]"|NO|
|"BlkioDeviceReadIOps"|Limit read rate (IO per second) from a device in the form of: "BlkioDeviceReadIOps": [{"Path": "device_path", "Rate": rate}], for example: "BlkioDeviceReadIOps": [{"Path": "/dev/sda", "Rate": "1000"}]|NO|
|BlkioDeviceWiiteIOps|Limit write rate (IO per second) to a device in the form of: "BlkioDeviceWriteIOps": [{"Path": "device_path", "Rate": rate}], for example: "BlkioDeviceWriteIOps": [{"Path": "/dev/sda", "Rate": "1000"}]|NO|
|CapAdd|A list of kernel capabilities to add to the container.|NO|
|CapDrop|A list of kernel capabilities to drop from the container.|NO|
|ContainerIDFile|-|NO|
|CpusetCpus|CPUs in which to allow execution (e.g., 0-3, 0,1).|NO|
|CpusetMems|Memory nodes (MEMs) in which to allow execution (0-3, 0,1). Only effective on NUMA systems.|NO|
|CpuShares|CPU shares (relative weight).|NO|
|CpuPeriod|The length of a CPU period in microseconds.|NO|
|Devices|A list of devices to add to the container specified as a JSON object in the form { "PathOnHost": "/dev/deviceName", "PathInContainer": "/dev/deviceName", "CgroupPermissions": "mrw"}|NO|
|Dns|A list of DNS servers for the container to use.|YES|
|DnsOptions|A list of DNS options|NO|
|DnsSearch|A list of DNS search domains|YES|
|ExtraHosts|A list of hostnames/IP mappings to add to the container’s /etc/hosts file. Specified in the form ["hostname:IP"]|NO|
|IpcMode|-|NO|
|Links|A list of links for the container. Each link entry should be in the form of container_name:alias|NO|
|LxcConf|-|NO|
|Memory|Memory limit in bytes|NO|
|MemorySwap|Total memory limit (memory + swap); set -1 to enable unlimited swap. You must use this with memory and make the swap value larger than memory.|NO|
|MemoryReservation|Memory soft limit in bytes|NO|
|KernelMemory|Kernel memory limit in bytes.|NO|
|OomKillDisable|Boolean value, whether to disable OOM Killer for the container or not.|NO|
|OomScoreAdj|An integer value containing the score given to the container in order to tune OOM killer preferences.|NO|
|NetworkMode|Sets the networking mode for the container. Supported standard values are: bridge, host, none, and container:<name|id>. Any other value is taken as a custom network’s name to which this container should connect to.|*diff*|
|PortBindings|A map of exposed container ports and the host port they should map to. A JSON object in the form { <port>/<protocol>: [{ "HostPort": "<port>" }] } Take note that port is specified as a string and not an integer value.|NO|
|Privileged|Gives the container full access to the host. Specified as a boolean value.|NO, vms are the isolation unit in VIC and the commands inside of the container has access to the host's vm|
|ReadonlyRootfs|Mount the container’s root filesystem as read only. Specified as a boolean value.|NO|
|PublishAllPorts|Allocates a random host port for all of a container’s exposed ports. Specified as a boolean value.|NO|
|RestartPolicy|The behavior to apply when the container exits.|NO|
|LogConfig|Log configuration for the container, specified as a JSON object in the form { "Type": "<driver_name>", "Config": {"key1": "val1"}}. Available types: json-file, syslog, journald, gelf, awslogs, splunk, none. json-file logging driver.|YES|
|SecurityOpt|A list of string values to customize labels for MLS systems, such as SELinux.|NO|
|VolumesFrom|A list of volumes to inherit from another container. Specified in the form <container name>[:<ro|rw>]|NO|
|Ulimits|A list of ulimits to set in the container, specified as { "Name": <name>, "Soft": <soft limit>, "Hard": <hard limit> }, for example: Ulimits: { "Name": "nofile", "Soft": 1024, "Hard": 2048 }|NO|
|VolumeDriver|Driver that this container users to mount volumes.|NO, VIC has its own volume driver|
|ShmSize|Size of /dev/shm in bytes. The size must be greater than 0. If omitted the system uses 64MB|NO|

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

This is the structure of the `Mounts` field in the request.

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

###Response

The response for the create call is very simple in nature

it will return these two fields as json object in the body of an HTTP call.

* __Id__ : this is the id of the created containers.
* __Warnings__ : an array of strings that contain warning messages.

### Status Codes

* 201 : no error
* 404 : no such container
* 506 : impossible to attach
* 500 : server error

## Get container logs

```GET /containers/(id)/logs```

This operation will return StdOut and StdErr logs from the target container.

### Query parameters

|Parameter|Description|support|
|---|---|---|
| __follow__ |1/true/True or 0/false/False, defaults to false. This determines if a stream is returned. | YES |
| __stdout__ |1/true/True or 0/false/False, defaults to false. Indicates whether stdout logs are desired. | YES |
| __stderr__ | 1/true/True or 0/false/False, defaults to false. Indicates whether stderr logs are desired.| YES |
| __since__ | integer based UNIX timestamp with which to filter logs. Only logs since the timestamp will be returned. Default: 0. | *maybe* |
| __timestamps__ | 1/true/True or 0/false/False, defaults to false. timestamped log lines. Defaults to false.| *maybe* |
| __tail__ | all or <number>. indicates how many lines to output since the end of the logs. | *maybe* |

### Status Codes

* 101 : no error, indicates proxy about hijacking
* 200 : no error, no upgrade header from attach found
* 404 : no such container
* 500 : server error

## Export a container

```GET /containers/(id)/export```

This exports the contents of a container as a flattened file for other than docker use.
__VIC plans to support this in the future__

### Response

HTTP response that returns a binary stream of the flattened file in a tarball. The response is of type application/octet-stream.

### Response Status Codes

* 200 : no error
* 404 : no such container
* 500 : server error

## Attach to a container

```POST /containers/(id or name)/attach```

This call will return a stream of the targeted containers common io streams(stdout, stdin, stderr). It is important to note that this involves a `HTTP UPGRADE` response.

### Query Parameters

|parameter| description| supported |
|---|---|---|
|detachKeys|Override key for detaching a container. e.g. a-z| *maybe*|
|logs| 1/true/True or 0/false/False. if `true` returns logs defaults to `false`| *maybe*|
|stream| 1/true/True or 0/false/False. if `true` returns a stream defaults to `false`| YES|
|stdin| 1/true/True or 0/false/False. if `true` and `stream=1` returns stdin as a stream. defaults to `false`| YES|
|stdout| 1/true/True or 0/false/False. if `true` and `stream=1` returns stdout as a stream. defaults to `false`| YES|
|stderr| 1/true/True or 0/false/False. if `true` and `stream=1` returns stderr as a stream. defaults to `false`| YES|


###Status codes

* 101 : no error, indicates proxy about hijacking
* 200 : no error, no upgrade header found
* 400 : bad parameter
* 404 : no such container
* 500 : server error

###Stream Implementation

__TBD__

## Attach to a container (websocket)

```GET /containers/(id or name)/attach/ws```

Handshake according to `RFC 6455`

### Query Parameters

|parameter| description| supported |
|---|---|---|
|detachKeys|Override key for detaching a container. e.g. a-z| *maybe*|
|logs| 1/true/True or 0/false/False. if `true` returns logs defaults to `false`| *maybe*|
|stream| 1/true/True or 0/false/False. if `true` returns a stream defaults to `false`| YES|
|stdin| 1/true/True or 0/false/False. if `true` and `stream=1` returns stdin as a stream. defaults to `false`| YES|
|stdout| 1/true/True or 0/false/False. if `true` and `stream=1` returns stdout as a stream. defaults to `false`| YES|
|stderr| 1/true/True or 0/false/False. if `true` and `stream=1` returns stderr as a stream. defaults to `false`| YES|

###Status codes

* 101 : no error, indicates proxy about hijacking
* 200 : no error, no upgrade header found
* 400 : bad parameter
* 404 : no such container
* 500 : server error

## Remove a container

```Delete /containers/(id or name)```

removes indicated container.

###Query Parameter

|parameter|description|supported?|
|---|---|---|
|__v__| 1/true/True or 0/false/False. Removes associated volumes of container.| *maybe*|
|__force__| 1/true/True or 0/false/False. Perform a kill then remove operation | YES |

###Status Codes

* 204 : no error
* 400 : bad parameter
* 404 : no container
* 500 : server error

# Image APIs

## Build image from a Dockerfile

__Supported in the Future__

```POST /build```

build an image from a Dockerfile.

### Request

the request should probide a tar stream of the file to be used as the docker file.

### Query Parameters

|parameter| description | Supported? |
|---|---|---|
|dockerfile|Path within build context to the dockerfile.|YES|
| t | tag name for the build. `name:tag`. Default: `latest`. More than one can be provided. |YES|
|remote| GIT repository URI or HTTP/HTTPS URI. if the target is a file, the contents are put into a file called Dockerfile. | YES |
| q | suppress verbose build output| YES |
|nocache| Do not use build cache | *maybe* |
| pull | Attempt to pull latest image even if one is cached | YES |
|rm| remove intermediate containers after build | *maybe* |
|forcerm|always remove intermediate containers| *maybe* |
|memory| memory limit for build|  *YES but different*|
|memswap|total memory|NO|
|cpushares|CPU shares|NO|
|cpusetcpus|CPUs in which to allow execution| NO|
|cpuperiod|the length of a cpu period|NO|
|cpuquota|micreseconds of cpu time use in each cpu period|NO|
|buildargs|the key values in this json object are used as the build context for commands run using the `Run` docker command. |YES|
|shmsize| size of `/dev/shm` in bytes. Default 64MB|NO|

### Request Headers

|header| description|supported?|
|---|---|---|
|Content-type| application/tar | YES |
|x-Registry-Config| JSON config object that allows for providing mapped credentials to different registries that may be needed to build a specific image| YES|

### Status Codes

* 200 : no error
* 500 : server error

# Network APIs

### List networks
### Inspect network
### Create a network
### Connect a container to a network
### Disconnect a container from a network
### Remove a network


# Volume APIs

### List volumes
### Create a volume
### Inspect a volume
### Remove a volume
