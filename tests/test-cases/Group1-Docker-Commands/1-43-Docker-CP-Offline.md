Test 1-43 - Docker CP Offline
=======

# Purpose:
To verify that docker cp command for offline containers is supported by VIC appliance

# References:
[1 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/cp/)

# Environment:
This test requires that a vSphere server is running and available

# Test Steps:
1. Deploy VIC appliance to vSphere server, set up test files, directories and volumes
2. Pull busybox image and create a container named offline
3. Issue docker cp foo.txt offline:/ to the new VIC appliance
4. Start offline and verify that foo.txt exists in the offline:/
5. Clean up offline:/foo.txt
6. Create a directory named testdir in offline:/ and create a fakefile in offline:/testdir
7. Stop offline and issue docker cp offline:/testdir . to the new VIC appliance
8. Verify that testdir and its content is successfully copied to the host
9. Issue docker cp bar offline:/bar to the new VIC appliance
10. Start offline and verify that bar exists in offline:/
11. Stop offline
12. Issue docker cp offline:/dne . to the new VIC appliance
13. Issue docker cp offline:/dne/. . to the new VIC appliance
14. Issue docker cp dne/ offline:/ to the new VIC appliance
15. Remove offline container
16. Create a container named offline with a small volume attached to it
17. Issue docker cp largefile.txt offline:/small to the new VIC appliance
18. Remove offline container
19. Create a container named offline with a single volume attached to it
20. Issue docker cp foo.txt offline:/vol1 to the new VIC appliance
21. Start offline and verify that the copy operation succeeded
22. Remove offline container
23. Create a container named offline with two volumes attached to it, vol2 nested in vol1
24. Issue docker cp foo.txt offline:/vol1/vol2 to the new VIC appliance
25. Start offline and verify that the copy operation succeeded
26. Remove offline container
27. Create a container named offline with three volumes attached to it, vol2 nested in vol1, vol3 nested in vol2
28. Issue docker cp foo.txt offline:/vol1/vol2/vol3 to the new VIC appliance
29. Start offline and verify that the copy operation succeeded
30. Remove offline container
31. Create a container named offline with a single volume attached to it
32. Issue docker cp foo.txt offline:/foo1 to the new VIC appliance
33. Issue docker cp foo.txt offline:/foo2 to the new VIC appliance
34. Issue docker cp foo.txt offline:/foo3 to the new VIC appliance
35. Start offline and verify that the copy operation succeeded
36. Stop offline
37. Issue docker cp largefile.txt offline:/vol1/lg1 to the new VIC appliance
38. Issue docker cp largefile.txt offline:/vol1/lg2 to the new VIC appliance
39. Issue docker cp largefile.txt offline:/vol1/lg3 to the new VIC appliance
40. Start offline and verify that the copy operation succeeded
41. Stop offline
42. Issue docker cp offline:/vol1/lg1 ./large1 to the new VIC appliance
43. Issue docker cp offline:/vol1/lg1 ./large2 to the new VIC appliance
44. Issue docker cp offline:/vol1/lg1 ./large3 to the new VIC appliance
45. Verify that the copy operation succeeded and clean up all the files copied to the host
46. Remove offline container

# Expected Outcome:
* Step 1-11 should all succeed
* Step 12-14 should all fail with no such directory error
* Step 15-16 should all succeed
* Step 17 should fail
* Step 18-30 should all succeed
* Step 31-46 should all succeed

# Possible Problems:
12-13 are skipped until we properly propagate errors back to user
17 should fail with not enough disk space, but currently returns a broken pipe
31-46 are skipped until we properly handle concurrent requests for offline cp