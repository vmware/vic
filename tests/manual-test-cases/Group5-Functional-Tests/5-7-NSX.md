Test 5-7 - NSX
=======

#Purpose:
To verify the VIC appliance works when the vCenter is using NSX networking

#References:
[1 - VMware NSX](http://www.vmware.com/products/nsx.html)

#Environment:
This test requires access to VMWare Nimbus cluster for dynamic ESXi and vCenter creation

#Test Steps:
1. Deploy a new vCenter in Nimbus with NSX configured  
2. Deploy VCH Appliance to the new vCenter
3. Run a variety of docker commands on the VCH appliance

#Expected Outcome:
The VCH appliance should deploy without error and each of the docker commands executed against it should return without error

#Possible Problems:
None
