Test 1-43 - Docker CP
=======

# Purpose:
To verify that docker cp command is supported by VIC appliance

# References:
[1 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/cp/)

# Environment:
This test requires that a vSphere server is running and available

# Test Steps:
1. Deploy VIC appliance to vSphere server and set up test files and directories
2. Pull busybox image and create a container named offline
3. Create a file named foo.txt and a directory named bar
4. Issue docker cp foo.txt offline:/ to the new VIC appliance
5. Issue docker cp bar offline:/ to the new VIC appliance
6. Start offline and issue docker exec offline ls / to the new VIC appliance to verify the copy operations succeeded
7. Clean up copied files and directories from offline
8. Issue docker exec offline sh -c 'mkdir testdir && echo "file content" > /testdir/fakefile' to the new VIC appliance
9. Stop offline and Issue docker cp offline:/testdir ${CURDIR}/ to the new VIC appliance
10. Verify that testdir is copied over successfully and clean up testdir on host
11. Issue docker cp bar offline:/bar to the new VIC appliance
12. Start offline and issue docker exec offline ls / to the new VIC appliance to verify the copy operation succeeded
13. Issue docker rm -f offline to the new VIC appliance to clean up the container
14. Pull busybox image and create a container named online
15. Start online and issue docker exec online sh -c 'mkdir newdir && echo "testing" > /newdir/test.txt' to the new VIC appliance
16. Issue docker cp online:/newdir newdir to the new VIC appliance
17. Verify that newdir is copied over successfully and clean up newdir on host
18. Issue docker cp online:/newdir/. /bar to the new VIC appliance
19. Verify that content in newdir are copied over successfully and remove the content on host
20. Issue docker cp online:/newdir/test.txt foo.txt to the new VIC appliance
21. Verify that test.txt is copied over to host successfully and overwrites foo.txt
22. Issue docker cp foo.txt online:/doesnotexist/ to the new VIC appliance
23. Issue docker cp foo.txt online:/ to the new VIC appliance
24. Issue docker cp bar online:/ to the new VIC appliance
25. Issue docker exec online ls / to the new VIC appliance to verify the copy operation succeeded
26. Issue docker rm -f online to the new VIC appliance to remove online
27. Issue docker create -i --name online -v vol1:/vol1 ${busybox} to the new VIC appliance
28. Start online and issue docker cp bar online:/vol1/ to the new VIC appliance
29. Issue docker exec online ls /vol1 to the new VIC appliance to verify the copy operation succeeded
30. Issue docker rm -f online to the new VIC appliance to clean up online
31. Issue docker create -i --name offline -v vol1:/vol1 ${busybox} to the new VIC appliance
32. Issue docker cp foo.txt offline:/vol1 to the new VIC appliance
33. Start offline and issue docker exec offline ls /vol1 to the new VIC appliance to verify the copy operation succeeded
34. Issue docker rm -f offline to the new VIC appliance to clean up offline
35. Issue docker create -i --name offline -v vol1:/vol1 -v vol2:/vol1/vol2 ${busybox} to the new VIC appliance
36. Issue docker cp foo.txt offline:/vol1/vol2 to the new VIC appliance
37. Start offline and issue docker exec offline ls /vol1/vol2 to the new VIC appliance to verify the copy operation succeeded
38. Issue docker rm -f offline to the new VIC appliance to clean up offline
39. Issue docker create -i --name offline -v vol1:/vol1 -v vol2:/vol1/vol2 -v vol3:/vol1/vol2/vol3 ${busybox} to the new VIC appliance
40. Issue docker cp foo.txt offline:/vol1/vol2/vol3 to the new VIC appliance
41. Start offline and issue docker exec offline ls /vol1/vol2/vol3 to the new VIC appliance to verify the copy operation succeeded
42. Issue docker rm -f offline to the new VIC appliance to clean up offline
43. Issue docker create -i --name online -v vol1:/vol1 ${busybox} to the new VIC appliance
44. Issue docker create -i --name offline -v vol1:/vol1 ${busybox} to the new VIC appliance
45. Start online and issue docker cp content offline:/vol1 to the new VIC appliance
46. Issue docker exec online ls /vol1 to the new VIC appliance to verify the copy operation succeeded
47. Issue docker cp offline:/vol1 . to the new VIC appliance
48. Verify that all contents in /vol1 are copied over successfully
49. Issue docker rm -f offline to the new VIC appliance to clean up offline
50. Generate a file of 4 MiB and issue docker cp largefile.txt online:/vol1/ to the new VIC appliance
51. Issue docker exec online ls -l /vol1/ to the new VIC appliance to verify the copy operation succeeded
52. Issue docker rm -f online to the new VIC appliance to clean up online
53. Clean up files and directories and volumes created at setup time


# Expected Outcome:
* Step 1-7 should all succeed
* Step 8-13 should all succeed
* Step 14-21 should all succeed
* Step 22 should fail with no such directory error
* Step 23-26 should all succeed
* Step 27-34 should all succeed
* Step 35-42 should all succeed
* Step 43-53 should all succeed


# Possible Problems: