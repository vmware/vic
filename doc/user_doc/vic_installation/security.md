# vSphere Integrated Containers Engine Security Overview #

The security requirements for vSphere Integrated Containers Engine relate to connections from virtual container hosts to different elements in your environment.

- [Connections to vSphere](#vsphere)
- [Connections to Docker clients](#docker)
- [Connections to Registry Servers for Docker Images](#registries)
- [Connections to and Downloads from the VIC Admin Portal](#vicadmin)

<a name="vsphere"></a>
## Connections to vSphere ##

If your vSphere environment uses default certificates that are generated, signed, and managed by the VMware Certificate Authority (VMCA), vSphere Integrated Containers uses auto-generated VMCA certificates to authenticate connections to vSphere.

If your vSphere environment uses custom certificates that are signed by a Certificate Authority (CA), you must configure virtual container hosts to use these certificates. 
- For information about configuring virtual container hosts to use custom CA certificates when connecting to vSphere, see the section on the [`thumbprint` option](vch_installer_options.md) in Virtual Container Host Deployment Options.
- For information about using custom certificates with vSphere, see [vSphere Security Certificates](https://pubs.vmware.com/vsphere-60/topic/com.vmware.vsphere.security.doc/GUID-779A011D-B2DD-49BE-B0B9-6D73ECF99864.html) in the vSphere documentation.

<a name="docker"></a>
## Connections to Docker Clients ##

Docker uses TLS authentication to allow for different levels of security when connecting Docker clients to container hosts:

- No TLS authentication.
- TLS certificate authentication between the Docker client and the container host, with untrusted, self-signed certificates that are not verified by a CA.
- TLS certificate authentication between the Docker client and the container host, with trusted certificates that are verified by a CA.

For more information about Docker certificates and security, see [the Docker security documentation](https://docs.docker.com/engine/security/https/).

By default, virtual container hosts automatically generate TLS certificates and use Docker `TLSVERIFY` to verify the certificates with a CA.

- You can configure a virtual container host to authenticate client connections by using a custom trusted certificates that are signed by a Certificate Authority (CA) instead of using an auto-generated  self-signed certificates. This corresponds to using `docker --tlsverify --tlscacert -- tlscert --tlskey`.
- You can configure a virtual container host so that it uses either auto-generated or custom certificates that are untrusted and are not verified by a CA.  This corresponds to using `docker --tls --tlscert --tlskey`.
- You can also disable authentication completely, to allow any Docker client to connect to the virtual container host.

<a name="registries"></a>
## Connections to Registry Servers for Docker Images ##

xccx

<a name="vicadmin"></a>
## Connections to and Downloads from the VIC Admin Portal ##

xccx