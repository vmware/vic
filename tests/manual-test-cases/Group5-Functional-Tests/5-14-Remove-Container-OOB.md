Test 5-14 - Remove Container OOB
=======

# Purpose:
To verify that VIC works properly when a container is removed OOB in VC

# References:
[1 - VMware vCenter Server Availability Guide](http://www.vmware.com/files/pdf/techpaper/vmware-vcenter-server-availability-guide.pdf)

# Environment:
This test requires access to VMware Nimbus cluster for dynamic ESXi and vCenter creation

# Test Steps:
1. Deploy a new vCenter with a simple cluster
2. Install the VIC appliance into one of the clusters
3. Run docker run -itd busybox /bin/top
4. Remove the container vm OOB
5. Run docker run -itd busybox /bin/top

# Expected Outcome:
Step 5 should result in success and the container created should behave as expected

# Possible Problems:
None
