# Remove a Virtual Container Host #

In the current builds of vSphere Integrated Containers, you remove virtual container hosts manually.

Removing a virtual container host involves several operations in the vSphere Web Client.

**Prerequisites**

- You have a virtual container host that you no longer require.
- Connect a vSphere Web Client to the vCenter Server instance on which you deployed the virtual container host.

**Procedure**

1. In the vSphere Web Client, go to **Hosts and Clusters** and select the cluster or host on which you deployed the virtual container host.
2. Select the resource pool for the virtual container host and click **Related Objects** > **Top Level Objects**.
3. Select the VM for the virtual container host and any container VMs, right-click, and select **Power** > **Power Off**.
4. Select the VM for the virtual container host and the container VMs, right-click, and select **Delete from Disk**.
5. Right-click the resource pool for the virtual container host and select **Delete**.
6. Go to **Storage**, select the datastore that you designated as the image datastore, and delete the folder for the virtual container host.
  
  If other virtual container hosts on this vCenter Server instance use this datastore, do not delete the `VIC` folder. The `VIC` folder contains container images that are shared between all virtual container hosts.
7. Select the datastore that you designated as the container datastore, and delete the folders for the container VMs that this virtual container host manages.
  
  In the current builds, container VM files are stored in the image datastore. 
8. (Optional) Go to **Hosts and Clusters**, select the host on which you configured a private port group for the virtual container host, and select **Manage** > **Networking** > **Virtual switches**.
9. (Optional) Select the private port group and click the delete button.
  
  Do not delete the port group if other virtual container hosts on this vCenter Server instance use it.
 
 