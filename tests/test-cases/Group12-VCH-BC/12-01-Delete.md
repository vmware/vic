Test 12-01 - Delete
=======

#Purpose:
To verify vic-machine delete can delete VCH created by vic 0.6.0

#Environment:
This test requires that a vSphere server is running and available

#Test Steps:
1. Download vic_0.6.0.tar.gz from bintray
2. Deploy VIC 0.6.0 to vsphere server
3. Create container
3. Using latest version vic-machine to delete this VCH

#Expected Outcome:
* All steps should result in success
