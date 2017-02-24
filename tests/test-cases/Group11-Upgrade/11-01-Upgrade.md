Test 11-01 - Upgrade
=======

#Purpose:
To verify vic-machine upgrade can upgrade VCH from a certain version

#Environment:
This test requires that a vSphere server is running and available

#Test Steps:
1. Download vic_5470.tar.gz from bintray
2. Deploy VIC 5470 to vsphere server
3. Issue docker network create bar, creating a new network called "bar"
4. Create container with port mapping
5. Upgrade VCH to latest version with short timeout 1s
6. Upgrade VCH to latest version
6. Check the previous created container and image are still there

#Expected Outcome:
* Step 5 should fail with timeout
* All other steps should result in success

#Possible Problems:
Upgrade test will upgrade VCH from build 5470, because that build has VCH restart and configuration restart features done.
Before GA, if there is any VCH configuration change, please bump upgrade from version, and be sure to add cases to cover those changes.
After GA, the upgrade from version will be GA release version.
