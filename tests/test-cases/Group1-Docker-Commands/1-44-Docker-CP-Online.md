Test 1-44 - Docker CP Online
=======

# Purpose:
To verify that docker cp command for online containers is supported by VIC appliance

# References:
[1 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/cp/)

# Environment:
This test requires that a vSphere server is running and available

# Test Steps:
1. Deploy VIC appliance to vSphere server and set up test files, directories and volumes
2. Pull busybox image and run a container named online
3. Create directory online:/newdir and file online:/newdir/test.txt
4. Issue docker cp online:/newdir newdir to the new VIC appliance
5. Inspect the host cwd to verify that the copy operation succeeded and clean up copied files
6. Issue docker cp online:/newdir/. bar to the new VIC appliance
7. Inspect bar on the host to verify that the copy operation succeeded and clean up copied files
8. Issue docker cp online:/newdir/test.txt foo.txt to the new VIC appliance
9. Verify that the copy operation succeeded
10. Issue docker cp foo.txt online:/doesnotexist/ to the new VIC appliance
11. Issue docker cp ./foo.txt online:/ to the new VIC appliance
12. Issue docker cp ./bar online:/ to the new VIC appliance
13. Inspect online:/ to verify that the copy operations succeeded
14. Remove online
15. Run a container called online with a single volume attached to it
16. Issue docker cp ./bar online:/vol1/ to the new VIC appliance
17. Inspect online:/vol1 to verify that the copy operation succeeded
18. Create a container callecd offline that shares a volume with online
19. Issue docker cp content offline:/vol1 to the new VIC appliance
20. Inspect online:/vol1 to verify that the copy operation succeeded
21. Issue docker cp offline:/vol1 . to the new VIC appliance
22. Verify that /vol1 and its content are copied over to host successfully and clean up copied files
23. Remove offline
24. Issue docker cp largefile.txt online:/vol1/ to the new VIC appliance
25. Inspect online:/vol1 to verify that the large file is copied successfully
26. Issue docker cp online:/dne . to the new VIC appliance
27. Issue docker cp online:/dne/. . to the new VIC appliance
28. Run a container called subVol with 2 volumes attached to it
29. Issue docker cp ./mnt subVol:/ to the new VIC appliance
30. Inspect subVol:/mnt, subVol:/mnt/vol1 and subVol:/mnt/vol2 to verify that the copy operation succeeded
31. Issue docker cp subVol:/mnt ./result to the new VIC appliance
32. Inspect ./result on the host to verify that copy succeeded and remove it afterwards
33. Remove subVol
34. Create a container called subVol_off with 2 volumes attached to it, one volume is shared with an online container
35. Issue docker cp ./mnt subVol_off:/ to the new VIC appliance
36. Stop the online container
37. Start subVol_off to inspect subVol_off:/mnt, subVol_off:/mnt/vol1 and subVol_off:/mnt/vol2 to verify the copy operation succeeded
38. Stop subVol and start online
39. Issue docker cp subVol_off:/mnt ./result to the new VIC appliance
40. Inspect ./result on the host to verify that copy succeeded and remove it afterwards
41. Remove subVol_off and online
42. Run a container called online with one volume attached to it
43. Start 10 background processes that issues docker cp foo.txt online:/foo-${idx} to the new VIC appliance
44. Wait for these processes to finish
45. Inspect online:/ to verify that copy operation succeeded
46. Start 10 background processes that issues docker cp largefile.txt online:/vol1/lg-${idx} to the new VIC appliance
47. Wait for these processes to finish
48. Inspect online:/vol1 to verify that copy operation succeeded
49. Start 10 background processes that issues docker cp online:/vol1/lg-${idx} . to the new VIC appliance
50. Wait for these processes to finish
51. Verify that the copy operation succeeded and clean up all the files copied to the host
52. Remove online

. Remove online
. Clean up created files, directories and volumes

# Expected Outcome:
* Step 1-9 should all succeed
* Step 10 should fail with no such directory
* Step 11-14 should all succeed
* Step 15-25 should all succeed
* Step 26-27 should both fail with no such directory
* Step 28-33 should all succeed
* Step 34-41 should all succeed
* Step 42-52 should all succeed

# Possible Problems:
Step 6-9 won't pass until we support cid:/. as source path