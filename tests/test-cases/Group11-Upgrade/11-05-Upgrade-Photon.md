Test 11-05 - Upgrade-Photon
=======

# Purpose:
To verify photon upgrade with vic-machine upgrade

# Environment:
This test requires that a vSphere server is running and available

# Test Steps:
1. Download vic_1.4.3.tar.gz from gcp
2. Deploy VCH using vic_1.4.3 and create some volumes with data on them.
3. Create a container VM with one of the created volumes.
4. Do VIC upgrade and check whether the previous running container still works.
5. Create a new container using another volume and check whether the container can successfully read the volume data.
6. Create a new volume with data on it.
7. Do upgrade --rollback on the VCH
8. Create a new container with the volume created in Step6. 
9. Poweroff the container which is created in Step5 and do VCH upgrade again 

# Expected Outcome:
* With VIC upgrade, the photon should be upgraded to 2.0
* Volumes created in the vic_1.4.3 still works well in vic_1.5.0

# Possible Problems:
* When VCH is downgraded to vic_1.4.3 and do re-upgrade, we need to make sure that the container created in previous vic_1.5.0 are in power-off state.
