# vSphere Integrated Containers Engine Plug-In Does Not Appear #

After you have installed the vSphere Client plug-in for vSphere Integrated Containers Engine, the plug-in does not appear in the vSphere Client.

## Problem ##

The UI plug-in installer reported success, but the plug-in does not appear in the vSphere Client. Logging out of the vSphere Client and logging back in again does not resolve the issue.

## Cause ##

- If a previous attempt at installing the vSphere Integrated Containers Engine plug-in failed, the failed installation state is retained in the vSphere Web Client cache.
- You installed a new version of the vSphere Integrated Containers Engine plug-in that has the same version number as the previous version, for example a hot patch.

## Solution ##

Restart the client service.

### Restart the Flex Client on vCenter Server 6.0 on Windows ###

1. Open Server Manager on the Windows system on which vCenter Server is running.
2. Select **Configuration** > **Services**.
3. Select **VMware vSphere Web Client** and click **Restart**.

### Restart the Flex Client on a vCenter Server Appliance ###

1. Use SSH to log in to the vCenter Server Appliance as root.
2. Stop the vSphere Web Client service by running one of the following commands.
   - vCenter Server 6.0: <pre>service vsphere-client stop</pre>
   - vCenter Server 6.5: <pre>service-control --stop vsphere-client</pre>
3. Restart the vSphere Web Client service by running one of the following commands.
   - vCenter Server 6.0:<pre>service vsphere-client start</pre>
   - vCenter Server 6.5: <pre>service-control --start vsphere-client</pre>

### Restart the HTML5 or Flex Client for vSphere 6.5 ###

1. Log into the Flex-based vSphere Web Client.
2. Got to **Home** > **System Configuration** > **Services** > **VMware vSphere Web Client**.
3. Click the **Restart** button.


