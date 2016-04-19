# Docker Remote API Compatibility
VIC attempts to be compatible with the Docker remote API; however, there are some differences.  The API differences are documented below, grouped together by functional areas.

# support Matrix

## Container API Support

|API|support?|
|---|---|
|List Containers|YES|
|Create a Container|YES|
|inspect a container| YES|
|List processes running inside a container|YES|
|Get Container Logs|YES|
|Inspect changes on a containers filesystem|YES|
|Export a container|YES - In the future|
|Get Container stats based on resource usage|__NO__|
|Resize a container TTY|__NO__|
|Start a Container|YES|
|Stop a container|YES|
|Restart a Container|YES|
|kill a container|YES|
|Update a Container|*Maybe*|
|Rename a Container|YES|
|Pause a Container|*maybe*|
|Unpause a Container|*maybe*|
|Attach a Container|YES|
|Attach a Container(websocket)|YES|
|wait a container|*maybe*|
|Remove a Container|YES|
|Copy files or folders from a Container|YES|
|Retrieve information about files and folders in a container|YES|
|Get and archive of a filesystem resource in a container|YES|
|Extract an archive of files or folders to a driectory in a container|YES|

## Image API Support

|API|support?|
|---|---|
|List Images|YES|
|Build image from a Dockerfile|YES - In the Future|
|Create an Image|YES|
|Inspect and Image |YES|
|Get the History of an Image|YES|
|Push an image on the registry|YES - In the Future|
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
|Create a New Image from a containers changes|YES - in the future|
|Monitor Docker Events|YES|
|Get a tarball containing all images in a repository|YES - in the future|
|Get a tarball containing all images|YES - in the future|
|Load a Tarball with a set of images and tags into docker|YES|
|Exec Create |*maybe*|
|Exec Start |*maybe*|
|Exec Resize|*maybe*|
|Exec Inspect|*maybe*|


# Container APIs
>>>>>>> origin/docker-support-docs

### List containers
```
GET /containers/json

```
No differences


## Create a container

```
POST /containers/create
```

### Json Fields

**misc**
|field|description|supported|
|---|---|---|
| AttachStdin | Attach to Stdin | YES |
| AttachStdout | Attach to STDOUT | YES |
| AttachStderr | Attach to STDERR | YES |
|Cmd|Array of strings which specify commands to run|*maybe*|
|Domainname| string to use for the containers domain | YES |
|Env|Set environment variables|YES|
|Entrypoint=""|Overwrite the default ENTRYPOINT of the image|YES|
| Hostname |Container host name|YES|
|Image|Image name to be used for the container| *maybe* |
| OpenStdin | Keep STDIN open even if not attached | *maybe* |
| StdinOnce| close stdin once single attached client disconnects | *maybe* |
|NetworkDisabled|boolean value that controls the presence of a network for the container.|YES|
|ExposedPorts|Object mapping to an empty object |YES|
|StopSignal|Signal to stop the container. can be a string or unsigned integer| *maybe*|
| Tty |Allocate a pseudo-TTY| *diff* |
| User | Username or UID | *diff* |
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

## Inspect a container

``` GET /containers/(id)/json ```

This remote API returns low level metadata for the container. This includes almost all of what is previded at the time a create call
is invoked. This is a __supported API__.

### Status Codes

* 200 : no error
* 404 : no such container
* 500 : no error

## List processes running inside a container

```GET containers/(id)/top```

__Not Supported By VIC__

## Get container logs

```GET /containers/(id)/logs```

This operation will return StdOut and StdErr logs from the target container. __Supported in VIC__

### Query parameters

|Paramter|Description|support|
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

## Inspect changes on a container’s filesystem

```GET /containers/(id)/changes```

inspects the r/w layer of a container for changes to it's filesystem.



### Response

the response for this api involves a json based HTTP response that is a json array of objects.

__JSON OBJECT STRUCTURE__

```
{
    "Path" : <string, path to file>,
    "Kind" : <integer, kind code>
}
```

__Kind Codes__

|code number| description|
|---|---|
|0|Modify|
|1|Add|
|2|Delete|

### Response Status Codes

* 200 : no error
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

## Get container stats based on resource usage
```GET /containers/(id)/stats```

This api returns a stream of the containers resource usage.

__Not Supported BY VIC__

the isolation unit of a containerVM is a vm so these stats are actively available from vSphere.

## Resize a container TTY

__Not Supported by VIC__

## Start a container

```POST /containers/(id or name)/start```

