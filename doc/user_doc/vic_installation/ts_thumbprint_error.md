# VCH Deployment Fails with a Certificate Verification Error #

When you use `vic-machine create` to deploy a virtual container host (VCH), the installation fails with a certificate verification error, noting that it `failed to create validator`.

## Problem ##
Deployment of the VCH fails during the validation of the configuration that you provided: 

<pre>
Failed to verify certificate for target=<i>vcenter_server_or_esxi_host</i>
(thumbprint=<i>vc_or_esxi_cert_thumbprint</i>)
Create cannot continue: failed to create validator
vic-machine-<i>platform</i>.exe failed: x509: certificate signed by unknown authority
</pre>

## Cause ##
The certificate on the vCenter Server or ESXi host that you specified in the `--target` option cannot be validated on the client system.
 
## Solution ##

If the certificate was signed by a certificate authority (CA), add that CA to the trusted roots for the client system.

If the CA should not be generally trusted, or the certificate is self-signed:
*  If the server is trusted and you did not specify the certificate thumbprint when you ran `vic-machine create`, specify the `--thumbprint` option, using the thumbprint from the error message.
*  If the thumbprint that you specified in `--thumbprint` does not match the server certificate reported in the error message: 
    1.  Remove the thumbprint from the `vic-machine create` command. 
    **WARNING:** A thumbprint mismatch could mean the server you have connected to is not the intended target and might have be spoofed.
    2. Validate that the change in server certificate is legitimate
    3. Re-run `vic-machine create`, specifying the new thumbprint in the `--thumbprint` option.
