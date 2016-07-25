# Install the vSphere Integrated Containers Plug-In on vCenter Server For Windows by Using a Web Server #

If your vCenter Server instance runs on Windows, you can use a Web server to host the vSphere Web Client plug-in for vSphere Integrated Containers.

**Prerequisites**

- You deployed at least one virtual container host to a vCenter Server instance.
- Download the latest build of the vSphere Integrated Containers plug-in bundle from https://bintray.com/vmware/vic-repo/build-deps/installer-onsite-beta#files.
- You are running a Web server that your vCenter Server instance can access.

**Procedure**

1. Unpack the `vic-ui-installer.zip` bundle on the Windows machine on which vCenter Server is running.
2. Upload the plug-in bundle to your Web server.
  <pre><i>unpack_dir</i>\installer\vsphere-client-serenity\com.vmware.vicui.Vicui-0.0.1.zip</pre>
3. Open the <code><i>unpack_dir</i>\installer\vCenterForWindows\configs</code> file in a text editor.
4. Enter the IPv4 address or FQDN of the vCenter Server instance on which to install the plug-in.<pre>target_vcenter_ip=<i>vcenter_server_address</i></pre>
5. Enter the URL of the ZIP file on your Web server.<pre>vic_ui_host_url=<i>vic_web_server_location</i></pre>
6. (Optional) If you used an HTTPS address in `vic_ui_host_url`, provide the SHA-1 thumbprint of the Web server.<pre>vic_ui_host_thumbprint=<i>thumbprint</i></pre> 
6. Save and close the `configs` file.
7. Open a command prompt and navigate to <code><i>unpack_dir</i>\installer\vCenterForWindows</code>.
8. Run the installer.<pre>install.bat</pre>
  Make sure that you use the correct account to run `install.bat`. 
  - If vCenter Server uses the local system account, run `install.bat` with the local system account.
  - If vCenter Server uses a different user account, run `install.bat` with that account.
9. Enter the password for the vCenter Server administrator account.
10. When installation finishes, if you are logged into the vSphere Web Client, log out then log back in again.