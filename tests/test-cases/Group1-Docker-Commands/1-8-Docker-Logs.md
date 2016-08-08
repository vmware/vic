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
2. Issue docker create busybox /bin/sh -c 'a=0; while [ $a -lt 5 ]; do echo "line $a"; a=`expr $a + 1`; sleep 1; done;' to the VIC appliance
3. Issue docker start <containerID> to the VIC appliance
4. Sleep 6 seconds
5. Issue docker logs <containerID> to the VIC appliance
6. Issue docker logs --tail=all <containerID> to the VIC appliance
7. Issue docker logs --tail=2 <containerID> to the VIC appliance
8. Issue docker logs --tail=0 <containerID> to the VIC appliance
9. Issue docker create -t busybox /bin/sh -c 'a=0; while [ $a -lt 5 ]; do echo "line $a"; a=`expr $a + 1`; sleep 1; done;' to the VIC appliance
10. Issue docker start <containerID> to the VIC appliance
11. Sleep 2 seconds
12. Issue docker logs --follow <containerID> to the VIC appliance
13. Issue docker create busybox /bin/sh -c 'a=0; while [ $a -lt 5 ]; do echo "line $a"; a=`expr $a + 1`; sleep 1; done;' to the VIC appliance
14. Issue docker start <containerID> to the VIC appliance
15. Sleep 2.5 seconds
16. Issue docker logs --tail=1 --follow <containerID> to the VIC appliance
17. Issue docker logs --since=1s <containerID> to the VIC appliance
18. Issue docker logs --timestamps <containerID> to the VIC appliance
19. Issue docker logs fakeContainer

#Expected Outcome:
* Steps 2-16 should all complete without error
* Step 6 should output 5 lines
* Step 7 should output 2 lines
* Step 8 should output 0 lines
* Step 12 should have last line be
```
line 5
```
* Step 16 should output 3 lines
* Step 17 and 18 should result in an error with the following message:
```
Error: vSphere Integrated Containers does not yet support timestampped logs.
```
* Step 19 should result in an error with the following message:
```
Error: No such container: fakeContainer
```

#Possible Problems:
None