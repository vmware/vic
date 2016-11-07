# Upgrade the vSphere Integrated Containers Engine Plug-In #

If you download a new version of vSphere Integrated Containers Engine, upgrade the vSphere Web Client plug-in for vSphere Integrated Containers Engine.

**Prerequisites**

- You deployed an older version of the vSphere Web Client plug-in for vSphere Integrated Containers Engine.
- You downloaded a new version of vSphere Integrated Containers Engine.

For information about updating the `configs` file and where to copy the <code>com.vmware.vicui.Vicui-<i>version</i>.zip</code> file or <code>com.vmware.vicui.Vicui-<i>version</i></code> folder, see the topic that corresponds to your type of deployment in [Installing the vSphere Web Client Plug-in for vSphere Integrated Containers Engine](../vic_installation/install_vic_plugin.html) in *vSphere Integrated Containers Engine Installation*.

**Procedure**

1. Copy the new versions of the <code>com.vmware.vicui.Vicui-<i>version</i>.zip</code> file or <code>com.vmware.vicui.Vicui-<i>version</i></code> folder to the appropriate location on your Web server or vCenter Server system.
2. Update the new version of the `ui/VCSA/configs` or `ui/vCenterForWindows/configs` file.
3. Run the `ui/VCSA/upgrade.sh` or `ui/vCenterForWindows/upgrade.bat` script.
4. When the upgrade finishes, if you are logged into the vSphere Web Client, log out then log back in again.


