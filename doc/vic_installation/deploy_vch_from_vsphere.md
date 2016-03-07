# Deploy a Virtual Container Host from the vSphere Web Client #

After you have deployed the vSphere Integrated Containers management server and installed the vSphere Web Client extension for vSphere Integrated Containers, you can start deploying virtual container hosts directly from the vSphere Web Client.

**Prerequisites**

* You deployed the vSphere Integrated Containers management server by using the OVA installer.
* You installed the vSphere Web Client extension for vSphere Integrated Containers.
* You have created a private port group on the ESXi host on which to deploy the virtual container host, as described in [Create a Private Port Group for Virtual Container Hosts](create_a_private_port_group_for_vch.md)

**Procedure**

1. In the **Home** view of the vSphere Web Client, click **vSphere Integrated Containers**.
2. Select **vSphere Integrated Containers Home** and click the **Getting Started** tab.
4. Click **Create a Virtual Container Host**.
5. Enter a name for the virtual container host and click **Next**.
6. Select the resource on which to run the virtual container host and click **Next**.
7. Select the datastore in which to store the virtual container host VM files.
8. Select the networks for the virtual container host to use, and click **Next**.

  * **External network**: The network to use for traffic between the virtual container host and vSphere. For example, VM Network.
  * **Internal network**: The network to use for all container traffic. Select  the network that you created when you created a private port group. 

9. Optionally set a static IP address for the virtual container host, and click **Next**.
10. Optionally increase the number of vCPUs and memory for the virtual container host VM, and click **Next**.
11. Leave the checkbox selected to participate in the VMware Customer Experience Inprovement Program, or deselect the checkbox, and click **Next**.
12. Review the details of the virtual container host deployment, leave the checkbox selected to perform a post-deployment functional verification, or deselect the checkbox, and click **Finish**.