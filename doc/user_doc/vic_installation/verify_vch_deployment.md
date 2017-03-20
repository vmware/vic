# Verify the Deployment of a VCH #

After you have deployed a virtual container host (VCH), you can verify the deployment by connecting a Docker client to the VCH and running Docker operations. You can check the results in the vSphere Client or vSphere Web Client.

**IMPORTANT**: Do not use the vSphere Client or vSphere Web Client to perform operations on VCH appliances or container VMs. Specifically, using the vSphere Client or vSphere Web Client to power off, power on, or delete VCH appliances or container VMs can cause vSphere Integrated Containers Engine to not function correctly. Always use `vic-machine` to perform operations on VCHs. Always use Docker commands to perform operations on containers.

**Prerequisites**

- You followed the instructions in [Deploy a VCH to an ESXi Host with No vCenter Server](deploy_vch_esxi.md) or [Deploy a VCH to a Basic vCenter Server Cluster](deploy_vch_vcenter.md) to deploy a VCH to either an ESXi host or to a vCenter Server instance.
- You ran `vic-machine create` with the `--no-tlsverify` option.
- You have installed a Docker client.
- If you deployed the VCH to vCenter Server, connect a vSphere Web Client to that vCenter Server instance.
- If you deployed the VCH to an ESXi host, connect a vSphere Client to that host.


**Procedure**    

1. View the VCH appliance in the vSphere Web Client or vSphere Client.
 
   - vCenter Server: Go to **Hosts and Clusters** in the vSphere Web Client and select the cluster or host on which you deployed the VCH. You should see a vApp with the name that you set for the VCH.
   - ESXi host: Go to **Inventory** in the vSphere Client and select the host on which you deployed the VCH. You should see a resource pool with the name that you set for the VCH.

   The vApp or resource pool contains the VCH endpoint VM.   

3.  Run the `docker info` command to confirm that you can connect to the VCH.<pre>docker -H <i>vch_address</i>:2376 --tls info</pre>

     You should see confirmation that the Storage Driver is ``` vSphere Integrated Containers Backend Engine```.

1.  Pull a Docker container image into the VCH, for example, the `BusyBox` container.<pre>docker -H <i>vch_address</i>:2376 --tls pull busybox</pre>
1. View the container image files in the vSphere Web Client or vSphere Client.

    - vCenter Server: Go to **Storage**, select the datastore that you designated as the image store, and click **Manage** > **Files**. 
    - ESXi host: Click the **Summary** tab for the ESXi host, right-click the datastore that you designated as the image store, and select **Browse Datastore**. 

    vSphere Integrated Containers Engine creates a folder a folder that has the same name as the VCH, that contains a folder named `VIC` in which to store container image files.
  
1. Expand the `VIC` folder to navigate to the `images` folder.  The `images` folder contains a folder for every container image that you pull into the VCH. The folders contain the container image files.
  
1. In your Docker client, run the Docker container that you pulled into the VCH.<pre>docker -H <i>vch_address</i>:2376 --tls run --name test busybox</pre>

1. View the container VMs in the vSphere Web Client or vSphere Client.

    - vCenter Server: Go to **Hosts and Clusters** and expand the VCH vApp.
    - ESXi host: Go to **Inventory** and expand the VCH resource pool.
 
    You should see a VM for every container that you run, including a VM named <code>test-<i>container_id</i></code>.

1. View the container VM files in the vSphere Web Client or vSphere Client.

    - vCenter Server: Go to **Storage** and select the datastore that you designated as the image store. 
    - ESXi host: Click the **Summary** tab for the ESXi host, right-click the datastore that you designated as the image store, and select **Browse Datastore**. 
 
     At the top-level of the datastore, you should see a folder for every container that you run. The folders contain the container VM files.