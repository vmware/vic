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
11. Configure VCH http proxy
12. Verify http proxy is set correctly through govc
13. Configure the VCH by adding a container network
14. Run docker network ls
15. Run vic-machine inspect --conf
16. Configure the VCH by adding a new container network without specifying the previous network
17. Configure the VCH by adding a new container network while specifying the previous network
18. Run docker network ls
19. Run vic-machine inspect --conf
20. Configure the VCH by attempting to change an existing container network

# Expected Outcome
* Steps 1 to 13 should succeed
* Steps 14 and 15 should succeed and the output should contain the new container network
* Step 16 should fail with an error message saying that the existing container network must be specified
* Steps 17 should succeed
* Steps 18 and 19 should succeed and the output should contain both container networks
* Step 20 should fail with an error message saying that changes to existing container networks are not supported
