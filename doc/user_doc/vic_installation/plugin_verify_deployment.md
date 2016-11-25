# Verify the Deployment of the vSphere Integrated Containers Engine Plug-In #

After you have installed the vSphere Web Client plug-in for vSphere Integrated Containers Engine, verify the deployment of the plug-in in the vSphere Web Client.

**Prerequisites**

- You deployed a virtual container host (VCH).
- You installed the vSphere Web Client plug-in for vSphere Integrated Containers Engine.
- You logged out of the vSphere Web Client after deploying the plug-in, and logged back in.
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

**What to Do Next**

If the VCH portlet still does not appear in the **Summary** tab for the VCH endpoint VM, restart the vSphere Web Client service. For instructions about how to restart the vSphere Web Client service, see [vSphere Integrated Containers Engine Plug-In Does Not Appear in the vSphere Web Client](ts_ui_not_appearing.md).