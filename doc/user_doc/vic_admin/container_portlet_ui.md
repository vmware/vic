# Find Container Information in the vSphere Web Client #

After you have installed the vSphere Web Client plug-in for vSphere Integrated Containers Engine, you can use the vSphere Web Client to find information about containers that are running in virtual container hosts (VCHs).

**IMPORTANT**: Do not use the vSphere Web Client to perform operations on VCH appliances or container VMs. Specifically, using the vSphere Web Client to power off, power on, or delete VCH appliances or container VMs can cause vSphere Integrated Containers Engine to not function correctly. Always use `vic-machine` to perform operations on VCHs. Always use Docker commands to perform operations on containers.

**Prerequisites**

- You deployed a VCH and pulled and ran at least one container.
- You installed the vSphere Web Client plug-in for vSphere Integrated Containers Engine.
- If you deployed the VCH to a vCenter Server 6.5 instance, use the Flash-based vSphere Web Client to view the vSphere Web Client plug-in for vSphere Integrated Containers Engine. vSphere Integrated Containers Engine does not currently provide a plug-in for the new HTML5 vSphere Client.

**Procedure**

1. In the vSphere Web Client Home page, select **Hosts and Clusters**.
2. Expand the hierarchy of vCenter Server objects to navigate to the VCH vApp.
3. Expand the VCH vApp and select a container VM.
4. Click the **Summary** tab for the container VM and scroll down to the **Container** portlet.

**Result**

Information about the container appears in the Container portlet in the **Summary** tab:

-  The name of the running container. If the container developer used <code>docker run -name <i>container_name</i></code> to run the container, <code><i>container_name</i></code> appears in the portlet.
-  The image from which the container was deployed.
-  If the container developer used <code>docker run -p <i>port</i></code> to map a port when running the container, the port number and the protocol appear in the portlet.