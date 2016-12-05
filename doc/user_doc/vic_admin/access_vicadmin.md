# Access the Administration Portal for a VCH #

vSphere Integrated Containers Engine provides a Web-based administration portal for virtual container hosts (VCHs), called VCH Admin.

If you deployed the VCH with `--no-tls` or `--no-tlsverify`, you can only log in to VCH Admin by specifying the username and password of the ESXi host or vCenter Server on which you deployed the VCH. If you deployed the VCH with client and server authentication by using `--tls-cname` or by specifying a static IP address on the client network, you can use the generated `*.pfx` certificate to authenticate with the VCH Admin portal. For information about using the `*.pfx` certificate to log into VCH admin, see [Browser-Based Certificate Login](browser_login.md) and [Command Line Certificate Login](cmdline_login.md).

**Prerequisites**

- You deployed a VCH.
- Obtain the address of the VCH:
  - Copy the address from the output of `vic-machine create` or `vic-machine inspect`. 
  - If you deployed the VCH to vCenter Server, copy the address from the **Summary** tab for the vSphere Integrated Containers Engine endpoint VM in the vSphere Web Client.
  - If you deployed the VCH to an ESXi host, copy the address from the **Summary** tab for the vSphere Integrated Containers Engine endpoint VM in the desktop vSphere Client.

**Procedure**

1. Go to https://<i>vch_address</i>:2378.

   If prompted about an insecure or not private connection, click Advanced and follow the prompts to proceed to the portal.   

2. Enter the username and password for the vCenter Server instance or ESXi host.

**Result**

The VCH Admin portal displays information about the VCH and the environment in which is running:

- Status information about the VCH, registry and Internet connections,  firewall configuration, and license. For information about these statuses and how to remedy error states, see the [VCH Status Reference](vicadmin_status_ref.md).
- The address of the Docker endpoint.
- The system time of the VCH. This is useful to know because clock skews between VCHs and client systems can cause TLS authentication to fail. For information about clock skews, see [Connections Fail with Certificate Errors when Using Full TLS Authentication with Trusted Certificates](ts_clock_skew.md). 
- The remaining capacity of the datastore that you designated as the image store. If the VCH is unable to connect to vSphere, the datastore information is not displayed.
- Live logs and log bundles for different aspects of the VCH. For information about the logs, see [Access vSphere Integrated Containers Engine Log Bundles](log_bundles.md).
