# Introduction to vSphere Integrated Containers

vSphere Integrated Containers provide developers the portability, speed, and agility of using containers and provide IT Ops the management, security, and visibility they require to run workloads in production. 

VIC enables IT teams to run traditional and container workloads side-by-side on existing infrastructure seamlessly.

vSphere Integrated Containers (VIC) uses Docker containers on top of vSphere infrastructure. These containers are provisioned as virtual machines, offering the same security and functionality of virtual machines in ESX. 

From a developer's perspective, VIC is a seamless Docker interface for containers with a vSphere back end. Developers can deploy, test, and run container processes faster in the same environment as traditional applications.

From an operations perspective, VIC leverages existing technology such as vCenter and vRealize for health monitoring, performance analysis, compliance, and security to move businesses faster. IT administrators can see containers while they are running, not as black boxes.

VIC allows you to create and run multiple containers rapidly with minimal overhead using vSphere 6 Instant Clone technology, which provisions child VMs forked directly from a parent VM template running a Linux kernel. VIC creates the kernel and a few supporting resources to run containers.

You can use the Docker command line interface or the VIC plugin in the vSphere Web Client to manage containers. 

VIC consists of various components to manage, execute and monitor containers.
* Virtual Container Host
* vSphere Web Client Plugin
* Docker Client




