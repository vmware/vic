Test 1-23 - Docker Inspect
=======

#Purpose:
To verify that docker inspect command is supported by VIC appliance

#References:
[1 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/inspect/)

#Environment:
This test requires that a vSphere server is running and available

#Test Steps:
1. Deploy VIC appliance to vSphere server
2. Issue docker pull busybox to the VIC appliance
3. Issue docker inspect busybox to the VIC appliance
4. Issue docker inspect --type=image busybox to the VIC appliance
5. Issue docker inspect --type=container busybox to the VIC appliance
6. Issue docker create busybox to the VIC appliance
7. Issue docker inspect <containerID> to the VIC appliance
8. Issue docker inspect --type=container <containerID> to the VIC appliance
9. Issue docker inspect <containerID> to the VIC appliance and verify the Cmd and Image fields
10. Issue docker inspect --type=image <containerID> to the VIC appliance
11. Issue docker inspect fake to the VIC appliance

#Expected Outcome:
* Step 3,4,7,8 should result in success and a properly formatted JSON response
* Step 5 should result in an error with the following message:  
```
Error: No such container: busybox
```
* Step 9 should result in success with the correct values in the Cmd and Image fields
* Step 10 should result in an error with the following message:
```
Error: No such image: <containerID>
```
* Step 11 should result in an error with the following message:
```
Error: No such image or container: fake
```

#Possible Problems:
None