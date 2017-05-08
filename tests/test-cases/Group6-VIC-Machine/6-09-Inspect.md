Test 6-09 - Verify vic-machine inspect
=======

# Purpose:
Verify vic-machine inspect functionality

# References:
* vic-machine-linux inspect -h

# Environment:
This test requires that a vSphere server is running and available

# Test Steps:
1. Install VCH
2. Issue vic-machine inspect --conf command
3. Create a VCH with tlsverify
4. Inspect the VCH without specifying --cert-path
5. Inspect the VCH with a valid --cert-path
6. Inspect the VCH with an invalid --cert-path
7. Create a VCH with --no-tls
8. Inspect the VCH without specifying --cert-path
9. Create a VCH with --no-tlsverify
10. Inspect the VCH without specifying --cert-path

# Expected Outcome:
* Steps 1 should succeed, and output from step 2 should contain expected flags & values
* Steps 3-10 should complete successfully, however, step 6 should show a warning in the output (see below)
* The output of steps 4 and 5 should contain the correct `DOCKER_CERT_PATH`.
* The output of step 6 should not contain a `DOCKER_CERT_PATH` and should contain:
```
Unable to find valid client certs
DOCKER_CERT_PATH must be provided in environment or certificates specified individually via CLI arguments
```
* The outputs of steps 8 and 10 should not contain a `DOCKER_CERT_PATH` and should not contain:
```
Unable to find valid client certs
DOCKER_CERT_PATH must be provided in environment or certificates specified individually via CLI arguments
```

