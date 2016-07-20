# Container Operations and Limitations in vSphere Integrated Containers

vSphere Integrated Containers currently include the following limitations:

- Container VMs only support root user
- You can resolve the symbolic names of a container from within another container except for the following:
	- aliases
	- IPv6 support
	- service discovery
- Containers are capable of acquiring DHCP addresses if they are on a network that has DHCP.