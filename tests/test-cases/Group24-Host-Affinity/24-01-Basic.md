Suite 24-01 - Basic
===================

# Purpose:
To verify basic VM-Host Affinity functionality

# References:
1. [The design document](../../../doc/design/host-affinity.md)

# Environment:
This suite requires a vCenter Server environment where VCHs can be deployed and container VMs created.

Note that because these basic tests do not test the behavior of DRS in the presence of rules, but just the management of
VM groups, these tests do not require an environment where DRS is enabled.


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


### 3. Removing containers cleans up the VM group

#### Test Steps:
1. Verify that no DRS VM Group exists by the expected name.
2. Create a VCH.
3. Create a variety of containers.
4. Verify that a DRS VM Group was created and that the endpoint VM and containers were added to it.
5. Delete the containers.
6. Verify that the DRS VM Group still exists, but does not include the removed containers.

#### Expected Outcome:
* Containers are removed from the DRS VM Group when they are deleted.
