# Verify the Deployment of a Virtual Container Host to vCenter Server #

After you have deployed a virtual container host to vCenter Server, you can verify the deployment by connecting a Docker client to the virtual container host and running Docker operations. You can check the results in the vSphere Web Client.

**Prerequisites**

- You used `vic-machine` `create` to deploy a virtual container host to vCenter Server. For information about running `vic-machine`, see [Deploy a Virtual Container Host](install_vic_cli.md) and [Virtual Container Host Deployment Options](vch_installer_options.md).
- You have installed a Docker client.
- If you did not explicitly disable certificate generation when you deployed the virtual container host, and if you are not running the Docker client in the same location as the one on which you ran `vic-machine`, copy the <code><i>vch_name</i>-cert.pem</code> and <code><i>vch_name</i>-key.pem</code> files to the system on which you you are running the Docker client. Copy the certificate and key files either to the folder in which you run Docker, or to the `~/.docker` folder. Alternatively, specify the following options when you run Docker commands:<pre>docker -H <i>vch_address</i>:2376 --tlscacert --tlscert=<i>path</i> --tlskey=<i>path</i></pre>
- Connect a vSphere Web Client to the vCenter Server instance on which you deployed the virtual container host.

**Procedure**    

1. In the vSphere Web Client, go to **Hosts and Clusters** and select the cluster or host on which you deployed the virtual container host.
 
    You should see a vApp with the name that you set for the virtual container host. The vApp contains the vSphere Integrated Containers Engine endpoint VM.   

1.  In your Docker client, run the `docker info` command to confirm that you can connect to the virtual container host.
 - With TLS authentication: <pre>$ docker -H <i>vch_address</i>:2376 --tls info</pre>
 - Without TLS authentication: <pre>$ docker -H <i>vch_address</i>:2375 info</pre>

 You should see confirmation that the Storage Driver is ``` vSphere Integrated Containers Backend Engine```.
1.  Pull a Docker container image into the virtual container host, for example, the `BusyBox` container.
 
   - With TLS authentication: <pre>$ docker -H <i>vch_address</i>:2376 --tls pull busybox:latest</pre>
   - Without TLS authentication: <pre>$ docker -H <i>vch_address</i>:2375 pull busybox:latest</pre>

1. In the vSphere Web Client, go to **Storage**, select the datastore that you designated as the image store, click **Manage** > **Files**, and navigate to the **VIC** folder. 

  If you specified a specific datastore folder as the image store, the **VIC** folder appears inside that folder.
 
    You should see a folder for every container image that you have pulled into the virtual container host. The folders contain the container image files.

1. In your Docker client, run the Docker container that you pulled into the virtual container host.
 
   - With TLS authentication: <pre>$ docker -H <i>vch_address</i>:2376 --tls run busybox --name test</pre> 
   - Without TLS authentication: <pre>$ docker -H <i>vch_address</i>:2375 run busybox --name test</pre> 

1. In the vSphere Web Client, go to **Hosts and Clusters** and expand the virtual container host vApp
 
    You should see a VM for every container that you run, including a VM named `test`.

1. In the vSphere Web Client, go to **Storage** and select the datastore that you designated as the container store.  
 
    At the top-level of the datastore, you should see a folder for every container that you run. The folders contain the container VM files.