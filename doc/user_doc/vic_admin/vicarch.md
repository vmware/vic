# vSphere Integrated Containers Engine Architecture

vSphere Integrated Containers Engine exists in a vSphere environment, allowing you to manage containers like virtual machines. The architecture consists of these components:

- vCenter Server management tools: Monitor and manage container virtual machines alongside regular virtual machines.

- Trusted networks: Deploy and use vSphere Integrated Containers Engine and create connections between Docker clients and virtual container hosts.

- Virtual SAN datastores: Specify Virtual SAN datastores in which to store container images, container VM files, container volumes, and the virtual container host vApp.

- Docker API appliance virtual machine: The vSphere Integrated Containers Engine installer deploys a vApp, referred to as the virtual container host. You point Docker clients to this appliance for use as the Docker endpoint.

- Docker container virtual machines: Using Photon OS technology, you create and provision multiple container virtual machines directly from a template. The Docker daemon runs outside the container virtual machine. The container is a x86 hardware virtualized virtual machine with a process ID, container interfaces and mounts.
 
![vSphere Integrated Containers Engine Architecture](vSphereContainerArch.png)

## Virtual Container Host 

The virtual container host appliance is backed by a Photon OS kernel that provides a virtual container endpoint backed by a vSphere vApp that allows you to control and consume container services.

You can access a Docker API endpoint for development and map ports for client connections to run containers as required.

vSphere resource management handles container placement within the virtual container host, so that a virtual container host can be served by an entire vSphere cluster or by a fraction of the same cluster. The only resources consumed by a container host in the cluster are the resources consumed by the container VMs that run in it.

You can reconfigure the virtual container host with no impact to containers running in it. The virtual container host is not limited by the kernel version or by the operating system that the containers are running.

You can deploy multiple virtual container hosts in an environment, depending on your business needs, including allocating separate resources for development, testing, and production.

You can configure virtual container hosts, giving your development team access to a large virtual container host, or sub-allocate smaller virtual container hosts for individual developers.

Each virtual container host maintains a cache of container images, which you download from either the public Docker Hub or a private registry.

The virtual container host maintains filesystem layers inherent in container images by mapping to discrete VMDK files, all of which are housed in vSphere datastores on VSAN, NFS, or local disks.

You deploy a virtual container host using the CLI installer, then access virtual container host endpoints remotely through a Docker command line interface or other API client.

## vSphere Web Client Plugin

You can monitor virtual container hosts and container VMs by using the vSphere Integrated Containers Engine plug-in for the vSphere Web Client

## Docker Client

Docker clients communicate with the virtual container host, not with each container, so you can see aggregated pools of vSphere resources, including storage and memory allocations.

You can pull standard container images from the Docker hub or from a private registry.

You can create, run, stop, and delete containers using standard docker commands and verify these actions in the vSphere Web Client.