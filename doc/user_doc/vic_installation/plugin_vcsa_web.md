# Install the Flex Plug-In on a vCenter Server Appliance by Using a Web Server #

If you are running the vCenter Server Appliance 6.0 or 6.5, you can use a Web server to host the Flex-based vSphere Web Client plug-in for vSphere Integrated Containers.

**Prerequisites**

- You deployed at least one virtual container host (VCH) to a vCenter Server  Appliance instance.
- You are running a Web server that the vCenter Server Appliance can access over HTTPS.
- On the system on which the Web server is running, if you do not use port 443 for secure HTTP traffic, make sure that the appropriate port is open for both inbound and outbound connections.
- Go to the vCenter Server Appliance Management Interface (VAMI) at https://<i>vcsa_address</i>:5480, click **Access**, and make sure that Bash Shell is enabled.

**Procedure**

1. On the system on which you have downloaded and unpacked vSphere Integrated Containers Engine, open the `\vic\ui\vsphere-client-serenity\` folder.
2. Upload the ZIP bundle for the Flex plug-in to your Web server.<pre>com.vmware.vic.ui-<i>version</i>.zip</pre>
3. Open the `/vic/ui/VCSA/configs` file in a text editor.
4. Enter the IPv4 address or FQDN of the vCenter Server instance on which to install the plug-in.<pre>VCENTER_IP="<i>vcenter_server_address</i>"</pre>
5. Enter the path to the folder on your Web server that contains the <code>com.vmware.vic.ui-<i>version</i>.zip</code> file. For example, if you saved the ZIP file on your Web server in a folder named `vsphere-plugins`, set one of the following URLs, depending on whether you use port 443 or a different port: <pre>SET vic_ui_host_url="https://<i>web_server_address</i>/vsphere-plugins"</pre><pre>SET vic_ui_host_url="https://<i>web_server_address</i>:<i>port</i>/vsphere-plugins"</pre>

      **NOTE**: Only include the path to the ZIP file. Do not include <code>com.vmware.vic.ui-<i>version</i>.zip</code> in the path. 

6. Provide the SHA-1 thumbprint of the Web server.<pre>VIC_UI_HOST_THUMBPRINT="<i>thumbprint</i>"</pre>**NOTE**: Use colon delimitation in the thumbprint. Do not use space delimitation. 
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
10. When installation finishes, if you are logged into the vSphere Web Client, log out then log back in again.

**What to Do Next**

Check that the deployment has succeeded by following the procedure in [Verify the Deployment of the vSphere Integrated Containers Plug-In](plugin_verify_deployment.md).