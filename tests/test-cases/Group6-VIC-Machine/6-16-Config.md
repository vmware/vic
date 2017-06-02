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

2. Configure VCH http proxy
3. Verify http proxy is set correctly through govc

2. Configure VCH
3. Run a variety of docker commands
4. Check the debug state of the VCH
5. Check the debug state of an existing containerVM
6. Configure the VCH by setting the debug state to 0
7. Check the debug state of the VCH
8. Check the debug state of the existing containerVM
9. Create a new container and check the debug state of it
10. Configure the debug state of the VCH again and check whether there is only a single snapshot left
11. check whether the output of vic-machine inspect contains the desired debug state


# Expected Outcome
* All steps should succeed