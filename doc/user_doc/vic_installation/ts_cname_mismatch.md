# VCH Deployment Fails with Certificate `cname` Mismatch #

When you use `vic-machine create` to deploy a virtual container host (VCH), the deployment fails with an error about the certificate `cname` value.

## Problem ##
Deployment fails during the validation of the configuration that you provided:

<pre>Provided cname does not match that in existing server certificate: <i>cname</i>
Unable to load certificates: cname option doesn't match existing server certificate 
in certificate path <i>path_to_certificate</i></pre>

## Cause ##
`vic-machine create` attempts to re-use certificates that it finds in `--cert-path`. The default value of `--cert-path` derives from the value that you specify in `--name`. If you are deploying a VCH from the same location and with the same name as a previous VCH, `vic-machine create` reuses the old certificates. This behavior is intentional, to allow you to easily redeploy a VCH without requiring you to re-issue client certificates to users.

Before reusing the existing certificates, `vic-machine` confirms that the existing certificate is valid given the options supplied for the new deployment. The options that influence this in order of priority are:
* `--tls-cname` if specified, or
* `--client-ip-address`, or 
* `--public-ip-address` if the client and public network roles share an interface.

The error message means that the existing certificate has a Common Name attribute that differs from the value derived from the options detailed above.

## Solution ##

- To reuse the certificates directly, change `--tls-cname`, `--client-ip-address`, or `--public-ip-address` to match the Common Name in the existing certificate.

- If you want to reuse the Certificate Authority so that client certificates remain valid, but you need to provide a different IP address:

  1. Manually generate the server certificates by using `openssl`, signing them with the existing CA.
  2.  Use the `--cert` and `--key` options to pass the newly generated certificates to `vic-machine create`.

- If you do not want to reuse the certificates, choose one of the following options:
  - Change the location from which you run `vic-machine`. This alters the default `--cert-path`.
  - Change the value of `--name`. This alters the default `--cert-path`.
  - Specify `--cert-path` explicitly.
  - Delete the existing certificates from `--cert-path`