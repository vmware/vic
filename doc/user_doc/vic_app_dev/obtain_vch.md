# Obtain a Virtual Container Host #

vSphere Integrated Containers Engine does not provide an automated means of obtaining virtual container hosts to container developers who are not vSphere administrators.

When you deploy a virtual container host, the virtual container host obtains an IP address from DHCP. This IP address is required by whoever will use the virtual container host as their Docker end-point. Depending on the nature of your organization, the consumer of the Docker end-point might be the person or team who deployed the virtual container host, or might be a different person or team. 