This call starts the indicated container, if it exists. The body of this post call is empty.

### Query Parameter

This call also has one query parameter, __detachKeys__ that allows for overriding the key sequence for detaching a container.       

###Response

```HTTP/1.1 <error code and message>```

###Response Codes

* 204 : no error
* 304 : container already started
* 404 - no such container
* 500 : server error

## Stop a container

```POST /containers/(id or name)/stop```

This call stops the indicated container, if it exists. The body of this post call is empty.

### Query Parameter

This call also has one query parameter, __t__ that specifies how long to wait before killing the indicated container.       

###Response

```HTTP/1.1 <error code and message>```

###Response Codes

* 204 : no error
* 304 : container already started
* 404 - no such container
* 500 : server error

## Restart a container
```POST /containers/(id or name)/restart```

This call restarts the indicated container, if it exists. The body of this post call is empty.

### Query Parameter

This call also has one query parameter, __t__ that specifies how long to wait before restarting the indicated container. 
      
###Response

```HTTP/1.1 <error code and message>```

###Response Codes

* 204 : no error
* 304 : container already started
* 404 - no such container
* 500 : server error

## Kill a container
```POST /containers/(id or name)/restart```

This call restarts the indicated container, if it exists. The body of this post call is empty.

### Query Parameter

This call also has one query parameter, __signal__ that specifies what signal to send with the kill command. e.g. `SIGKILL` 
      
###Response

```HTTP/1.1 <error code and message>```

###Response Codes

* 204 : no error
* 304 : container already started
* 404 - no such container
* 500 : server error

## Update a container

__Might Be Supported By VIC__

## Rename a container

```POST /containers/(id or name)/rename```

rename the target container to the name relayed in the query parameter.

### Query Parameter(s)

* name : new name for the indicated container.

### Status Codes

* 204 : no error
* 404 : no such container found
* 409 : name collision
* 500 : server error

## Pause a container

__Not Supported By VIC__

## Unpause a container

__Not Supported By VIC__

## Attach to a container
 
```POST /containers/(id or name)/attach```

This call will return a stream of the targeted containers common io streams(stdout, stdin, stderr). It is important to not that this involves a `HTTP UPGRADE` response.

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

## Wait a container

__Not Supported By VIC__

## Remove a container

```Delete /containers/(id or name)```

removes indicated container.

###Query Paramter

|parameter|description|supported?|
|---|---|---|
|__v__| 1/true/True or 0/false/False. Removes associated volumes of container.| *maybe*|
|__force__| 1/true/True or 0/false/False. Perform a kill then remove operation | YES |

###Status Codes

* 204 : no error
* 400 : bad parameter
* 404 : no container
* 500 : server error

## Copy files or folders from a container

```POST /containers/(id or name)/copy```

copies files from the target container.

### Request
The request will have one paramter in it's json object.

* resource : the file or folder you want to copy from docker.

### Response

The response will return a tar stream of the requested files.

###Status COdes

* 200 : no error
* 404 : no such container
* 500 : server error 

## Retrieving information about files and folders in a container

```HEAD /containers/(id or name)/archive```

## Get an archive of a filesystem resource in a container

```Get /containers/(id or name)/archive```

get a tar archive of a file or directory in the filesystem of the target container.



### Query Parameters

* __Path__ : required, if it is not an absolute path, it should be a path from the root directory. path should end in `/` or `/.` to assure that the target is a directory. A symlink is always resolved to it's target. When using `/.` this should indicate that only the contents of the directory targeted are copied. 

### Status Codes:

* 200 : success, returns resource as tar stream.
* 400 : client error, bad parameter, details in json response body must specify path, not a directory error.
* 404 : client error, resource not found, no such container found. 
* 500: server error.

## Extract an archive of files or folders to a directory in a container

```PUT /containers/(id or name)/archive```

upload a tar archive to be extracted onto the target container.

### Query Parameter

|parameter|description|supported?|
|---|---|---|
|path|path to extact to. Thepath must exist.| YES|
|noOverwriteDirNonDir|1/true/True. if true, then an error will be returned if the extraction would replace an existing directory with a non directory|YES|

### Status Codes

* 200 : Success
* 400 : client error, bad parameter, must specify path error, not a directory, unable to overwrite existing directory with non directory, unable to overwrite existing non directory with directory
* 403 : client error, permission denied, the volume or container rootfs is marked read-only
* 404 : client error, resource not found, no such container, no such file or directory
* 500 : server error

# Image APIs

