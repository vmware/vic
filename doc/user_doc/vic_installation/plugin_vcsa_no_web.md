# Install the vSphere Integrated Containers Plug-In on a vCenter Server Appliance Without Access to a Web Server #

If you are running the vCenter Server Appliance and you do not have access to a Web server, you can manually install the vSphere Web Client plug-in for vSphere Integrated Containers.

**Prerequisites**

You deployed at least one virtual container host to a vCenter Server Appliance instance.

**Procedure**

1. On the system on which you run `vic-machine`, navigate to the folder that contains the `vic-machine` utility and open the `ui` folder.
3. Open the  <code><i>vic_unpack_dir</i>/vic/ui/VCSA/configs</code> file in a text editor.
4. Enter the IPv4 address or FQDN of the vCenter Server instance on which to install the plug-in. <pre>VCENTER_IP="<i>vcenter_server_address</i>"</pre>
6. Save and close the `configs` file.
7. Open a command prompt and navigate to <code><i>vic_unpack_dir</i>/vic/ui/VCSA</code>.
8. Run the installer.<pre>./install.sh</pre>Make sure that `install.sh` is executable by running `chmod` before you run it.
9. Enter the user name and password for the vCenter Server administrator account.
10. Enter the root password for the vCenter Server Appliance twice.

  The installer requires the root password of the vCenter Server Appliance twice: once to copy the files to the appliance over SSH, and once to set the correct ownership on the files and folders.
10. When installation finishes, if you are logged into the vSphere Web Client, log out then log back in again.