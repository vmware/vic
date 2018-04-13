Test 5-25 - OPS User Grant
=======

# Purpose:
To verify that VIC works properly when a VCH is installed with the option to create the proper permissions for the OPS-user

# Environment:
This test requires access to VMware Nimbus cluster for dynamic ESXi and vCenter creation

# Test Steps:
1. Deploy a new vCenter with a simple cluster
2. Create Local OPS User on VC
3. Install the VIC appliance into the cluster with the --ops-grant-perms option
4. With the ops-user, use govc to attempt to change the DRS settings on the cluster
5. Run a variety of docker operations on the VCH

# Expected Outcome:
* Steps 1-3 should succeed
* Step 4 should fail since the ops-user does not have enough permissions for the operation
* Step 5 should succeed

# Possible Problems:
None
