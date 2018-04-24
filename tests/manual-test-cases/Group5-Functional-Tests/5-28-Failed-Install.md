Test 5-28 - Failed Install
=======

# Purpose:
To verify that when VIC install fails, all artifacts are cleaned up properly on the machine

# References:

# Environment:
This test requires access to VMware Nimbus cluster for dynamic ESXi and vCenter creation

# Test Steps:
1. Deploy a new vCenter with a simple cluster
2. Install VIC into an invalid portgroup

# Expected Outcome:
After step 2 fails, verify that the system is properly cleaned up

# Possible Problems:
None
