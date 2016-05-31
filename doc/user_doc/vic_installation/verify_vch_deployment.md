# Verify the Deployment of a Virtual Container Host to vCenter Server #

After you have deployed a virtual container host to vCenter Server, you can verify the deployment by connecting a Docker client to the virtual container host and running Docker operations. You can check the results in the vSphere Web Client.

**Prerequisites**

- You used `vic-machine` to deploy a virtual container host to vCenter Server. For information about running `vic-machine`, see [Deploy a Virtual Container Host](install_vic_cli.md).
- You have installed a Docker client.
- If you did not explicitly disable certificate generation when you deployed the virtual container host, copy the <code><i>vch_name</i>-cert.pem</code> and <code><i>vch_name</i>-key.pem</code> files to the system on which you you are running the Docker client. 
- Connect a vSphere Web Client to the vCenter Server instance on which you deployed the virtual container host.

**Procedure**    

1. In the vSphere Web Client, go to **Hosts and Clusters** and select the cluster or host on which you deployed the virtual container host.
 
    You should see a resource pool with the name that you set for the virtual container host. The resource pool contains the virtual container host appliance VM.   

1.  In your Docker client, run the `docker info` command to confirm that you can connect to the virtual container host. 
 
 - Without TLS authentication:<pre>$ docker -H <i>vch_address</i>:2375 info</pre>
 - With TLS authentication:<pre>$ docker -H <i>vch_address</i>:2376 
--tls 
--tlscert='<i>path_to_cert_file</i>/<i>vch_name</i>-cert.pem' 
--tlskey='<i>path_to_key_file</i>/<i>vch_name</i>-key.pem' 
info</pre>

 You should see confirmation that the Storage Driver is ```Portlayer Storage``` and the `Name` is `VIC`.
1.  Pull a Docker container image into the virtual container host, for example, the `BusyBox` container.

 - Without TLS authentication:<pre>$ docker -H <i>vch_address</i>:2375 pull busybox:latest</pre>
 - With TLS authentication:<pre>$ docker -H <i>vch_address</i>:2376 
--tls 
--tlscert='<i>path_to_cert_file</i>/<i>vch_name</i>-cert.pem' 
--tlskey='<i>path_to_key_file</i>/<i>vch_name</i>-key.pem' 
pull busybox:latest</pre>

1. In the vSphere Web Client, go to **Storage**, select the datastore that you designated as the image store, click **Manage** > **Files**, and navigate to the **VIC** > **Photon** folder. 
 
    You should see a folder for every container image that you have pulled into the virtual container host. The folders contain the container image files.

1. In your Docker client, run the Docker container that you pulled into the virtual container host.

  - Without TLS authentication:<pre>$ docker -H <i>vch_address</i>:2375 run busybox</pre>
  - With TLS authentication:<pre>$ docker -H <i>vch_address</i>:2376 
--tls 
--tlscert='<i>path_to_cert_file</i>/<i>vch_name</i>-cert.pem' 
--tlskey='<i>path_to_key_file</i>/<i>vch_name</i>-key.pem' 
run busybox</pre> 

1. In the vSphere Web Client, go to **Hosts and Clusters** and select the resource pool that contains the virtual container host.
 
    You should see a VM for every container that you run.

1. In the vSphere Web Client, go to **Storage** and select the datastore that you designated as the container store.  
 
    At the top-level of the datastore, you should see a folder for every container that you run. The folders contain the container VM files.