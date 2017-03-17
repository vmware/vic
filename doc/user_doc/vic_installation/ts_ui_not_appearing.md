# vSphere Integrated Containers Plug-In Does Not Appear #

After you have installed either of the HTML5 or Flex-based plug-ins for vSphere Integrated Containers, the plug-ins do not appear in the HTML5 vSphere Client or the Flex-based vSphere Web Client.

## Problem ##

The UI plug-in installer reported success, but the plug-in does not appear in the client. Logging out of the client and logging back in again does not resolve the issue.

## Cause ##

- If a previous attempt at installing the vSphere Integrated Containers plug-in failed, the failed installation state is retained in the client cache.
- You installed a new version of the vSphere Integrated Containers plug-in that has the same version number as the previous version, for example a hot patch.

## Solution ##

Restart the client service.

### Restart the HTML5 or Flex Client on vSphere 6.5 on Windows ###

1. Log into the Flex-based vSphere Web Client.
2. Go to **Home** > **System Configuration** > **Services**.
3. Select the service for either the HTML5 or Flex client:

   - HTML5: **VMware vSphere Client**.
   - Flex: **VMware vSphere Web Client**
3. Click the **Restart** button.

### Restart the Flex Client on vCenter Server 6.0 on Windows ###

1. Open Server Manager on the Windows system on which vCenter Server is running.
2. Select **Configuration** > **Services**.
3. Select **VMware vSphere Web Client** and click **Restart**.

### Restart the Flex Client on a vCenter Server Appliance ###

1. Use SSH to log in to the vCenter Server Appliance as `root`.
2. Stop the vSphere Web Client service by running one of the following commands.
   - vCenter Server 6.0: <pre>service vsphere-client stop</pre>
   - vCenter Server 6.5: <pre>service-control --stop vsphere-client</pre>
3. Restart the vSphere Web Client service by running one of the following commands.
   - vCenter Server 6.0:<pre>service vsphere-client start</pre>
   - vCenter Server 6.5: <pre>service-control --start vsphere-client</pre>


