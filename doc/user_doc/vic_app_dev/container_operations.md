# Supported Docker Commands

 vSphere Integrated Containers Engine 0.9 supports Docker 1.13. The supported version of the Docker API is 1.25. If you are using a more recent version of the Docker client, see [Docker Commands Fail with a Docker API Version Error](ts_docker_version_error.md). 

- [Docker Management Commands](#mgmt)
- [Image Commands](#image)
- [Container Commands](#container)
- [Hub and Registry Commands](#registry)
- [Network and Connectivity Commands](#network)
- [Shared Data Volume Commands](#volume)
- [Docker Compose Commands](#compose)
- [Swarm Commands](#swarm)

<a name="mgmt"></a>
## Docker Management Commands ##

| **Command** | **Docker Reference** | **Supported** |
| --- | --- | --- |
|`dockerd`|[Launch the Docker daemon](https://docs.docker.com/engine/reference/commandline/dockerd/)|No|
|`info`|[Docker system information](https://docs.docker.com/engine/reference/commandline/info/)|Yes. Provides Docker-specific data, basic capacity information, lists configured volume stores, and virtual container host information. Does not reveal vSphere datastore paths that might contain sensitive vSphere information.|
|`inspect`|[Inspect a container or image](https://docs.docker.com/engine/reference/commandline/inspect/)|Yes. Includes information about the container network.|
|`version`|[Docker version information](https://docs.docker.com/engine/reference/commandline/version/)|Yes. vSphere Integrated Containers Engine version provided.|

<a name="image"></a>
## Image Commands ##

| **Command** | **Docker Reference** | **Supported** |
| --- | --- | --- |
|`build`|[Build an image from a Dockerfile](https://docs.docker.com/engine/reference/commandline/build/)|No|
|`commit`|[Create a new image from a container’s changes](https://docs.docker.com/engine/reference/commandline/commit/)|No|
|`history`|[Show the history of an image](https://docs.docker.com/engine/reference/commandline/history/)|No|
|`images`|[Images](https://docs.docker.com/engine/reference/commandline/images/)|Yes. Supports `--filter`, `--no-trunc`, and `--quiet`|
|`import`|[Import the contents from a tarball to create a filesystem image](https://docs.docker.com/engine/reference/commandline/import/)|No|
|`load`|[Load an image from a tar archive or STDIN](https://docs.docker.com/engine/reference/commandline/load/)|No|
|`rmi`|[Remove a Docker image](https://docs.docker.com/engine/reference/commandline/rmi/)|Yes|
|`save`|[Save images](https://docs.docker.com/engine/reference/commandline/save/)|No|
|`tag`|[Tag an image into a repository](https://docs.docker.com/engine/reference/commandline/tag/)|Yes|

<a name="container"></a>
## Container Commands ##

| **Command** | **Docker Reference** | **Supported** |
| --- | --- | --- |
|`attach`|[Attach to a container](https://docs.docker.com/engine/reference/commandline/attach/)|Yes|
|`container list`|[List Containers](https://docs.docker.com/engine/reference/api/docker_remote_api_v1.22/#list-containers)|Yes|
|`container resize`|[Resize a container](https://docs.docker.com/engine/reference/api/docker_remote_api_v1.23/#resize-a-container-tty)|Yes|
|`cp`|[Copy files or folders between a container and the local filesystem](https://docs.docker.com/engine/reference/commandline/cp/)|No|
|`create`|[Create a container](https://docs.docker.com/engine/reference/commandline/create/)|Yes. <br>`--cpuset-cpus`in Docker specifies CPUs the container is allowed to use during execution (0-3, 0,1). In vSphere Integrated Containers Engine, this parameter specifies the number of virtual CPUs to allocate to the container VM. Minimum CPU count is 1, maximum is unlimited. Default is 2.<br>`--ip` allows you to set a static IP on the container. By default, the virtual container host  manages the container IP.<br>`--memory` Minimum memory is 512MB, maximum unlimited. If unspecified, default is 2GB. Supports the `--attach`, `--cpuset-cpus`,`--env`, `--ip`, `--memory`, `--interactive`, `--link`, `--label`, `--network`, `--tty`, and `--volume` options.|
|`diff`|[Inspect changes on a container's filesystem](https://docs.docker.com/engine/reference/commandline/diff/)|No|
|`events`|[Get real time events from the server](https://docs.docker.com/engine/reference/commandline/events/)|Yes. Supports passive Docker events for containers and images. Does not yet support events for volumes or networks.|
|`exec`|[Run a command in a running container](https://docs.docker.com/engine/reference/commandline/exec/)|No|
|`export`|[Export a container](https://docs.docker.com/engine/reference/commandline/export/)|No|
|`kill`|[Kill a running container](https://docs.docker.com/engine/reference/commandline/kill/)|Yes. Docker must wait for the container to shut down.|
|`logs`|[Get container logs](https://docs.docker.com/engine/reference/commandline/logs/)|Yes. Does not support `docker logs --timestamps` (`-t`) and `--since` options.|
|`pause`|[Pause processes in a container](https://docs.docker.com/engine/reference/commandline/pause/)|No|
|`port`|[Obtain port data](https://docs.docker.com/engine/reference/commandline/port/)|Yes. Displays port mapping data. <br>Supports mapping a random host port to the container when the host port is not specified.|
|`ps`|[Show running containers](https://docs.docker.com/engine/reference/commandline/ps/)|Yes. Supports the `-a/--all`, `-f/--filter`, `--no-trunc`, and `-q/--quiet` options. Filtering by network name is supported, but filtering by network ID is not supported.|
|`rename`|[Rename a container](https://docs.docker.com/engine/reference/commandline/rename/)|No|
|`restart`|[Restart a container](https://docs.docker.com/engine/reference/commandline/restart/)|Yes|
|`rm`|[Remove a container](https://docs.docker.com/engine/reference/commandline/rm/)|Yes. Removes associated anonymous and regular volumes. Supports the `--force` option and the `name` parameter.  Does not support `docker rm -v`. To view volumes attached to a container that is removed, use `docker volume ls` and `docker volume inspect <id>`. If you continually invoke `docker create` to make more anonymous volumes, those volumes are left behind after each subsequent removal of that container.|
|`run`|[Run a command in a new container](https://docs.docker.com/engine/reference/commandline/run/)|<a id="docker_run"></a>Yes. Supports container search by using prettyname-ID with `docker run --name`. Supports the `--detach`, `--detach-keys`, and `--dns` options.<br> Supports mapping a random host port to the container when the host port is not specified. <br>Supports running images from private and custom registries.<br>`docker run --net=host` is not supported. You can specify a container network by using the [`--container-network` option](../vic_installation/vch_installer_options.html#container-network) when you deploy a virtual container host.|
|`start`|[Start a container](https://docs.docker.com/engine/reference/commandline/start/)|Yes|
|`stats`|[Get container stats based on resource usage](https://docs.docker.com/engine/reference/commandline/stats/)|No|
|`stop`|[Stop a container](https://docs.docker.com/engine/reference/commandline/stop/)|Yes. Attempts to politely stop the container. If that fails, powers down the VM.|
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
|`pull`|[Pull an image or repository from a registry](https://docs.docker.com/engine/reference/commandline/pull/)| Yes. Supports pulling from  secure or insecure public and private registries.|
|`push`|[Push an image or a repository to a registry](https://docs.docker.com/engine/reference/commandline/push/)|No|
|`search`|[Search the Docker hub for images](https://docs.docker.com/engine/reference/commandline/search/)|No|

<a name="network"></a>
## Network and Connectivity Commands ##

For more information about network operations with vSphere Integrated Containers Engine, see [Container Networking with vSphere Integrated Containers Engine](network_use_cases.md).

| **Command** | **Docker Reference** | **Supported** |
| --- | --- | --- |
|`network connect`|[Connect to a network](https://docs.docker.com/engine/reference/commandline/network_connect/)|Yes. Not supported for running containers.|
|`network create`|[Create a network](https://docs.docker.com/engine/reference/commandline/network_create/)|Yes. See the use case to connect to an external network in [Container Networking with vSphere Integrated Containers Engine](network_use_cases.md). Bridge is also supported.<br>The `--label` and `--internal` options are currently not supported.|
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
|`volume inspect`|[Information about a volume](https://docs.docker.com/engine/reference/commandline/volume_inspect/)|Yes|
|`volume ls`|[List volumes](https://docs.docker.com/engine/reference/commandline/volume_ls/)|Yes|
|`volume rm`|[Remove or delete a volume](https://docs.docker.com/engine/reference/commandline/volume_rm/)|Yes|

<a name="compose"></a>
## Docker Compose Commands ##

vSphere Integrated Containers Engine 0.9 supports Docker Compose version 1.9.0.

For more information about using Docker Compose with vSphere Integrated Containers Engine, see [Creating a Containerized Application with vSphere Integrated Containers Engine](creating_containerized_app_with_vic.md).

For information about Docker Compose file support, see [Supported Docker Compose File Options](docker_compose_file_options.md).

| **Command** | **Docker Reference** | **Supported** |
| --- | --- | --- |
| `build`  | [Build or rebuild service](https://docs.docker.com/compose/reference/build/)  | No. Depends on `docker build`.|
| `bundle`  | [Generate a Distributed Application Bundle (DAB) from the Compose file](https://docs.docker.com/compose/reference/bundle/)| No |
| `config`  | [Validate and view the compose file](https://docs.docker.com/compose/reference/config/)  | Yes  |
| `create`  | [Create services](https://docs.docker.com/compose/reference/create/)  | Yes  |
| `down`  | [Stop and remove containers, networks, images, and volumes](https://docs.docker.com/compose/reference/down/)  | Yes  |
| `events`  |[Receive real time events from containers](https://docs.docker.com/compose/reference/events/)  | Yes. Supports passive Docker events for containers and images. Does not yet support events for volumes or networks.|
| `exec`  | [Run commands in services](https://docs.docker.com/compose/reference/exec/) | No. Depends on `docker exec`. |
| `help`  | [Get help on a command](https://docs.docker.com/compose/reference/help/)  | Yes  |
| `kill`  | [Kill containers](https://docs.docker.com/compose/reference/kill/)  | No, but `docker kill` works. |
| `logs`  | [View output from containers](https://docs.docker.com/compose/reference/logs/)  | Yes |
| `pause`  | [Pause services](https://docs.docker.com/compose/reference/pause/)  | No. Depends on `docker pause`.  |
| `port`  | [Print the public port for a port binding](https://docs.docker.com/compose/reference/port/)  | Yes |
| `ps`  | [List containers](https://docs.docker.com/compose/reference/ps/)  |Yes |
| `pull`  | [Pulls service images](https://docs.docker.com/compose/reference/pull/)  | Yes  |
| `push`  | [Pushes images for service](https://docs.docker.com/compose/reference/push/)  | No. Depends on `docker push`  |
| `restart`  |	[Restart services](https://docs.docker.com/compose/reference/restart/)  | Yes  |
| `rm`  | [Remove stopped containers](https://docs.docker.com/compose/reference/rm/)  | Yes  |
| `run`  | [Run a one-off command](https://docs.docker.com/compose/reference/run/)  | Yes  |
| `scale`  | [Set number of containers for a service](https://docs.docker.com/compose/reference/scale/)  | Yes |
| `start`  | [Start services](https://docs.docker.com/compose/reference/start/)  | Yes  |
| `stop`  | [Stop services](https://docs.docker.com/compose/reference/stop/)  | Yes  |
| `unpause`  | [Unpause services](https://docs.docker.com/compose/reference/unpause/)  | No. Depends on `docker unpause`.  |
| `up`  | [Create and start containers]()  | Conditionally supported. Does not work if there are orphaned containers. Depends on `docker rename` and `docker ps --filter`.  |
| `version`  | Show Docker Compose version information  | Yes  |


<a name="swarm"></a>
## Swarm Commands ##

This version of vSphere Integrated Containers Engine does not support Docker Swarm.