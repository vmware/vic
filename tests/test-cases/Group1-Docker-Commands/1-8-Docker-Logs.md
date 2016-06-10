Test 1-8 - Docker Logs
=======

#Purpose:
To verify that docker logs command is supported by VIC appliance

#References:
[1 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/logs/)

#Environment:
This test requires that a vSphere server is running and available

#Test Steps:
1. Deploy VIC appliance to vSphere server
2. Issue docker create busybox dmesg to the VIC appliance
3. Issue docker start <containerID> to the VIC appliance
4. Issue docker logs <containerID> to the VIC appliance
5. Issue docker logs -t <containerID> to the VIC appliance
6. Issue docker logs --since="<timeStamp from Step 5>" <containerID> to the VIC appliance
7. Issue docker logs --tail="3" <containerID> to the VIC appliance
8. Issue docker create busybox /bin/top to the VIC appliance
9. Issue docker start <containerID> to the VIC appliance
10. Issue docker logs -f <containerID> to the VIC appliance
11. Issue docker logs fakeContainer

#Expected Outcome:
* Steps 2-10 should all complete without error
* Step 4 should return the result of the dmesg output
* Step 5 should return the result of the dmesg output with timestamps on each of the lines
* Step 6 should return a set of results limited to only messages after that timestamp
* Step 7 should return only the last 3 lines of output from the dmesg command
* Step 10 should show the live running output from the top command
* Step 11 should result in an error with the following message:    
```
Error: No such container: fakeContainer
```

#Possible Problems:
None