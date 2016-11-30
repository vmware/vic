# Use and Limitations of Containers in vSphere Integrated Containers Engine

vSphere Integrated Containers Engine currently includes the following  capabilities and limitations:

- Container VMs only support root user.
- You can resolve the symbolic names of a container from within another container, except in the following cases:
	- Aliases
	- IPv6
	- Service discovery
- Containers can acquire DHCP addresses only if they are on a network that has DHCP.
- This version of vSphere Integrated Containers Engine does not support pulling images via image digest.

For limitations of using vSphere Integrated Containers with volumes, see [Using Volumes with vSphere Integrated Containers Engine](using_volumes_with_vic.md).
