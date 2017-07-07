Test 1-41 - Docker CP
=======

# Purpose:
To verify that docker cp command is supported by VIC appliance

# References:
[1 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/cp/)

# Environment:
This test requires that a vSphere server is running and available

# Test Steps:
1. Deploy VIC appliance to vSphere server
2. Pull busybox image and create an interactive docker container named test1
3. Create a file named foo.txt and a directory named bar
4. Issue docker cp foo.txt test1:/ to the new VIC appliance
5. Issue docker cp bar test1:/ to the new VIC appliance
6. Start test1 and issue docker exec test1 ls / to the new VIC appliance to verify the copy operations succeeded
7. Issue docker exec test1 sh -c 'rm /foo.txt && rmdir /bar' to the new VIC appliance to clean up
8. Issue docker cp bar test1:/bar to the new VIC appliance
9. Issue docker exec test1 ls / to the new VIC appliance to verify bar is there
10. Issue docker exec test1 sh -c 'mkdir newdir && echo "testing" > /newdir/test.txt' to the new VIC appliance
11. Issue docker cp test1:/newdir newdir to the new VIC appliance
12. Check the current directory to verify newdir and test.txt are copied successfully
13. Issue docker cp test1:/newdir/. bar to the new VIC appliance
14. Verify that bar/test.txt exists
15. Issue docker cp test1:/newdir/test.txt foo.txt to the new VIC appliance
16. Verify that the content foo.txt is testing
17. Issue docker cp foo.txt test1:/doesnotexist/ to the new VIC appliance
18. Issue docker exec test1 sh -c 'rm -rf /newdir && rmdir /bar' to the new VIC appliance to clean up
19. Issue docker stop -t 0 test1 to the new VIC appliance
20. Issue docker cp foo.txt test1:/ to the new VIC appliance
21. Issue docker start test1 to the new VIC appliance
22. Issue docker exec test1 ls / to the new VIC appliance to verify the copy operation succeeded
23. Issue docker rm -f test1 to the new VIC appliance to clean up
24. Issue docker volume create --name vol1 to the new VIC appliance to create a volume
25. Issue docker create -i --name test1 -v vol1:/vol1 to the new VIC appliance using vol1 as its volume
26. Issue docker cp $foo.txt test1:/vol1 to the new VIC appliance
27. Issue docker start test1 to the new VIC appliance
28. Issue docker exec test1 ls /vol1 to the new VIC appliance to verify the copy operation succeeded
29. Issue docker cp bar test1:/vol1/ to the new VIC appliance
30. Issue docker exec test1 ls /vol1 to the new VIC appliance to verify the copy operation succeeded
31. Issue docker rm -f test1 to the new VIC appliance to clean up the container
32. Issue docker volume create --name vol2 to the new VIC appliance
33. Issue docker create -i --name test1 -v vol1:/vol1 -v vol2:/vol1/vol2 busybox to the new VIC appliance
34. Issue docker cp $foo.txt test1:/vol1/vol2 to the new VIC appliance
35. Start the container and issue docker exec test1 ls /vol1/vol2 to the new VIC appliance to verify the copy operation succeeded
36. Issue docker rm -f test1 to the new VIC appliance to delete the container
37. Issue docker volume create --name vol3 to the new VIC appliance
38. Issue docker create -i --name test1 -v vol1:/vol1 -v vol2:/vol1/vol2 -v vol3:/vol1/vol2/vol3 busybox to the new VIC appliance
39. Issue docker cp foo.txt test1:/vol1/vol2/vol3 to the new VIC appliance
40. Start the container and issue docker exec test1 ls /vol1/vol2/vol3 to the new VIC appliance to verify the copy operation succeeded
41.

# Expected Outcome:
* Step 2-7 should all succeed
* Step 8-16 should all succeed
* Step 17 should fail with no such directory error
* Step 18-23 should all succeed
* Step 24-32 should all succeed
* Step 33-40 should all succeed






# Possible Problems:

