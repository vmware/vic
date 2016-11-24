# Install the vSphere Integrated Containers Engine Plug-In on vCenter Server for Windows Without Access to a Web Server #

You can install the vSphere Web Client plug-in for vSphere Integrated Containers Engine on a vCenter Server instance for Windows that does not have access to a Web Server.

**Prerequisites**

- You deployed at least one virtual container host to a vCenter Server instance that runs on Windows.
- You must use a Windows system to run the script to install the plug-in on a vCenter Server that runs on Windows. If you used a Linux or Mac OS system to deploy the virtual container host, download and unpack the vSphere Integrated Containers Engine package on a Windows system. For example, download the package to the system on which vCenter Server is running.
- If you deployed the virtual container host to a vCenter Server 6.5 instance, use the Flash-based vSphere Web Client to view the vSphere Web Client plug-in for vSphere Integrated Containers Engine. vSphere Integrated Containers Engine does not currently provide a plug-in for the new HTML5 vSphere Client.

**Procedure**

1. On the Windows system on which you have downloaded and unpacked vSphere Integrated Containers Engine, navigate to the folder that contains the `vic-machine` utility and open the `ui` folder.
2. Copy the <code>com.vmware.vicui.Vicui-<i>version</i></code> folder into the folder on the vCenter Server system that contains the vSphere Web Client packages.
  
  - Source location on `vic-machine` system: <pre><i>vic_unpack_dir</i>\vic\ui\vsphere-client-serenity</pre>
  - Destination location on vCenter Server Windows system: <pre><i>instl_dir</i>\vCenterServer\cfg\vsphere-client\vc-packages\vsphere-client-serenity</pre>

    <code><i>instl_dir</i></code> is the location in which vCenter Server is installed. If the `vc-packages\vsphere-client-serenity` folders do not exist under the <code>vsphere-client</code> folder, create them manually.
3. On the `vic-machine` system, open the <code><i>vic_unpack_dir</i>\vic\ui\vCenterForWindows\configs</code> file in a text editor.
4. Enter the IPv4 address or FQDN of the vCenter Server instance on which to install the plug-in.<pre>SET target_vcenter_ip=<i>vcenter_server_address</i></pre>
6. Save and close the `configs` file.
7. Open a command prompt, navigate to <code><i>vic_unpack_dir</i>\vic\ui\vCenterForWindows</code>, and run the installer.<pre>install.bat</pre>
9. Enter the user name and password for the vCenter Server administrator account.
10. When installation finishes, if you are logged into the vSphere Web Client, log out then log back in again.

**What to Do Next**
Check that the deployment has succeeded by following the procedure in [Verify the Deployment of the vSphere Integrated Containers Engine Plug-In](plugin_verify_deployment.md).