Test 19-7 New Install of VIC 1.2
=======

# Purpose:
To verify that the install process of VIC 1.2 works and provides a properly functioning system

# References:
[1- VIC 1.2 Install](TBD - waiting on official docs to link)

# Environment:
This test requires that a vCenter server is running and available

# Test Steps:
1. Download the VIC 1.2 OVA from:
[OVA Releases](https://console.cloud.google.com/storage/browser/vic-product-ova-releases/?project=eminent-nation-87317&authuser=1)
2. Walk through the standard OVA install process - setting the appliance password, registry passwords and engine admin username
3. Finish the OVA install process and power on the new VM
4. Log into the admiral portal
5. Add the harbor instance as a registry within admiral
6. Download the root certificate from harbor
7. Navigate to the demo VCH page
8. Configure and install a basic VCH
9. Add the harbor and VCH instances into the default cluster within admiral and enable trusted content in order to whitelist the harbor registry
10. Using a mix of command line and admiral, deploy a variety of containers
11. Verify the VIC UI portlet is installed by default and works as expected in the vSphere UI
12. Remove each of the containers created
13. Delete the VCH
14. Delete the OVA VM

# Expected Outcome:
* All steps should succeed without error
* Step 14 should result in no lingering artifacts in the environment, nothing in the datastore, networks, host, resource pool, etc.

# Possible Problems:
None
