Test 11-06- Upgrade-StorageQuota
=======

# Purpose:
To verify vic-machine configure for storage quota work well after upgrade

# Environment:
This test requires that a vSphere server is running and available

# Test Steps:
1. Download vic_1.4.3.tar.gz from gcp
2. Deploy a VCH_1 using vic_1.4.3 and create several containers and volumes.
3. Deploy another 2 VCH(VCH_2 and VCH_3) using vic_1.4.3 using the same image store and VCH_3 is using 10G as base image size and VCH_2 is using default 8G as base image size. Create containers and volumes in the two VCHs as well.
4. Upgrade VCH_1/VCH_2/VCH_3 to the latest vic build
5. Using docker -H VCH_1 info to check VCH_1 info's correctness
6. Set a lower storage quota for VCH_1
7. Create a new container in VCH_1, which should be failure.
8. Reset a bigger storage quota for VCH_1 and new container should be successfully created
9. Set a bigger storage quota for VCH_2 and create a new container in VCH_2
10. Check storage quota is showing in VCH_2
11. Pick VCH_1 to do upgrade revert back
12. Re-do upgrade on VCH_1

# Expected Outcome:
* When VCH is upgraded to the latest, container and images used size can be collected by docker info
* docker info can be responsed in acceptable time.
* When setting a lower storage quota in VCH and new container can not be successfully created due to Storage quota exceed.
* When setting a bigger storage quota in VCH and new container can be successfully created.
* When doing upgrade revert on a VCH which is configured storage quota, VCH could be reverted back to 1.4.3

# Possible Problems:
* May update the baseline of acceptable value of docker info response time.
* When a VCH is configured with storage quota, upgrade --rollback can not revert the VCH to 1.4.3 version.
