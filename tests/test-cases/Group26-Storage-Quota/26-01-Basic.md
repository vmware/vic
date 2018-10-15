Suite 26-01 - Basic
===================

# Purpose:
To verify storage quota functionality

# Environment:
This suite requires a vCenter Server environment where VCHs can be deployed and container VMs created.


### 1. Create a VCH with storage quota and docker info shows the storage quota and storage usage

#### Test Steps:
1. Create a VCH with a storage quota of 15GB
2. Verify that docker info shows the storage quota and storage usage is 0

#### Expected Outcome:
* The VCH with storage quota is created.


### 2. Create a container and docker info shows the storage usage has changed

#### Test Steps:
1. Create a busybox container
2. Verify that docker info shows the storage usage is larger than 9GB(the default container vm will create a 2GB swap vmdk and an about maximum 8GB delta vmdk)

#### Expected Outcome:
* The container is created and storage usage has changed.


### 3. Pull a debian image and docker info shows the storage usage has changed

#### Test Steps:
1. Pull a debian image.
2. Verify that docker info show the storage usage has increased(the new image layers contribute to the storage usage)

#### Expected Outcome:
* The image is downloaded and storage usage has increased.


### 4. Create second container and get storage quota exceeding failure

#### Test Steps:
1. Verify that the container creation fails with storage quota exceeding error message.
2. Verify that docker info show the storage usage is the same

#### Expected Outcome:
* Contiainer creation fails because storage quota exceeds.


### 5. Configure VCH with a larger storage quota of 35GB

#### Test Steps:
1. Configure VCH with a larger storage quota of 35GB.
2. Verify that docker info show the storage quota is set to 35GB and storage usage does not change.

#### Expected Outcome:
* VCH configure succeeds and storage quota is updated.


### 6. Create second container successfully

#### Test Steps:
1. Create second busybox container.
2. Verify that docker info show the storage usage is larger than 18GB

#### Expected Outcome:
* The second container creation succeeds.


### 7. Remove a container successfully with storage usage changes

#### Test Steps:
1. Delete a busybox container.
2. Verify that docker info show the storage usage decreases

#### Expected Outcome:
* The container deletion decreases storage usage.


### 8. Create a busybox container with memory of 4GB successfully

#### Test Steps:
1. Create a busybox container with memory of 4GB
2. Verify that docker info show the storage usage is larger than 20GB(swap vmdk is equal to memory size)

#### Expected Outcome:
* The container creation succeeds with correct storage usage.


### 9. Create a debian container and commit to an image successfully

#### Test Steps:
1. Create a debian container
2. Verify that the storage usage is larger than 29GB
3. Commit it to an image
4. Verify that the storage usage increases

#### Expected Outcome:
* The container commits to an image and storage usage increases.


### 10. Delete an image successfully with storage usage decreased

#### Test Steps:
1. Delete the debian container committed image
2. Verify that the storage uasge decreases

#### Expected Outcome:
* The image deletion makes storage usage decrease.


### 11. Create a busybox continainer afer unsetting storage quota

#### Test Steps:
1. Verify that the container creation fails with storage quota exceeding error message
2. Configure VCH storage quota to 0(unlimited)
3. Verify that the container creation succeeds

#### Expected Outcome:
* The container creation succeeds after unsetting storage quota
