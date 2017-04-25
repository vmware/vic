Test 11-04 - Upgrade UpdateInProgress
=======

#Purpose:
To verify that vic-machine inspect could detect the upgrade status of a VCH

#Environment:
This test requires that a vSphere server is running and available

#Test Steps:
1. Download vic_7315.tar.gz from bintray
2. Deploy VIC 7315 to vsphere server
3. Set UpdateInProgress to true using govc
4. Upgrade VCH
5. Run vic-machine upgrade --resetInProgressFlag to reset UpdateInProgress to false
6. Upgrade VCH
7. Run vic-machine inspect to check the upgrade status of the VCH (this should run in parallel with step 6)
8. After step 3 finishes, run step 4 again.

#Expected Outcome:
* In step 4, output should contain "Upgrade failed: another upgrade/configure operation is in progress"
* In step 5, output should contain "Reset UpdateInProgress flag successfully"
* In step 6, output should contain "Completed successfully"
* In step 7, output should contain "Upgrade/configure in progress"
* In step 8, output should not contain "Upgrade/configure in progress"
