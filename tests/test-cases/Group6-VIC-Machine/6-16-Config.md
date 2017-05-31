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
<<<<<<< Updated upstream
2. Configure VCH
3. Run a variety of docker commands
=======
2. Configure VCH http proxy
3. Verify http proxy is set correctly through govc
>>>>>>> Stashed changes

# Expected Outcome
* All steps should succeed