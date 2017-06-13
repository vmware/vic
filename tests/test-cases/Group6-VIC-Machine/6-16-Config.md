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
24. Run vic-machine inspect --conf

# Expected Outcome
* Steps 1 to 11 should succeed
* Steps 12 and 13 should succeed and the output should contain the new container network
* Step 14 should succeed
* Step 15 should fail with an error message saying that the existing container network must be specified
* Step 16 should succeed
* Steps 17 and 18 should succeed and the output should contain both container networks
* Step 19 should succeed
* Step 20 should fail with an error message saying that changes to existing container networks are not supported
* Steps 21 and 22 should succeed
* Steps 23 and 24 should succeed
* Step 24's output should contain the operations user's name and the host thumbprint
