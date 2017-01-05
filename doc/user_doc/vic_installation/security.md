# Securing VCH Connections #

vSphere Integrated Containers Engine provides different possibilities for securing the connections between virtual container hosts (VCHs) and different elements in your environment.

- [Connections to Docker clients](#docker)
- [Connections to Registry Servers for Container Images](#registries)
- [Connections to and Downloads from the VCH Admin Portal](#vicadmin)
- [Connections to vSphere](#vsphere)

<a name="docker"></a>
## Connections to Docker Clients ##

Docker uses TLS authentication to allow for different levels of security when connecting Docker clients to container hosts. For information about Docker certificates and security, see [the Docker security documentation](https://docs.docker.com/engine/security/https/).

vSphere Integrated Containers Engine allows you to deploy VCHs that implement the following standard Docker authentication options:

- Mutual authentication with client and server certificates
- Server authentication with no client authentication
- No TLS authentication

You specify the type of authentication to use when you deploy a VCH.

### Two-Way TLS Authentication with Trusted Certificates ###

Docker allows for two-way TLS certificate authentication between the Docker client and the container host, with trusted server and client certificates that are both signed by a Certificate Authority (CA). 

With vSphere Integrated Containers Engine, you can use auto-generated trusted certificates that the VCH deployment utility, `vic-machine`, creates for you. You can also use custom trusted certificates that you create yourself.

If you implement two-way authentication with trusted certificates, container developers must have copies of the certificate and key files and must use the `--tlsverify --tlscacert --tlscert --tlskey` options when running Docker commands against the VCH. 

**Auto-Generated Trusted Certificates**

To deploy a VCH to use auto-generated, trusted certificates that are signed by a CA, you run `vic-machine create` with  options that populate the different fields of the certificate application form. The `vic-machine create` command uses this information when it generates the certificates and sends them to a CA for signature. For information about deploying VCHs that implement two-way TLS authentication with auto-generated trusted certificates, see the description of the `vic-machine create` [`--tls-cname` option](vch_installer_options.md#tls-cname).

Alternatively, if you assign a static address on the client network to a VCH, `vic-machine create` uses this address to auto-generate trusted CA certificates. In this case, two-way TLS authentication with trusted certificates is implemented by default, and you do not need to perform any additional TLS configuration when you deploy the VCH. For information about creating a VCH with a static IP address, see [Options for Specifying a Static IP Address for the VCH Endpoint VM](vch_installer_options.md#static-ip).

**Custom Trusted Certificates**
 
To deploy a VCH to use trusted certificates that you created yourself and have had signed by a CA, run `vic-machine create` with the `--cert` and `--key` options. The `--cert` and `--key` options allow you to specify the paths to existing certificate and key files, that `vic-machine create` uploads to the VCH. 

For information about creating a VCH that uses custom trusted certificates, see the descriptions of the `vic-machine create` [`--cert`](vch_installer_options.md#cert) and [`--key`](vch_installer_options.md#key) options.

### Server-Side Authentication with Untrusted Self-Signed Certificates ###

Docker allows for TLS authentication of the container host with untrusted, self-signed certificates that are not verified by a CA. Container hosts are verified, but Docker clients are not.

If you implement TLS authentication of the VCH with untrusted, self-signed certificates, container developers do not require copies of the certificate and key files. They must use the `--tls` option when running Docker commands against the VCH.

To deploy a VCH that uses auto-generated, self-signed untrusted certificates to verify the VCH, you disable TLS client verification when you run `vic-machine create`. 

For information about creating a VCH that uses auto-generated, self-signed untrusted certificates for verification of the VCH only, see the description of the `vic-machine create` [`--no-tlsverify` option](vch_installer_options.md#no-tlsverify).  

### No TLS Authentication ###

You can disable TLS authentication completely when you deploy a VCH. In this case, neither the Docker client nor the VCH are verified and no certificate files are created.

For information about creating a VCH that does not authenticate connections, see the description of the `vic-machine create` [`--no-tls` option](vch_installer_options.md#no-tls).

<a name="registries"></a>
## Connections to Private Registry Servers for Container Images ##

You can configure VCHs to connect to insecure  private registry servers to download container images. The current version of vSphere Integrated Containers Engine does not support connections to private registry servers that are secured by TLS authentication.

<a name="vicadmin"></a>
## Connections to and Downloads from the VCH Admin Portal ##

If you configured a VCH to use full TLS authentication with trusted certificates, connections to the VCH Admin Portal are authenticated by a browser-friendly PFX client certificate, `cert.pfx` that `vic-machine create` generates during deployment of the VCH. To access the VCH Admin Portal for a VCH that implements full TLS authentication with trusted certificates, you must provide the PFX certificate to the browser that you are using to access the portal. If you deployed the VCH with anything less than full TLS authentication, anyone can access the VCH Admin Portal for that VCH.

<a name="vsphere"></a>
## Connections to vSphere ##

Connections from VCHs to vCenter Server and ESXi hosts are authenticated by using the vCenter Server and ESXi host certificates. If your vSphere environment uses untrusted, self-signed certificates you must provide the thumbprint of the vCenter Server or ESXi host certificate when you deploy a VCH. If your vSphere environment uses trusted certificates signed by a known CA, no action is required when you deploy a VCH. All  connections between the VCH and the vCenter Server instance or ESXi host are authenticated by using the vCenter Server or ESXi host certificate.