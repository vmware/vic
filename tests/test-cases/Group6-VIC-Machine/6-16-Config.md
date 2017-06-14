Test 6-16 - Verify vic-machine configure
=======

# Purpose:
Verify vic-machine configure

# References:
* vic-machine-linux create -h

# Environment:
This test requires that a vSphere server is running and available

# Test Steps
1. Deploy VCH
2. Configure VCH
3. Check the debug state of the VCH
4. Check the debug state of an existing containerVM
5. Configure the VCH by setting the debug state to 0
6. Check the debug state of the VCH
7. Check the debug state of the existing containerVM
8. Create a new container and check the debug state of it
9. Configure the debug state of the VCH again and check whether there is only a single snapshot left
10. check whether the output of vic-machine inspect contains the desired debug state
11. Configure the VCH by adding a container network
12. Run docker network ls
13. Run vic-machine inspect config
14. Run a container with the new container network
15. Configure the VCH by adding a new container network without specifying the previous network
16. Configure the VCH by adding a new container network while specifying the previous network
17. Run docker network ls
18. Run vic-machine inspect config
19. Run a container with the new container network
20. Configure the VCH by attempting to change an existing container network
21. Configure VCH http proxy
22. Verify http proxy is set correctly through govc
23. Configure the VCH's operations user credentials
24. Run vic-machine inspect config
26. Reset VCH http proxy using VCH ID
26. Verify http proxy is reset correctly through govc
27. Run vic-machine inspect config
28. Configure VCH dns server to 10.118.81.1 and 10.118.81.2
29. Run vic-machine inspect config

# Expected Outcome
* Step 15 should fail with an error message saying that the existing container network must be specified
* Step 20 should fail with an error message saying that changes to existing container networks are not supported
* Step 24's output should contain the operations user's name and the host thumbprint
* All other steps should succeed

# Configure VCH resources
1. Configure VCH resources
2. Verify VCH configuration through vic-machine inspect
3. Configure VCH resources with too small values
4. Verify VCH configuration is rollback to old value

# Expected Outcome
* All other steps should succeed
