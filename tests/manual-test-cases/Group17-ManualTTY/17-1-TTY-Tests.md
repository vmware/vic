Test 17-1 - TTY Tests
=======

# Purpose:
To verify that docker commands using TTY work with VIC

# References:


# Environment:
This test requires that a vSphere server is running and available

# Test Cases
1. Deploy VIC appliance to vSphere server
2. Issue docker run -it busybox date to the new VCH
3. Issue docker run -it busybox df to the new VCH
4. Issue docker run -it busybox top to the new VCH

### Expected Outcome:
* Steps 1-4 should all succeed and return the expected output from those commands