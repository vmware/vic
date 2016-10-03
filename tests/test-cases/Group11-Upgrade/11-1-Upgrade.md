Test 11-1 - Upgrade
=======

#Purpose:
To verify vic-machine upgrade can upgrade VCH from a certain version

#Environment:
This test requires that a vSphere server is running and available

#Test Steps:
1. Download vic_5316.tar.gz from bintray
2. Deploy VIC 5287 to the vSphere server
3. Issue docker network create bar, creating a new network called "bar"
4. Create container with port mapping
5. Upgrade VCH to latest version
6. Check the previous created container and image are still there

#Expected Outcome:
* All steps should result in success

#Possible Problems:
None
