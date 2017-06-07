# Create a Private Port Group for Virtual Container Hosts

If you install vSphere Integrated Containers in a vCenter Server cluster, you might need to create a private port group for each virtual container host. You must create a private port group in the following circumstances:

* You are installing vSphere Integrated Containers by using the OVA deployment.
* You are installing vSphere Integrated Containers by using the command line installer and you want to use the `containerNetwork` option to assign a non-default port group and network for containers to use.

If you install vSphere Integrated Containers by using the command line installer and you do not use the `containerNetwork` option, the installer creates a default private port group named `container_port_group` and a network with the same name as the vSphere Integrated Containers appliance.

You create private port groups in the vSphere Web Client.

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

