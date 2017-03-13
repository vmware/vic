# Supported Docker Compose File Options #

vSphere Integrated Containers Engine 0.9 supports [Docker Compose file version 2 and 2.1](https://docs.docker.com/compose/compose-file/compose-file-v2/).

This topic provides information about the Docker Compose file options that vSphere Integrated Containers Engine 0.9 supports.

- [Service Configuration Options](#service)
- [Volume Configuration Options](#volume)
- [Network Configuration Options](#network)

<a id="service"></a>
## Service Configuration Options ##

| **Option** | **Compose File Reference** | **Supported** |
| --- | --- | --- |
|`build`|[Options applied at build time](https://docs.docker.com/compose/compose-file/compose-file-v2/#build)|No|
|`cap_add`, `cap_drop`|[Add or drop container capabilities](https://docs.docker.com/compose/compose-file/compose-file-v2/#capadd-capdrop)|No. Depends on  `docker run --cap-add` and `docker run --cap-drop`|
|`command`|[Override the default command](https://docs.docker.com/compose/compose-file/compose-file-v2/#command)|Yes|
|`cgroup_parent` |[Specify an optional parent `cgroup` for the container.](https://docs.docker.com/compose/compose-file/compose-file-v2/#cgroupparent) | No; need docker run --cgrop_parent |
|`container_name` |[Specify a custom container name](https://docs.docker.com/compose/compose-file/compose-file-v2/#containername)| Yes |
|`devices` |[List of device mappings](https://docs.docker.com/compose/compose-file/compose-file-v2/#devices) | No. Depends on `docker create --device`.|
|`depends_on`|[Express dependency between services](https://docs.docker.com/compose/compose-file/compose-file-v2/#dependson)| Yes |
|`dns` |[Custom DNS servers](https://docs.docker.com/compose/compose-file/compose-file-v2/#dns)| Yes |
|`dns_search` |[Custom DNS search domains](https://docs.docker.com/compose/compose-file/compose-file-v2/#dnssearch) | No. Depends on `docker run --dns-search`. |
|`tmpfs` |[Mount a temporary file system inside the container](https://docs.docker.com/compose/compose-file/compose-file-v2/#tmpfs) | No. Depends on `docker run --tmpfs`.|
|`entrypoint` |[Override the default entry point](https://docs.docker.com/compose/compose-file/compose-file-v2/#entrypoint) | No. Depends on `docker run --entrypoint`.|
|`env_file` |[Add environment variables from a file](https://docs.docker.com/compose/compose-file/compose-file-v2/#envfile)| Yes |
|`environment` |[Add environment variables](https://docs.docker.com/compose/compose-file/compose-file-v2/#environment)| Yes |
|`expose` |[Expose ports without publishing them to the host machine](https://docs.docker.com/compose/compose-file/compose-file-v2/#expose) | No. Depends on `docker run --expose`.|
|`extends` |[Extend another service](https://docs.docker.com/compose/compose-file/compose-file-v2/#extends)| Yes |
|`external_links` |[Link to containers started outside this YML](https://docs.docker.com/compose/compose-file/compose-file-v2/#externallinks) | No. Depends on `docker rename`.|
|`extra_hosts` |[Add hostname mappings](https://docs.docker.com/compose/compose-file/compose-file-v2/#extrahosts) | No. Depends on `docker run --add-host`.|
|`group_add` |[Specify additional groups for the user inside the container](https://docs.docker.com/compose/compose-file/compose-file-v2/#groupadd)| Yes |
|`healthcheck` |[Check container health](https://docs.docker.com/compose/compose-file/compose-file-v2/#healthcheck) | No. Depends on `docker run --health-cmd`.|
|`image` |[Specify container image](https://docs.docker.com/compose/compose-file/compose-file-v2/#image)| Yes |
|`isolation` |[Specify isolation technology](https://docs.docker.com/compose/compose-file/compose-file-v2/#isolation) | No. Depends on `docker run --isolation`.|
|`labels` |[Add metadata by using labels](https://docs.docker.com/compose/compose-file/compose-file-v2/#labels)| Yes |
|`links` |[Link to containers in another service](https://docs.docker.com/compose/compose-file/compose-file-v2/#links)| Yes |
|`logging`, `log_driver`, `log_opt` |[Logging configuration](https://docs.docker.com/compose/compose-file/compose-file-v2/#logging) | No. Depends on `docker run --log-driver` and `--log-opt`.|
|`net` |[Network mode (version 1)](https://docs.docker.com/compose/compose-file/compose-file-v1/#net)| Yes |
|`network_mode` |[Network mode (version 2)](https://docs.docker.com/compose/compose-file/compose-file-v2/#networkmode)| Yes |
|`networks` |[Networks to join](https://docs.docker.com/compose/compose-file/compose-file-v2/#networks)| Yes |
|`aliases` |[Aliases for this service on the network](https://docs.docker.com/compose/compose-file/compose-file-v2/#aliases)| Yes |
|`ipv4_address`,`ipv6_address` |[Static IP address for containers](https://docs.docker.com/compose/compose-file/compose-file-v2/#ipv4address-ipv6address) | Yes for IPv4. vSphere Integrated Containers Engine does not support IPv6. |
|`link_local_ips` |[List of link-local IPs](https://docs.docker.com/compose/compose-file/compose-file-v2/#linklocalips) | No. Depends on `docker run --link-local-ip` |
|`pid` |[Sets PID mode](https://docs.docker.com/compose/compose-file/compose-file-v2/#pid) | No. Depends on `docker run --pid`. |
|`ports` |[Expose ports](https://docs.docker.com/compose/compose-file/compose-file-v2/#ports)| Yes |
|`security-opt` |[Override the default labeling scheme for containers](https://docs.docker.com/compose/compose-file/compose-file-v2/#securityopt) | No. This option only applies to Windows containers, which vSphere Integrated Containers Engine does not support. |
|`stop-signal` |[Sets an alternative signal to stop the container.](https://docs.docker.com/compose/compose-file/compose-file-v2/#stopsignal) | Yes |
|`stop-grace-period` |[Specify how long to wait stopping a container](https://docs.docker.com/compose/compose-file/compose-file-v2/#stopgraceperiod) | No |
|`sysctls` |[Kernel parameters to set in the container](https://docs.docker.com/compose/compose-file/compose-file-v2/#sysctls) | No |
|`ulimits` |[Override the default ulimits for a container](https://docs.docker.com/compose/compose-file/compose-file-v2/#ulimits) | No |
|`userns_mode` |[Disables the user namespace](https://docs.docker.com/compose/compose-file/compose-file-v2/#usernsmode) | No |
|`volumes`,`volume_driver` |[xxx](https://docs.docker.com/compose/compose-file/compose-file-v2/#volumes-volumedriver) | Yes |
|`volumes_from` |[Mount volumes from another service or container](https://docs.docker.com/compose/compose-file/compose-file-v2/#volumesfrom)| No |


 The following [Docker run options](https://docs.docker.com/compose/compose-file/compose-file-v2/#cpushares-cpuquota-cpuset-domainname-hostname-ipc-macaddress-memlimit-memswaplimit-memswappiness-oomscoreadj-privileged-readonly-restart-shmsize-stdinopen-tty-user-workingdir) are supported if their [`docker run`](container_operations.md#docker_run) counterpart is supported: `security_opt`, `stop_grace_period`, `stop_signal`, `sysctls`, `ulimits`, `userns_mode`, `cpu_shares`, `cpu_quota`, `cpuset`, `domainname`, `hostname`, `ipc`, `mac_address`, `mem_limit`, `memswap_limit`, `oom_score_adj`, `privileged`, `read_only`, `restart`, `shm_size`, `stdin_open`, `tty`, `user`, `working_dir`.

<a id="volume"></a>
## Volume Configuration Options ##

NOTE: vSphere Integrated Containers 0.9 does support shared volumes. You can use these options for containers that do not share volumes. 

| **Option** | **Compose File Reference** | **Supported** |
| --- | --- | --- |
|`driver` | [Specify driver to use for this volume](https://docs.docker.com/compose/compose-file/compose-file-v2/#driver) | Yes |
|`driver_opts` | [Specify options to pass to the driver for this volume](https://docs.docker.com/compose/compose-file/compose-file-v2/#driveropts) | Yes |
|`labels` | [Add metadata to containers](https://docs.docker.com/compose/compose-file/compose-file-v2/#labels-1) | Yes |
|`external` | [Specify that volume has been created outside of Compose](https://docs.docker.com/compose/compose-file/compose-file-v2/#external) | Yes |

<a id="network"></a>
## Network Configuration Options ##

| **Option** | **Compose File Reference** | **Supported** |
| --- | --- | --- |
|`driver` | [Specify driver to use for this network](https://docs.docker.com/compose/compose-file/compose-file-v2/#driver-1) | Yes |
|`driver_opts` | [Specify options to pass to the driver for this network](https://docs.docker.com/compose/compose-file/compose-file-v2/#driveropts-1) | No |
|`enable_ipv6` | [Enables IPv6](https://docs.docker.com/compose/compose-file/compose-file-v2/#enableipv6) | No. vSphere Integrated Containers Engine does not support IPv6. |
|`ipam`| [Specify custom IPAM configuration](https://docs.docker.com/compose/compose-file/compose-file-v2/#ipam) | No. Depends on `docker network create --ipam`. |
|`internal` | [Create an externally isolated overlay network](https://docs.docker.com/compose/compose-file/compose-file-v2/#internal) | No. Depends on `docker network create --internal`.|
|`labels` | [Add metadata to containers](https://docs.docker.com/compose/compose-file/compose-file-v2/#labels-2) | No. Depends on `docker network --label`. |
|`external` | [Specify that network has been created outside of Compose](https://docs.docker.com/compose/compose-file/compose-file-v2/#external-1) | Yes |
