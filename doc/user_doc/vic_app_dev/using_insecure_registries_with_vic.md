# Using Insecure Private Registry Servers with vSphere Integrated Containers Engine #

An insecure private registry server is a private registry server that is secured by self-signed certificates rather than by TLS.

If your Docker environment stores Docker images in an insecure private registry server, you or the vSphere administrator must have set the `vic-machine  create --docker-insecure-registry` option when creating the virtual container host. Setting the `docker-insecure-registry` option on a virtual container host informs that virtual container host that it is permitted to pull images from the designated insecure registry server. 

For information about how to use the `vic-machine create --docker-insecure-registry` option, see the section on `docker-insecure-registry` in [Virtual Container Host Deployment Options](../vic_installation/vch_installer_options.html#registry) in *vSphere Integrated Containers Engine Installation*. 

## Pull a Container Image from an Insecure Private Registry Server ##

To pull a container image from an insecure private registry server, run the following Docker command. 
<pre>docker -H <i>vch_address</i>:2376 --tls 
pull <i>registry_server_address</i>/path/to/image/<i>image_name</i>:<i>image_version</i></pre>

If the private registry server listens for connections on a specific port, include the port number in the registry server URL.

<pre>docker -H <i>vch_address</i>:2376 --tls 
pull <i>registry_server_address</i>:<i>port_number</i>/path/to/image/<i>image_name</i>:<i>image_version</i></pre>

**NOTE**: The current builds of vSphere Integrated Containers do not yet support private registry servers that you secure by using TLS certificates.
