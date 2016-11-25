# Connections Fail with Certificate Errors when Using Full TLS Authentication with Trusted Certificates #

Connections to a virtual container host (VCH) that uses full TLS authentication with trusted Certificate Authority (CA) certificates fail with certificate errors.

## Problem ##
 
- `vic-machine` operations on a VCH result in a "bad certificate" error:
  <pre>Connection failed with TLS error "bad certificate" 
check for clock skew on the host
Collecting host-227 hostd.log
vic-machine-windows.exe failed: tls: bad certificate
</pre>
- Connections to the VCH Admin portal for the VCH fail with an `ERR_CERT_DATE_INVALID` error.
- Connections to the VCH from Docker clients fail with a `bad certificate` error.

## Cause ##
There is potentially a clock skew between the VCH and the system from which you are connecting to the VCH.

## Solution ##

1. Go to the VCH Admin portal for the VCH at https://<i>vch_address</i>:2378 and check the System Time under **VCH Info**.
2. If the system time of the VCH is wrong, run `vic-machine debug` to enable SSH access to the VCH. 

  For information about enabling SSH on a VCH, see [Authorize SSH Access to the VCH Endpoint VM](vch_ssh_access.md).
2.  Connect to the VCH endpoit VM by using SSH.
3.  Use the `date --set` Linux command to set the system clock to the correct date and time.

To prevent this issue recurring on VCHs that you deploy in the future, verify that the host time is correct on the ESXi host on which you deploy VCHs. For information about verifying time synchronization on ESXi hosts, see [VMware KB 1003736](https://kb.vmware.com/kb/1003736).