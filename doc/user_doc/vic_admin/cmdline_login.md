# Command Line Certificate Login #

You can use certificate-based authentication with tools such as `curl` or `wget` to access the VCH Admin log server.

## With TLS Client Authentication ##

If you deployed the VCH with client authentication by using `--tls-cname` or by specifying a static IP address on the client network, you can point `curl` to the `cert.pem` and `key.pem` files for the VCH. The following example authenticates connections to the `port-layer.log` file.

<pre>
curl https://<i>vch_address</i>:2378/logs/port-layer.log 
--key ./<i>cert_folder</i>/key.pem 
--certificate ./<i>cert_folder</i>/cert.pem
</pre>


**NOTE**: If your certificates are self-signed, you might also need to specify the `curl -k` flag.

In the example above, <i>cert_folder</i> is either of the following locations:

  - The folder with the same name as the VCH, in the directory from which you ran `vic-machine create`.
  - A folder that you specified in the `vic-machine create --cert-path` option.

## Without Client Authentication ##

If you deployed the VCH without client authentication by using either of `--no-tls` or `--no-tlsverify`, you can use curl to access the logs but you must first authenticate connections to VCH Admin by using the  vSphere username and password. 

1. Log in to VCH Admin to gather an authentication cookie for subsequent access:
   <pre>curl -sk https://<i>vch_address</i>:2378/authentication 
   -XPOST -F username=<i>vsphere_username</i> 
   -F password=<i>vsphere_password</i> 
   -D cookies_file
</pre>

2. Use the cookie from Step 1 in a `curl` command to access the logs.

    <pre>curl -sk https://<i>vch_address</i>:2378/logs/port-layer.log 
   -b cookies_file</pre>