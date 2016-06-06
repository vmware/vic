Test 3-1 - Force Install vic-machine
=======

#Purpose:
To verify that vic-machine -force parameter works to re-install VCH Appliance when it is already installed

#References:
* vic-machine --help

#Environment:
This test requires that an ESXi or vCenter server is running and available.

#Test Steps:
1. Deploy VIC appliance to the test server using vic-machine
2. Re-deploy VIC appliance to the test server using vic-machine -force

#Expected Outcome:
vic-machine should succeed at both installs, no errors thrown and provide the IP address of the VCH appliance at the end.

#Possible Problems:
None