# Constraints of Using vSphere Integrated Containers Engine to Build Applications # 

There are some constraints on the types of containerized applications that you can deploy with this release of vSphere Integrated Containers Engine. For the lists of Docker features that this release of vSphere Integrated Containers Engine supports and does not support, see [Use and Limitations of Containers in vSphere Integrated Containers Engine](container_limitations.md). 

##  Building Container Images ##

This release of vSphere Integrated Containers Engine does not support  the `docker build` or `push` commands. As a consequence, you must use regular Docker without vSphere Integrated Containers Engine to build a container and to push it to the global hub or to your private registry server. 

## Sharing Configuration ##

This release of vSphere Integrated Containers Engine does not support  data volume sharing or `docker copy`. As a consequence, providing configuration to a containerized application has some constraints. 

An example of a configuration is the configuration files for a Web server. To pass configuration to a container when using vSphere Integrated Containers Engine, you can use the following workaround:

- Use command line arguments or environment variables. 
- Add a script to the container image that ingests the command line argument or environment variable and passes the configuration to the contained application. 

A benefit of using environment variables to transfer configuration is the containerized application closely follows the popular [12-factor application model](https://12factor.net/).

Since vSphere Integrated Containers Engine does not support sharing volumes between containers, you have the following options for processes that must share files:

- Build the files into the same image and run them in the same container.
- When containers are on the same network, add a script to the container that mounts an NFS share:
	- Run the container with an NFS server that shares a data volume.
	- Mount the NFS share on the containers that need to share files.



