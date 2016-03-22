# Using TLS Authentication with vSphere Integrated Containers #

You can use vSphere Integrated Containers in environments that require  trusted TLS certificates to authenticate connections between Docker clients and the Docker daemon in the virtual container hosts. 

**Prerequisites**

Familiarize yourself with the instructions for protecting the Docker daemon socket at https://docs.docker.com/engine/security/https/. 

**Procedure**

1. Use SSH to log in to the vSphere Integrated Containers appliance.
2. Follow the instructions at https://docs.docker.com/engine/security/https/ to create TLS certificates on the virtual container host.
3. Log into your Docker client.
4. Run the following command to copy the certificate files into the `.docker` folder in your docker client. <pre>scp ca.pem server-cert.pem server-key.pem 
root@<i>VIC_appliance_address</i>:/root/.docker/ 
&& ssh root@vch /opt/dockerd.sh</pre>
5. Restart the Docker daemon.

After you restart the Docker daemon, Docker will use TLS authentication for all connections between the client and the virtual container host. 


