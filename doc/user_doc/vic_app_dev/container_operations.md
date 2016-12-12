# Supported Docker Commands

 vSphere Integrated Containers Engine supports Docker 1.11.2. The supported version of the Docker API is 1.23. If you are using a more recent version of the Docker client, see [Docker Commands Fail with a Docker API Version Error](ts_docker_version_error.md). 

- [Docker Management Commands](#mgmt)
- [Image Commands](#image)
- [Container Commands](#container)
- [Hub and Registry Commands](#registry)
- [Network and Connectivity Commands](#network)
- [Shared Data Volume Commands](#volume)
- [Swarm Commands](#swarm)

<a name="mgmt"></a>
## Docker Management Commands ##

| **Command** | **Docker Reference** | **Supported** |
| --- | --- | --- |
|`dockerd`|[Launch the Docker daemon](https://docs.docker.com/engine/reference/commandline/dockerd/)|No|
|`info`|[Docker system information](https://docs.docker.com/engine/reference/commandline/info/)|Yes, docker-specific data, basic capacity information, list of configured volume stores, virtual container host information. Does not reveal vSphere datastore paths that might contain sensitive vSphere information.|
|`inspect`|[Inspect a container or image](https://docs.docker.com/engine/reference/commandline/inspect/)|Yes. Includes information about the container network.|
|`version`|[Docker version information](https://docs.docker.com/engine/reference/commandline/version/)|Yes. vSphere Integrated Containers Engine version provided.|

<a name="image"></a>
## Image Commands ##

| **Command** | **Docker Reference** | **Supported** |
| --- | --- | --- |
|`build`|[Build an image from a Dockerfile](https://docs.docker.com/engine/reference/commandline/build/)|No|
|`commit`|[Create a new image from a container’s changes](https://docs.docker.com/engine/reference/commandline/commit/)|No|
|`history`|[Show the history of an image](https://docs.docker.com/engine/reference/commandline/history/)|No|
|`images`|[Images](https://docs.docker.com/engine/reference/commandline/images/)|Yes. vSphere Integrated Container Engine currently supports `--no-trunc` and `-q/--quiet`|
|`import`|[Import the contents from a tarball to create a filesystem image](https://docs.docker.com/engine/reference/commandline/import/)|No|
|`load`|[Load an image from a tar archive or STDIN](https://docs.docker.com/engine/reference/commandline/load/)|No|
|`rmi`|[Remove a Docker image](https://docs.docker.com/engine/reference/commandline/rmi/)|Yes. vSphere Integrated Container Engine does not yet support tag or untag operations, or any options.|
|`save`|[Save images](https://docs.docker.com/engine/reference/commandline/save/)|No|
|`tag`|[Tag an image into a repository](https://docs.docker.com/engine/reference/commandline/tag/)|No|

<a name="container"></a>
## Container Commands ##

| **Command** | **Docker Reference** | **Supported** |
| --- | --- | --- |
|`attach`|[Attach to a container](https://docs.docker.com/engine/reference/commandline/attach/)|Yes|
|`container list`|[List Containers](https://docs.docker.com/engine/reference/api/docker_remote_api_v1.22/#list-containers)|Yes|
|`container resize`|[Resize a container](https://docs.docker.com/engine/reference/api/docker_remote_api_v1.23/#resize-a-container-tty)|Yes|
|`cp`|[Copy files or folders between a container and the local filesystem](https://docs.docker.com/engine/reference/commandline/cp/)|No|
|`create`|[Create a container](https://docs.docker.com/engine/reference/commandline/create/)|Yes. <br>`--cpuset-cpus`in Docker specifies CPUs the container is allowed to use during execution (0-3, 0,1). In vSphere Integrated Containers Engine, this parameter specifies the number of virtual CPUs to allocate to the container VM. Minimum CPU count is 1, maximum is unlimited. Default is 2.<br>`--ip` allows you to set a static IP on the container. By default, the virtual container host  manages the container IP.<br>`-m --memory` Minimum memory is 512MB, maximum unlimited. If unspecified, default is 2GB. vSphere Integrated Container Engine supports these options: `-a, --attach`, `--cpuset-cpus`,`-e, --env`, `--ip`, `-m, --memory`, `-i, --interactive`, `--link`, `-l, --label`, `--network`, `--tty`, `-v, --volume`|
|`diff`|[Inspect changes on a container's filesystem](https://docs.docker.com/engine/reference/commandline/diff/)|No|
|`events`|[Get real time events from the server](https://docs.docker.com/engine/reference/commandline/events/)|No|
|`exec`|[Run a command in a running container](https://docs.docker.com/engine/reference/commandline/exec/)|No|
|export|[Export a container](https://docs.docker.com/engine/reference/commandline/export/)|No|
|`kill`|[Kill a running container](https://docs.docker.com/engine/reference/commandline/kill/)|Yes. Docker must wait for the container to shut down.|
|`logs`|[Get container logs](https://docs.docker.com/engine/reference/commandline/logs/)|Yes,  vSphere Integrated Container Engine does not yet support `docker logs --timestamps` (`-t`) and `--since` options.|
|`pause`|[Pause processes in a container](https://docs.docker.com/engine/reference/commandline/pause/)|No|
|`port`|[Obtain port data](https://docs.docker.com/engine/reference/commandline/port/)|Yes. Displays port mapping data. <br> vSphere Integrated Container Engine supports mapping a random host port to the container when the host port is not specified.|
|`ps`|[Show running containers](https://docs.docker.com/engine/reference/commandline/ps/)|Yes, vSphere Integrated Containers Engine supports only `-a/--all`, `--no-trunc`, and `-q/--quiet` options.|
|`rename`|[Rename a container](https://docs.docker.com/engine/reference/commandline/rename/)|No|
|`restart`|[Restart a container](https://docs.docker.com/engine/reference/commandline/restart/)|Yes|
|`rm`|[Remove a container](https://docs.docker.com/engine/reference/commandline/rm/)|Yes, vSphere Integrated Containers Engine also removes associated anonymous and regular volumes, supports `-f/--force` option and the <code>name</code> parameter.  `docker rm -v` and `docker rm -f` are not supported. To view volumes attached to a container that is removed, use `docker volume ls` and `docker volume inspect <id>`. If you continually invoke `docker create` to make more anonymous volumes, those volumes are left behind after each subsequent removal of that container.|
|`run`|[Run a command in a new container](https://docs.docker.com/engine/reference/commandline/run/)|Yes. vSphere Integrated Container Engine supports container search using prettyname-ID `docker run --name` and these options: `-d, --detach`. `--detach-keys`, `--dns` <br> vSphere Integrated Container Engine supports mapping a random host port to the container when the host port is not specified. <br>vSphere Integrated Container Engine supports running images from private and custom registries.<br> `docker run --net=host` is not supported. You can specify a container network by using the [`--container-network` option](../vic_installation/vch_installer_options.html#container-network) when you deploy a virtual container host.|
|`start`|[Start a container](https://docs.docker.com/engine/reference/commandline/start/)|Yes|
|`stats`|[Get container stats based on resource usage](https://docs.docker.com/engine/reference/commandline/stats/)|No|
|`stop`|[Stop a container](https://docs.docker.com/engine/reference/commandline/stop/)|Yes. vSphere Integrated Containers Engine attempts to politely stop the container. If that fails, it powers down the VM.|
|`top`|[Display the running processes of a container](https://docs.docker.com/engine/reference/commandline/top/)|No|
|`unpause`|[Unpause processes within a container](https://docs.docker.com/engine/reference/commandline/unpause/)|No|
|`update`| [Update a container](https://docs.docker.com/engine/reference/commandline/update/)|No|
|`wait`|[Wait for a container](https://docs.docker.com/engine/reference/commandline/wait/)|Yes|

<a name="registry"></a>
## Hub and Registry Commands ##

| **Command** | **Docker Reference** | **Supported** |
| --- | --- | --- |
|`login`|[Log into a registry](https://docs.docker.com/engine/reference/commandline/login/)|Yes|
|`logout`|[Log out from a registry](https://docs.docker.com/engine/reference/commandline/logout/)|Yes|
|`pull`|[Pull an image or repository from a registry](https://docs.docker.com/engine/reference/commandline/pull/)| Yes, vSphere Integrated Container Engine supports pulling from public and private registries, secure or insecure.|
|`push`|[Push an image or a repository to a registry](https://docs.docker.com/engine/reference/commandline/push/)|No|
|`search`|[Search the Docker hub for images](https://docs.docker.com/engine/reference/commandline/search/)|No|

<a name="network"></a>
## Network and Connectivity Commands ##

For more information about network operations with vSphere Integrated Containers Engine, see [Container Networking with vSphere Integrated Containers Engine](network_use_cases.md).

| **Command** | **Docker Reference** | **Supported** |
| --- | --- | --- |
|`network connect`|[Connect to a network](https://docs.docker.com/engine/reference/commandline/network_connect/)|Yes, However, this is not supported for running containers.|
|`network create`|[Create a network](https://docs.docker.com/engine/reference/commandline/network_create/)|Yes. See the use case to connect to an external network in [Container Networking with vSphere Integrated Containers Engine](network_use_cases.md). Bridge is also supported.|
|`network disconnect`|[Disconnect a network](https://docs.docker.com/engine/reference/commandline/network_disconnect/)|No|
|`network inspect`|[Inspect a network](https://docs.docker.com/engine/reference/commandline/network_inspect/)|Yes|
|`network ls`|[List networks/](https://docs.docker.com/engine/reference/commandline/network_ls/)|Yes|
|`network rm`|[Remove a network](https://docs.docker.com/engine/reference/commandline/network_rm/)|Yes. Network name and network ID are supported.|

<a name="volume"></a>
## Shared Data Volume Commands ##

For more information about volume operations with vSphere Integrated Containers Engine, see [Using Volumes with vSphere Integrated Containers Engine](using_volumes_with_vic.md).

| **Command** | **Docker Reference** | **Supported** |
| --- | --- | --- |
|`volume create`|[Create a volume](https://docs.docker.com/engine/reference/commandline/volume_create/)|Yes|
|`volume inspect`|[Information about a volume](https://docs.docker.com/engine/reference/commandline/volume_inspect/)|Yes, use with docker compose.|
|`volume ls`|[List volumes](https://docs.docker.com/engine/reference/commandline/volume_ls/)|Yes. `--filter` is not supported.|
|`volume rm`|[Remove or delete a volume](https://docs.docker.com/engine/reference/commandline/volume_rm/)|Yes|

<a name="swarm"></a>
## Swarm Commands ##

This version of vSphere Integrated Containers Engine does not support Docker Swarm.