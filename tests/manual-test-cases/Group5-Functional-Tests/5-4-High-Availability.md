Test 5-4 - High Availability
=======

#Purpose:
To verify the VIC appliance works in when the vCenter appliance is using high availability

#References:
[1 - VMware vCenter Server Availability Guide](http://www.vmware.com/files/pdf/techpaper/vmware-vcenter-server-availability-guide.pdf)

#Environment:
This test requires access to VMWare Nimbus cluster for dynamic ESXi and vCenter creation

#Test Steps:
1. Deploy a new vCenter in Nimbus
2. Add two different ESXi hosts to the new vCenter
3. Turn on high availability between the ESXi hosts
4. Deploy VCH Appliance to the new vCenter  
5. Run a variety of docker commands on the VCH appliance
6. Power off the ESXi host that the VCH is currently running on
7. Run a variety of docker commands on the VCH appliance

#Expected Outcome:
The VCH appliance should deploy without error and each of the docker commands executed against it should return without error

#Possible Problems:
None