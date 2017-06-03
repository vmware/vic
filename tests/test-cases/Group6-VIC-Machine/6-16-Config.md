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
4. Configure VCH
5. Check the debug state of the VCH
6. Check the debug state of an existing containerVM
7. Configure the VCH by setting the debug state to 0
8. Check the debug state of the VCH
9. Check the debug state of the existing containerVM
10. Create a new container and check the debug state of it
11. Configure the debug state of the VCH again and check whether there is only a single snapshot left
12. check whether the output of vic-machine inspect contains the desired debug state

# Expected Outcome
* All steps should succeed