# Verify the Deployment of a Virtual Container Host to an ESXi Host #

After you have deployed a virtual container host to a standalone ESXi host, you can verify the deployment by connecting a Docker client to the virtual container host and running Docker operations. You can check the results in the vSphere Client.

**Prerequisites**

- You used `vic-machine` `create` to deploy a virtual container host to a standalone ESXi host. For information about running `vic-machine`, see [Deploy a Virtual Container Host](install_vic_cli.md).
- You have installed a Docker client.
- If you did not explicitly disable certificate generation when you deployed the virtual container host, and if you are not running the Docker client in the same location as the one on which you ran `vic-machine`, copy the <code><i>vch_name</i>-cert.pem</code> and <code><i>vch_name</i>-key.pem</code> files to the system on which you you are running the Docker client. 
- Connect a vSphere Client to the ESXi host on which you deployed the virtual container host.

**Procedure**    

1. In the vSphere Client, go to **Inventory** and select the host on which you deployed the virtual container host.
 
    You should see a resource pool with the name that you set for the virtual container host. The resource pool contains the vSphere Integrated Containers endpoint VM.   

1.  In your Docker client, run the `docker info` command to confirm that you can connect to the virtual container host. 
 
 <pre>$ docker -H <i>vch_address</i>:2376 --tls info</pre>

 You should see confirmation that the Storage Driver is `vSphere Integrated Containers Backend Engine`.
1.  Pull a Docker container image into the virtual container host, for example, the `BusyBox` container.

  <pre>$ docker -H <i>vch_address</i>:2376 --tls pull busybox:latest</pre>

1. In the vSphere Client, click the **Summary** tab for the ESXi host, right-click the datastore that you designated as the image store, select **Browse Datastore**, and navigate to the **VIC** folder. 

  If you specified a specific datastore folder as the image store, the **VIC** folder appears inside that folder.
 
    You should see a folder for every container image that you have pulled into the virtual container host. The folders contain the container image files.

1. In your Docker client, run the Docker container that you pulled into the virtual container host.<pre>$ docker -H <i>vch_address</i>:2376 --tls run busybox</pre> 

1. In the vSphere Client, go to **Inventory** and expand the virtual container host resource pool.
 
    You should see a running VM for every container that you run.

1. In the vSphere Client, click the **Summary** tab for the ESXi host, right-click the datastore that you designated as the container store, and select **Browse Datastore**.  
 
    At the top-level of the datastore, you should see a folder for every container that you run. The folders contain the container VM files.