Test 5-22 - NFS Volume
=======

# Purpose:
To verify that NFS shared volumes work with currently supported docker commands

# References:
[1 - Best practices for running VMware vSphere on NFS](http://www.vmware.com/content/dam/digitalmarketing/vmware/en/pdf/techpaper/vmware-nfs-bestpractices-white-paper-en.pdf)

[2 - Docker Command Line Reference - Volume Create](https://docs.docker.com/engine/reference/commandline/volume_create/)

[3 - Docker Command Line Reference - Exec](https://docs.docker.com/engine/reference/commandline/exec/)

[4 - Docker Command Line Reference - Volume Inspect](https://docs.docker.com/engine/reference/commandline/volume_inspect/)

[5 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/volume_ls/)


# Environment:
This test requires access to VMware Nimbus for dynamic ESXi and NFS server creation

# Test Steps:
1. Deploy VIC appliance to an ESX and use fake NFS mount point
2. Deploy VIC appliance to an ESX and use valid NFS mount point
3. Issue docker volume create using no name for the volume (unnamed) on VolumeStore=nfsVolumeStore (NFS volume)
4. Run docker mount on container using unnamed volume, run docker wait, then docker rm on temp container
5. Issue docker volume create --name=nfs_default_%{VCH-NAME} on VolumeStore=nfsVolumeStore (NFS volume)
6. Run docker mount on container using named volume, run docker wait, then docker rm on temp container
7. Issue docker volume create on unnamed volume
8. Issue docker volume create on named volume
9. Issue docker volume create --name="test!@\#$%^&*()"
10. Create container ${createFileContainer} using named nfs volume
11. Issue docker exec -i ${createFileContainer} echo # --> to write contents to a file (created by this echo command) on NFS volume
12. Issue docker exec -i ${createFileContainer} ls   # to verify file is created in the correct directory on NFS volume
13. Issue docker exec -i ${createFileContainer} cat  # to verify contents of the file
14. Create a container using named nfs volume and echo # to append more contents to the file used by earlier container
15. Create a container using named nfs volume and echo # to append to the same file
16. Create a container using named nfs volume and echo # to append to the same file
17. Create a container using named nfs volume and cat # verify contents of the file
18. Create a detached container using named nfs volume using named nfs volume to cat the file from last test
19. Issue docker logs to see the results of the cat command # verify contents of the file
20. Create a container using named nfs volume and rm the file just used.
21. Create a container using named nfs volume and cat the file that was just removed
22. Issue docker start on detached container from earlier
23. Issue docker logs on detached container
24. Spin up on container per item in a list to write once a sec to a file the value passed in from the list and save the container ids
25. Create container using named nfs volume and cat the contents of the file from the previous step
26. Check output from each container that was writing to the file.
27. Stop all the running write containers.
28. Issue docker volume inspect ${nfsNamedVolume}
29. Issue docker volume ls
30. Issue docker volume rm ${nfsDefaultVolume}
31. Issue docker volume rm ${nfsNamedVolume}
32. Create a detached container using named nfs volume and write to file every second
33. Create a container using named nfs volume and tail the file from previous step
34. Kill the NFS Server from Nimbus
35. Create a container using named nfs volume from killed NFS server and tail the file from previous step
36. Create a container using named nfs volume from killed NFS server and write to file from previous step
37. Create a container using named nfs volume from killed NFS server and ls the mydata directory



# Expected Outcome:
* Step 1 will succeed VCH creation but should fail in mounting the fake NFS mount point
* Step 2 should complete successfully; VCH should be created/installed
* Step 3 should complete successfully and return a long string name for the volume created
* Step 4 should verify that the NFS volume is mounted on a temp container; container rm should succeed
* Step 5 should complete successfully and return named volume
* Step 6 same as step 4 but using the named volume instead
* Step 7 should result in error with the following error message:
```
Error response from daemon: A volume named ${nfsDefaultVolume} already exists. Choose a different volume name.
```
* Step 8 should result in error with the following error message:
 ```
 Error response from daemon: A volume named ${nfsNamedVolume} already exists. Choose a different volume name.
 ```
* Step 9 should result in error with the following message:
```
Error response from daemon: create test???: "test???" includes invalid characters for a local volume name, only "\[a-zA-Z0-9][a-zA-Z0-9_.-]" are allowed
```
* Step 10 - 20 should result in success
* Step 21 should result in error with the following error message:
```
cat: can't open 'mydata/test_nfs_file.txt': No such file or directory
```
* Step 22 and 23 should succeed, however step 23 will show the same error as above in the logs
* Step 24 - 27 should result in success
* Step 28 should result in a properly formatted JSON response
* Step 29 should result in each nfs volume being listed with both driver and volume name
* Step 30 should result in success and the volume should not be listed anymore
* Step 31 should result in error with the following message:  
```
Error response from daemon: volume ${nfsNamedVolume} in use by
```
* Step 32 - 34 should result in success; step 34 should kill/drop the server
* Step 35 should result in error with the following message:
```
Server error from portlayer: unable to wait for process launch status:
```
* Step 36 - 37 should result in error with the rc = 125.


# Possible Problems:
Mount command may be affected by Nimbus's performance returning '' when the volume was successfully created/mounted