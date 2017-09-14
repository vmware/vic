Test 5-23 - Resource Pool Install
=======

# Purpose:
To verify that installing a VCH into a resource pool instead of a vApp works as expected

# References:

# Environment:
This test requires access to VMware Nimbus for dynamic ESXi and NFS server creation

# Test Steps:
1. Create a simple vCenter cluster in Nimbus
2. Install a VCH with the --use-rp option
3. Run a variety of docker commands

# Expected Outcome:
* Each step should result in success

# Possible Problems:
