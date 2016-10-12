# Install the vSphere Integrated Containers Engine Plug-In on a vCenter Server Appliance by Using a Web Server #

If you are running the vCenter Server Appliance, you can use a Web server to host the vSphere Web Client plug-in for vSphere Integrated Containers Engine.

**Prerequisites**

- You deployed at least one virtual container host to a vCenter Server Appliance instance.
- You are running a Web server that the vCenter Server Appliance can access.

**Procedure**

1. On the system on which you run `vic-machine`, navigate to the folder that contains the `vic-machine` utility and open the `ui` folder.
2. Upload the plug-in bundle to your Web server.
  <pre><i>vic_unpack_dir</i>/vic/ui/vsphere-client-serenity/com.vmware.vicui.Vicui-0.0.1.zip</pre>
3. Open the  <code><i>vic_unpack_dir</i>/vic/ui/VCSA/configs</code> file in a text editor.
4. Enter the IPv4 address or FQDN of the vCenter Server instance on which to install the plug-in.<pre>VCENTER_IP="<i>vcenter_server_address</i>"</pre>
5. Enter the path to the folder on your Web server that contains the `com.vmware.vicui.Vicui-0.0.1.zip`  file.<pre>VIC_UI_HOST_URL="<i>vicui_zip_location</i>"</pre>
6. (Optional) If you used an HTTPS address in `VIC_UI_HOST_URL`, provide the SHA-1 thumbprint of the Web server.<pre>VIC_UI_HOST_THUMBPRINT="<i>thumbprint</i>"</pre> 
6. (Optional) If you are deploying the plug-in to a vCenter Server 5.5 instance, change the  value of `IS_VCENTER_5_5` from 0 to 1. 

  **IMPORTANT**: Deploying vSphere Integrated Containers Engine to vSphere 5.5 environments works but is unsupported.
  <pre>IS_VCENTER_5_5=1</pre>
6. Save and close the `configs` file.
7. (Optional) If you run `vic-machine` on a Windows system, open  the <code><i>vic_unpack_dir</i>/vic/ui/VCSA/install.sh</code> file in a text editor and point `PLUGIN_MANAGER_BIN` to the Windows UI executable.

 - Before:<pre>if [[ $(echo $OS | grep -i "darwin") ]] ; then
    PLUGIN_MANAGER_BIN="../../vic-ui-darwin"
else
    PLUGIN_MANAGER_BIN="../../vic-ui-linux"</pre>
  - After:<pre>if [[ $(echo $OS | grep -i "darwin") ]] ; then
    PLUGIN_MANAGER_BIN="../../vic-ui-darwin"
else
    PLUGIN_MANAGER_BIN="../../vic-ui-windows"</pre>

7. Open a command prompt, navigate to <code><i>vic_unpack_dir</i>/vic/ui/VCSA</code>, and run the installer.
   <pre>./install.sh</pre>
  - Make sure that `install.sh` is executable by running `chmod` before you run it.
  - On Windows systems, run `install.sh` in a UNIX shell that supports SSH and SCP, for example Cygwyn or Git Bash. Do not use Windows 10 native Bash.
9. Enter the user name and password for the vCenter Server administrator account.
10. When installation finishes, if you are logged into the vSphere Web Client, log out then log back in again.

**What to Do Next**
Check that the deployment has succeeded by following the procedure in [Verify the Deployment of the vSphere Integrated Containers Engine Plug-In](plugin_verify_deployment.md).