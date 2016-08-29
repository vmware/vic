# Install the vSphere Integrated Containers Engine Plug-In on vCenter Server For Windows by Using SFTP #

If your vCenter Server instance runs on a Windows system that has access to an SFTP service, you can use SFTP to install the vSphere Web Client plug-in for vSphere Integrated Containers Engine.

**Prerequisites**

- You deployed at least one virtual container host to a vCenter Server instance that runs on Windows.
- Your vCenter Server instance is running on a Windows system that has an SFTP service.

**Procedure**

1. On the system on which you run `vic-machine`, navigate to the folder that contains the `vic-machine` utility and open the `ui` folder.
3. Open the <code><i>vic_unpack_dir</i>\vic\ui\vCenterForWindows\configs</code> file in a text editor.
4. Enter the IPv4 address or FQDN of the vCenter Server instance on which to install the plug-in.<pre>SET target_vcenter_ip=<i>vcenter_server_address</i></pre>
5. Make sure that `sftp_supported` is set to 1.<pre>SET sftp_supported=1</pre>
6. Provide the SFTP user name and password.
  <pre>SET sftp_username=<i>username</i>
SET sftp_password=<i>password</i></pre> 
6. Provide the location of the `vsphere-client-serenity` folder on the vCenter Server Windows system relative to the root folder of the SFTP connection.

  For example, if the root folder of the SFTP connection is <code><i>vcenter_server_install_dir</i>\vCenterServer\cfg</code>, set the relative path as follows:  

  <pre>target_vc_packages_path=/vsphere-client/vc-packages/vsphere-client-serenity/</pre>

7. Save and close the `configs` file.
7. Open a command prompt and navigate to <code><i>vic_unpack_dir</i>\vic\ui\vCenterForWindows</code>.
8. Run the installer.<pre>install.bat</pre>
  Make sure that you use the correct account to run `install.bat`. 
  - If vCenter Server uses the local system account, run `install.bat` with the local system account.
  - If vCenter Server uses a different user account, run `install.bat` with that account.
9. Enter the user name and password for the vCenter Server administrator account.
10. When installation finishes, if you are logged into the vSphere Web Client, log out then log back in again.