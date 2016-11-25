# Obtain a VCH #

vSphere Integrated Containers Engine does not currently provide an automated means of obtaining virtual container hosts (VCHs).

When you or the vSphere Administrator use `vic-machine create` to deploy a VCH, the VCH endpoint VM obtains an IP address. The IP address can either be static or be obtained from DHCP. As a container developer, you require the IP address of the VCH endpoint VM when you run Docker commands. 

Depending on the nature of your organization, you might deploy VCHs yourself, or you might request a VCH from a different person or team. If you do not run `vic-machine create` yourself, your organization must define the process by which you obtain VCH addresses. This process can be as simple as an exchange of emails with a vSphere Administrator, or as advanced as a custom self-provisioning portal or API end-point. For example, your organization could use VMware vRealize&reg; Automation&trade; to provide a self-provisioning service. In this case, you would use the vRealize Automation interface or APIs to request a VCH. At the end of the provisioning process, vRealize Automation would communicate the VCH endpoint VM address to you.

## Using Docker Environment Variables ##

If you or the vSphere Administrator deploy the VCHs   with TLS authentication, either with trusted certificates or with untrusted self-signed certificates, `vic-machine create` generates a  file named `vch_address.env`. The `env` file contains Docker environment variables that are specific to the VCH. You can use the contents of the `env` file to set environment variables in your Docker client so that it connects to the correct VCH and uses the appropriate level of authentication. A self-provisioning service such as vRealize Automation could potentially provide the `env` file at the end of the provisioning process for VCHs.