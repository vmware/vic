# Add Pre- and Post-Initialization Scripts to the vSphere Integrated Containers Appliance #

The options that the vSphere Integrated Containers installer provides do not allow for all possible vSphere configurations. After installation, your vSphere environment might not allow you to use vSphere Integrated Containers in the way that you require. 

To allow vSphere Integrated Containers to run correctly in your environment, you might have to reconfigure the vSphere Integrated Containers appliance. If the reconfiguration must be persistent, you can provide scripts that perform additional configuration on the vSphere Integrated Containers appliance during its initialization process. The scripts can be any arbitrary executable that can run on the appliance, for example UNIX shell scripts. The scripts can run either before or after the initialization of the appliance, or both.

- Pre-initialization scripts run immediately after the mount of the Docker metadata disk, before the initialization and setup of the appliance.
- Post-initialization scripts run after the initialization and setup of the appliance, but before the Docker daemon starts.

**NOTE** Any pre- or post-initialization configuration that you perform via scripting is unsupported. This functionality is provided purely so that you can test the vSphere Integrated Containers technical preview in your specific vSphere environment.

For example, you might require the vSphere Integrated Containers appliance to have an additional network interface with a static IP address, so that you can expose the virtual container host to users who do not have access to the vCenter Server management network. This is not possible with the current installer options, which do not allow for multiple network interfaces. Achieving such a configuration requires you to reconfigure the vSphere Integrated Containers appliance after deployment. 

Similarly, your vSphere environment might also require TLS authentication to be persistent across reboots of the vSphere Integrated Containers appliance.   

**Prequisites**

- Deploy the vSphere Integrated Containers appliance.

 To use pre- and post-initialization scripts, you must use build 58 or later of the vSphere Integrated Containers command line installer.
- Verify that the deployment was successful and obtain the address of the appliance.

**Procedure**

1. Use SSH to log in to the vSphere Integrated Containers appliance as `root`.
3. Create folders named `pre`, or `post`, or both, under `/var/lib/docker/`.
4. Copy script files into the `pre` or `post` folders, or create the scripts directly in those folders.

 - Place pre-initialization scripts in the `pre` folder.
 - Place post-initialization scripts in the `post` folder.
 - Use `#!/bin/ash` as the script interpreter
 - Ensure that the scripts are executable.
6. Restart the vSphere Integrated Containers appliance.

**Example: Implementing Persistent TLS Authentication**

To implement TLS authentication that is persistent across reboots of the vSphere Integrated Containers appliance, you must provide a post-initialization script to create the certificates and provide them to the Docker daemon. 

1. Use SSH to log in to the vSphere Integrated Containers appliance as `root`.
2. Create a script file named `cert.sh`.
3. Copy and paste the following script into `cert.sh`.

 <pre>#! /bin/ash
HOST=$(hostname)
IPAddr=$(ip addr show eth0 | awk '$1 == "inet" {gsub(/\/.*$/, "", $2); print $2}') echo "Creating Certs for $HOST with IP: $IPAddr This script will automatically REBOOT!"
openssl genrsa -aes256 -out ca-key.pem 4096 openssl req -new -x509 -days 365 -key ca-key.pem -sha256 -out ca.pem openssl genrsa -out server-key.pem 4096 openssl req -subj "/CN=$HOST" -sha256 -new -key server-key.pem -out server.csr echo subjectAltName = IP:$IPAddr,IP:127.0.0.1 > extfile.cnf openssl x509 -req -days 365 -sha256 -in server.csr -CA ca.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile extfile.cnf openssl genrsa -out key.pem 4096 openssl req -subj '/CN=client' -new -key key.pem -out client.csr echo extendedKeyUsage = clientAuth > extfile.cnf openssl x509 -req -days 365 -sha256 -in client.csr -CA ca.pem -CAkey ca-key.pem -CAcreateserial -out cert.pem -extfile extfile.cnf echo "Setting pem permissions"
chmod -v 0400 ca-key.pem key.pem server-key.pem chmod -v 0444 ca.pem server-cert.pem cert.pem echo "Copying certs to persistant location /var/lib/docker/"
cp ca.pem /var/lib/docker/
cp server-cert.pem /var/lib/docker/
cp server-key.pem /var/lib/docker/
echo "Creating persistant CertCopy.sh script to copy certs to /root/.docker"
echo "The docker startup script checks the .docker directory and enables TLS if certs are present"
cat > /var/lib/docker/post << EOL
#! /bin/ash
mkdir /root/.docker
cp /var/lib/docker/ca.pem /root/.docker/ca.pem cp /var/lib/docker/server-cert.pem /root/.docker/server-cert.pem cp /var/lib/docker/server-key.pem /root/.docker/server-key.pem EOL echo "Modifying script to make it executable"
chmod +x /var/lib/docker/post
echo "Done!  Now copy your certs to your client and then reboot this server."
</pre>

 This script performs the following actions:
 - Creates the TLS certificates.
 - Copies the certificates into the boot-persistent `/var/lib/docker` folder.
 - Creates the `/var/lib/docker/post` folder.
 - Copies the certificates into the non-persistent `/root/.docker` folder at initialization time.
 - Makes the scripts in the `/var/lib/docker/post` folder executable.
4. Save `cert.sh`.
5. Log into your Docker client.
4. Create a folder named `.docker` in the `/root` folder.<pre>mkdir /root/.docker</pre>
4. Run the following command to copy the certificate files into the `/root/.docker` folder in your Docker client. <pre>scp ca.pem server-cert.pem server-key.pem 
root@<i>VIC_appliance_address</i>:/root/.docker/</pre>
5. Run the following command to restart the Docker daemon.<pre>ssh root@<i>VIC_appliance_address</i> /opt/dockerd.sh</pre>
6. Run the following command in a Docker client terminal to connect the Docker client to the virtual container host with TLS authentication.<pre>docker --tlsverify -H tcp://<i>VIC_appliance_address</i>:2376 ps</pre>