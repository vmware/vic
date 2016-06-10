Test 1-6 - Docker Run
=======

#Purpose:
To verify that docker run command is supported by VIC appliance

#References:
[1 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/run/)

#Environment:
This test requires that a vSphere server is running and available

#Test Steps:
1. Deploy VIC appliance to vSphere server
2. Issue docker run busybox dmesg to the VIC appliance
3. Issue docker run busybox -i dmesg to the VIC appliance
4. Issue docker run -it busybox /bin/top to the VIC appliance
5. Issue 'q' command to the container
6. Issue docker run busybox /bin/top to the VIC appliance
7. Issue docker run busybox fakeCommand to the VIC appliance
8. Issue docker run fakeImage /bin/bash to the VIC appliance
9. Issue docker run -d --name busy3 busybox /bin/top to the VIC appliance
10. Issue docker run --link busy3:busy3 busybox ping -c2 busy3 to the VIC appliance

#Expected Outcome:
* Step 2 and 3 should result in success and print the dmesg of the container
* Step 4 should result in the top command starting and printing it's results to the screen
* Step 5 should result in top stopping and the container exiting
* Step 6 should result in the top command starting and printing it's results to the screen, as it is not interactive you will need to issue ctrl-c to stop the container with a KILL signal
* Step 7 should result in an error and the following message:  
```
exec: "fakeCommand": executable file not found in $PATH
docker: Error response from daemon: Container command not found or does not exist..
```
* Step 8 should result in an error and the following message:  
```
docker: Error parsing reference: "fakeImage" is not a valid repository/tag.
```
* Step 10 should result in success and the output should indicate that the ping succeeded across containers just using the linked name

#Possible Problems:
None