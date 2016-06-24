Test 5-5 - Heterogeneous ESXi
=======

#Purpose:
To verify the VIC appliance works when the vCenter appliance is using multiple different ESXi versions

#References:
[1 - VMware vCenter Server Availability Guide](http://www.vmware.com/files/pdf/techpaper/vmware-vcenter-server-availability-guide.pdf)

#Environment:
This test requires access to VMWare Nimbus cluster for dynamic ESXi and vCenter creation

#Test Steps:
1. Deploy a new vCenter in Nimbus
2. Deploy two different ESXi hosts with build numbers(6.0.0u2 and 5.5u3):
```3620759``` and ```3029944```
3. Add each host to the vCenter
4. Deploy a VCH appliance to each of the ESXi hosts in the vCenter
5. Run a variety of docker commands on each of the VCH appliances.

#Expected Outcome:
The VCH appliance should deploy without error in both scenarios and each of the docker commands executed against it should return without error

#Possible Problems:
This case is not supported with version 1.0 of VIC, as VIC only supports 6.0.0u2 on release.