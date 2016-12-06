# Obtain a VCH #

vSphere Integrated Containers Engine does not currently provide an automated means of obtaining virtual container hosts (VCHs).

When you or the vSphere Administrator use `vic-machine create` to deploy a VCH, the VCH endpoint VM obtains an IP address. The IP address can either be static or be obtained from DHCP. As a container developer, you require the IP address of the VCH endpoint VM when you run Docker commands. 

Depending on the nature of your organization, you might deploy VCHs yourself, or you might request a VCH from a different person or team. If you do not run `vic-machine create` yourself, your organization must define the process by which you obtain VCH addresses. This process can be as simple as an exchange of emails with a vSphere Administrator, or as advanced as a custom self-provisioning portal or API end-point. For example, your organization could use VMware vRealize&reg; Automation&trade; to provide a self-provisioning service, by using the vRealize Automation interface or APIs to request VCHs. At the end of the provisioning process, vRealize Automation would communicate the VCH endpoint VM address to you.