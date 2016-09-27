# vSphere Integrated Containers Engine Plug-In Does Not Appear in the vSphere Web Client #

After you have installed the vSphere Web Client plug-in for vSphere Integrated Containers Engine, the plug-in does not appear in the vSphere Web Client.

## Problem ##

The UI plug-in installer reported success, but the Virtual Container Host portlet does not appear in the **Summary** tab for the virtual container host endpoint VM. Logging out of the vSphere Web Client and logging back in again does not resolve the issue.

## Cause ##

If a previous attempt at installing the vSphere Integrated Containers engine plug-in failed, the failed installation state is retained in the vSphere Web Client cache.


## Solution ##

Restart the vSphere Web Client service.

### vCenter Server on Windows ###

1. Open Server Manager on the Windows system on which vCenter Server is running.
2. Select **Configuration** > **Services**.
3. Select **VMware vSphere Web Client** and click **Restart**.

### vCenter Server Appliance ###

1. Use SSH to log in to the vCenter Server Appliance as root.
2. Stop the vSphere Web Client service by running the following command:<pre>service vsphere-client stop</pre>
3. Restart the vSphere Web Client service by running the following command:<pre>service vsphere-client start</pre>


