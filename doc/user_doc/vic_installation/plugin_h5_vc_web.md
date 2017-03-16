# Install the HTML5 Plug-In on vCenter Server for Windows by Using a Web Server #

If your vCenter Server instance runs on Windows, you can use a Web server to host the HTML5 vSphere Client plug-in for vSphere Integrated Containers Engine.

**Prerequisites**

- You are running a Web server that your vCenter Server instance can access.
- Use a Windows system to run the script to install the plug-in on a vCenter Server that runs on Windows. If you used a Linux or Mac OS system to deploy the VCH, download and unpack the vSphere Integrated Containers Engine package on a Windows system. For example, download the package to the system on which vCenter Server is running.

**Procedure**

1. On the Windows system on which you have downloaded and unpacked vSphere Integrated Containers Engine, navigate to the folder that contains the `vic-machine` utility and open the `ui` folder.
2. Upload the plug-in bundle to your Web server.
    <pre><i>unpack_dir</i>\vic\ui\plugin-packages\com.vmware.vic-<i>version</i>.zip</pre>
3. On the system on which the Web server is running, make sure that port 3443 is open for both inbound and outbound connections.
3. On the `vic-machine` system, open the <code><i>vic_unpack_dir</i>\vic\ui\vCenterForWindows\configs</code> file in a text editor.
4. Enter the IPv4 address or FQDN of the vCenter Server instance on which to install the plug-in.<pre>SET target_vcenter_ip=<i>vcenter_server_address</i></pre>
5. Enter the path to the folder on your Web server that contains the <code>com.vmware.vic-<i>version</i>.zip</code> file. For example, if you saved the ZIP file in a folder named `vsphere-plugins`, set the following URL: <pre>SET vic_ui_host_url="<i>web_server_address</i>:3443/vsphere-plugins"</pre>

      **NOTE**: Only include the path to the ZIP file. Do not include <code>com.vmware.vic-<i>version</i>.zip</code> in the path. 

6. (Optional) If you used an HTTPS address in `vic_ui_host_url`, provide the SHA-1 thumbprint of the Web server.<pre>SET vic_ui_host_thumbprint="<i>thumbprint</i>"</pre>**NOTE**: Use colon delimitation in the thumbprint. Do not use space delimitation. 
6. Save and close the `configs` file.
7. Open a command prompt, navigate to <code><i>vic_unpack_dir</i>\vic\ui\vCenterForWindows</code>, and run the installer.<pre>install.bat</pre>
9. Enter the user name and password for the vCenter Server administrator account.
10. When installation finishes, if you are logged into the vSphere Client, log out then log back in again.
11. When you have logged back in to the vSphere Client, select **Home**.

   The **vSphere Integrated Containers** icon should appear in the Inventories section.

**What to Do Next**
If the vSphere Integrated Containers icon does not appear in the **Home** view, restart the vSphere Client service. For instructions about how to restart the vSphere Client service, see [vSphere Integrated Containers Engine Plug-In Does Not Appear](ts_ui_not_appearing.md).