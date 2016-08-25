# Find Container Information in the vSphere Web Client #

After you have installed the vSphere Web Client plug-in for vSphere Integrated Containers, you can use the vSphere Web Client to find information about containers that are running in virtual container hosts.

**Prerequisites**

- You deployed a virtual container host and pulled and ran at least one container. For information about how to pull and run a container, see 
- You installed the vSphere Web Client plug-in for vSphere Integrated Containers.

**Procedure**

1. In the vSphere Web Client Home page, select **Hosts and Clusters**.
2. Expand the hierarchy of vCenter Server objects to navigate to the virtual container host vApp.
3. Expand the virtual container host vApp and select a container VM.
4. Click the **Summary** tab for the virtual container host VM and scroll down to the **Container** portlet.

**Result**

Information about the container appears in the Container portlet in the **Summary** tab:

-  The name of the running container
-  The image from which the container was deployed