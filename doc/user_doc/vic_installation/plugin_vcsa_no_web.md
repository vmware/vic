# Install the vSphere Integrated Containers Plug-In on a vCenter Server Appliance Without Access to a Web Server #

If you are running the vCenter Server Appliance and you do not have access to a Web server, you can manually install the vSphere Web Client plug-in for vSphere Integrated Containers.

**Prerequisites**

- You deployed at least one virtual container host to a vCenter Server instance.
- Download the latest build of the vSphere Integrated Containers plug-in bundle from https://bintray.com/vmware/vic-repo/build-deps/installer-onsite-beta#files.

**Procedure**

1. Unpack the `vic-ui-installer.zip` bundle on any system.
2. (Optional) Connect to the vCenter Server Appliance via SSH and copy the contents of the `vic-ui-installer.zip` bundle to a temporary location on the vCenter Server Appliance. 

  Perform this step if you unpacked the `vic-ui-installer.zip` bundle on a Windows system. If you unpacked `vic-ui-installer.zip` on a Mac OS or Linux system, you do not need to copy the files to the vCenter Server Appliance.
3. Open the <code><i>unpack_dir</i>\installer\VCSA\configs</code> file in a text editor.
4. Enter the IPv4 address or FQDN of the vCenter Server instance on which to install the plug-in. <pre>VCENTER_IP="<i>vcenter_server_address</i>"</pre>
6. Save and close the `configs` file.
7. Navigate to <code><i>unpack_dir</i>/installer/VCSA</code>.
8. Run the installer.<pre>./install.sh</pre>Make sure that `install.sh` is executable by running `chmod` before you run it.
9. Enter the password for the vCenter Server administrator account.
10. Enter the root password for the vCenter Server Appliance twice.

  The installer requires the root password of the vCenter Server Appliance twice: once to copy the files to the appliance over SSH, and once to set the correct ownership on the files and folders.
10. When installation finishes, if you are logged into the vSphere Web Client, log out then log back in again.