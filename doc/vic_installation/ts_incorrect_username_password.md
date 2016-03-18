# Login to Target Fails Due to an Incorrect Username or Password #
When you run the command line installer with the `-passwd` option, the installation fails due to an incorrect username or password.

## Problem ##
Command line installation with the `-passwd` option fails with an error when connecting to vCenter Server or the ESXi host. 

<pre>
INFO[<em>timestamp</em>] Connecting to target
      
ERROR[<em>timestamp</em>] Failed to create client: 
Failed to log in to https://<em>vcenter_server_or_esxi_address</em>/sdk: 
ServerFaultCode: Cannot complete login due to an incorrect user name 
or password.
FATAL[<em>timestamp</em>] Connecting to target failed. Exiting...
</pre>

## Cause ##

- You are installing on Mac OS X or Linux and you wrapped the password in double quotation marks (") or omitted the quotations marks.
- You are installing on Windows and you wrapped the password in single quotation marks (').
- You specified the incorrect username or password for the vCenter Server instance or ESXi host.

## Solution ##
- Verify that you are using the correct username and password for the vCenter Server instance or ESXi host.
- If you are installing on Mac OS X or Linux, wrap the password in single quotation marks (').
- You are installing on Windows, wrap the password in double quotation marks (") or omit the quotations marks.
- Run the installer command without the `-passwd` option and enter the password when the installer prompts you to do so.