# Upgrade the Flex-Based vSphere Web Client Plug-In #

If you download a new version of vSphere Integrated Containers, upgrade the Flex-based vSphere Web Client plug-in for vSphere Integrated Containers.

**Prerequisites**

- You deployed an older version of the vSphere Integrated Containers plug-in for the Flex-based vSphere Web Client.
- You downloaded a new version of vSphere Integrated Containers Engine.
- For information about updating the `configs` and `install.sh` files, and where to copy the <code>com.vmware.vic.ui-<i>version</i>.zip</code> file or <code>com.vmware.vic.ui-<i>version</i></code> folder, see the topic that corresponds to your type of deployment in [Installing the vSphere Client Plug-ins](../vic_installation/install_vic_plugin.html) in *vSphere Integrated Containers Engine Installation*.

**Procedure**

1. Copy the new versions of the <code>com.vmware.vic.ui-<i>version</i>.zip</code> file or <code>com.vmware.vic.ui-<i>version</i></code> folder to the appropriate location on your Web server or vCenter Server system.
2. Update the new version of the `vic/ui/VCSA/configs` or `vic/ui/vCenterForWindows/configs` file.
3. (Optional) If you are upgrading the plug-in on a vCenter Server Appliance and you are working on a Windows system, update the `vic/ui/VCSA/install.sh` file to point `PLUGIN_MANAGER_BIN` to the Windows UI executable.
3. Run the `vic/ui/VCSA/upgrade.sh` or `vic/ui/vCenterForWindows/upgrade.bat` script. 
4. When the upgrade finishes, if you are logged into the vSphere Web Client, log out then log back in again.


