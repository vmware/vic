# VCH Deployment Fails with Certificate `cname` Mismatch #

When you use `vic-machine create` to deploy a virtual container host, the deployment fails with an error about the certificate `cname` value.

## Problem ##
Deployment fails during the validation of the configuration that you provided:

<pre>Provided cname does not match that in existing server certificate: <i>cname</i>
Unable to load certificates: cname option doesn't match existing server certificate 
in certificate path <i>path_to_certificate</i></pre>

## Cause ##
This error can occur in the following circumstances when you run `vic-machine create`:
- You specified the Common Name attribute to use in auto-generated CA certificates in one of the following ways:
  - You specified the `--tls-cname` option.
  - You specified a static IP address in the `--client-network-ip` option.
  - You specified a static IP address in the `--public-network-ip` option and the public network shares a port group with the client network.
- `vic-machine create` finds and attempts to use an existing auto-generated certificate. 
- The existing certificate includes a Common Name attribute that is different to the address that you specified in `--tls-cname`, `--client-network-ip`, or `--public-network-ip`. 

`vic-machine create` is attempting to use an existing certificate for one of the following reasons:

- You specified a virtual container host name in the `--name` option that is the same as that of an existing virtual container host. The certificate folder for the existing virtual container host has the same name as the one you specified in the `--name` option for the new virtual container host.
- You specified a virtual container host name in the `--name` option that is the same as that of a virtual container host that has been deleted, but for which the default certificate folder still exists. The certificate folder for the deleted virtual container host has the same name as the one you specified in the `--name` option for the new virtual container host.
- You used the `--cert-path` option to specify a certificate folder that already contains certificates for another virtual container host.
- You intentionally attempted to reuse an existing certificate, but the value that you provided in `--tls-cname`, `--client-network-ip`, or `--public-network-ip` does not match the Common Name attribute in the existing certificate.

## Solution ##

Run `vic-machine create` again with one of the following modifications:

- If another virtual container host of the same name already exists, specify a different name for the new virtual container host in the `--name` option.
- If a certificate folder still exists for a virtual container host that has been deleted, if that folder has the same name as the one you are specifying in `--name`, and if you do not intend to reuse the existing certificate, delete the existing certificates.
- If you used the `--cert-path` option, delete the existing certificate if it is no longer required, or specify a different certificate folder in  `--cert-path`.
- If you do intend to reuse the existing certificate, update the `--tls-cname` option or `--client-network-ip` option to match the `cname` that the error message included. 