## List Images

```Get /images/json```

### Query Parameters

|parameter|description|Supported?|
|---|---|---|
|all|1/true/True or 0/false/False. If true all images are listed| YES |
|filters|json encoded filters. options: `dangling=true` and `label="key=value`"|YES|
|filter| only return the image by the name specified| YES |

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

## Create an image

```POST /images/create```

create image by pulling or importing it. 

### Query Parameters

|parameter|description|supported?|
|---|---|---|
|fromImage|Name of the image to pull|YES|
|fromSrc|Source to import|YES|
|repo|the repo name attached to the image when it is imported. May include a tag.|YES|
|tag|tag or digest|YES|

### Request Headers
|header| description|supported?|
|---|---|---|
|x-Registry-Config| JSON config object that allows for providing mapped credentials to different registries that may be needed to build a specific image| YES|

### Status Codes

* 200 : no error
* 500 : server error

## Inspect an image

```GET /images/(name)/json```

returns low level metadata on the indicated image.

### Status Codes

* 200 : no error
* 404 : no such image
* 500 : server error

## Get the history of an image

```GET /images/(name)/history```

returns the history of the indicated image.

### Response 

The response will have several json array where each element in an image manifest for an image that is in the history of the target image. they are in historical order. newest first.

### Status Codes

* 200 : no error
* 404 : no such image
* 500 : server error

## Push an image on the registry
__Supported at a Future Date__

## Tag an image into a repository

```POST /images/(name)/tag```

tag the target image with the indicated tag.

### Query Parameters

|parameter|description|supported?|
|---|---|---|
|repo|target repo for the tag|YES|
|force|1/true/True or 0/false/False|YES|
|tag|The tag for the target image|YES|

### Status Codes

* 201 : no error(note different than the noremall 200)
* 400 : bad parameter
* 404 : no such image
* 409 : conflict 
* 500 : server error

## Remove an image

```DELETE /images(name)```

### Query Parameters

|parameter|description|supported?|
|---|---|---|
|force|1/true/True or 0/false/False|YES|
|noprune|1/true/True or 0/false/False. default false.|YES|

### Status Codes

* 200 : no error(note different than the noremall 200)
* 404 : no such image
* 409 : conflict 
* 500 : server error

## Search images

```GET /images/search```

### Query Parameters

|parameter|description|supported?|
|---|---|---|
|term|substring to search for|YES|

### Status Codes

* 200 : no error(note different than the noremall 200)
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

# MISC APIs

## Check auth configuration

```POST /auth```

default credentials.

### Request 

`Content-type : application/json` 

|parameters|
|---|
|username|
|password|
|email|
|serveraddress|

### Status Codes
note: the 1.22 api docs do not explain a difference between 200 and 204. so standard difference is assumed.

* 200 : no error
* 204 : no error
* 500 : server error

## Display system-wide information

```GET /info```

Gets system-wide information about the machine running the daemon.

### Status Codes 

* 200 : no error
* 500 : server error

## Show the docker version information

```Get /version```

responds with the version of docker the daemon is currently running. The response
comes as a json object.

### Status Codes

* 200 : no error
* 500 : server error

## Ping the docker server

```GET /_ping```

makes a simple ping against the docker server. 

### Status Codes

* 200 : no error
* 500 : server error

## Monitor Docker’s events

```GET /events```

Fetches the events for containers, images, volumes, and networks.

*Note: this needs more information*

### Query Parameters

|parameters|description|supported?|
|---|---|---|
|since|timestamp for all events since a certain time.|YES|
|until|Timestamp indicating events that are occuring up until a specified time|YES|
|filters| a json object that contains the following key value |YES|

NOTE:
this is the format of the filters object
```
{
    containers=<string>,//container name to look for.
    event=<string>,//event to look for.
    image=<string>,//image to look for.
    label=<string>,//image and container label to look for.
    type=<string>,//this will be one of: container, image, volume, or network
    volume=<string>,//the volume to look for.
    network=<string>//the network to look for.
}
```

## Create a new image from a container’s changes
__Supported at a Future Date__
## Get a tarball containing all images in a repository
__Supported at a Future Date__
## Get a tarball containing all images.
__Supported at a Future Date__
## Load a tarball with a set of images and tags into docker
__Supported at a Future Date__
## Exec Create
__Not Supported By VIC__
## Exec Start
__Not Supported By VIC__
## Exec Resize
__Not Supported By VIC__
## Exec Inspect
__Not Supported By VIC__
