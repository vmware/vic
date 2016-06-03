# Create a Private Port Group for Virtual Container Hosts

If you install vSphere Integrated Containers in an environment that includes more than one ESXi host, you must create a private port group for container VMs to use to communicate with each other.

You must create the private port group in advance of deploying a virtual container host, and use the `-bridge-network` option to assign the network to the virtual container host.

In an environment with a single ESXi host, the `-bridge-network` option is optional and you do not need to create a port group in advance. For information about when to specify the `-bridge-network` option, see [Virtual Container Host Deployment Options](vch_installer_options.md).

You create private port groups in the vSphere Web Client.

**IMPORTANT** vSphere Integrated Containers does not support tagging on private port groups. Do not apply tags to the private port groups that you create. 

1. Log into the vSphere Web Client for the vCenter Server instance that manages your host or cluster.
2. In the **Home** view, click **Hosts and Clusters**. 
3. Select an ESXi host and click the **Manage** tab.
4. Select **Networking** > **Virtual switches**. 
5. Click the icon to **Add host networking**.
6. Select **Virtual Machine Port Group for a Standard Switch** and click **Next**.
7. Select an existing standard switch, or create a new standard switch, and click **Next**.
8. Rename the network for ease identification, and click **Next**.

    For example, rename the network to ```VIC_Network```.    
    
9. Click **Finish**.

