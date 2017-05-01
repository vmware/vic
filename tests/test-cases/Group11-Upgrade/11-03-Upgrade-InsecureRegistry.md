Test 11-03 - Upgrade InsecureRegistry
=======

#Purpose:
To verify InsecureRegistries are correctly migrated

#Environment:
This test requires that a vSphere server is running and available

#Test Steps:

1. Deploy an insecure registry on http
2. Create a test project on insecure registry
3. Start a docker daemon with insecure registry configured
4. Push test image to insecure registry through above docker daemon(VCH docker push not implemented)
5. Download vic_7315.tar.gz from bintray if it does not exist locally
6. Deploy VIC 7315 to vsphere server with above insecure registry
7. Make sure pull given test image through VCH successfully
8. Upgrade VCH to latest version
9. Make sure pull given test image through VCH successfully

10. Deploy an insecure registry on https with self-signed certificate
11. Create a test project on insecure registry
12. Start a docker daemon with insecure registry configured
13. Push test image to insecure registry through above docker daemon(VCH docker push not implemented)
14. Download vic_7315.tar.gz from bintray if it does not exist locally
15. Deploy VIC 7315 to vsphere server with above insecure registry
16. Make sure pull given test image through VCH successfully
17. Upgrade VCH to latest version
18. Make sure pull given test image through VCH successfully

#Expected Outcome:
* Able to pull given test image through VCH successfully both before and after upgrade
