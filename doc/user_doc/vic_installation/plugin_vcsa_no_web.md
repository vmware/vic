# Install the Flex Plug-In on a vCenter Server Appliance Without Access to a Web Server #

If you are running the vCenter Server Appliance 6.0 or 6.5 and you do not have access to a Web server, you can install the Flex-based vSphere Web Client plug-in for vSphere Integrated Containers by running an installation script.

**Prerequisites**

- You deployed at least one virtual container host (VCH) to a vCenter Server Appliance instance.
- Go to the vCenter Server Appliance Management Interface (VAMI) at https://<i>vcsa_address</i>:5480, click **Access**, and make sure that Bash Shell is enabled.

**Procedure**

1. On the system on which you have downloaded and unpacked vSphere Integrated Containers Engine, open the `/vic/ui/VCSA/configs` file in a text editor.
4. Enter the IPv4 address or FQDN of the vCenter Server instance on which to install the plug-in. <pre>VCENTER_IP="<i>vcenter_server_address</i>"</pre>
6. Save and close the `configs` file.
7. (Optional)  If you unpacked vSphere Integrated Containers Engine on a Windows system, open  the `/vic/ui/VCSA/install.sh` file in a text editor and point `PLUGIN_MANAGER_BIN` to the Windows UI executable.

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

10. Enter the root password for the vCenter Server Appliance.

    The installer requires the root password of the vCenter Server Appliance three times: 
     - Once to check whether the Bash shell is enabled on the vCenter Server Appliance. If the Bash shell is not enabled, the installation fails and the installer provides remedial instructions.
     - Once to copy the files to the appliance over SSH.
     - Once to set the correct ownership on the files and folders.
10. When installation finishes, if you are logged into the vSphere Web Client, log out then log back in again.

**What to Do Next**

Check that the deployment has succeeded by following the procedure in [Verify the Deployment of the vSphere Integrated Containers Plug-In](plugin_verify_deployment.md).