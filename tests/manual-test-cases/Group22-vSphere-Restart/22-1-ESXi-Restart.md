Test 22-1 - ESXi Restart
=======

# Purpose:
To verify the VIC appliance works when running on an ESXi that gets restarted during operation

# References:

# Environment:
This test requires access to VMWare Nimbus cluster for dynamic ESXi and vCenter creation

# Test Steps:
1. Deploy a new ESXi in nimbus
2. Deploy a VCH on the new ESXi server
3. Run a variety of docker commands on the new VCH
4. Restart the ESXi server
5. Power on the VCH after the reboot and run a variety of docker commands again

# Expected Outcome:
Verify after the restart that the VCH continues to operate as expected and any containers from before the restart are still functional

# Possible Problems:
* Since the containers are not in a vApp on an ESXi server, it is unclear whether they should come back when you reboot VCH or not
* Need a reliable way to determine that the vSphere infrastructure is back up