# VCH Deployment Fails with a Certificate Verification Error #

When you use `vic-machine create` to deploy a virtual container host, the installation fails with a certificate verification error.

## Problem ##
Deployment of the virtual container host fails during the validation of the configuration that you provided: 

<pre>
Failed to verify certificate for target=<i>vcenter_server_or_esxi_host</i>
(thumbprint=<i>vc_or_esxi_cert_thumbprint</i>)
Create cannot continue: failed to create validator
vic-machine-<i>platform</i>.exe failed: x509: certificate signed by unknown authority
</pre>

## Cause ##
The vCenter Server or ESXi host on which you are deploying the virtual container host uses untrusted certificates that have not been signed by a Certificate Authority (CA).
 
## Solution ##
If you cannot use trusted certificates:

1. Copy the thumbprint of the untrusted certificate from the `vic-machine create` error message.
2. Run `vic-machine create` again, specifying the certificate thumbprint in the `--thumbprint` option.