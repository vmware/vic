# Find VCH Information in the vSphere Web Client #

After you have installed the vSphere Web Client plug-in for vSphere Integrated Containers Engine, you can find information about virtual container hosts (VCHs) in the vSphere Web Client.

**IMPORTANT**: Do not use the vSphere Web Client to perform operations on VCH appliances or container VMs. Specifically, using the vSphere Web Client to power off, power on, or delete VCH appliances or container VMs can cause vSphere Integrated Containers Engine to not function correctly. Always use `vic-machine` to perform operations on VCHs. Always use Docker commands to perform operations on containers.

**Prerequisites**

- You deployed a VCH.
- You installed the vSphere Web Client plug-in for vSphere Integrated Containers Engine.
- If you deployed the VCH to a vCenter Server 6.5 instance, use the Flash-based vSphere Web Client to view the vSphere Web Client plug-in for vSphere Integrated Containers Engine. vSphere Integrated Containers Engine does not currently provide a plug-in for the new HTML5 vSphere Client.

**Procedure**

1. In the vSphere Web Client Home page, select **Hosts and Clusters**.
2. Expand the hierarchy of vCenter Server objects to navigate to the VCH vApp.
3. Expand the VCH vApp and select the VCH endpoint VM.
4. Click the **Summary** tab for the VCH endpoint VM and scroll down to the VCH portlet.

**Result**

Information about the VCH appears in the VCH portlet in the **Summary** tab:

-  The address of the Docker API endpoint for this VCH
-  A link to the vic-admin portal for the VCH, from which you can obtain health information and download log bundles for the VCH.

