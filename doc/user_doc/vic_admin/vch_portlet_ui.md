# Find Virtual Container Host Information in the vSphere Web Client #

After you have installed the vSphere Web Client plug-in for vSphere Integrated Containers Engine, you can find information about virtual container hosts in the vSphere Web Client.

**IMPORTANT**: Do not use the vSphere Web Client to perform operations on virtual container host appliances or container VMs. Specifically, using the vSphere Web Client to power off, power on, or delete virtual container host appliances or container VMs can cause vSphere Integrated Containers Engine to not function correctly. Always use `vic-machine` to perform operations on virtual container hosts. Always use Docker commands to perform operations on containers.

**Prerequisites**

- You deployed a virtual container host.
- You installed the vSphere Web Client plug-in for vSphere Integrated Containers Engine.
- If you deployed the virtual container host to a vCenter Server 6.5 instance, use the Flash-based vSphere Web Client. The vSphere Web Client plug-in for vSphere Integrated Containers Engine does not support the new HTML5 vSphere Client.

**Procedure**

1. In the vSphere Web Client Home page, select **Hosts and Clusters**.
2. Expand the hierarchy of vCenter Server objects to navigate to the virtual container host vApp.
3. Expand the virtual container host vApp and select the virtual container host endpoint VM.
4. Click the **Summary** tab for the virtual container host endpoint VM and scroll down to the Virtual Container Host portlet.

**Result**

Information about the virtual container host appears in the Virtual Container Host portlet in the **Summary** tab:

-  The address of the Docker API endpoint for this virtual container host
-  A link to the vic-admin portal for the virtual container host, from which you can obtain health information and download log bundles for the virtual container host.

