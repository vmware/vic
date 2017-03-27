Test 11-03 - Upgrade InsecureRegistry
=======

#Purpose:
To verify InsecureRegistries are correctly migrated

#Environment:
This test requires that a vSphere server is running and available

#Test Steps:

# Note ideally, we should test against a deployed insecure registry instead of depending on guestinfo
1. Download vic_7315.tar.gz from bintray
2. Deploy VIC 7315 to vsphere server with an insecure registry
3. Get guestinfo from VCH and verify Path is set as hostname while host is empty
4. Upgrade VCH to latest version
5. Get guestinfo from upgraded VCH and verify host is set as hostname and path is empty


#Expected Outcome:
* Guestinfo for insecure registry in upgraded VCH should contain valid host
