# Install the HTML5 Plug-In on a vCenter Server Appliance #

If you are running version 6.5 of the vCenter Server Appliance, you can install the HTML5 vSphere Client plug-in for vSphere Integrated Containers.

**Prerequisites**

Go to the vCenter Server Appliance Management Interface (VAMI) at https://<i>vcsa_address</i>:5480, click **Access**, and make sure that Bash Shell is enabled.

**Procedure**

1. On the system on which you have downloaded and unpacked vSphere Integrated Containers Engine, open the `/vic/ui/HTML5Client/configs` file in a text editor.
4. Enter the IPv4 address or FQDN of the vCenter Server Appliance on which to install the plug-in. <pre>VCENTER_IP="<i>vcenter_server_address</i>"</pre>
6. Save and close the `configs` file.
7. (Optional) If you unpacked vSphere Integrated Containers Engine on a Windows system, open  the `/vic/ui/HTML5Client/install.sh` file in a text editor and point `PLUGIN_MANAGER_BIN` to the Windows UI executable.

   Before:
     <pre>if [[ $(echo $OS | grep -i "darwin") ]] ; then
       PLUGIN_MANAGER_BIN="../../vic-ui-darwin"
     else
       PLUGIN_MANAGER_BIN="../../vic-ui-linux"</pre>
   After:
      <pre>if [[ $(echo $OS | grep -i "darwin") ]] ; then
       PLUGIN_MANAGER_BIN="../../vic-ui-darwin"
      else
       PLUGIN_MANAGER_BIN="../../vic-ui-windows"</pre>

7. Open a command prompt, navigate to `/vic/ui/VCSA`, and run the installer.
   <pre>./install.sh</pre>

    - Make sure that `install.sh` is executable by running `chmod` before you run it.
    - On Windows systems, run `install.sh` in a UNIX shell that supports SSH and SCP, for example Cygwyn or Git Bash. Do not use Windows 10 native Bash.
  
9. Enter the user name and password for the vCenter Server administrator account.

10. Enter the `root` password for the vCenter Server Appliance.

    The installer requires the `root` password of the vCenter Server Appliance three times: 
     - Once to check whether the Bash shell is enabled on the vCenter Server Appliance. If the Bash shell is not enabled, the installation fails.
     - Once to copy the files to the appliance over SSH.
     - Once to set the correct ownership on the files and folders.
10. When installation finishes, if you are logged into the HTML5 vSphere Client, log out then log back in again.
11. When you have logged back in to the vSphere Client, select **Home**.

**Result**

The **vSphere Integrated Containers** icon should appear in the Inventories section.

**What to Do Next**

If the vSphere Integrated Containers icon does not appear in the **Home** view, restart the vSphere Client service. For instructions about how to restart the vSphere Client service, see [vSphere Integrated Containers Plug-In Does Not Appear](ts_ui_not_appearing.md).