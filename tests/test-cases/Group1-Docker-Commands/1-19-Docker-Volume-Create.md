Test 1-19 - Docker Volume Create
=======

#Purpose:
To verify that docker volume create command is supported by VIC appliance

#References:
[1 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/volume_create/)

#Environment:
This test requires that a vSphere server is running and available

#Test Steps:
1. Deploy VIC appliance to vSphere server
2. Issue docker volume create to the VIC appliance
3. Issue docker volume create --name=test to the VIC appliance
4. Issue docker volume create -d remote --name=test to the VIC appliance

#Expected Outcome:
* Steps 2 and 3 should complete successfully and return the name of the volume created, you should then be able to see the volume has been created
* Step 4 should result in error with the following error message:  
```
Error response from daemon: A volume named test already exists. Choose a different volume name.
```

#Possible Problems:
None