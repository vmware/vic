Test 23-13 - VCH Tdnf 
======

# Purpose:
To test whether tdnf in VCH work well 

# References:
[1 - tdnf](https://github.com/vmware/tdnf/wiki)

# Environment:
* requires a working VCSA setup with VCH installed
* Target VCSA has Bash enabled for the root account

# Test Steps:
1. Enable VCH SSh
2. SSH into VCH and check gpgcheck=1 in photon.repo
3. SSH into VCH and run tdnf install which -y

# Expected Outcome:
* Each step should return success

# Possible Problems:
None
