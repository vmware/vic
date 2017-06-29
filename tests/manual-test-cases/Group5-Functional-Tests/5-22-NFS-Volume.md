Test 5-22 - NFS Volume
=======

# Purpose:
To verify that NFS shared volumes work with current supported docker commands

# References:
[1 - Best practices for running VMware vSphere on NFS](http://www.vmware.com/content/dam/digitalmarketing/vmware/en/pdf/techpaper/vmware-nfs-bestpractices-white-paper-en.pdf)

# Environment:
This test requires access to VMware Nimbus cluster for dynamic ESXi and vCenter creation

# Test Steps:
1. Deploy a new vCenter with a simple cluster
2. Deploy an NFS server
3. Create a new datastore out of a NFS share on the NFS server
4. Install the VIC appliance in cluster


# Expected Outcome:

# Possible Problems:
None