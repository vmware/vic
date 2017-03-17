# Find VCH Information in the vSphere Client #

After you have installed either or both of the HTML5 or Flex-based client plug-ins for vSphere Integrated Containers, you can find information about virtual container hosts (VCHs) in the HTML5 vSphere Client or the Flex-based vSphere Web Client.

**IMPORTANT**: Do not use the vSphere Client to perform operations on VCH appliances or container VMs. Specifically, using the vSphere Client to power off, power on, or delete VCH appliances or container VMs can cause vSphere Integrated Containers Engine to not function correctly. Always use `vic-machine` to perform operations on VCHs. Always use Docker commands to perform operations on containers.

**Prerequisites**

- You deployed a VCH.
- You installed either or both of HTML5 or Flex-based plug-ins for vSphere Integrated Containers.

**Procedure**

1. Log in to either the HTML5 vSphere Client or the Flex-based vSphere Web Client.
2. On the **Home** page, select **Hosts and Clusters**.
2. Expand the hierarchy of vCenter Server objects to navigate to the VCH vApp.
3. Expand the VCH vApp and select the VCH endpoint VM.
4. Click the **Summary** tab for the VCH endpoint VM and scroll down to the Virtual Container Host portlet.

**Result**

Information about the VCH appears in the VCH portlet in the **Summary** tab:

-  The address of the Docker API endpoint for this VCH
-  A link to the VCH Admin portal for the VCH, from which you can obtain health information and download log bundles for the VCH.