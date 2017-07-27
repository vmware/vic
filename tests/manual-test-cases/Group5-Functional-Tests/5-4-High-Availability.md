Test 5-4 - High Availability
=======

# Purpose:
To verify the VIC appliance works in when the vCenter appliance is using high availability

# References:
[1 - VMware vCenter Server Availability Guide](http://www.vmware.com/files/pdf/techpaper/vmware-vcenter-server-availability-guide.pdf)
[2 - Managing HA Clusters](https://pubs.vmware.com/vsphere-50/index.jsp#com.vmware.wssdk.pg.doc_50/PG_Ch13_Resources.15.9.html)

# Environment:
This test requires access to VMWare Nimbus cluster for dynamic ESXi and vCenter creation

# Test Steps:
1. Deploy a new vCenter with 3 ESXi hosts in a cluster
2. Enable HA on the cluster:  
```govc cluster.change -drs-enabled -ha-enabled /ha-datacenter/host/cls```
3. Deploy a new VCH Appliance to the cluster  
4. Run a variety of docker commands on the VCH appliance

5. Power off the ESXi host that the VCH is currently running on

6. Run a variety of docker commands on the VCH appliance

# Expected Outcome:
The VCH appliance should deploy without error and each of the docker commands executed against it should return without error

# Possible Problems:
None
