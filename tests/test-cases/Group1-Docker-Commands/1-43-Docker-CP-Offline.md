Test 1-43 - Docker CP Offline
=======

# Purpose:
To verify that docker cp command for offline containers is supported by VIC appliance

# References:
[1 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/cp/)

# Environment:
This test requires that a vSphere server is running and available

# Test Steps:
1. Deploy VIC appliance to vSphere server and set up test files, directories and volumes
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
32. Start 10 background processes that issues docker cp foo.txt offline:/foo-${idx} to the new VIC appliance
33. Wait for these processes to finish
34. Start and inspect offline:/ to verify that copy operation succeeded
35. Stop offline
36. Start 10 background processes that issues docker cp largefile.txt offline:/vol1/lg-${idx} to the new VIC appliance
37. Wait for these processes to finish
38. Start and inspect offline:/vol1 to verify that copy operation succeeded
39. Stop offline
40. Start 10 background processes that issues docker cp offline:/vol1/lg-${idx} . to the new VIC appliance
41. Wait for these processes to finish
42. Verify that the copy operation succeeded and clean up all the files copied to the host
43. Remove offline container
44. Clean up created files, directories and volumes

# Expected Outcome:
* Step 1-11 should all succeed
* Step 12-14 should all fail with no such directory error
* Step 15-16 should all succeed
* Step 17 should fail
* Step 18-30 should all succeed
* Step 31-44 should all succeed

# Possible Problems:
12-13 are skipped until we properly propagate errors back to user
17 should fail with not enough disk space, but currently returns a broken pipe
31-43 are skipped until we properly handle concurrent requests for offline cp