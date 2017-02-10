# Deploy a VCH with vSphere Integrated Containers Registry (Harbor)

This example uses vSphere Integrated Containers Engine 0.9.0, Harbor 0.5.0, and Ubuntu on the user machine. Harbor requires 60GB or more free space on your datastore.

If no server certificate and private key are provided during installation, Harbor will create a self-generated CA (certificate authority) certificate, a server certificate, and a server private key.  The self-generated CA certificate will be available for download from the Harbor web client.

See the OVA installation guide for Harbor [Harbor docs](https://github.com/vmware/harbor/blob/master/docs/installation_guide_ova.md). Harbor requires both an IP address and FQDN (fully qualified domain name) for the server. A DHCP install method is available for debugging purposes, but it is not a recommended production deployment model.

This example assumes a Harbor instance was installed without a server certificate or a private key and the CA cert downloaded using the Harbor instructions. For Harbor to work with vSphere Integrated Containers Engine you need to update standard docker with the Harbor CA cert and deploy a new VCH with the CA cert.

## Update the User Working Machine's Standard Docker with the Harbor CA Certificate

Update the machine with standard Docker so it recognizes the CA certificate. Docker can look for additional CA certificates outside of the operating system's CA bundle folder if the new CA certificates are in the correct location. See [Verify repository client with certificates](https://docs.docker.com/engine/security/certificates/).

Create the necessary folder, copy the CA cert file to the folder. Restart Docker, then verify that you can log onto the Harbor server.

**Note** This example has Ubuntu-specific commands.

    user@Devbox:~/mycerts$ sudo su
    [sudo] password for user: 
    root@Devbox:/home/user/mycerts# mkdir -p /etc/docker/certs.d/<Harbor FQDN>
    root@Devbox:/home/user/mycerts# mkdir -p /etc/docker/certs.d/<Harbor IP>
    root@Devbox:/home/user/mycerts# cp ca.crt /etc/docker/certs.d/<Harbor FQDN>/
    root@Devbox:/home/user/mycerts# cp ca.crt /etc/docker/certs.d/<Harbor IP>/
    root@Devbox:/home/user/mycerts# exit
    exit
    user@Devbox:~/mycerts$ sudo systemctl daemon-reload
    user@Devbox:~/mycerts$ sudo systemctl restart docker

    user@Devbox:~$ docker logout <Harbor FQDN>
    Remove login credentials for <Harbor FQDN>
    user@Devbox:~$ docker logout <Harbor IP>
    Remove login credentials for <Harbor IP>

    user@Devbox:~$ docker login <Harbor FQDN>
    Username: user
    Password: 
    Login Succeeded

    user@Devbox:~$ docker login <Harbor IP>
    Username: user
    Password: 
    Login Succeeded

    user@Devbox:~$ docker logout <Harbor FQDN>
    Remove login credentials for <Harbor FQDN>

    user@Devbox:~$ docker logout <Harbor IP>
    Remove login credentials for <Harbor IP>

This example creates folders for both FQDN and IP in the docker cert folder and copies the CA cert to both folders. You can then log into Harbor from Docker using both FQDN and IP address.

## Install a VCH with the New CA Certificate

Deploy a VCH and specify the CA cert with `--registry-ca` parameter in vic-machine.  This parameter is a list, and you can easily add multiple CA certs by specifying multiple `--registry-ca` parameters.

For simplicity, this example installs a VCH with the `--no-tls` flag, so you do not need TLS from a docker CLI to the VCH. However, it does not imply that access to Harbor is performed without TLS.
<pre>
./vic-machine-linux create 
--target=<vCenter_IP> 
--image-store="vsanDatastore" 
--name=vic-docker 
--user=root 
--password=<vCenter_password> 
--compute-resource="/dc1/host/cluster1/Resources" 
--bridge-network DPortGroup 
--force 
--no-tls 
--registry-ca=ca.crt

 WARN[2016-11-11T11:46:37-08:00] Configuring without TLS 
- all communications will be insecure
 ...

 INFO[2016-11-11T11:47:57-08:00] Installer completed successfully</pre>            