# Using Private Registry Servers with vSphere Integrated Containers Engine #

If your development environment includes private registry servers for container images, you or the vSphere Administrator must correctly configure virtual container hosts (VCHs) to allow them to connect to the private registry servers.

You can use vSphere Integrated Containers Engine with either secure or insecure private registry servers.

## Secure Private Registry Servers ##

If the private registry server is secured by a Certificate Authority (CA) certificate and private key, you or the vSphere Administrator must have deployed the VCH with the appropriate certificates. You deploy a VCH with private registry server CA certificates by specifying the `--registry-ca` option when you deploy the VCH. 

For information about how to configure a VCH to use private registry server CA certificates, see the section on `--registry-ca` in [VCH Deployment Options](../vic_installation/vch_installer_options.html#registry-ca) in *vSphere Integrated Containers Engine Installation*. 

## Insecure Private Registry Servers ##

If you set up a private registry that does not use certificates, you or the vSphere administrator must deploy the VCH with the `vic-machine create --insecure-registry` option. Setting the `insecure-registry` option on a VCH informs that VCH that it is authorized to pull images from the designated insecure private registry server.

If you authorize a VCH to connect to an insecure private registry server, the VCH attempts to access the registry server via HTTP if access via HTTPS fails. VCHs always use HTTPS when connecting to registry servers for which you have not authorized insecure access. Insecure private registries are not recommended in production environments.

For information about how to use the `vic-machine create --insecure-registry` option, see the section on `insecure-registry` in [VCH Deployment Options](../vic_installation/vch_installer_options.html#insecure-registry) in *vSphere Integrated Containers Engine Installation*. 

## Pull a Container Image from a Private Registry Server ##

To pull a container image from a private registry server, run the following Docker command. 

<pre>docker -H <i>vch_address</i>:2376 --tls 
pull <i>registry_server_address</i>/path/to/image/<i>image_name</i>:<i>image_version</i></pre>

If the private registry server listens for connections on a specific port, include the port number in the registry server URL.

<pre>docker -H <i>vch_address</i>:2376 --tls 
pull <i>registry_server_address</i>:<i>port_number</i>/path/to/image/<i>image_name</i>:<i>image_version</i></pre>

These commands will only work in the following circumstances:

- The private registry server at <i>registry_server_address</i> is secured by CA certificates, and you or the vSphere Administrator passed the appropriate certificates to the VCH during deployment by using the `--registry-ca` option.
- The private registry server at <i>registry_server_address</i> is not secured by certificates, and you or the vSphere Administrator authorized access to this registry server by using the `--insecure-registry` option during VCH deployment.

**NOTE**: In the examples, the Docker commands specify `--tls`. This is to specify that the connection between the Docker client and the VCH is secured by TLS. The level of security of the connection between the Docker client and the VCH is completely independent from the level of security of the connection between the VCH and the private registry server. The connection to the private registry server can be insecure when the connection between the client and the VCH is secure, and the reverse.
