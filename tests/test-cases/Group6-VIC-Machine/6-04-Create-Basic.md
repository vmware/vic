Test 6-04 - Verify vic-machine create basic function
=======

# Purpose:
Verify vic-machine create basic connection variables, certificates, timeout, and all arguments after appliance-iso

# References:
* vic-machine-linux create -h

# Environment:
This test requires that a vSphere server is running and available



Image size
=======

## Create VCH - custom base disk
1. Issue the following command:
```
vic-machine-linux create --name=${vch-name} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --bridge-network=%{BRIDGE_NETWORK} --external-network=%{EXTERNAL_NETWORK} --image-store=%{TEST_DATASTORE} --password=%{TEST_PASSWORD} --base-image-size=6GB ${vicmachinetls}
```

### Expected Outcome
* VCH is deployed successfully
* Container has correct disk size
* Regression tests pass



Connection
=======

## Create VCH - URL without user and password
1. Create with vSphere URL in --target parameter, without --user and --password

### Expected Outcome
* Command should fail for no user password available


## Create VCH - URL without password
1. Create with vSphere URL in --target parameter and --user provided, but without --password

### Expected Outcome
* Command should promote interactive password input


## Create VCH - target URL
1. Create with vSphere URL and user password encoded in the same --target parameter
```
vic-machine-linux create --name=<VCH_NAME> --target="<TEST_USERNAME>:<TEST_PASSWORD>@<TEST_URL>" \
    --image-store=<TEST_DATASTORE>
```
2. Run regression tests

### Expected Outcome
* Deployment succeed
* Regression test pass


## Create VCH - specified datacenter
1. Prepare test env with multiple DC exists
2. Create with vSphere URL with correct DC appended as <ip>/DC1

### Expected Outcome
* Verify deployed successfully
* Verify VCH is in correct DC through govc



vic-machine create Parameters
=======

## Create VCH - defaults
1. Issue the following command:
```
vic-machine create --name=<VCH_NAME> --target=<TEST_URL> \
    --user=<TEST_USERNAME> --image-store=<TEST_DATASTORE> --password=<TEST_PASSWORD> \
    --bridge-network=<NETWORK> --compute-resource=<TEST_RESOURCE>
```
2. Run regression tests

### Expected Outcome
* Deployment succeed
* Regression test pass


## Create VCH - full params
1. Issue the following command:
```
vic-machine-linux create --name=<VCH_NAME> --target=<TEST_URL> \
    --user=<TEST_USERNAME> --image-store=<TEST_DATASTORE> \
    --appliance-iso=bin/appliance.iso --bootstrap-iso=bin/bootstrap.iso \
    --password=<TEST_PASSWORD> --force=true --bridge-network=network \
    --compute-resource=<TEST_RESOURCE> --timeout <TEST_TIMEOUT> \
    --volume-store=<TEST_DATASTORE>/test:default
```
2. Run regression tests

### Expected Outcome
* Deployment succeed
* Regression test pass


## Create VCH - custom image store directory
1. Issue the following command:
```
vic-machine-linux create --name=${vch-name} --target=%{TEST_URL} \
    --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} \
    --image-store %{TEST_DATASTORE}/vic-machine-test-images \
    --appliance-iso=bin/appliance.iso --bootstrap-iso=bin/bootstrap.iso \
    --password=%{TEST_PASSWORD} --force=true --bridge-network=%{BRIDGE_NETWORK} \
    --external-network=%{EXTERNAL_NETWORK} --compute-resource=%{TEST_RESOURCE} \
    --timeout %{TEST_TIMEOUT} ${vicmachinetls}
```
2. Run regression tests

### Expected Outcome
* Deployment succeeds
* Regression tests pass


## Create VCH - long VCH name
1. Provide long name to create VCH, e.g. 100 characters

### Expected Outcome
* Command failed for name is too long


## Create VCH - Existing VCH name
1. Create with same name with existing VCH

### Expected Outcome
* Command failed for VCH is found


## Create VCH - Existing VM name
1. Create with existing VM name
2. Run regression tests

### Expected Outcome
* Deployment succeeds
* Regression tests pass


## Create VCH - Existing RP on ESX
1. Create resource pool on ESX
2. Create VCH with the same (already existing) name

