# Browser-Based Certificate Login #

If you deployed the VCH with client and server authentication by using `--tls-cname` or by specifying a static IP address on the client network, you can use browser-based certificate authentication to access the VCH Admin Portal. In this way, you do not need to provide the  vSphere credentials each time that you log in to VCH Admin.

**Prerequisites**

- You deployed a VCH with `--tls-cname` or a static IP address for the VCH on the client network.
- Use Firefox. Currently, this feature is only supported with Firefox.
- Locate the file named `cert.pfx` on the system on which you ran `vic-machine create`. The `cert.pfx` is located in either of the following locations:

  - In the folder with the same name as the VCH, in the directory from which you ran `vic-machine create`.
  - In a folder that you specified in the `vic-machine create --cert-path` option.

**Procedure** 

1. In Firefox, select `Tools` > `Options` and select  `Advanced`.
3. Click `View Certificates`.
4. Click `Import`.
5. Browse to the `cert.pfx` file and click `Open`.
6. Click `OK`. 

   Do not enter a password when prompted. 

**Result**

You see a message stating that the certificate was successfully installed.
With the VCH certificate installed in your browser, you can navigate to https://<i>vch_address</i>:2378/ or to one of the log pages without having to enter the vSphere credentials.


