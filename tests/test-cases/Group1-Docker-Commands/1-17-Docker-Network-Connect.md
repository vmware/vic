Test 1-17 - Docker Network Connect
=======

#Purpose:
To verify that docker network connect command is supported by VIC appliance

#References:
[1 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/network_connect/)

#Environment:
This test requires that a vSphere server is running and available

#Test Steps:
1. Deploy VIC appliance to vSphere server
2. Issue docker network create test-network
3. Issue docker create busybox ifconfig
4. Issue docker network connect test-network <containerID>
5. Issue docker start <containerID>
6. Issue docker logs <containerID>
7. Issue docker network connect test-network fakeContainer
8. Issue docker network connect fakeNetwork <containerID>

#Expected Outcome:
* Step 4 should complete successfully
* Step 6 should print the results of the ifconfig command and there should be two network interfaces in the container(eth0, eth1)
* Step 7 should result in an error with the following message:  
```
Error response from daemon: No such container: fakeContainer
```
* Step 8 should result in an error with the following message:  
```
Error response from daemon: network fakeNetwork not found
```

#Possible Problems:
None