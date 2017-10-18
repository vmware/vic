Test 1-11 - Docker RM
=======

# Purpose:
To verify that docker rm command is supported by VIC appliance

# References:
[1 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/rm/)

# Environment:
This test requires that a vSphere server is running and available

# Test Steps:
1. Deploy VIC appliance to vSphere server
2. Issue docker create busybox dmesg to the VIC appliance
3. Issue docker rm <containerID> to the VIC appliance
4. Issue docker create busybox ls to the VIC appliance
5. Issue docker start <containerID> to the VIC appliance
6. Issue docker rm <containerID> to the VIC appliance
7. Issue docker create busybox /bin/top to the VIC appliance
8. Issue docker start <containerID> to the VIC appliance
9. Issue docker rm <containerID> to the VIC appliance
10. Issue docker rm -f <containerID> to the VIC appliance
11. Issue docker rm fakeContainer to the VIC appliance
12. Issue docker create --name test busybox to the VIC appliance
13. Remove the containerVM out-of-band using govc
14. Issue docker rm test to the VIC appliance
15. Issue docker rm to container created with an unknown executable

# Expected Outcome:
* Steps 2-8,12,15 should complete without error
* Step 3,6,10 should result in the container being removed from the VIC appliance
* Step 9 should result in the following error:  
```
Error response from daemon: Conflict, You cannot remove a running container. Stop the container before attempting removal or use -f
```
* Step 11 should result in the following error:  
```
Error response from daemon: No such container: fakeContainer
```
* Step 13 should succeed on ESXi and fail on vCenter with the following error:
```
govc: ServerFaultCode: The method is disabled by 'VIC'
```
* When run on standalone ESXi, step 14 should result in the following error:  
```
Error response from daemon: No such container: test
```

# Possible Problems:
None
