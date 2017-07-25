Test 1-42 - Docker Push
=======

# Purpose:
To verify that `docker push` is supported and works as expected.

# References:
[1 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/push/)

# Environment:
This test requires that a vSphere server is running and available and a local image repository

# Test Steps:
1. Pull latest busybox image, busybox:1.26.0, ubuntu, alpine from docker hub
2. Tag each of the images to the local image repository
3. Issue docker push <local-repo>/busybox
4. Issue docker push <local-repo>/busybox:1.26.0
5. Issue docker push <local-repo>/ubuntu
6. Issue docker push --disable-content-trust alpine
7. Make a small change to the busybox image, commit it, then re-tag it as busybox:test
8. Issue docker push <local-repo>/busybox:test
9. Issue docker push <local-repo>/fakeimage
10. Issue docker push fakeRepo/busybox

# Expected Outcome:
* Steps 1-8 should all succeed without error
* Step 8 should result in only the new layer being uploaded
* Steps 9 and 10 should cause an error with an known error message

# Possible Problems:
* None