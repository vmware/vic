Test 1-10 - Docker PS
=======

#Purpose:
To verify that docker ps command is supported by VIC appliance

#References:
[1 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/ps/)

#Environment:
This test requires that a vSphere server is running and available

#Test Steps:
1. Deploy VIC appliance to vSphere server
2. Issue docker ps to the VIC appliance
3. Issue docker create busybox /bin/top to the VIC appliance
4. Issue docker start <containerID> to the VIC appliance
5. Issue docker create busybox ls to the VIC appliance
6. Issue docker start <containerID> to the VIC appliance
7. Issue docker create busybox dmesg to the VIC appliance
8. Issue docker ps to the VIC appliance
9. Issue docker ps -a to the VIC appliance
10. Issue docker ps -l to the VIC appliance
11. Issue docker ps -n=2 to the VIC appliance
12. Issue docker ps -ls to the VIC appliance
13. Issue docker ps -aq to the VIC appliance
14. Issue docker ps -f status=created to the VIC appliance

#Expected Outcome:
* Steps 2-13 should all return without error
* Step 2 should return with only the printed ps command header and no containers
* Step 8 should return with only the information for the /bin/top container
* Step 9 should return with the information for all 3 containers
* Step 10 should return with the information for only the 'dmesg' container
* Step 11 should return with the information for both the 'ls' and the 'dmesg' containers
* Step 12 should return with the information in addition to the size information of the 'dmesg' container
* Step 13 should return with only the three container IDs
* Step 14 should return with only the information for the 'dmesg' container

#Possible Problems:
None