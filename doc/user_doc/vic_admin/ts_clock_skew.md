# Connections Fail with Certificate Errors when Using Full TLS Authentication with Trusted Certificates #

Connections to a virtual container host that uses full TLS authentication with trusted Certificate Authority (CA) certificates fail with certificate errors.

## Problem ##
 
- `vic-machine` operations on a virtual container host result in a "bad certificate" error:
  <pre>Connection failed with TLS error "bad certificate" 
check for clock skew on the host
Collecting host-227 hostd.log
vic-machine-windows.exe failed: tls: bad certificate
</pre>
- Connections to the VCH Admin portal for the virtual container host fail with an `ERR_CERT_DATE_INVALID` error.
- Connections to the virtual container host from Docker clients fail with a `bad certificate` error.

## Cause ##
There is potentially a clock skew between the virtual container host and the system from which you are connecting to the virtual container host.

## Solution ##

1. Go to the VCH Admin portal for the virtual container host at https://<i>vch_address</i>:2378 and check the System Time under **Virtual Container Host Info**.
2. If the system time of the virtual container host is wrong, run `vic-machine debug` to enable SSH access to the virtual container host. 

  For information about enabling SSH on a virtual container host, see [Authorize SSH Access to the Virtual Container Host Endpoint VM](vch_ssh_access.md).
2.  Connect to the virtual container host endpoit VM by using SSH.
3.  Use the `date --set` Linux command to set the system clock to the correct date and time.

To prevent this issue recurring on virtual container hosts that you deploy in the future, verify that the host time is correct on the ESXi host on which you deploy virtual container hosts. For information about verifying time synchronization on ESXi hosts, see [VMware KB 1003736](https://kb.vmware.com/kb/1003736).