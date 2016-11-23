# vSphere Integrated Containers Engine Security Reference
Use the Security Reference to learn about the security features of vSphere Integrated Containers Engine.

- [List of External Interfaces, Ports, and Services](#list_open_ports)
- [vSphere Integrated Containers Service Accounts and Privileges](#accounts)
- [Apply Security Updates and Patches](#patches)
- [Security Related Log Messages](#logs)

<a name="list_open_ports"></a>
### List of External Interfaces, Ports, and Services

The following ports must be open on the virtual container host appliance:<br>

- 2375 if no TLS is deployed
- 2376 for TLS
- 2377 to configure or disable the firewall
- 2378 vic-admin
- 53 DNS server
- 22 for SSH for debugging purposes

<a name="accounts"></a>
### vSphere Integrated Containers Engine Service Accounts and Privileges
vSphere Integrated Containers Engine does not create service accounts and does not assign privileges. Instead, it creates a vSphere Extension and authenticates against it.

<a name="patches"></a>
### Apply Security Updates and Patches
Download a new version of vSphere Integrated Containers and run `vic-machine upgrade`.

<a name="logs)"></a>
### Security Related Log Messages
vSphere Integrated Containers security related information appears in docker-personality.log and vicadmin.log that you can access from the VCH Admin portal.