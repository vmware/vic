Test 5-14 - Remove Container OOB
=======

# Purpose:
To verify that VIC works properly when a container is removed OOB in VC

# References:
[1 - VMware vCenter Server Availability Guide](http://www.vmware.com/files/pdf/techpaper/vmware-vcenter-server-availability-guide.pdf)

# Environment:
This test requires access to VMware Nimbus cluster for dynamic ESXi and vCenter creation

# Docker run a container and verify cannot destroy the VM Out of Band - test steps:
1. Deploy a new vCenter with a simple cluster
2. Install the VIC appliance into one of the clusters
3. Run docker run -itd busybox /bin/top
4. Attempt to remove the container vm OOB

# Expected Outcome:
Step 4 should result in and error and a message stating that OOB deletion is disabled on VIC

# Possible Problems:
None

# Docker run a container destroy VM Out Of Band verify container is cleaned up - test steps:
1. Deploy a new vCenter with a simple cluster
2. Install the VIC appliance into one of the clusters
3. Deploy a busybox anchor container to prevent losing images/scrachdisk when the main containerVM is removed
4. Run docker run -itd busybox /bin/top as the main containerVM
5. Enable removal of the container vm OOB
6. Remove the container vm OOB
7. Verify the associated container has been removed
8. Remove anchor container

# Expected Outcome:
Step 7 should result in not finding the container in the list

# Possible Problems:
None

