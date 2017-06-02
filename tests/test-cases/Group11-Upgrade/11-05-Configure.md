Test 11-05 - Configure
=======

# Purpose:
To verify vic-machine configure can upgrade with --upgrade specified

#Environment:
This test requires that a vSphere server is running and available

# Test Steps:
1. Download vic_7315.tar.gz from gcp
2. Deploy VIC 7315 to vsphere server
3. Using latest version vic-machine to configure this VCH

# Expected Outcome:
* Step 3 should get expected error
