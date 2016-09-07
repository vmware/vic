Test 13-1 vMotion VCH Appliance
=======

#Purpose:
To verify the VCH appliance continues to function properly after being vMotioned to a new host

#References:
[1- vMotion A Powered On Virtual Machine](http://pubs.vmware.com/vsphere-4-esx-vcenter/index.jsp?topic=/com.vmware.vsphere.dcadmin.doc_41/vsp_dc_admin_guide/migrating_virtual_machines/t_migrate_a_powered-on_virtual_machine_with_vmotion.html)

#Environment:
This test requires that a vCenter server is running and available

#Test Steps:
1. Install a new VCH appliance onto one of the hosts within the vCenter server
2. Power down the VCH appliance
3. vMotion the VCH appliance to a new host
4. Power on the VCH appliance and run a variety of docker commands
5. Install a new VCH appliance onto one of the hosts within the vCenter server
6. While the VCH appliance is powered on, vMotion the VCH appliance to a new host
7. Run a variety of docker commands on the VCH appliance after it has moved
8. Install a new VCH appliance onto on the hosts within the vCenter server
9. Create several containers on the new VCH appliance that are in the following states: created but not started, started and running, started and stopped
10. vMotion the VCH appliance to a new host
11. Complete the life cycle of the containers created in Step 9.

#Expected Outcome:
In each scenario, the VCH appliance should continue to work as expected after being vMotioned and all docker commands should return without error

#Possible Problems:
None