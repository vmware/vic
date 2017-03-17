# Find Container Information in the vSphere Client #

After you have installed either or both of the plug-ins for vSphere Integrated Containers, you can find information about the container VMs that are running in virtual container hosts (VCHs).

**IMPORTANT**: Do not use the vSphere Client to perform operations on VCH appliances or container VMs. Specifically, using the vSphere Client to power off, power on, or delete VCH appliances or container VMs can cause vSphere Integrated Containers Engine to not function correctly. Always use `vic-machine` to perform operations on VCHs. Always use Docker commands to perform operations on containers.

**Prerequisites**

- You deployed a VCH and pulled and ran at least one container.
- You installed either or both of HTML5 or Flex-based plug-ins for vSphere Integrated Containers.

**Procedure**

1. Log in to either the HTML5 vSphere Client or the Flex-based vSphere Web Client.
2. On the **Home** page, select **Hosts and Clusters**.
2. Expand the hierarchy of vCenter Server objects to navigate to the VCH vApp.
3. Expand the VCH vApp and select a container VM.
4. Click the **Summary** tab for the container VM and scroll down to the **Container** portlet.

**Result**

Information about the container VM appears in the Container portlet in the **Summary** tab:

-  The name of the running container. If the container developer used <code>docker run -name <i>container_name</i></code> to run the container, <code><i>container_name</i></code> appears in the portlet.
-  The image from which the container was deployed.
-  If the container developer used <code>docker run -p <i>port</i></code> to map a port when running the container, the port number and the protocol appear in the portlet.