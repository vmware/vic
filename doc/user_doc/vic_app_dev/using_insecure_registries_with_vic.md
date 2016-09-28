# Using Insecure Registry Servers with vSphere Integrated Containers Engine #

An insecure registry server is a registry server that is secured by self-signed certificates rather than by TLS.

If your Docker environment stores Docker images in an insecure registry server, you or the vSphere administrator must have set the `vic-machine  create --docker-insecure-registry` option when creating the virtual container host. Setting the `docker-insecure-registry` option on a virtual container host informs that virtual container host that it is permitted to pull images from the designated insecure registry server. 

For information about how to use the `vic-machine create --docker-insecure-registry` option, see the section on `docker-insecure-registry` in [Virtual Container Host Deployment Options](../vic_installation/vch_installer_options.html#registry) in *vSphere Integrated Containers Engine Installation*. 

## Pull a Container Image from an Insecure Registry ##

To pull a container image from an insecure registry server, run the following Docker command.
<pre>docker -H <i>vch_address</i>:2376 --tls 
pull <i>vch_address</i>/path/to/image/<i>image_name</i>:<i>image_version</i></pre>
