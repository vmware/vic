Test 1-05 - Docker Start
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
6. Issue docker create vmware/photon
7. Issue docker start vmware/photon <containerID>
8. Issue docker start fakeContainer
9. Create a container, remove it's ethernet adapter, then start the container 
10. Create and start 5 busybox containers running /bin/top serially
11. Create and start 5 ubuntu containers running /bin/top serially
12. Create and start 5 busybox containers running /bin/top all at once
13. Run a container with a test-network, stop the container, remove the test-network, then start the container again

#Expected Outcome:
* Commands 1-7 should all return without error and respond with the container ID
* After commands 3, 5, and 7 verify that the containers are running
* Step 8 should result in the VIC appliance returning the following error:
```
Error response from daemon: No such container: fakeContainer
Error: failed to start containers: fakeContainer
```
* Step 9 should result in an error message stating unable to wait for process launch status
* Steps 10-12 should all result in all containers succeeding and not throwing any errors
* Step 13 should result in the VIC appliance returning the following error:
```
Error response from daemon: Server error from portlayer: network test-network not found
Error: failed to start containers: containerID
```
#Possible Problems:
None
