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
This test requires access to VMware Nimbus cluster for dynamic ESXi, vCenter, and NFS server creation

# Test Steps:
1. Deploy VIC appliance to vCenter with a simple cluster and one NFS server; nfs share point added here.
2. Issue docker volume create using no name volume (unnamed) on VolumeStore=nfsVolumeStoreDefault (NFS volume)
3. Run docker mount on container using unnamed volume, run docker wait, then docker rm on temp container
4. Issue docker volume create --name=nfs_default_%{VCH-NAME} on VolumeStore=nfsVolumeStoreDefault (NFS volume)
5. Run docker mount on container using named volume, run docker wait, then docker rm on temp container
6. Issue docker volume create on unnamed volume
7. Issue docker volume create on named volume
8. Issue docker volume create --name=test???
9. Create container ${createFileContainer} using named nfs volume
10. Issue docker exec -i ${createFileContainer} echo # --> to write contents to a file (created by this echo command) on NFS volume
11. Issue docker exec -i ${createFileContainer} ls   # to verify file is created in the correct directory on NFS volume
12. Issue docker exec -i ${createFileContainer} cat  # to verify contents of the file
13. Create container ${addToFileContainer} using named nfs volume
14. Issue docker exec -i ${addToFileContainer} echo # to append to contents of the file used by earlier container
15. Issue docker exec -i ${createFileContainer} echo # to append to the same file
16. Issue docker exec -i ${addToFileContainer} cat # verify contents of the file
17. Issue docker exec -i ${addToFileContainer} echo # write to the file some more
18. Issue docker exec -i ${createFileContainer} cat # verify contents of the file
19. Issue docker volume inspect ${nfsNamedVolume}
20. Issue docker volume ls
21. Issue docker volume rm ${nfsDefaultVolume}
22. Issue docker volume rm ${nfsNamedVolume}


# Expected Outcome:
* Step 2 should complete successfully and return a long string name for the volume created
* Step 3 should verify that the NFS volume is mounted on a temp container; we then rm this container for cleanup
* Step 4 should complete successfully and return named volume
* Step 5 same as step 3 but using the named volume instead
* Step 6 should result in error with the following error message:
```
Error response from daemon: A volume named ${nfsDefaultVolume} already exists. Choose a different volume name.
```
* Step 7 should result in error with the following error message:
 ```
 Error response from daemon: A volume named ${nfsNamedVolume} already exists. Choose a different volume name.
 ```
* Step 8 should result in error with the following message:
```
Error response from daemon: create test???: "test???" includes invalid characters for a local volume name, only "\[a-zA-Z0-9][a-zA-Z0-9_.-]" are allowed
```
* Step 9 should complete successfully creating ${createFileContainer} using ${nfsNamedVolume} (the named nfs volume from before)
* Step 10 should result in success and create/write contents to the /mydata/test_nfs_file.txt file using the ${createFileContainer}
* Step 11 should result in success and verify that the 'test_nfs_file.txt' file is in the mydata directory using the ${createFileContainer}
* Step 12 should result in success and verify the contents of the 'test_nfs_file.txt' file using the ${createFileContainer}
* Step 13 should complete successfully creating ${addToFileContainer} using ${nfsNamedVolume} (the named nfs volume from before)
* Step 14 should result in success and append a line to the end of the 'test_nfs_file.txt' using ${addToFileContainer}
* Step 15 should result in success and append a line to the end of the 'test_nfs_file.txt' using ${createFileContainer}
* Step 16 should result in success and verify the contents of the 'test_nfs_file.txt' file using the ${addToFileContainer}
* Step 17 same as Step 14
* Step 18 same as Step 12
* Step 19 should result in a properly formatted JSON response
* Step 20 should result in each nfs volume being listed with both driver and volume name
* Step 21 should result in success and the volume should not be listed anymore
* Step 22 should result in error with the following message:  
```
Error response from daemon: volume ${nfsNamedVolume} in use by
```

# Possible Problems:
Mount command may be affected by Nimbus's performance returning '' when the volume was successfully created/mounted