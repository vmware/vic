Test 11-02 - Upgrade Exec
=======

#Purpose:
To verify that exec does not work in VIC version 0.9.0 and lower

#Environment:
This test requires that a vSphere server is running and available

#Test Steps:
1. Download vic_8351.tar.gz from gcp
2. Deploy VIC 8351 to vsphere server
3. Create a busybox container running the background
4. Upgrade VCH to latest version
5. Run docker exec on container created in step (3.)
6. Create new container
7. Run docker exec on new container in (6.)

#Expected Outcome:
* Step 5 should fail
* All other steps should result in success
