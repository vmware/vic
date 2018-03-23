Suite 24-01 - Basic
===================

# Purpose:
To verify basic VM-Host Affinity functionality

# References:
1. [The design document](../../../doc/design/host-affinity.md)

# Environment:
This test requires a vCenter Server environment where VCHs can be deployed and container VMs created.


Positive Testing
----------------

### 1. Creating a VCH creates a VM group and container VMs get added to it

#### Test Steps:
1. Verify that no DRS VM Group exists by the expected name.
2. Create a VCH.
3. Verify that a DRS VM Group was created and that the endpoint VM was added to it.
4. Create a variety of containers.
5. Verify that the container VMs were added to the DRS VM Group.

#### Expected Outcome:
* The DRS VM Group is created and the VCH endpoint VM and all container VMs are added to it.


### 2. Deleting a VCH deletes its VM group

#### Test Steps:
1. Verify that no DRS VM Group exists by the expected name.
2. Create a VCH.
3. Verify that a DRS VM Group was created and that the endpoint VM was added to it.
4. Delete the VCH.
5. Verify that the DRS VM Group no longer exists.

#### Expected Outcome:
* The DRS VM Group is deleted when the VCH is deleted.
