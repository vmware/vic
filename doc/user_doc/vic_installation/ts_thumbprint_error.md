# VCH Deployment Fails with a Certificate Verification Error #

When you use `vic-machine create` to deploy a virtual container host (VCH), the installation fails with a certificate verification error.

## Problem ##
Deployment of the VCH fails during the validation of the configuration that you provided: 

<pre>
Failed to verify certificate for target=<i>vcenter_server_or_esxi_host</i>
(thumbprint=<i>vc_or_esxi_cert_thumbprint</i>)
Create cannot continue: failed to create validator
vic-machine-<i>platform</i>.exe failed: x509: certificate signed by unknown authority
</pre>

## Cause ##

This error can occur for the following reasons:

- The vCenter Server or ESXi host on which you are deploying the VCH uses untrusted certificates that have not been signed by a Certificate Authority (CA).
- The vCenter Server or ESXi host on which you are deploying the VCH uses trusted certificates that have been signed by a Certificate Authority (CA) and you are using an incorrect certificate on the machine on which you are running `vic-machine create`.
- There is a security problem and the connection has possibly been subject to a spoofing attack.
 
## Solution ##
If you are using untrusted certificates and you cannot use trusted certificates:

1. Copy the thumbprint of the untrusted certificate from the `vic-machine create` error message.
2. Run `vic-machine create` again, specifying the certificate thumbprint in the `--thumbprint` option.

If you are using trusted certificates, check that you have the correct client certificate.

