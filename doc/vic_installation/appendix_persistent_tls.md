# Appendix: Example of Implementing Persistent TLS Authentication #

To implement TLS authentication that is persistent across reboots of the vSphere Integrated Containers appliance, you provide a post-initialization script to create the certificates and provide them to the Docker daemon. 

For information about how to implement TLS authentication in the Docker daemon, see https://docs.docker.com/engine/security/https/.

**NOTE**: The scripting in this appendix is for example purposes and is unsupported.

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
 - Creates the `/var/lib/docker/post` executable.
 - Copies the certificates into the non-persistent `/root/.docker` folder at initialization time.
 - Makes the `/var/lib/docker/post` file executable.
4. Save `cert.sh`.
5. Log into your Docker client.
4. Create a folder named `.docker` in the `/root` folder.<pre>mkdir /root/.docker</pre>
4. Run the following commands to copy the certificate files into the `/root/.docker` folder in your Docker client. <pre>scp root@<i>VIC_appliance_address</i>:/root/.docker/ca.pem /root/.docker</pre>
<pre>scp root@<i>VIC_appliance_address</i>:/root/.docker/server-cert.pem /root/.docker</pre>
<pre>scp root@<i>VIC_appliance_address</i>:/root/.docker/server-key.pem  /root/.docker</pre>
5. Run the following command to restart the Docker daemon.<pre>ssh root@<i>VIC_appliance_address</i> /opt/dockerd.sh</pre>
6. Run the following command in a Docker client terminal to connect the Docker client to the virtual container host with TLS authentication.<pre>docker --tlsverify -H tcp://<i>VIC_appliance_address</i>:2376 ps</pre>