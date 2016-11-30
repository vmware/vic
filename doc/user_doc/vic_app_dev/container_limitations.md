# Use and Limitations of Containers in vSphere Integrated Containers Engine

vSphere Integrated Containers Engine currently includes the following  capabilities and limitations:

- Container VMs only support root user.
- You can resolve the symbolic names of a container from within another container except for the following:
	- aliases
	- IPv6 support
	- service discovery
- Containers are capable of acquiring DHCP addresses if they are on a network that has DHCP.
- This version of vSphere Integrated Containers Engine does not support pulling images via image digest.
