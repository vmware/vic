# Verify the Deployment of a Virtual Container Host #

After you have deployed a virtual container host, you can verify the deployment by connecting a Docker client to the virtual container host and running Docker operations. You can check the results in the vSphere Client or vSphere Web Client.

**Prerequisites**

- You used `vic-machine` `create` to deploy a virtual container host to either a vCenter Server instance or an ESXi host. For information about running `vic-machine`, see [Deploy a Virtual Container Host](install_vic_cli.md) and [Virtual Container Host Deployment Options](vch_installer_options.md).
- You have installed a Docker client.
- Configure your Docker client according to the type of authentication that the virtual container host uses:

  - If you deployed the virtual container host with either of the `--no-tlsverify` or `--no-tls` options, disable TLS in the Docker client.<pre>set DOCKER_TLS_VERIFY=0</pre> 
  - If you deployed the virtual container host with full TLS autentication with trusted certificates, enable TLS in the Docker client.<pre>set DOCKER_TLS_VERIFY=1</pre>
  - If you deployed the virtual container host with full TLS autentication with trusted certificates, and if you are not running the Docker client in the same location as the one on which you ran `vic-machine`, copy the <code>*.pem</code> certificate files to the system on which you you are running the Docker client. Copy the certificate files either to the folder in which you run Docker commands, or to the `~/.docker` folder.
- If you deployed the virtual container host to vCenter Server, connect a vSphere Web Client to that vCenter Server instance.
- If you deployed the virtual container host to an ESXi host, connect a vSphere Client to that host.

**Procedure**    

1. View the virtual container host appliance in the vSphere Web Client or vSphere Client.
 
   - vCenter Server: Go to **Hosts and Clusters** in the vSphere Web Client and select the cluster or host on which you deployed the virtual container host. You should see a vApp with the name that you set for the virtual container host.
   - ESXi host: Go to **Inventory** in the vSphere Client and select the host on which you deployed the virtual container host. You should see a resource pool with the name that you set for the virtual container host.

   The vApp or resource pool contains the virtual container host endpoint VM.   

2.  In your Docker client terminal, run the `docker info` command to confirm that you can connect to the virtual container host.

  How you run Docker commands depends on the level of authentication that the virtual container host requires:
 - Full TLS authentication with trusted CA certificates: <pre>docker -H <i>vch_address</i>:2376 --tlscacert --tlscert=<i>path</i> --tlskey=<i>path</i> info</pre>
 - TLS authentication with untrusted self-signed certificates: <pre>docker -H <i>vch_address</i>:2376 --tls info</pre>
 - With no TLS authentication: <pre>$ docker -H <i>vch_address</i>:2375 info</pre>

 You should see confirmation that the Storage Driver is ``` vSphere Integrated Containers Backend Engine```. If the connection fails with a Docker API version error, see [Docker Commands Fail with a Docker API Version Error](ts_docker_version_error.md).
1.  Pull a Docker container image into the virtual container host, for example, the `BusyBox` container.
 
  - Full TLS authentication with trusted CA certificates: <pre>docker -H <i>vch_address</i>:2376 --tlscacert --tlscert=<i>path</i> --tlskey=<i>path</i> pull busybox:latest</pre>
  - TLS authentication with untrusted self-signed certificates: <pre>docker -H <i>vch_address</i>:2376 --tls pull busybox:latest</pre>
  - With no TLS authentication: <pre>$ docker -H <i>vch_address</i>:2375 pull busybox:latest</pre>

1. View the container image files in the vSphere Web Client or vSphere Client.

  - vCenter Server: Go to **Storage**, select the datastore that you designated as the image store, and click **Manage** > **Files**. 
  - ESXi host: Click the **Summary** tab for the ESXi host, right-click the datastore that you designated as the image store, and select **Browse Datastore**. 

  vSphere Integrated Containers Engine creates a folder named `VIC` in which to store container image files. 

  -  If you specified a specific datastore folder as the image store when you deployed the virtual container host, the `VIC` folder appears inside that folder.
  -  If you did not specify a specific datastore folder as the image store when you deployed the virtual container host, the `VIC` folder appears in a folder that has the same name as the virtual container host.
  
  Expand the `VIC` folder to navigate to the `images` folder.  The `images` folder contains a folder for every container image that you pull into the virtual container host. The folders contain the container image files.
  
1. In your Docker client, run the Docker container that you pulled into the virtual container host.
 
   - Full TLS authentication with trusted CA certificates: <pre>docker -H <i>vch_address</i>:2376 --tlscacert --tlscert=<i>path</i> --tlskey=<i>path</i> run --name test busybox</pre>
  - TLS authentication with untrusted self-signed certificates: <pre>docker -H <i>vch_address</i>:2376 --tls run --name test busybox</pre>
  - With no TLS authentication: <pre>$ docker -H <i>vch_address</i>:2375 run --name test busybox</pre> 

1. View the container VMs in the vSphere Web Client or vSphere Client.

  - vCenter Server: Go to **Hosts and Clusters** and expand the virtual container host vApp.
  - ESXi host: Go to **Inventory** and expand the virtual container host resource pool.
 
  You should see a VM for every container that you run, including a VM named `test`.

1. View the container VM files in the vSphere Web Client or vSphere Client.

  - vCenter Server: Go to **Storage** and select the datastore that you designated as the image store. 
  - ESXi host: Click the **Summary** tab for the ESXi host, right-click the datastore that you designated as the image store, and select **Browse Datastore**. 
 
   At the top-level of the datastore, you should see a folder for every container that you run. The folders contain the container VM files.