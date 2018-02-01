Test 22-2 - vCenter Restart
=======

# Purpose:
To verify the VIC appliance works when running on an vCenter that gets restarted during operation

# References:

# Environment:
This test requires access to VMWare Nimbus cluster for dynamic ESXi and vCenter creation

# Test Steps:
1. Deploy a new simple vCenter cluster in nimbus
2. Deploy a VCH on the new vCenter cluster
3. Run a variety of docker commands on the new VCH
4. Restart the vCenter server
5. Power on the VCH after the reboot and run a variety of docker commands again

# Expected Outcome:
Verify after the restart that the VCH continues to operate as expected and any containers from before the restart are still functional

# Possible Problems:
* Need a reliable way to determine that the vSphere infrastructure is back up