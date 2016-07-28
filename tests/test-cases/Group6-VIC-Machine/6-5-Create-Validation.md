Test 6-4 - Verify vic-machine create validation function
=======

#Purpose:
Verify vic-machine create validation functions, this does not include validation for network, datastore, and compute resources

#References:
* vic-machine-linux create -h

#Environment:
This test requires that a vSphere server is running and available

#Test Cases: - suggestion resources

#Test Steps
1. Create with wrong compute-resource: not exist resource pool, not existed vc cluster, not existed datacenter.
2. Create with wrong compute-resource format

#Expected Outcome:
* Verify resource suggestion successfully show available values

#Test Cases: - validate license
#Test steps
1. Prepare env with different license level
2. Verify license validation works for different license
3. If license verification passed, verify VCH deployment succeed.

#Test Cases: - firewall
#Test steps
1. Prepare env with firewall disabled
2. Verify deployment failed for firewall is not enabled with user-friendly error message

#Test Cases: - drs
#Test steps
1. Prepare env with drs disabled
2. Verify deployment failed for drs disabled with user-friendly error message

#Test Cases: - resource accessibility
#Test steps
1. Prepare env with datastore not connected to hosts
2. Verify deployment failed for host/datastore connectability with user-friendly error message
