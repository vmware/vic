Test 1-2 - Docker Pull
=======

#Purpose:
To verify that docker pull command is supported by VIC appliance

#References:
[1 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/pull/)

#Environment:
This test requires that an vSphere server is running and available.

#Test Steps:
1. Deploy VIC appliance to vSphere server
2. Issue a docker pull command to the new VIC appliance for each of the top 3 most popular images in hub.docker.com
    * nginx, busybox, ubuntu
3. Issue a docker pull command to the new VIC appliance using a tag that isn't the default latest
    * ubuntu:14.04
4. Issue a docker pull command to the new VIC appliance using a digest
    * ubuntu@sha256:45b23dee08af5e43a7fea6c4cf9c25ccf269ee113168c19722f87876677c5cb2
5. Issue a docker pull command to the new VIC appliance using a different repo than the default
    * myregistry.local:5000/testing/test-image
6. Issue a docker pull command to the new VIC appliance using all tags option
    * --all-tags fedora
7. Issue a docker pull command to the new VIC appliance using an image that doesn't exist
8. Issue a docker pull command to the new VIC appliance using a non-default repository that doesn't exist

#Expected Outcome:
VIC appliance should respond with a properly formatted pull response to each command issued to it. No errors should be seen, except in the case of step 7 and 8.

#Possible Problems:
None
