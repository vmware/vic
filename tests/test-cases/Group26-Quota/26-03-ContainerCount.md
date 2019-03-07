Suite 26-02 - Reconfigure
===================

# Purpose:
To verify container count functionality

# Environment:
This suite requires a vCenter Server environment where VCHs can be deployed and container VMs created.


### 1. Create a VCH with container count and inspect shows the container count

#### Test Steps:
1. Create a VCH with container count
2. Verify that vch inspect config output has container count information
3. Create a busybox container

#### Expected Outcome:
* The VCH and container creation succeeds within container count.

### 2. Create second container and get container count exceeding failure

#### Test Steps:
1. Create a busybox container
2. Verify that the container creation fails with container count exceeding error message

#### Expected Outcome:
* Container creation fails because container count exceeds.

### 3. Configure VCH with a larger container count of 2

#### Test Steps:
1. Configure VCH with a larger container count of 2
2. Verify that vch inspect config the container count is set to 2

#### Expected Outcome:
* VCH configure succeeds and container count is updated.

### 4. Create second container and succeed

#### Test Steps:
1. Verify that the second container creation succeeds
2. Verify that the third container creation fails because container count exceeds

#### Expected Outcome:
* The third container creation fails because container count exceeds.

### 5. Configure VCH with unlimited container count of 0

#### Test Steps:
1. Configure VCH with unlimited container count
2. Verify that the third container creation succeeds

#### Expected Outcome:
* The third container creation succeeds.
