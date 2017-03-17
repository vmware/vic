# Upgrade the HTML5 vSphere Client Plug-In #

If you download a new version of vSphere Integrated Containers Engine, upgrade the HTML5 vSphere Client plug-in for vSphere Integrated Containers.

**Prerequisites**

- You deployed an older version of the vSphere Integrated Containers plug-in for the HTML5 vSphere Client.
- You downloaded a new version of vSphere Integrated Containers Engine.
- For information about updating the `configs` and `install.sh` files, and where to copy the <code>com.vmware.vic-<i>version</i>.zip</code> file, see [Install the HTML5 Plug-In on a vCenter Server Appliance](../vic_installation/plugin_h5_vcsa.html) or [Install the HTML5 Plug-In on vCenter Server for Windows by Using a Web Server](../vic_installation/plugin_h5_vc_web.html) in *vSphere Integrated Containers Engine Installation*.

**Procedure**

1. If you run vCenter Server on Windows, copy the new version of the <code>com.vmware.vic-<i>version</i>.zip</code> file to the appropriate location on your Web server.
2. Update the new version of the `configs` file.

   - vCenter Server Appliance: `vic/ui/HTML5Client/configs`
   - vCenter Server on Windows: `vic/ui/vCenterForWindows/configs`
3. (Optional) If you are upgrading the plug-in on a vCenter Server Appliance and you are working on a Windows system, update the `vic/ui/HTML5Client/install.sh` file to point `PLUGIN_MANAGER_BIN` to the Windows UI executable. 
4. Run the `vic/ui/HTML5Client/upgrade.sh` or `vic/ui/vCenterForWindows/upgrade.bat` script. 
4. When the upgrade finishes, if you are logged into the vSphere Client, log out then log back in again.