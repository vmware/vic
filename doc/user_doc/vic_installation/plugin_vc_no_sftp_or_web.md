# Install the vSphere Integrated Containers Plug-In on vCenter Server For Windows, without SFTP or a Web Server #

You can install the vSphere Web Client plug-in for vSphere Integrated Containers on a vCenter Server instance that has neither access to an SFTP service nor to a Web Server.

**Prerequisites**

- You deployed at least one virtual container host to a vCenter Server instance.
- Download the latest build of the vSphere Integrated Containers plug-in bundle from https://bintray.com/vmware/vic-repo/build-deps/installer-onsite-beta#files.

**Procedure**

1. Unpack the `vic-ui-installer.zip` bundle on the Windows machine on which vCenter Server is running.
2. Copy the `com.vmware.vicui.Vicui-0.0.1` folder into the folder that contains the vSphere Web Client packages.
  
  - Source folder: <pre><i>unpack_dir</i>\installer\vsphere-client-serenity</pre>
  - Destination folder: <pre><i>instl_dir</i>\vCenterServer\cfg\vsphere-client\vc-packages\vsphere-client-serenity</pre>

    <code><i>instl_dir</i></code> is the location in which vCenter Server is installed. If the `vc-packages\vsphere-client-serenity` folders do not exist under the <code>vsphere-client</code> folder, create them manually.
3. Open the <code><i>unpack_dir</i>\installer\vCenterForWindows\configs</code> file in a text editor.
4. Enter the IPv4 address or FQDN of the vCenter Server instance on which to install the plug-in.<pre>SET target_vcenter_ip=<i>vcenter_server_address</i></pre>
5. Deactivate SFTP by changing the value of `sftp_supported` to 0.<pre>SET sftp_supported=0</pre>
6. Save and close the `configs` file.
7. Open a command prompt and navigate to <code><i>unpack_dir</i>\installer\vCenterForWindows</code>.
8. Run the installer.<pre>install.bat</pre>
  Make sure that you use the correct account to run `install.bat`. 
  - If vCenter Server uses the local system account, run `install.bat` with the local system account.
  - If vCenter Server uses a different user account, run `install.bat` with that account.
9. Enter the password for the vCenter Server administrator account.
10. When installation finishes, if you are logged into the vSphere Web Client, log out then log back in again.