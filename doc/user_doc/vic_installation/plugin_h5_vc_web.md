# Install the HTML5 Plug-In on vCenter Server for Windows by Using a Web Server #

If you are running vCenter Server 6.5 on Windows, you must use a Web server to deploy the HTML5 vSphere Client plug-in for vSphere Integrated Containers.

**Prerequisites**

- You are running a Web server that your vCenter Server instance can access over HTTPS.
- On the system on which the Web server is running, if you do not use port 443 for secure HTTP traffic, make sure that the appropriate port is open for both inbound and outbound connections.
- Use a Windows system to run the script to install the plug-in on a vCenter Server that runs on Windows. If you usually use a Linux or Mac OS system to run `vic-machine`, download and unpack the vSphere Integrated Containers Engine package on a Windows system. For example, download the package to the system on which vCenter Server is running.

**Procedure**

1. On the Windows system on which you have downloaded and unpacked vSphere Integrated Containers Engine, open the `\vic\ui\plugin-packages\` folder.
2. Upload the ZIP bundle for the HTML5 plug-in to your Web server.<pre>com.vmware.vic-<i>version</i>.zip</pre>
3. Open the `\vic\ui\vCenterForWindows\configs` file in a text editor.
4. Enter the IPv4 address or FQDN of the vCenter Server instance on which to install the plug-in.<pre>SET target_vcenter_ip=<i>vcenter_server_address</i></pre>
5. Enter the path to the folder on your Web server that contains the <code>com.vmware.vic-<i>version</i>.zip</code> file. For example, if you saved the ZIP file on your Web server in a folder named `vsphere-plugins`, set one of the following URLs, depending on whether you use port 443 or a different port: <pre>SET vic_ui_host_url="https://<i>web_server_address</i>/vsphere-plugins"</pre><pre>SET vic_ui_host_url="https://<i>web_server_address</i>:<i>port</i>/vsphere-plugins"</pre>

      **NOTE**: Only include the path to the ZIP file. Do not include <code>com.vmware.vic-<i>version</i>.zip</code> in the path. 

6. Provide the SHA-1 thumbprint of the Web server.<pre>SET vic_ui_host_thumbprint="<i>thumbprint</i>"</pre>**NOTE**: Use colon delimitation in the thumbprint. Do not use space delimitation. 
6. Save and close the `configs` file.
7. Open a command prompt, navigate to `\vic\ui\vCenterForWindows`, and run the installer.<pre>install.bat</pre>
9. Enter the user name and password for the vCenter Server administrator account.
10. When installation finishes, if you are logged into the vSphere Client, log out then log back in again.
11. When you have logged back in to the vSphere Client, select **Home**.

**Result**

   The **vSphere Integrated Containers** icon should appear in the Inventories section.

**What to Do Next**

If the vSphere Integrated Containers icon does not appear in the **Home** view, restart the vSphere Client service. For instructions about how to restart the vSphere Client service, see [vSphere Integrated Containers Plug-In Does Not Appear](ts_ui_not_appearing.md).