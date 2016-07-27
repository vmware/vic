Test 6-4 - Verify vic-machine create basic function
=======

#Purpose:
Verify vic-machine create basic connection variables, certificates, timeout, and all arguments after appliance-iso

#References:
* vic-machine-linux create -h

#Environment:
This test requires that a vSphere server is running and available



#Test Cases: - Connection

#Test Steps
1. Create with vSphere URL in --target parameter, without --user and --password

#Expected Outcome:
* Command should fail for no user password available


#Test Steps
1. Create with vSphere URL in --target parameter and --user provided, but without --password

#Expected Outcome:
* Command should promote interactive password input


#Test Steps
1. Create with vSphere URL with wrong DC appended as <ip>/WrongDC

#Expected Outcome:
* Command should fail for WrongDC is not found


#Test Steps
1. Prepare test env with multiple DC exists
2. Create with vSphere URL with correct DC appended as <ip>/DC1

#Expected Outcome:
* Verify deployed successfully
* Verify VCH is in correct DC through govc


#Test Steps
## Create VCH - defaults
1. Issue the following command:
```
vic-machine create --name=<VCH_NAME> --target=<TEST_URL> \
    --user=<TEST_USERNAME> --image-datastore=<TEST_DATASTORE> --password=<TEST_PASSWORD> \
    --bridge-network=<NETWORK> --compute-resource=<TEST_RESOURCE>
```
2. Run regression tests

#Expected Outcome:
* Deployment succeed
* Regression test pass


#Test Steps
## Create VCH - target URL
1. Create with vSphere URL and user password encoded in the same --target parameter
2. Issue the following command:
```
vic-machine-linux create --name=<VCH_NAME> --target="<TEST_USERNAME>:<TEST_PASSWORD>@<TEST_URL>" \
    --image-datastore=<TEST_DATASTORE>
```
3. Run regression tests

#Expected Outcome:
* Deployment succeed
* Regression test pass


#Test Steps
## Create VCH - full params
1. Issue the following command:
```
vic-machine-linux create --name=<VCH_NAME> --target=<TEST_URL> --user=<TEST_USERNAME> --image-datastore=<TEST_DATASTORE> --appliance-iso=bin/appliance.iso --bootstrap-iso=bin/bootstrap.iso --password=<TEST_PASSWORD> --force=true --bridge-network=network --compute-resource=<TEST_RESOURCE> --timeout <TEST_TIMEOUT> --volume-store=<TEST_DATASTORE>/test:default
```
2. Run regression tests

#Expected Outcome:
* Deployment succeed
* Regression test pass


#Test Steps
1. Provide long name to create VCH, e.g. 100 characters

#Expected Outcome:
* Command failed for name is too long


#Test Steps
1. Create with same name with existing VCH

#Expected Outcome:
* Command failed for VCH is found


#Test Steps
1. Create with existing VM name, but not VCH

#Expected Outcome:
* Command failed for VM name is used


#Test Steps
1. Create with existing rp or vApp name, not VCH

#Expected Outcome:
* Command failed for name is used



#Test Cases: - Image files

#Test Steps
1. Provide wrong iso files

#Expected Outcome:
* Command failed for no iso files found



#Test Cases: - Certificate

#Test Steps
1. Create with --no-tls is true

#Expected Outcome:
* Command should succeed
* Docker is accessed through non-tls connection
* Regression test pass


#Test Steps
1. Generate key/cert files
2. Specify key, cert files during creation

#Expected Outcome:
* Command should succeed, 
* Docker is accessed through tls connection, and the key/cert previous generated can be used to access docker
* Regression test pass


#Test Steps
1. Specify key, cert files with mal-format files

#Expected Outcome:
* Command fail for wrong key/cert file



#Test Cases: - Timeout

#Test Steps
1. Specify short timeout to 2s

#Expected Outcome:
* Command fail for timeout error #1557



#Test Cases: - Appliance size

#Test Steps
1. Specify appliance size to 4cpu, 4096MB

#Expected Outcome:
* Deployed successfully
* Appliance VM size is set correctly in vsphere
* Regression test pass


#Test Steps
1. Specify appliance size to 1cpu, 256MB

#Expected Outcome:
* Deployment failed for no enought resource
* Should have user-friendly error message


#Test Steps
1. --use-rp=true

#Expected Outcome:
* Deployed successfully
* VCH is created under resource pool against VC
* Regression test pass



#Test Cases: - extension

#Test Steps
1. Verify VC extension is registered correctly against VC
2. Pre register one extension before VCH deployment, check deployment failed with user-friendly error message

#Test Cases: - VCH size


#Test Steps
1. Specify VCH CPU size to reservation: 4, limit: 8, shares: wrong

#Expected Outcome:
* Deployment failed for wrong shares format


#Test Steps
1. Specify VCH CPU size to reservation: 4, limit: 2, shares: normal

#Expected Outcome:
* Deployment failed for user-friendly error message


#Test Steps
1. Specify VCH CPU size to reservation: 4, limit: 8, shares: high

#Expected Outcome:
* Deployed successfully
* Check vApp or rp resource settings is correct through govc
* Integration test passed


#Test Steps
1. Specify VCH Memory size to reservation: 4096, limit: 8192, shares: wrong

#Expected Outcome:
* Deployment failed for wrong shares format


#Test Steps
1. Specify VCH Memory size to reservation: 4096, limit: 2048, shares: normal

#Expected Outcome:
* Deployment failed for user-friendly error message


#Test Steps
1. Specify VCH Memory size to reservation: 4096, limit: 8192, shares: high

#Expected Outcome:
* Deployed successfully
* Check vApp or rp resource settings is correct through govc
* Integration test passed


#Test Steps
1. Specify VCH Memory size to reservation: 256, limit: 256, shares: high

#Expected Outcome:
* Deployment failed with user-friendly error message


#Test Steps
1. Specify VCH Memory size to reservation: 200, limit: 200, shares: high

#Expected Outcome:
* Deployment failed with user-friendly error message
