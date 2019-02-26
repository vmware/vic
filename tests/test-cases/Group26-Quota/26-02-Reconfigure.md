Suite 26-02 - Reconfigure
===================

# Purpose:
To verify adding storage quota settings during configure functionality

# Environment:
This suite requires a vCenter Server environment where VCHs can be deployed and container VMs created.


### 1. Create a VCH without storage quota and docker info shows the storage usage

#### Test Steps:
1. Create a VCH without storage quota
2. Verify that docker info output has no storage quota information, but has storage usage information
3. Create a busybox container
4. Verify that docker info shows the storage usage is larger than 9GB

#### Expected Outcome:
* The VCH and container creation succeeds without setting storage quota and keep storage usage track(this is required when VCH is reconfigured with storage quota).

### 2. Create a container and docker info shows the storage usage has changed

#### Test Steps:
1. Create a busybox container
2. Verify that docker info shows the storage usage is larger than 9GB(the default container vm will create a 2GB swap vmdk and an about maximum 8GB delta vmdk)

#### Expected Outcome:
* The container is created and storage usage has changed.

### 3. Configure VCH with a larger storage quota of 15GB

#### Test Steps:
1. Configure VCH with a larger storage quota of 15GB.
2. Verify that docker info show the storage quota is set to 35GB and storage usage does not change.

#### Expected Outcome:
* VCH configure succeeds and storage quota is updated.

### 4. Create second container and get storage quota exceeding failure

#### Test Steps:
1. Verify that the container creation fails with storage quota exceeding error message.
2. Verify that docker info show the storage usage is the same

#### Expected Outcome:
* Contiainer creation fails because storage quota exceeds.
