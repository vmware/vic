**VICe Overview for vSphere Admins**

**Introduction to Containers, Images and Volumes**

The word Container is overloaded one these days. It’s helpful to distinguish the *runtime* aspect from the *packaging* in understanding containers and how they relate to VIC.

*Runtime*

At its most basic, a container is simply a sandbox in which a process can be run. The sandbox isolates the process from other processes running on the same system. A container has a lifecycle which is typically tied to lifecycle of the process it’s designed to run. If you start a container, it starts its main process and when that process ends, the container stops. The container may have access to some storage and it will typically have an identity on a network.

Conceptually a container represents many of the same capabilities as a VM. The main difference between the two is the abstraction layer

* A software container is a sandbox within a guest OS and it is up to the guest to provide the container with its dependencies and to enforce isolation. Multiple containers share the guest kernel, networking and storage. A container does not boot - it is simply a slice of an already-running OS. The OS running the container is called its *host*.

* A VM on the other hand is a sandbox within a hypervisor and it is up to the hypervisor to provide it with dependencies, such as virtual disks and NICs. A VM has to boot an OS and its lifecycle is typically tied to that of the OS rather than any one process. A VM is strongly isolated from other VMs and its host by design.

One of the most interesting facets of containers is how they deal with state. Any data written by a container will be non-persistent by default and will be lost when a container is deleted. State however can be persisted beyond the lifespan of a container by attaching a *Volume* or sending it over a network. Binary dependencies needed by the container - such as OS libraries or application binaries - are encapsulated in *Images* which are immutable.

*Packaging*

One of the most significant benefits of containers is that they allow you to package up an entire environment needed for an application and run it anywhere. You can go to Docker Hub today, pick any one of hundreds of thousands of applications and run it anywhere that Docker is installed on a compatible OS. The packaging encapsulates the binary dependencies, environment variables, volumes even network configuration. 

The format of this packaging, is called an *Image.* An image is a template from which many containers can be instantiated. The Docker image format allows for images to be composed in a parent child-relationship, just like a disk snapshot. This image hierarchy allows for common dependencies to be shared between containers. For example, I may have a Debian 8 image which has a child image with Java installed. That Java image may itself have a child with Tomcat installed. The Debian 8 image may have other children, such as PHP, Python etc. 

The immutability of the image format means that you never modify an image, you always create a new one. The layered nature of the image format means that you can cache commonly-used layers and only have to download/upload layers you don’t already have. It also means that if you want to patch a particular image, you create a new one and then rebuild all of its children. 

The beauty of the image format is its portability. So long as you have a destination running a container engine (such as Docker) and a compatible kernel ABI, then you can download and run an image on it. This portability is facilitated by a *Registry* which is a service that indexes and stores images. You can run your own private image registry that forms part of a development pipeline: Images are *pushed* to the registry from development, *pulled* to a test environment, verified and can then be *pulled* to a production environment.

**What is VIC**

The VIC product portfolio includes a container engine (VICe), management portal (Admiral) and container registry (Harbor) that all currently support the Docker image format. VIC is entirely Open Source, free to use and support is included in the vSphere Enterprise Plus license.

This guide is focussed primarily on VIC engine. Admiral ties together image registries, container deployment endpoints and application templates into a UI portal. Harbor is a Docker image registry with additional capabilities such as RBAC, replication etc.

VIC engine is designed to integrate all the packaging and runtime benefits of containers with the enterprise capabilities of your vSphere environment. 

It is designed to solve many of the challenges associated with putting containerized applications into production by directly leveraging the clustering, dynamic scheduling and virtualized infrastructure in vSphere and bypassing the need to maintain discrete Linux VMs as container hosts.

It allows you as a vSphere administrator to provide a container management endpoint to a user as a service, while remaining in complete control over the infrastructure it depends on.

* vSphere is the container host, not Linux

    * Containers are spun up *as* VMs not *in* VMs

    * Every container is fully isolated from the host and from each other

    * Provides per-tenant dynamic resource limits within a vCenter cluster

* vSphere is the infrastructure, not Linux

    * vSphere networks may be selected to appear in the Docker client as container networks

    * Images, volumes, and container state are provisioned directly to VMFS

* vSphere is the control plane

    * Use the Docker client to directly control selected elements of vSphere infrastructure

    * A container endpoint is Service-as-a-Servicepresents as a service abstraction, not as IaaS

VICe is designed to be the fastest and easiest way to provision any Linux-based workload to vSphere, provided that workload can be serialized as a Docker image.

**What VICe Does**

VICe gives you, the vSphere admin, the tools to easily make your vSphere infrastructure accessible to users for provisioning container workloads into production.

Imagine this scenario. A user raises a ticket and says, "I need Docker". Your response to that is to provision a large Linux VM and send them the IP address. Your user is now on the hook for installing Docker, patching the OS, configuring in-guest network and storage virtualization, securing the guest, isolating containers, packing containers efficiently, managing upgrades and downtime. You’ve just given your user something akin to a nested hypervisor that they have to manage and which is opaque to you. If you scale that up to one large Linux VM per tenant, you’ll end up creating a large distributed silo for containers.

