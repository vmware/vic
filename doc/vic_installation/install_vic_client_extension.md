# Install the vSphere Web Client Extension for vSphere Integrated Containers #

After you have deployed the vSphere Integrated Containers management server, you can install a vSphere Web Client extension for vSphere Integrated Containers. The vSphere Web Client extension allows you to deploy virtual container hosts and to manage containers and images in your vSphere environment.

**Prerequisites**

You deployed the vSphere Integrated Containers management server by using the OVA installer.

**Procedure**

1. In the vSphere Web Client, select the vSphere Integrated Containers management server VM, click the **Summary** tab, and take a note of the external IP address of the VM.
2. In a Web browser, go to https://*ip_address_of_management_server*/
3. Select the **Install** radio button.
4. Enter the address, username, and password for the vCenter Server instance on which you installed the management server, and click **Submit**.
5. When the installation finishes, log out of the vSphere Web Client and log back in again.

 The vSphere Integrated Containers extension appears in **Home** view of the vSphere Web Client. 