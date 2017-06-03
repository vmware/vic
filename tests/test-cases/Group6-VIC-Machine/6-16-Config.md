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

# Expected Outcome
* All steps should succeed