Test 23-03 - VCH Create
=======

# Purpose:
To verify vic-machine-server can create a VCH with a specified configuration

# References:
1. [The design document](../../../doc/design/vic-machine/service.md)

# Environment:
This test requires a vSphere system where VCHs can be deployed

# Test Steps:
1. Create a VCH with as minimal a configuration as possible
2. Inspect that VCH using the CLI
3. Create a VCH with a more complex configuration
4. Inspect that VCH using the CLI

# Expected Outcome:
* The results of 2 should contain the same information as was supplied when the VCH was created in 1.
* The results of 4 should contain the same information as was supplied when the VCH was created in 3.

# Possible Problems:
None known
