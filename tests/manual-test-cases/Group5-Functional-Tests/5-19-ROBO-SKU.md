Test 5-19 - ROBO SKU
=======

# Purpose:
To verify that VIC works properly when a VCH is installed in a remote office branch office (ROBO) version of vSphere

# References:
[1 - vSphere Remote Office and Branch Office](http://www.vmware.com/products/vsphere/remote-office-branch-office.html)

# Environment:
This test requires access to VMware Nimbus cluster for dynamic ESXi and vCenter creation

# Test Steps:
1. Deploy a new vCenter with stand alone hosts
2. Add the ROBO SKU license to the vCenter appliance
3. Assign the ROBO SKU license to each of the hosts within the vCenter
4. Install the VIC appliance onto one of the hosts in the vCenter
5. Run a variety of docker operation on the VCH

# Expected Outcome:
All test steps should complete without error

# Possible Problems:
None
