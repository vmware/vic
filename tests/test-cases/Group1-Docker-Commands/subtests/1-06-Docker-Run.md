Test 1-06 - Docker Run
=======

# Purpose:
To verify that docker run command is supported by VIC appliance

# References:
[1 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/run/)

# Environment:
This test requires that a vSphere server is running and available

# Test Steps:
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
11. Issue docker run -it busybox /bin/df to the VIC appliance
12. Issue docker run -d -p 6379 redis:alpine to the VIC appliance
13. Issue docker run -it busybox /bin/true
14. Issue docker run -it busybox /bin/false
15. Issue docker run -it busybox /bin/fakeCommand
16. Issue docker run busybox date
17. Create container1 with id1 and then create container2 with name = id1
18. Run a short lived container with autoremove specified
19. Issue docker run -d -v vol:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=pw --name test-mysql mysql to the VIC appliance
20. Issue docker ps to the VIC appliance to verify that test-mysql is running and clean up afterward
21. Issue docker run -d -e MYSQL_ROOT_PASSWORD=pw --name test-mariadb mariadb to the VIC appliance
22. Issue docker ps to the VIC appliance to verify that test-mariadb is running and clean up afterward
23. Issue docker run -d --name test-postgres postgres to the VIC appliance
24. Issue docker ps to the VIC appliance to verify that test-postgres is running and clean up afterward

# Expected Outcome:
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
* Step 11 should result in success with exit code 0 and show the output of the df command
* Step 12 should result in success with exit code 0
* Step 13 should result in success with exit code 0
* Step 14 should result in success with exit code 1
* Step 15 should result in success with exit code 127
* Step 16 should result in success and the output should contain the current date
* Step 17 should result in no conflicts
* Step 18 should result in the same container count at beginning and end
* Step 19-24 should result in success with exit code 0

# Possible Problems:
None
