## Virtual Container Host (VCH)

The Virtual Container Host (VCH) provides a virtual container endpoint backed by a vSphere resource pool that allows you to control and consume container services. 

You can access a Docker API endpoint for development and map ports for client connections to run containers as required.

vSphere resource management handles container placement within VCH, so that a virtual Docker host can serve as an entire vSphere cluster or a fraction of the same cluster. The only resource consumed by a container host in the cluster is the resource consumed by running containers.

You can reconfigure the VCH with no impact to containers running in it. VCH imposes no conceptual limitations on the kernel version or on the operating system the containers are running.  

You can deploy multiple VCHs in an environment, depending on your business needs, including allocating separate resources for development, testing, and production. 
You can also nest VCHs, giving your team access to a large VCH, or sub-allocate smaller VCHs for individuals.

Each VCH maintains a cache of container images, which you download from either the public Docker Hub or a private registry. VCH maintains filesystem layers inherent in container images by mapping to discrete VMDK files, all of which are housed in vSphere datastores on VSAN, NFS, or local disks.

You deploy a VCH using the CLI installer, then access VCH endpoints remotely through a Docker command line interface or other API client.




