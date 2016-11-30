# Install the vSphere Integrated Containers Engine Plug-In on a vCenter Server Appliance Without Access to a Web Server #

If you are running the vCenter Server Appliance and you do not have access to a Web server, you can manually install the vSphere Web Client plug-in for vSphere Integrated Containers Engine.

**Prerequisites**

- You deployed at least one virtual container host (VCH) to a vCenter Server Appliance instance.
- If you deployed the VCH to a vCenter Server 6.5 instance, use the Flash-based vSphere Web Client to view the vSphere Web Client plug-in for vSphere Integrated Containers Engine. vSphere Integrated Containers Engine does not currently provide a plug-in for the new HTML5 vSphere Client.

**Procedure**

1. On the system on which you run `vic-machine`, open the <code><i>vic_unpack_dir</i>/vic/ui/VCSA/configs</code> file in a text editor.
4. Enter the IPv4 address or FQDN of the vCenter Server instance on which to install the plug-in. <pre>VCENTER_IP="<i>vcenter_server_address</i>"</pre>
6. Save and close the `configs` file.
7. (Optional) If you run `vic-machine` on a Windows system, open  the <code><i>vic_unpack_dir</i>/vic/ui/VCSA/install.sh</code> file in a text editor and point `PLUGIN_MANAGER_BIN` to the Windows UI executable.

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

7. Open a command prompt, navigate to <code><i>vic_unpack_dir</i>/vic/ui/VCSA</code>, and run the installer.
   <pre>./install.sh</pre>
    - Make sure that `install.sh` is executable by running `chmod` before you run it.
    - On Windows systems, run `install.sh` in a UNIX shell that supports SSH and SCP, for example Cygwyn or Git Bash. Do not use Windows 10 native Bash.
  
9. Enter the user name and password for the vCenter Server administrator account.

10. Enter the root password for the vCenter Server Appliance.

    The installer requires the root password of the vCenter Server Appliance three times: 
     - Once to check whether the Bash shell is enabled on the vCenter Server Appliance. If the Bash shell is not enabled, the installation fails and the installer provides remedial instructions.
     - Once to copy the files to the appliance over SSH.
     - Once to set the correct ownership on the files and folders.
10. When installation finishes, if you are logged into the vSphere Web Client, log out then log back in again.

**What to Do Next**
Check that the deployment has succeeded by following the procedure in [Verify the Deployment of the vSphere Integrated Containers Engine Plug-In](plugin_verify_deployment.md).