# Supported Docker Commands

 vSphere Integrated Containers Engine supports Docker 1.11.2. The supported version of the Docker API is 1.23. If you are using a more recent version of the Docker client, see [Docker Commands Fail with a Docker API Version Error](ts_docker_version_error.md).

## General Container Operations ##

| **Container** | **Docker Reference** | **Supported** |
| --- | --- | --- |
|attach|[Attach to a container](https://docs.docker.com/engine/reference/api/docker_remote_api_v1.22/#attach-to-a-container)<br> [Attach to a container websocket](https://docs.docker.com/engine/reference/api/docker_remote_api_v1.22/#attach-to-a-container-websocket)|Yes|
|container list|[List Containers](https://docs.docker.com/engine/reference/api/docker_remote_api_v1.22/#list-containers)|Yes|
|container resize|[Resize a container](https://docs.docker.com/engine/reference/api/docker_remote_api_v1.23/#resize-a-container-tty)|Yes|
|cp|[Copy files or folders between a container and the local filesystem](https://docs.docker.com/engine/reference/commandline/cp/)|No|
|create|[Create a container](https://docs.docker.com/engine/reference/commandline/create/)|Yes. <br>`--cpuset-cpus`in Docker specifies CPUs the container is allowed to use during execution (0-3, 0,1). In vSphere Integrated Containers, this parameter specifies the number of virtual CPUs to allocate to the container VM. Minimum CPU count is 1, maximum is unlimited. Default is 2.<br>`--ip` allows you to set a static IP on the container. By default, the virtual container host  manages the container IP.<br>`-m --memory` Minimum memory is 512MB, maximum unlimited. If unspecified, default is 2GB.|
|diff|[Inspect changes on a container's filesystem](https://docs.docker.com/engine/reference/commandline/diff/)|No|
|events|[Get real time events from the server](https://docs.docker.com/engine/reference/commandline/events/)|No|
|exec|[Run a command in a running container](https://docs.docker.com/engine/reference/commandline/exec/)|No|
|export|[Export a container](https://docs.docker.com/engine/reference/api/docker_remote_api_v1.22/#export-a-container)|No|
|images|[Images](https://docs.docker.com/engine/reference/commandline/images/)<br>[list-images](https://docs.docker.com/engine/reference/api/docker_remote_api_v1.22/#list-images)|Yes|
|info|[Docker system information](https://docs.docker.com/engine/reference/commandline/info/)|Yes, docker-specific data, basic capacity information, list of configured volume stores, virtual container host information. Does not reveal vSphere datastore paths that might contain sensitive vSphere information|
|inspect|[Inspect a container](https://docs.docker.com/engine/reference/api/docker_remote_api_v1.22/#inspect-a-container) <br>[Inspect an image](https://docs.docker.com/engine/reference/api/docker_remote_api_v1.22/#inspect-an-image)|Yes. Includes information about the container network.|
|kill|[Kill a running container](https://docs.docker.com/engine/reference/api/docker_remote_api_v1.22/#kill-a-container)<br> [Kill](https://docs.docker.com/engine/reference/commandline/kill/)|Yes. Docker must wait for the container to shut down.|
|login| [Log into a Docker registry](https://docs.docker.com/engine/reference/commandline/login/)|Yes, use to access private repository|
|logs|[Get container logs](https://docs.docker.com/engine/reference/api/docker_remote_api_v1.22/#get-container-logs)|Yes, except for the `docker logs --timestamps` (`-t`) and `--since` options, which are not supported.|
|pause|[Pause processes in a container](https://docs.docker.com/engine/reference/api/docker_remote_api_v1.22/#pause-a-container)<br> [Pause](https://docs.docker.com/engine/reference/commandline/pause/)|No|
|port|[Obtain port data](https://docs.docker.com/engine/reference/commandline/port/)|Yes. Displays port mapping data. <br> Mapping a random host port to the container when the host port is not specified is supported.|
|ps|[Show running containers](https://docs.docker.com/engine/reference/commandline/ps/)|Yes|
|pull|[Pull an image or repository from a registry](https://docs.docker.com/engine/reference/commandline/pull/)| Yes, pulling from insecure and custom registries is supported|
|rename|[Rename a container](https://docs.docker.com/engine/reference/api/docker_remote_api_v1.22/#rename-a-container) [Rename](https://docs.docker.com/engine/reference/commandline/rename/)|No|
|restart|[Restart a container](https://docs.docker.com/engine/reference/api/docker_remote_api_v1.22/#restart-a-container)<br> [Restart](https://docs.docker.com/engine/reference/commandline/restart/)|Yes|
|rm|[Remove a container](https://docs.docker.com/engine/reference/api/docker_remote_api_v1.22/#remove-a-container)|Yes, only the <code>name</code> parameter is supported. <code>force</code> and <code>v</code> are a future implementation. Also removes associated volumes.|
|rmi|[Remove a Docker image](https://docs.docker.com/engine/reference/api/docker_remote_api_v1.22/#remove-an-image)|Yes|
|run|Composite command of create, start, inspect, attach, rm, resize, wait, kill|Yes. Container search using prettyname-ID `docker run --name` is supported. <br> Mapping a random host port to the container when the host port is not specified is supported. <br>Running images from private and custom registries is supported.<br> `docker run --net=host` is not supported. You can specify a container network by using the [`--container-network` option](../vic_installation/vch_installer_options.html#container-network) when you deploy a virtual container host.|
|save|[Save images](https://docs.docker.com/engine/reference/commandline/save/)|No|
|start|[Start a container](https://docs.docker.com/engine/reference/commandline/start/)|Yes|
|stats|[Get container stats based on resource usage](https://docs.docker.com/engine/reference/api/docker_remote_api_v1.22/#get-container-stats-based-on-resource-usage)<br> [Stats](https://docs.docker.com/engine/reference/commandline/stats/)|No|
|stop|[Stop a container](https://docs.docker.com/engine/reference/api/docker_remote_api_v1.22/#stop-a-container)<br> [Stop](https://docs.docker.com/engine/reference/commandline/stop/)|Yes. Powers down the VM.|
|unpause|[Unpause processes in a container](https://docs.docker.com/engine/reference/api/docker_remote_api_v1.22/#unpause-a-container)<br> [Unpause](https://docs.docker.com/engine/reference/commandline/unpause/)|No|
|update| [Update a container](https://docs.docker.com/engine/reference/api/docker_remote_api_v1.22/#update-a-container) <br> [Update](https://docs.docker.com/engine/reference/commandline/update/)|No|
|version|[Docker version information](https://docs.docker.com/engine/reference/commandline/version/)|Yes. vSphere Integrated Containers Engine version provided|
|wait|[Wait for a container](https://docs.docker.com/engine/reference/api/docker_remote_api_v1.22/#wait-a-container) <br>[Wait](https://docs.docker.com/engine/reference/commandline/wait/)|Yes|

## Network Operations ##

For more information about network operations, see [Network Port Use Cases](network_use_cases.md).

| **Network** | **Docker Reference** | **Supported** |
| --- | --- | --- |
|Network connect|[Connect to a network](https://docs.docker.com/engine/reference/commandline/network_connect/)|Yes|
|Network create|[Create a network](https://docs.docker.com/engine/reference/commandline/network_create/)|Yes. See the use case to connect to an external network in [Network Port Use Cases](network_use_cases.md). Bridge is also supported.|
|Network disconnect|[Disconnect a network](https://docs.docker.com/engine/reference/commandline/network_disconnect/)|No|
|Network inspect|[Inspect a network](https://docs.docker.com/engine/reference/commandline/network_inspect/)|Yes|
|Network ls|[List networks/](https://docs.docker.com/engine/reference/commandline/network_ls/)|Yes|
|Network rm|[Remove a network](https://docs.docker.com/engine/reference/commandline/network_rm/)|Yes. Network name and network ID are supported|

## Volume Operations ##

For more information about volume operations, see [Using Volumes with vSphere Integrated Containers Engine](using_volumes_with_vic.md).

| **Volume** | **Docker Reference** | **Supported** |
| --- | --- | --- |
|Volume create|[Create a volume](https://docs.docker.com/engine/reference/commandline/volume_create/)| The driver option is ignored even if you specify it. You must include <code>--opt VolumeStore=<name of volumestore> --Capacity=<size  in MB></code> as these are direct vSphere arguments. vSphere Integrated Containers does not assign random names during a volume create, but only for anonymous volumes.|
|Volume inspect|[Information about a volume](https://docs.docker.com/engine/reference/commandline/volume_inspect/)|Yes, use with docker compose|
|Volume ls|[List volumes](https://docs.docker.com/engine/reference/commandline/volume_ls/)|Yes|
|Volume rm|[Remove or delete a volume](https://docs.docker.com/engine/reference/commandline/volume_rm/)|Yes|