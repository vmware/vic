# Obtain a VCH #

vSphere Integrated Containers Engine does not currently provide an automated means of obtaining virtual container hosts (VCHs).

When you or the vSphere administrator use `vic-machine create` to deploy a VCH, the VCH endpoint VM obtains an IP address. The IP address can either be static or be obtained from DHCP. As a container developer, you require the IP address of the VCH endpoint VM when you run Docker commands. 

Depending on the nature of your organization, you might deploy VCHs yourself, or you might request a VCH from a different person or team. If you do not run `vic-machine create` yourself, your organization must define the process by which you obtain VCH addresses. This process can be as simple as an exchange of emails with a vSphere administrator, or as advanced as a custom self-provisioning portal or API end-point. For example, your organization could use VMware vRealize&reg; Automation&trade; to provide a self-provisioning service, by using the vRealize Automation interface or APIs to request VCHs. At the end of the provisioning process, vRealize Automation would communicate the VCH endpoint VM address to you.

## Using Docker Environment Variables ##

If you or the vSphere administrator deploy VCHs with TLS authentication, `vic-machine create` generates a file named `vch_address.env`. The `env` file contains Docker environment variables that are specific to the VCH. You can use the contents of the `env` file to set environment variables in your Docker client. A self-provisioning service such as vRealize Automation could potentially provide the env file at the end of the provisioning process for VCHs.

## Connecting to the VCH ##

How you connect to your VCH depends on the security options with which  you or the vSphere administrator deployed the VCH. 

- If the VCH uses TLS authentication, either by using server certificates or by using mutual authentication with client and server certificates, you connect to the VCH at *vch_address*:2376.  
- If the VCH uses mutual authentication with client and server certificates, you must configure the Docker client appropriately with one of the following options:
  - By using the following `--tlsverify`, `--tlscert`, and `--tlskey` Docker options, adding `tlscacert` if a custom CA was used to sign the server certificate.
  - By setting `DOCKER_CERT_PATH=/path/to/client/cert.pem` and `DOCKER_TLS_VERIFY=1`.
- If the VCH uses server certificates without client authentication, you run Docker commands with the `--tls` option. The `DOCKER_TLS_VERIFY` environment variable must not be set. Note that setting `DOCKER_TLS_VERIFY` to 0 or `false` has no effect. 
- If TLS is completely disabled on the VCH, you connect to the VCH at *vch_address*:2375 and do not need to specify any additional Docker options. 