### Expected Outcome
* Deployment succeeds
* Regression tests pass


## Create VCH - Existing vApp on vCenter
1. Create vApp on vCenter
2. Create VCH with the same (already existing) name

### Expected Outcome
* Command fails



Image files
=======

## Create VCH - wrong ISOs
1. Provide wrong iso files

### Expected Outcome
* Command failed for no iso files found



TLS Auth
=======

## Create VCH - defaults with --no-tls
1. Issue the following command:
```
vic-machine-linux create --name=${vch-name} --target=%{TEST_URL} \
    --user=%{TEST_USERNAME} --image-store=%{TEST_DATASTORE} --appliance-iso=bin/appliance.iso \
    --bootstrap-iso=bin/bootstrap.iso --password=%{TEST_PASSWORD} --no-tls --force=true \
    --bridge-network=%{BRIDGE_NETWORK} --external-network=%{EXTERNAL_NETWORK} \
    --compute-resource=%{TEST_RESOURCE} --timeout %{TEST_TIMEOUT} \
    --volume-store=%{TEST_DATASTORE}/test:${vol}
```
2. Run regression tests

### Expected Outcome
* Deployment succeeds
* Regression tests pass


## Create VCH - force accept target thumbprint
1. Issue the following command:
```
vic-machine-linux create --name=${vch-name} --target="%{TEST_USERNAME}:%{TEST_PASSWORD}@%{TEST_URL}" \
    --force --image-store=%{TEST_DATASTORE} --bridge-network=%{BRIDGE_NETWORK} \
    --external-network=%{EXTERNAL_NETWORK} ${vicmachinetls}
```
2. Run regression tests

### Expected Outcome
* Deployment succeed
* Regression test pass


## Create VCH - Specified keys
1. Generate key/cert files
2. Specify key, cert files during creation

### Expected Outcome
* Command should succeed,
* Docker is accessed through tls connection, and the key/cert previous generated can be used to access docker
* Regression test pass


## Create VCH - Invalid keys
1. Specify key, cert files with mal-format files

### Expected Outcome
* Command fail for wrong key/cert file



Timeout
=======

## Basic timeout
1. Specify short timeout to 2s

### Expected Outcome
* Command fail for timeout error #1557



Appliance size
=======

## Basic VCH resource config
1. Specify appliance size to 4cpu, 4096MB

### Expected Outcome
* Deployed successfully
* Appliance VM size is set correctly in vsphere
* Regression test pass


## Invalid VCH resource config
1. Specify appliance size to 1cpu, 256MB

### Expected Outcome
* Deployment failed for no enought resource
* Should have user-friendly error message


## Use resource pool
1. --use-rp=true

### Expected Outcome
* Deployed successfully
* VCH is created under resource pool against VC
* Regression test pass


## CPU reservation shares invalid
1. Specify VCH CPU size to reservation: 4, limit: 8, shares: wrong

### Expected Outcome
* Deployment failed for wrong shares format


## CPU reservation invalid
1. Specify VCH CPU size to reservation: 4, limit: 2, shares: normal

### Expected Outcome
* Deployment failed for user-friendly error message


## CPU reservation valid
1. Specify VCH CPU size to reservation: 4, limit: 8, shares: high

### Expected Outcome
* Deployed successfully
* Check vApp or rp resource settings is correct through govc
* Integration test passed


## Memory reservation shares invalid
1. Specify VCH Memory size to reservation: 4096, limit: 8192, shares: wrong

### Expected Outcome
* Deployment failed for wrong shares format


## Memory reservation invalid 1
1. Specify VCH Memory size to reservation: 4096, limit: 2048, shares: normal

### Expected Outcome
* Deployment failed for user-friendly error message


## Memory reservation invalid 2
1. Specify VCH Memory size to reservation: 256, limit: 256, shares: high

### Expected Outcome
* Deployment failed with user-friendly error message


## Memory reservation invalid 3
1. Specify VCH Memory size to reservation: 200, limit: 200, shares: high

### Expected Outcome
* Deployment failed with user-friendly error message


## Memory reservation valid
1. Specify VCH Memory size to reservation: 4096, limit: 8192, shares: high

### Expected Outcome
* Deployed successfully
* Check vApp or rp resource settings is correct through govc
* Integration test passed
