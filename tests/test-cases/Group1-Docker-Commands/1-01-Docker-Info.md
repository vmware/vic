Test 1-01 - Docker Info
=======

#Purpose:
To verify that docker info command is supported by VIC appliance

#References:
[1 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/info/)

#Environment:
This test requires that a vSphere server is running and available

#Test Steps:
1. Deploy VIC appliance to the vSphere server
2. Issue a docker info command to the new VIC appliance
3. Issue a docker -D info command to the new VIC appliance
4. Issue docker info command, docker create busybox, docker info, docker start <containerID>, docker info

#Expected Outcome:
* VIC appliance should respond with a properly formatted info response, no errors should be seen. Step 3 should result in additional debug information being returned as well.
* Verify in step 4 that the correct number of containers is reported

#Possible Problems:
None