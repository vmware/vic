Test 11-01 - Upgrade
=======

# Purpose:
To verify vic-machine upgrade can upgrade VCH from a certain version

# Environment:
This test requires that a vSphere server is running and available

# Test Steps:
1. Download vic_7315.tar.gz from bintray
2. Deploy VIC 7315 to vsphere server
3. Issue docker network create bar, creating a new network called "bar"
4. Create container with port mapping
5. Upgrade VCH to latest version with short timeout 1s
6. Upgrade VCH to latest version
7. Roll back to the previous version
8. Upgrade again to the upgraded version
9. Check the previous created container and image are still there
10. Attempt to rename an old container created with a VCH that doesn't support rename.
11. Rename a new container created with a VCH that supports rename.
12. Check the previous created container's display name and datastore folder name
13. Check the display name and datastore folder name of a new container created after VCH upgrade

# Expected Outcome:
* Step 5 should fail with timeout
* Step 10 should result in an error containing the following message:
```
does not support rename
```
* Step 11 should succeed and the container's new name should be present in ps, inspect and govc vm.info output.
* Step 12 should show that both the container's display name and datastore folder name are containerName-containerID
* Step 13 should show that (1) on a non-vsan setup, the container's display name is containerName-containerShortID while the datastore folder name is containerID, or (2) on a vsan setup, both the container's display name and datastore folder name are containerName-containerShortID
* All other steps should result in success

# Possible Problems:
Upgrade test will upgrade VCH from build 7315, because that build has VCH restart and configuration restart features done.
Before GA, if there is any VCH configuration change, please bump upgrade from version, and be sure to add cases to cover those changes.
After GA, the upgrade from version will be GA release version.
