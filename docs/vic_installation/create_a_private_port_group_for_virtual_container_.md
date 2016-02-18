# Create a Private Port Group for Virtual Container Hosts

If you use the experimental vSphere Integrated Containers UI to deploy virtual container hosts, you must create a private port group for the virtual container hosts to use. 

If you use the default command line installer to deploy a virtual container host, the installer creates the port group automatically.

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

