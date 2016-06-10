Test 1-3 - Docker Images
=======

#Purpose:
To verify that docker images command is supported by VIC appliance

#References:
[1 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/images/)

#Environment:
This test requires that a vSphere server is running and available

#Test Steps:
1. Deploy VIC appliance to vSphere server
2. Docker pull alpine
3. Docker pull alpine:3.2
4. Docker pull alpine:3.1
5. Issue docker images command to the new VIC appliance
6. Issue docker images -a command to the new VIC appliance
7. Issue docker images -q command to the new VIC appliance
8. Issue docker images --no-trunc command to the new VIC appliance
9. Issue docker images alpine:3.1

#Expected Outcome:
* Each of the commands issued should return error free with a properly formatted response
* The docker images and docker images -a command should return the 3 alpine images as expected
* The docker images -q command should return only the short hash of the three images
* The docker --no-trunc command should return the full non-truncated image ID of the three images
* The docker images alpine:3.1 command should return only the specific image specified

#Possible Problems:
None