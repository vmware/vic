Test 23-01 - Version
=======

# Purpose:
To verify vic-machine-server returns a valid version number

# References:
1. [The design document](../../../doc/design/vic-machine/service.md)

# Environment:
This test has no environmental requirements

# Test Steps:
1. Start the vic-machine-server
2. Wait 1 second
3. Use curl to issue a GET request for the version

# Expected Outcome:
* Step 3 should succeed with a 200 OK response containing a version number

# Possible Problems:
* Step 1 could take more than 1 second to complete causing step 3 to fail with a return code of 7. (Other tests in this suite may wait on a response from the version endpoint to determine whether the service is available, but this test should not as it is the endpoint under test.)

