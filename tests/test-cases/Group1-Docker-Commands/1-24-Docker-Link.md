Test 1-24 - Docker Link
=======

#Purpose:
To verify that docker --link/--net-alias commands are supported by VIC appliance

#References:
[1 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/run/)

#Environment:
This test requires that a vSphere server is running and available

#Test Steps:
1. Deploy VIC appliance to vSphere server
2. Issue docker volume pull busybox to the VIC appliance
3. Issue docker create --name first busybox to the VIC appliance
4. Issue docker start to the VIC appliance
5. Issue docker create busybox ping -c3 first to the VIC appliance
6. Issue docker start to the VIC appliance
7. Issue docker create --link first:1st busybox ping -c3 1st to the VIC appliance
8. Issue docker start to the VIC appliance
9. Issue docker create --net-alias 2nd busybox to the VIC appliance
10. Issue docker start to the VIC appliance
11. Issue docker create busybox ping -c3 2nd to the VIC appliance
12. Issue docker start to the VIC appliance

#Expected Outcome:
* Every step should result in success

#Possible Problems:
None