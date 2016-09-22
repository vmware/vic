Test 10-1 - VCH Restart
=======

#Purpose:
To verify that restarting a VCH results in a consistent and correct state

#Environment:
This test requires that a vSphere server is running and available

#Test Steps:
1. Deploy VIC appliance to the vSphere server
2. Issue docker network create bar, creating a new network called "bar"
3. Reboot VCH
4. Issue a docker network ls
5. Issue docker inspect
6. Issue docker stop, start and ls

#Expected Outcome:
* VCH should reboot within a reasonable amount of time
* After VCH restart, network ls should have the previously created network listed
* Step 5, 6 should result in success

#Possible Problems:
None
