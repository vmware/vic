# Access the Administration Portal for a VCH #

vSphere Integrated Containers Engine provides a Web-based administration portal for virtual container hosts (VCHs), called VCH Admin.

## Prerequisites ##

- You deployed a VCH.
- Obtain the address of the VCH:
  - Copy the address from the output of `vic-machine create` or `vic-machine inspect`. 
  - If you deployed the VCH to vCenter Server, copy the address from the **Summary** tab for the vSphere Integrated Containers Engine endpoint VM in the vSphere Web Client.
  - If you deployed the VCH to an ESXi host, copy the address from the **Summary** tab for the vSphere Integrated Containers Engine endpoint VM in the desktop vSphere Client.

## Procedure ##

### From a Web Browser ###

#### Username / Password Login ####
1. Go to https://<i>vch_address</i>:2378.

   If prompted about an insecure or not private connection, click Advanced and follow the prompts to proceed to the portal.   

2. Enter the username and password for the vCenter Server instance or ESXi host.


**Note:** If the VCH was installed with `--no-tls` or `--no-tlsverify` this is the only way to log in, but it is always available regardless of how the VCH was installed. 

#### Browser-based Certificate Login ####

If the VCH was installed with `--tls-cname` then browser-based certificate authentication may be used instead, if desired. 

Currently, only Firefox is known to work but support for this feature is expected to improve as support is added to the various browsers.

Look for the file named `cert.pfx` in the VCH folder (the folder with the same name as your VCH present in the directory from which you ran `vic-machine`) and install it in your browser. 

**Installing the `.pfx` certificate in Firefox:**

1. Navigate to `about:preferences` in the URL bar or click the menu and choose `Preferences`
2. Click `Advanced`
3. Click `View Certificates`
4. Click `Import`
5. Find the `cert.pfx` file in the file browser that opens and click `Open`
6. Do not enter a password when prompted; simply click `OK`. You should receive a message stating that the certificate was successfully installed.

Now that the certificate is installed, navigating to https://<i>vch_address</i>:2378/ or one of the log pages will no longer redirect you to the authentication page.


### From the Command Line ###

If you wish to access the log server using a tool such as `curl` or `wget`, authentication is most easily performed with certificates.

With `curl` this is as simple as specifying `--key vch-name/key.pem --certificate vch-name/cert.pem` as command line options to curl. The folder named `vch-name`, where `vch-name` is the same as provided to `vic-machine` at install time, will be present in the directory from which you ran `vic-machine`. `key.pem` and `cert.pem` will be present within.

**Example:**
<pre>curl https://<i>vch_address</i>:2378/logs/port-layer.log --key ./vch-name/key.pem --certificate vch-name/cert.pem</pre>
**Note:** If your certificates are self-signed, you may also need to provide the `-k` flag to `curl`.

If the VCH was *not* installed with full TLS (`--tls-cname`) then you must authenticate using the username/password form on the login page. 

1. First log in to gather an authentication cookie for subsequent access:<pre>curl -sk https://<i>vch_address</i>:2378/authentication -XPOST -F username=<i>your_username</i> -F password=<i>your_password</i> -D cookies_file</pre>
2. Then to actually access the logs, use the cookies from the last step along with a normal `curl` command, e.g.:<pre>curl -sk https://<i>vch_address</i>:2378/logs/port-layer.log -b cookies_file</pre>


## Result ##

The VCH Admin portal displays information about the VCH and the environment in which is running:

- Status information about the VCH, registry and Internet connections,  firewall configuration, and license. For information about these statuses and how to remedy error states, see the [VCH Status Reference](vicadmin_status_ref.md).
- The address of the Docker endpoint.
- The system time of the VCH. This is useful to know because clock skews between VCHs and client systems can cause TLS authentication to fail. For information about clock skews, see [Connections Fail with Certificate Errors when Using Full TLS Authentication with Trusted Certificates](ts_clock_skew.md). 
- The remaining capacity of the datastore that you designated as the image store. If the VCH is unable to connect to vSphere, the datastore information is not displayed.
- Live logs and log bundles for different aspects of the VCH. For information about the logs, see [Access vSphere Integrated Containers Engine Log Bundles](log_bundles.md).
