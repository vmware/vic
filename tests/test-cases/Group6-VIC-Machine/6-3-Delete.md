Test 6-3 - Verify delete clean up all resources
=======

#Purpose:
Verify vic-machine delete can delete vch installed by vic-machine create

#References:
* vic-machine-linux delete -h
#Environment:
This test requires that a vSphere server is running and available

#Test Steps:
1. Create VCH through vic-machine create
2. Create container
3. Make sure container VM is powered on
4. Delete VCH without --force, to get delete failure
5. Delete VCH with --force, to get delete success
6. Check vSphere through govc to make sure all resources are deleted correctly

#Expected Outcome:
* Step 4 should fail for VM is powered on
* Step 5 should success
* Step 6 should not find any resources created by vic-machine and VCH
