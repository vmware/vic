Test 1-11 - Docker RM
=======

#Purpose:
To verify that docker rm command is supported by VIC appliance

#References:
[1 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/rm/)

#Environment:
This test requires that a vSphere server is running and available

#Test Steps:
1. Deploy VIC appliance to vSphere server
2. Issue docker create busybox dmesg to the VIC appliance
3. Issue docker rm <containerID> to the VIC appliance
4. Issue docker create busybox /bin/top to the VIC appliance
5. Issue docker start <containerID> to the VIC appliance
6. Issue docker rm <containerID> to the VIC appliance
7. Issue docker rm -f <containerID> to the VIC appliance
8. Issue docker rm fakeContainer to the VIC appliance

#Expected Outcome:
* Steps 2-5 should complete without error
* Step 3 and 7 should result in the container being removed from the VIC appliance
* Step 6 should result in the following error:  
```
Failed to remove container (<containerID>): Error response from daemon: Conflict, You cannot remove a running container. Stop the container before attempting removal or use -f
```
* Step 8 should result in the following error:  
```
Failed to remove container (fakeContainer): Error response from daemon: No such container: fakeContainer
```

#Possible Problems:
None