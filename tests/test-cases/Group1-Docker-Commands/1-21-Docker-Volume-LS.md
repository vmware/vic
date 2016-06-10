Test 1-21 - Docker Volume LS
=======

#Purpose:
To verify that docker volume ls command is supported by VIC appliance

#References:
[1 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/volume_ls/)

#Environment:
This test requires that a vSphere server is running and available

#Test Steps:
1. Deploy VIC appliance to vSphere server
2. Issue docker create --name=test to the VIC appliance
3. Issue docker volume ls to the VIC appliance
4. Issue docker volume ls -q to the VIC appliance
5. Issue docker volume ls -f dangling=true to the VIC appliance
6. Issue docker volume ls -f name=test to the VIC appliance
7. Issue docker create -v test:/test busybox to the VIC appliance
8. Issue docker volume ls -f dangling=true to the VIC appliance

#Expected Outcome:
* Step 3 should result in each volume being listed with both driver and volume name
* Step 4 should result in each volume being listed with only the volume name being listed
* Step 5 should result in at least the test volume being listed
* Step 6 should result in the following error:  
```
Error response from daemon: Invalid filter 'name'
```
* Step 8 should result in a list of volumes that does not include test anymore

#Possible Problems:
None