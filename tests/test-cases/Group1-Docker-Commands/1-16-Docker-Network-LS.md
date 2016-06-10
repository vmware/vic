Test 1-16 - Docker Network LS
=======

#Purpose:
To verify that docker network ls command is supported by VIC appliance

#References:
[1 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/network_ls/)

#Environment:
This test requires that a vSphere server is running and available

#Test Steps:
1. Deploy VIC appliance to vSphere server
2. Issue docker network ls to the VIC appliance
3. Issue docker network ls -q to the VIC appliance
4. Issue docker network ls -f name=host to the VIC appliance
5. Issue docker network ls --no-trunc to the VIC appliance
6. Issue docker network ls -f name=fakeName to the VIC appliance

#Expected Outcome:
* Steps 2-6 should all complete successfully
* Step 2 should return at the least the host, null, and bridge networks that are default
* Step 3 should return the networks ID only
* Step 4 should return only the host network
* Step 5 should return all of the networks with their full IDs
* Step 6 should return no networks listed

#Possible Problems:
None