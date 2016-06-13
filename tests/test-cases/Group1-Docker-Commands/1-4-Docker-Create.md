Test 1-4 - Docker Create
=======

#Purpose:
To verify that docker create command is supported by VIC appliance

#References:
[1 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/create/)

#Environment:
This test requires that a vSphere server is running and available

#Test Steps:
1. Deploy VIC appliance to vSphere server
2. Docker pull busybox
3. Issue docker create busybox to the new VIC appliance
4. Issue docker create -t -i busybox to the new VIC appliance
5. Issue docker create --name test1 busybox to the new VIC appliance
6. Issue docker create -v /var/log:/var/log busybox ls /var/log to the new VIC appliance
7. Issue docker start <containerID> to the new VIC appliance
8. Issue docker logs <containerID> to the new VIC appliance
9. Issue docker create busybox /bin/top to the new VIC appliance
10. Issue docker create fakeimage to the new VIC appliance
11. Issue docker create fakeImage to the new VIC appliance
12. Issue docker create --name busy1 busybox /bin/top to the new VIC appliance
13. Issue docker start busy1 to the new VIC appliance
14. Issue docker create --link busy1:busy1 --name busy2 busybox ping -c2 busy1 to the new VIC appliance
15. Issue docker start busy2 to the new VIC appliance
16. Issue docker logs busy2 to the new VIC appliance

#Expected Outcome:
* Steps 3-7 should all return without error and printing the container ID on return
* Step 8 should show that the contents of the containers /var/log matches the contents of the hosts /var/log
* Step 10 should return with the error message - Error: image library/fakeimage not found
* Step 11 should return with the error message - Error parsing reference: "fakeImage" is not a valid repository/tag
* Step 14 should result in success and the busy2 container should exist
* Step 16 should show that busy2 was able to successfully ping busy1 just using the linked name 

#Possible Problems:
None