Test 6-4 - Verify vic-machine create validation function
=======

#Purpose:
Verify vic-machine create validation functions, this does not include validation for network, datastore, and compute resources

#References:
* vic-machine-linux create -h

#Environment:
This test requires that a vSphere server is running and available


#Test Cases: - suggest resources
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


#Test Cases: - networking

# Test Steps
## vDS contains all hosts in cluster
1. Prepare vCenter environment with a vDS that is connected to all hosts in the cluster
2. Issue the following command:
```
vic-machine create --name=<VCH_NAME> --target=<TEST_URL> \
    --user=<TEST_USERNAME> --image-datastore=<TEST_DATASTORE> --password=<TEST_PASSWORD> \
    --bridge-network=<NETWORK> --compute-resource=<TEST_RESOURCE>
```
3. Run regression tests

#Expected Outcome:
* Output contains message indicating vDS configuration OK
* Deployment succeeds
* Regression tests pass

## vDS does not contain all hosts in cluster
1. Prepare vCenter environment with a vDS that is not connected to all hosts in the cluster
2. Issue the following command:
```
vic-machine create --name=<VCH_NAME> --target=<TEST_URL> \
    --user=<TEST_USERNAME> --image-datastore=<TEST_DATASTORE> --password=<TEST_PASSWORD> \
    --bridge-network=<NETWORK> --compute-resource=<TEST_RESOURCE>
```

#Expected Outcome:
* Output contains message indicating vDS configuration is incorrect with user-friendly error message
* Deployment fails
