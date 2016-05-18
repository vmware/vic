# Create Private Port Groups for Virtual Container Hosts

If you deploy a vSphere Integrated Containers virtual container host in a vCenter Server cluster, you must create a private port group for each virtual container host.

You create private port groups in the vSphere Web Client.

**IMPORTANT** vSphere Integrated Containers does not support tagging on private port groups. Do not apply tags to the private port groups that you create. 

1. Log into the vSphere Web Client for the vCenter Server instance that manages your host or cluster.
2. In the **Home** view, click **Hosts and Clusters**. 
3. Select an ESXi host and click the **Manage** tab.
4. Select **Networking** > **Virtual switches**. 
5. Click the icon to **Add host networking**.
6. Select **Virtual Machine Port Group for a Standard Switch** and click **Next**.
7. Select an existing standard switch and click **Next**.
8. Rename the network for ease identification, and click **Next**.

    For example, rename the network to ```VIC_Network```.    
    
9. Click **Finish**.