Test 6-03 - Verify delete clean up all resources
=======

# Purpose:
Verify vic-machine delete can delete vch installed by vic-machine create

# References:
* vic-machine-linux delete -h
# Environment:
This test requires that a vSphere server is running and available

# Test Cases

## Delete VCH and verify
1. Create VCH through vic-machine create
2. Create container
3. Make sure container VM is powered on
4. Delete VCH without --force, to get delete failure
5. Delete VCH with --force, to get delete success
6. Check vSphere through govc to make sure all resources are deleted correctly

### Expected Outcome:
* Step 4 should fail for VM is powered on
* Step 5 should success
* Step 6 should not find any resources created by vic-machine and VCH


## Attach Disks and Delete VM
1. Create VCH normally
2. Pull a moderately sized image (ubuntu)
3. Check to see that disk images exist via govc
4. Use govc to attach disk images to appliance VM directly
5. Delete the VCH
6. Ensure that VCH folder is entirely gone after delete via govc

### Expected Outcome:
1. All steps should succeed
2. govc should return 1 when not finding images in step 6
