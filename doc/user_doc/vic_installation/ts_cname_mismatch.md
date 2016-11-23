# VCH Deployment Fails with Certificate `cname` Mismatch #

When you use `vic-machine create` to deploy a virtual container host, the deployment fails with an error about the certificate `cname` value.

## Problem ##
Deployment fails during the validation of the configuration that you provided:

<pre>Provided cname does not match that in existing server certificate: <i>cname</i>
Unable to load certificates: cname option doesn't match existing server certificate 
in certificate path <i>path_to_certificate</i></pre>

## Cause ##
This error occurs in the following circumstances:
- You specified the `--tls-cname` option or specified a static IP address in the `--client-network-ip` option, to serve as the Common Name attribute in auto-generated CA certificates.
- `vic-machine create` finds and attempts to use an existing auto-generated certificate. 
- The existing certificate includes a different Common Name attribute to the one that you specified in `--tls-cname` or `--client-network-ip`. 

`vic-machine create` is attempting to use the existing certificate for one of the following reasons:

- You specified a virtual container host name in the `--name` option that is the same as that of an existing virtual container host. The certificate folder for the existing virtual container host has the same name as the one you specified in the `--name` option.
- You specified a virtual container host name in the `--name` option that is the same as that of a virtual container host that has been deleted, but for which the default certificate folder still exists. The certificate folder for the deleted virtual container host has the same name as the one you specified in the `--name` option.
- You used the `--cert-path` option to specify a certificate folder that already contains certificates for another virtual container host.
- You intentionally attempted to reuse an existing certificate, but the value that you provided in `--tls-cname` or `--client-network-ip` does not match the Common Name attribute in the existing certificate.

## Solution ##

Run vic-machine create again with one of the following modifications:

- If another virtual container host of the same name already exists, specify a different name in the `--name` option.
- If a certificate folder for a virtual container host that has been deleted still exists, if that folder has the same name as the one you are specifying in `--name`, and if you do not intend to reuse the existing certificate, delete the existing certificates.
- If you used the `--cert-path` option, delete the existing certificate if it is no longer required, or specify a different folder in  `--cert-path`.
- If you do intend to reuse the existing certificate, update the `--tls-cname` option or `--client-network-ip` option to match the `cname` that the error message included. 