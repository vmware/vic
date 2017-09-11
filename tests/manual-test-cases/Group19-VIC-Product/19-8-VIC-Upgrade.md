Test 19-8 VIC 1.2 Upgrade from 1.1.0 or 1.1.1
=======

# Purpose:
To verify VIC 1.2 upgrade from 1.1.0 or 1.1.1 works and provides access to all of the vSphere Integrated Containers components.

# References:
[VIC 1.2 Upgrade](TBD - waiting on official docs to link)

# Environment:
* This test requires that a vCenter server is running and available.
* Existing vSphere Integrated Containers 1.1.x appliance that has deployed by using the official OVA installer.
* Configure and install VCH
* Deploy a variety of containers through command line and admiral

# Test Steps:
1. Deploy new 1.2 OVA of the vSphere Integrated Containers appliance to a location in which it can access the VMDK files of the previous appliance.
2. Log in to a vSphere Web Client instance to access both versions of the vSphere Integrated Containers appliance.
3. Shut down both the previous and new vSphere Integrated Containers appliances by selecting Shut Down Guest OS.
4. Click on the previous vSphere Integrated Containers appliance, and Remove data disk from VM
5. Click on the new vSphere Integrated Containers appliance.
6. Remove and delete empty data disk from new OVA VM
7. Attach old data disk to new OVA VM
8. Power on the new appliance to complete the upgrade
9. Run the upgrade script to migrate harbor and admiral
10. Run the vic-machine upgrade command to upgrade your existing VCH to the latest version
11. Verify that the data from your previous vSphere Integrated Containers installation has migrated successfully
12. Verify that VCH and containers from previous VIC appliance exists after upgrade and run docker commands.
13. Deploy a variety of containers through command line and admiral on new VIC appliance.
14. Remove each of the containers created
15. Delete VCH
16. Delete new and old OVA VM


# Expected Outcome:
* All steps should succeed without error

# Possible Problems:
* None
