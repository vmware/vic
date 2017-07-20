Test 1-20 - Docker Volume Inspect
=======

# Purpose:
To verify that docker volume inspect command is supported by VIC appliance

# References:
[1 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/volume_inspect/)

# Environment:
This test requires that a vSphere server is running and available

# Test Steps:
1. Deploy VIC appliance to vSphere server
2. Issue docker volume create --name=test to the VIC appliance
3. Issue docker volume inspect test to the VIC appliance
4. Issue docker volume inspect fakeVolume to the VIC appliance
5. Issue docker volume inspect test to the VIC appliance
6. Restart the VCH
7. Issue docker volume inspect test to the VIC appliance

# Expected Outcome:
* Step 3 should result in a properly formatted JSON response
* Step 4 should result in an error with the following message:
```
Error: No such volume: fakeVolume
```
* Step 5 should result in a properly formatted JSON response
* Step 6 VCH should restart successfully
* Step 7 should result in a properly formatted JSON response

# Possible Problems:
None