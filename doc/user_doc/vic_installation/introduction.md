# Overview of vSphere Integrated Containers Engine for vSphere Administrators

vSphere Integrated Containers Engine provides developers the portability, speed, and agility of using  enterprise-class containers, and provide IT Ops the management, security, and visibility they require to run workloads in production.

vSphere Integrated Containers Engine enables IT teams to run traditional and container workloads side-by-side on existing infrastructure seamlessly.

Using constructs from the Open Container Initiative to map Docker containers to vSphere infrastructure, vSphere Integrated Containers Engine containers are provisioned as virtual machines, offering the same security and functionality of virtual machines in VMware ESXi&trade; hosts or VMware vCenter Server&reg; instances.

A virtual container host (VCH) is compatible with standard Docker client tools and backed by a pool of resources to accommodate applications.

From a developer's perspective, vSphere Integrated Containers Engine is a seamless Docker interface for containers with a vSphere back end. Developers can deploy, test, and run container processes faster in the same environment as traditional applications.

You install vSphere Integrated Containers Engine by using a command line installer, `vic-machine`, that deploys VCHs to ESXi hosts or vCenter Server. You connect Docker clients to the VCHs and use the Docker clients to work with containers. You use your vSphere environment to manage the container VMs and container images.

## Comparing vSphere Integrated Containers Engine and Traditional Container Hosts
vSphere Integrated Containers Engine provisions containers as virtual machines, rather than in virtual machines. 

### Traditional Container Host

A traditional container host is a virtual machine running a Linux OS with the necessary libraries, kernel version, and daemon installed. The container host has a fixed amount of memory and vCPU resource used by the containers provisioned into it.

The hypervisor provides hardware virtualization of the entire container host VM, one or more VMDKs providing local disk for the OS, one or more vNICs to provide network connectivity for the OS and possibly paravirtualization capabilities allowing the containers to directly access hypervisor infrastructure.

### vSphere Integrated Containers Engine VCH

vSphere Integrated Containers Engine containers run as virtual machines. The VCH is not a VM, but a vApp, which is a kind of resource pool. It is an abstract dynamic resource boundary defined and controlled by vSphere into which you can provision container VMs. The VCH can be a subset of a physical host or a subset of a cluster of hosts.

A one to one coupling exists between a container and a virtual machine. A container image is attached to the VM as a disk, the VM is either booted or forked from the kernel ISO, then the containerVM chroots into the container filesystem, effectively becoming the container.

## VCH Deployment and Management ##

vSphere Integrated Containers Engine provides a command-line utility, `vic-machine`, that you use to deploy and manage VCHs. The different commands of the vic-machine utility allow you to perform the following actions:

- Deploy VCHs in configurations that are tailored to your vSphere and container development environments.
- List the VCHs that are running on a particular ESXi host or vCenter Server instance.
- Inspect, delete, and debug running VCHs.

## The Port Layer

The port layer augments the vSphere API with low level, platform-specific primitives to allow you to implement a simple container engine:

- Port Layer Execution: Handles container management, such as create, start, and stop. 
- Port Layer Interaction: Handles interaction with a running container. 
- Port Layer Networking: Handles specific vSphere network mappings into the Docker network namespace as well as mapping existing network entities such as database servers into the Docker container namespace with defined aliases.
- Port Layer Storage: Provides storage manipulation, including container image storage, layering with volume creation and manipulation. imagec, the docker registry client library, uses this component to translate registry images into a layered format that VMDK disk chains can use directly.


## Tether Process ##

The tether process is a minimal agent in the container VM that starts and stops processes and provides monitoring statistics.