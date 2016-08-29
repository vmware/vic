# Using TLS Authentication with vSphere Integrated Containers Engine #

You can use vSphere Integrated Containers Engine in environments that require  trusted TLS certificates to authenticate connections between Docker clients and the Docker daemon in the virtual container hosts. 

**Prerequisites**

Familiarize yourself with the instructions for protecting the Docker daemon socket at https://docs.docker.com/engine/security/https/. 

**Procedure**

1. Use SSH to log in to the vSphere Integrated Containers Engine appliance as `root`.
2. Change the shell interpreter to `ash`.<pre>sudo /bin/ash</pre>
2. Follow the instructions at https://docs.docker.com/engine/security/https/ to create TLS certificates on the virtual container host.
3. Copy the certificates into the `/root/.docker` folder on the vSphere Integrated Containers Engine appliance.
3. Log into your Docker client host machine.
4. Create a folder named `.docker` in the `/root` folder.<pre>mkdir /root/.docker</pre>
4. Run the following commands to copy the certificate files into the `/root/.docker` folder in your Docker client. <pre>scp root@<i>VIC_appliance_address</i>:/root/.docker/ca.pem /root/.docker</pre>
<pre>scp root@<i>VIC_appliance_address</i>:/root/.docker/server-cert.pem /root/.docker</pre>
<pre>scp root@<i>VIC_appliance_address</i>:/root/.docker/server-key.pem  /root/.docker</pre>
5. Run the following command in a Docker client terminal to restart the Docker daemon in the vSphere Integrated Containers Engine appliance.<pre>ssh root@<i>VIC_appliance_address</i> /opt/dockerd.sh</pre>
6. Run the following command in a Docker client terminal to restart the Docker client.<pre>docker-machine restart</pre>
6. Run the following command in a Docker client terminal to connect the Docker client to the virtual container host with TLS authentication.<pre>docker --tlsverify -H tcp://<i>VIC_appliance_address</i>:2376 ps</pre>

Docker now uses TLS authentication for all connections between the client and the virtual container host.

**NOTE**: The `/root/.docker` folder is rebuilt with each reboot of the vSphere Integrated Containers Engine appliance. As a consequence, the configuration in this procedure will not persist across a reboot of the  appliance. To make TLS authentication persistent across reboots, you must provide the certificates by using  post-initializaton scripting.