With VIC, the conversation is different. A user raises a ticket and says, "I need Docker". Your response is to set aside a certain amount of storage, networking and compute on your cluster and formalize that using a tool called vic-machine. This tool installs a small appliance that runs a secure remote Docker API and which represents authorization to the infrastructure you’ve assigned, into which the user can self-provision container workloads. Instead of sending your user a VM and saying “good luck”, you instead send them the IP address of the appliance, the port of the remote API and a certificate for secure access. This secure remote API is the only access the user has to the infrastructure - you’ve provided a service portal. This is great for the user because they don’t have to worry about isolation, patching, security, backup etc. It’s great for you because every container that’s spun up is a VM (“containerVM”) and can be vMotioned and monitored just like all your other VMs.

Let’s say the user needs more compute capacity. With the first option, the pragmatic choice is to power down the VM and reconfigure it, or give them a new VM and have them deal with the clustering implications. Both these scenarios are disruptive to the user. With VIC, it’s a simple reconfiguration that’s completely transparent to the user.

VICe allows you to select and dictate appropriate infrastructure for the task in hand. When it comes to networking, you can select multiple Port Groups for different types of network traffic, ensuring that all containers provisioned by a user get appropriate interfaces on the right networks. With storage, you can select different vSphere datastores for different types of state. For example, Container state is ephemeral and is unlikely to need to be backed up, but Volume state almost certainly should be backed up. VIC automatically ensures that state gets written to the appropriate datastore when the user provisions a container.

So, to summarize, VIC gives you a mechanism to allow users to self-provision VMs as containers into your virtual infrastructure.

**What Is VICe For?**

VICe currently offers a subset of the Docker API and is designed to specifically address the provisioning of containers into production, solving many of the problems highlighted above.

VICe exploits the portability of the Docker Image format to present itself as an enterprise deployment target. Containers get built on one system and pushed to a registry. They get tested by another system and blessed for production. VIC can then pull them out of the registry and deploy them to vSphere.

**VICe Concepts**

If we consider a Venn diagram of "What vSphere Does" in one circle and “What Docker Does” in another, the intersection is not insignificant. The vision of VIC is to take as much of vSphere as possible and layer on whatever Docker capabilities are missing, reusing as much of Docker’s own code as possible. The end result should not sacrifice portability of the Docker image format and should be completely transparent to a Docker client. Below are the key concepts and components that make this possible.

![VICe components](graphics/vice-components.png)

*The "ContainerVM"*

VMs created by VIC engine have all of the *characteristics* of software containers:

* Ephemeral storage layer with optionally attached persistent "volumes"

* Custom Linux guest - designed to be "just a kernel" - needs “images” to be functional

* Mechanism for persisting and attaching read-only binary "image layers"

* PID 1 guest agent "Tether" extends the control plane into the containerVM

* Various well-defined methods of configuration and state ingress / egress

* Automatically configured to various network topologies

To be clear, the provisioned VM does not contain any OS container abstraction. The VM boots from an ISO containing the Photon Linux kernel and is configured with container images mounted as a disk. Container image layers are represented as a read-only VMDK snapshot hierarchy on a vSphere datastore. At the top of this hierarchy is a read-write snapshot that stores ephemeral state. Container volumes are formatted VMDKs attached as disks and indexed on a datastore. Networks are distributed port groups attached as vNICs.

*The "Virtual Container Host"*

A VCH is the virtual functional equivalent of a Linux VM running Docker, but with some significant benefits. A VCH represents a clustered pool of resource into which containerVMs can be provisioned, a single-tenant container namespace and a secure API endpoint. It also represents authorization to use and configure pre-approved virtual infrastructure.

A VCH is functionally distinct from a traditional container host in the following ways:

* It naturally encapsulates clustering and dynamic scheduling by provisioning to vSphere targets

* The resource constraints are dynamically configurable with no impact to the containers

* The containers don’t share a kernel

* There is no local image cache. This is kept on a datastore somewhere in the cluster

* There is no read-write shared storage

A VCH is a multi-functional appliance deployed as a vApp. The vApp provides a useful visual parent/child relationship in the vSphere UI for containerVMs provisioned into the VCH and also has the capability to specify resource limits. Multiple VCHs can be provisioned onto a single ESXi host, to a vSphere resource pool or to a vSphere cluster.

*The VIC Appliance*

There is a 1:1 relationship between a VCH and a VIC appliance. It has the following functions:

* Run the services that are necessary for a VCH

* Provide a secure remote API to a client

* Provide network forwarding such that ports to containers can be opened on the appliance and the containers can access a public network

* Manage the lifecycle of containers, image store, volume store and container state

* Provide logging and monitoring of its own services and of its containers

The lifecycle of the VIC appliance is managed by a binary called VIC machine which installs, upgrades, deletes and enables debugging for a VIC appliance.

*VIC machine*

VIC machine is a binary built for Windows, Linux and OSX that manages the lifecycle of VCHs. VIC machine has been designed to be used by vSphere admins. It takes pre-existing compute, network, storage and a vSphere user as input and creates a VCH as output. It has the following additional functions:

* Creates certificates for Docker client TLS authentication

* Checks that prerequisites have been met on the cluster (firewall, licenses etc)

* Configures a running VCH for debugging

* Lists, upgrades and deletes VCHs



