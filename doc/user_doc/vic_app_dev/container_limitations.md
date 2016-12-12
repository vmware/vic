# Use and Limitations of Containers in vSphere Integrated Containers Engine

vSphere Integrated Containers Engine currently includes the following capabilities and limitations:

## Supported Docker Features
This version of vSphere Integrated Containers Engine supports these features:

- Docker Compose (basic)
- Registry pull from docker hub and private registry
- Named Data Volumes
- Anonymous Data Volumes
- Bridged Networks
- External Networks
- Port Mapping
- Network Links/Alias

## Limitations
vSphere Integrated Containers Engine includes these limitations:

- Container VMs only support root user.
- When you do not configure a PATH environment variable, or create a container from an image that does not supply a PATH, vSphere Integrated Containers Engine provides a default PATH.
- You can resolve the symbolic names of a container from within another container, except in the following cases:
	- Aliases
	- IPv6
	- Service discovery
- Containers can acquire DHCP addresses only if they are on a network that has DHCP.

## Unsupported Docker Features

This version of vSphere Integrated Containers Engine does not support these features:

- Pulling images via image digest 
- Pushing a registry
- Sharing concurrent data volume between containers
- Mapping a local host folder to a container volume
- Mapping a local host file to a container
- Docker build
- Docker copy files into a container, both running and stopped
- Docker container inspect does not return all container network for a container

For limitations of using vSphere Integrated Containers with volumes, see [Using Volumes with vSphere Integrated Containers Engine](using_volumes_with_vic.md).