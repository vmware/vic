# Verify the Deployment of the vSphere Integrated Containers Engine Plug-In #

After you have installed the vSphere Web Client plug-in for vSphere Integrated Containers Engine, verify the deployment of the plug-in in the vSphere Web Client.

**Prerequisites**

- You deployed a virtual container host.
- You installed the vSphere Web Client plug-in for vSphere Integrated Containers Engine.
- You logged out of the vSphere Web Client after deploying the plug-in, and logged back in.
- If you deployed the virtual container host to a vCenter Server 6.5 instance, use the Flash-based vSphere Web Client to view the vSphere Web Client plug-in for vSphere Integrated Containers Engine. vSphere Integrated Containers Engine does not currently provide a plug-in for the new HTML5 vSphere Client.

**Procedure**

1. In the vSphere Web Client Home page, select **Hosts and Clusters**.
2. Expand the hierarchy of vCenter Server objects to navigate to the virtual container host vApp.
3. Expand the virtual container host vApp and select the virtual container host endpoint VM.
4. Click the **Summary** tab for the virtual container host endpoint VM and scroll down to the Virtual Container Host portlet.

**Result**

Information about the virtual container host appears in the Virtual Container Host portlet in the **Summary** tab:

-  The address of the Docker API endpoint for this virtual container host
-  A link to the vic-admin portal for the virtual container host, from which you can obtain health information and download log bundles for the virtual container host.

**What to Do Next**

If the Virtual Container Host portlet still does not appear in the **Summary** tab for the virtual container host endpoint VM, restart the vSphere Web Client service. For instructions about how to restart the vSphere Web Client service, see [vSphere Integrated Containers Engine Plug-In Does Not Appear in the vSphere Web Client](ts_ui_not_appearing.md).