Test 1-04 - Docker Create
=======

# Purpose:
To verify that docker create command is supported by VIC appliance

# References:
[1 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/create/)

# Environment:
This test requires that a vSphere server is running and available

# Test Steps:
1. Deploy VIC appliance to vSphere server
2. Docker pull busybox
3. Issue docker create busybox to the new VIC appliance
4. Issue docker create -t -i busybox to the new VIC appliance
5. Issue docker create --name test1 busybox to the new VIC appliance
6. Issue docker create -v /var/log busybox ls /var/log to the new VIC appliance
7. Issue docker start <containerID> to the new VIC appliance
8. Issue docker logs <containerID> to the new VIC appliance
9. Issue docker create -v test-named-vol:/testdir busybox
10. Issue docker start <containerID>
11. Issue docker logs <containerID> to grab the disk size of the volume
12. Issue docker create -v /dir:/dir busybox
13. Issue docker create busybox /bin/top to the new VIC appliance
14. Issue docker create fakeimage to the new VIC appliance
15. Issue docker create fakeImage to the new VIC appliance
16. Issue docker create --name busy1 busybox /bin/top to the new VIC appliance
17. Issue docker start busy1 to the new VIC appliance
18. Issue docker create --link busy1:busy1 --name busy2 busybox ping -c2 busy1 to the new VIC appliance
19. Issue docker start busy2 to the new VIC appliance
20. Issue docker logs busy2 to the new VIC appliance
21. Create a container, rm the container, then create another container
22. Create a container directly without pulling the image first for an image that hasn't been pulled yet
23. Create a container without specifying a command
24. Create a container with a custom amount of CPUs
25. Create a container with a custom amount of Memory in GB
26. Create a container with a custom amount of Memory in MB
27. Create a container with a custom amount of Memory in KB
28. Create a container with a custom amount of Memory in Bytes
29. Create a container using a rest api call without HostConfig in the form data
30. Create a container, then check the vm display name in vsphere through govc
31. Create a container, then check the vm Destroy_Task method is disabled in VC through govc
32. Create two containers with the same name in parallel, then check that only one attempt is successful
33. Remove the container from Step 32 by name and create another container with the same name

# Expected Outcome:
* Steps 3-7 should all return without error and printing the container ID on return
* Step 8 should show that the contents of the containers /var/log matches the contents of the hosts /var/log
* Steps 9, 10 and 11 should return without errors and should successfully create a new volume called `test-named-vol` with disk size 975.9M
* Step 12 should return with the error message - Error response from daemon: vSphere Integrated Containers does not support mounting directories as a data volume.
* Step 14 should return with the error message - Error: image library/fakeimage not found
* Step 15 should return with the error message - Error parsing reference: "fakeImage" is not a valid repository/tag
* Step 18 should result in success and the busy2 container should exist
* Step 20 should show that busy2 was able to successfully ping busy1 just using the linked name
* Step 21 should result in success for all three parts
* Step 22 should return without error
* Step 23 should return with the following error message - Error response from daemon: No command specified
* Steps 24-28 should return without error.
* Step 29 should return without error.
* Step 30 should show that the VM display name equals to containerName-containerShortID and datastore folder name equal to containerID
* Step 31 should show that the VM Destroy_Task method is disabled in VC
* Step 32 should have one container create process succeed and the other fail with an error
* Step 33 should succeed

# Possible Problems:
None
