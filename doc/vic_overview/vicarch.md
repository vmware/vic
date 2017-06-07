# vSphere Integrated Containers Architecture

vSphere Integrated Containers exists in a vSphere environment, allowing you to use virtual machines like containers. The architecture consists of these components:

- vCenter management tools: monitor and manage virtual machines as well as container virtual machines.

- vCenter Server: manage a single ESXi host or cluster of ESXi hosts with DRS enabled. Specify and deploy datastores and network paths, define clusters, resource pools, and port groups.

- Trusted networks: Deploy and use vSphere Integrated Containers and connections between Docker clients and virtual container hosts.

- Virtual SAN datastores: specify a datastore and top level directory with the name of the virtual container host.

- Docker API appliance virtual machine: The installer deploys this appliance virtual machine, which is also referred to as the virtual container host. You set the docker client to this appliance.

- Docker container virtual machines: Using vSphere Instant Clone and Photon OS technology, you can create and provision multiple container virtual machines directly from a template. The Docker daemon runs outside the container virtual machine. The container is a x86 hardware virtualized virtual machine with a process ID, container interfaces and mounts.
 
![vSphere Integrated Containers Architecture](vSphereContainerArch.png)

# Virtual Container Host Appliance

The virtual container host appliance is backed by a Photon OS kernel that provides a virtual container endpoint backed by a vSphere resource pool that allows you to control and consume container services.

You can access a Docker API endpoint for development and map ports for client connections to run containers as required.

vSphere resource management handles container placement within the virtual container host, so that a virtual Docker host can serve as an entire vSphere cluster or a fraction of the same cluster. The only resource consumed by a container host in the cluster is the resource consumed by running containers.

You can reconfigure the virtual container host with no impact to containers running in it. The virtual container host is not limited by the kernel version or by the operating system the containers are running.

You can deploy multiple virtual container hosts in an environment, depending on your business needs, including allocating separate resources for development, testing, and production.

You can also nest virtual container hosts, giving your team access to a large virtual container host, or sub-allocate smaller virtual container hosts for individuals.

Each virtual container host maintains a cache of container images, which you download from either the public Docker Hub or a private registry.

The virtual container host maintains filesystem layers inherent in container images by mapping to discrete VMDK files, all of which are housed in vSphere datastores on VSAN, NFS, or local disks.

You deploy a virtual container host using the CLI installer, then access Virtual Container Host endpoints remotely through a Docker command line interface or other API client.

## vSphere Web Client Plugin

You can monitor and manage containers using the vSphere Integrated Containers plugin in the vSphere Web Client.

The plugin allows you to create virtual container hosts, perform administrative tasks on containers such as resource allocation, port mapping, and manage communications between administrators, developers, and application owners during troubleshooting.

You can create, run, stop, and delete containers using standard docker commands in a command line interface and verify these actions in the vSphere Web Client.

## Docker Client

Docker clients communicate with the virtual container host, not each container, so you can see aggregated pools of vSphere resources, including storage and memory allocations.

You can pull standard container images from the Docker hub or private registry.

You can create, run, stop, and delete containers using standard docker commands and verify these actions in the vSphere Web Client.


