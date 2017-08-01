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
2. Issue vic-machine inspect config command
3. Issue vic-machine inspect config --format raw command
4. Create a VCH with tlsverify
5. Inspect the VCH without specifying --tls-cert-path
6. Inspect the VCH with a valid --tls-cert-path
7. Inspect the VCH with an invalid --tls-cert-path
8. Create a VCH with --no-tls
9. Inspect the VCH without specifying --tls-cert-path
10. Create a VCH with --no-tlsverify
11. Inspect the VCH without specifying --tls-cert-path
12. Create a VCH with some container-network options
13. Verify the `inspect config` contains the correct options

# Expected Outcome:
* Steps 1 should succeed, and output from step 2 and 3 should contain expected flags & values
* Steps 4-13 should complete successfully, however, step 6 should show a warning in the output (see below)
* The output of steps 5 and 6 should contain the correct `DOCKER_CERT_PATH`.
* The output of step 7 should not contain a `DOCKER_CERT_PATH` and should contain:
```
Unable to find valid client certs
DOCKER_CERT_PATH must be provided in environment or certificates specified individually via CLI arguments
```
* The outputs of steps 9 and 11 should not contain a `DOCKER_CERT_PATH` and should not contain:
```
Unable to find valid client certs
DOCKER_CERT_PATH must be provided in environment or certificates specified individually via CLI arguments
```
