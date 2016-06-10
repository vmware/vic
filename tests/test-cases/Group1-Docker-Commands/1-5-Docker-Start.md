Test 1-5 - Docker Start
=======

#Purpose:
To verify that docker start command is supported by VIC appliance

#References:
[1 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/start/)

#Environment:
This test requires that a vSphere server is running and available

#Test Steps:
1. Deploy VIC appliance to vSphere server
2. Issue docker create -it busybox /bin/top to VIC appliance
3. Issue docker start <containerID>
4. Issue docker create -it busybox /bin/top to VIC appliance
5. Issue docker start -ai <containerID>
6. Issue docker start fakeContainer

#Expected Outcome:
* Commands 1-5 should all return without error and respond with the container ID
* After command 3 and 5 verify that the containers are running
* Step 6 should result in the VIC applaiance returning the following error:  
```
Error response from daemon: No such container: fakeContainer  
Error: failed to start containers: fakeContainer  
```

#Possible Problems:
None