# Find Virtual Container Host Information in the vSphere Web Client #

After you have installed the vSphere Web Client plug-in for vSphere Integrated Containers, you can find information about virtual container hosts in the vSphere Web Client.

**Prerequisites**

- You deployed a virtual container host.
- You installed the vSphere Web Client plug-in for vSphere Integrated Containers.

**Procedure**

1. In the vSphere Web Client Home page, select **Hosts and Clusters**.
2. Expand the hierarchy of vCenter Server objects to navigate to the virtual container host vApp.
3. Expand the virtual container host vApp and select the virtual container host VM.
4. Click the **Summary** tab for the virtual container host VM.

**Result**

Information about the virtual container host appears in the Virtual Container Host portlet in the Summary tab:

-  The address of the Docker API endpoint for this virtual container host
-  A link to the vic-admin portal for the virtual container host, from which you can obtain health information and download log bundles for the virtual container host.

