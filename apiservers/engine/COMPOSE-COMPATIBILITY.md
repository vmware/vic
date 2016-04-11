# Docker Compose Compatibility
Docker Compose is a tool to define an application consisting of multiple containers that interact with one another.  The concept has some similarity with RKT's and Kubernete's pods.

The way Compose works is by orchestrating the start of containers, attaching volumes, and attaching networks by calling the Docker Remote API.  Since Compose uses the Remote API, VIC should support nearly every CLI operation and Compose file options.

## Config File options

These are the compose yml config file options support provided by VIC.  See the following for Docker's compose file documentation: https://docs.docker.com/compose/compose-file/

|Compose yml config option|VIC Support|
|---|:---|
|build|No, *caveat|
|context|No, not until VIC support build|
|dockerfile|No, not until VIC support build|
|args|No, not until VIC support build|
|cap_add, cap_drop|?|
|command|Yes|
|cgroup_parent|No, VIC does not support cgroups|
|container_name|Yes|
|devices|No, a VIC container host is a VCH|
|depends_on|Yes|
|dns|Yes|
|dns_search|Yes|
|tmpfs|No, a VIC container host is a VCH|
|entrypoint|Yes|
|env_file|Yes|
|environment|Yes|
|expose|Yes|
|extends|Yes|
|external_links|Yes|
|extra_hosts|Yes|
|image|Yes|
|labels|Yes|
|links|Yes|
|logging|Yes|
|log_driver|Yes|
|log_opt|Yes|
|network_mode|Yes, *see below|
|networks|Yes|
|ipv4_address, ipv6_address|Yes|
|pid|No, a VIC container host is a VCH|
|ports|Yes|
|security_opt|Yes|
|stop_signal|Yes|
|ulimits|?|
|volumes, volume_driver|Yes, *see below|
|volumes_from|?|
|cpu_shares, cpu_quota, cpuset, domainname, hostname, ipc, mac_address, mem_limit, memswap_limit, privileged, read_only, restart, shm_size, stdin_open, tty, user, working_dir|?|
|Volume configuration reference: driver, driver_opts, external|Yes, *see below|
|Network configuration reference: driver, driver_opts, ipam, external|Yes, *see below|

### network_mode, network_driver

In VIC, all networks are currently *bridged* mode.  Policy around access within those networks are managed via vSphere.  This provides a much simpler model for administration.  Networks can span multiple ESX hosts within a VCH.  These networks can be applied to any container running within the VCH, regardless of which ESX host the container is currently running on.

The network_mode option is ignored.  network_driver option is also ignored.

### volume_driver

VIC provides only one volume driver.  This option is ignored.  Data volumes for VIC containers are VMDKs, and they can exist on any ESX host participating in a VCH.  During volume creation, VIC will handle the mapping of these VMDKs to the path specified in the docker CLI call.  The VIC container and data volume can be on the same ESX host or on different ESX host, as long as the ESX hosts participate in the VCH.  This is much more flexible than the standard docker volumes.
