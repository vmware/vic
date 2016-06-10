Test 1-13 - Docker Version
=======

#Purpose:
To verify that docker version command is supported by VIC appliance

#References:
[1 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/version/)

#Environment:
This test requires that a vSphere server is running and available.

#Test Steps:
1. Deploy VIC appliance to the vSphere server
2. Issue a docker version command to the new VIC appliance
3. Issue a docker version --format '{{.Client.Version}}' command to the new VIC appliance
4. Issue a docker version --format '{{.Client.APIVersion}}' command to the new VIC appliance
5. Issue a docker version --format '{{.Client.GoVersion}}' command to the new VIC appliance
3. Issue a docker version --format '{{.Server.Version}}' command to the new VIC appliance
4. Issue a docker version --format '{{.Server.APIVersion}}' command to the new VIC appliance
5. Issue a docker version --format '{{.Server.GoVersion}}' command to the new VIC appliance

#Expected Outcome:
* VIC appliance should respond with a properly formatted version response, it should be capable of returning each individual field as well without error.
* The server version field should indicate that it is VIC.

#Possible Problems:
None