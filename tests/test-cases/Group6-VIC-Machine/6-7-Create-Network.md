Test 6-4 - Verify vic-machine create network function
=======

#Purpose:
Verify vic-machine create external, management, bridge network and container networks functions

#References:
* vic-machine-linux create -h

#Environment:
This test requires that a vSphere server is running and available

#Test Cases: - external network

#Test Steps
1. Create without external network provided
2. Verify "VM Network" is connected in VCH VM
3. Integration test passed 

#Test Steps
1. Create with wrong network name provided for external network
2. Verify create failed for network is not found
3. Create with distribute virtual switch as external network name
4. Verify create failed for network type is wrong

#Test Steps
1. Create with network name no DHCP availabile for external network
2. Verify VCH created but without ip address
3. Verify VCH can be deleted without anything wrong through vic-machine delete

#Test Steps
1. Create with DPG as external network in VC and correct switch in ESXi
2. Verify create passed
3. Verify integration test passed

#Test Cases: - management network

#Test Steps
1. Create without management network provided, but external network correctly set
2. Verify warning message set for management network and client network sharing the same network
3. No multiple attachement in VM network to same vSphere virtual switch (or DPG)
4. Integration test passed 

#Test Steps
1. Create with wrong network name provided for management network
2. Verify create failed for network is not found
3. Create with distribute virtual switch as management network name
4. Verify create failed for network type is wrong

#Test Steps
1. Create with network unreachable for vSphere or VC as management network
2. Verify VCH created but VC or vSphere is unreachable
3. Make sure vic-machine failed with user-friendly error message

#Test Steps
1. Create with DPG as management network in VC and correct switch in ESXi
2. Verify create passed
3. Verify integration test passed

#Test Cases: - bridge network

#Test Steps
1. Create without bridge network provided in VC
2. Create failed for bridge network should be specified in VC

#Test Steps
1. Create without bridge network provided in ESXi
2. Integration test pass

#Test Steps
1. Create with wrong network name provided for bridge network
2. Verify create failed for network is not found
3. Create with distribute virtual switch as bridge network name
4. Verify create failed for network type is wrong

#Test Steps
1. Create with standard network in VC as bridge network
2. vic-machine failed for DPG is required for bridge network

#Test Steps
1. Create with DPG as management network in VC and correct switch in ESXi
2. Verify create passed
3. Verify integration test passed

#Test Steps
1. Create with same network for bridige and external network 
2. Verify create failed for same network with external network
3. Same case with management network
4. Same case with container network

#Test Steps
1. Create with bridge network correctly set
2. Set bridge network ip range with wrong format
3. Verify create failed with user-friendly error message

#Test Steps
1. Create with bridge network correctly set
2. Set bridge network ip range correctly
3. Verify create passed
4. Regression test passed
5. docker create container, with ip address correctly set in the above ip range

#Test Cases:

#Test Steps
1. Create with invalid container network: <WrongNet>:alias
2. Verify create failed with WrongNet is not found

#Test Steps
1. Create with container network: <stand switch network name>:alias in VC
2. Verify create failed with standard network is not supported

#Test Steps
1. Create with container network: <dpg name>:net1 in VC or <stand switch network name>:net1 in ESXi
2. Verify create passed
3. Regression test passed
4. Verify docker network ls command to show net1 network

#Test Steps
1. Create with container network: <dpg name> in VC or <stand switch network name> in ESXi
2. Verify create passed
3. Regression test passed
4. Verify docker network ls command to show the <vsphere network name> network

#Test Steps
1. Create with two container network map to same alias
2. Verify create failed with two different vsphere network map to same docker network

#Test Steps
1. Create with two container network map to same alias
2. Verify create failed with two different vsphere network map to same docker network

#Test Steps
1. Create with container network mapping
2. Set container network gateway as <dpg name>:1.1.1.1/24
3. Set container network gateway as <dpg name>:192.168.1.0/24
4. Set container network gateway as <wrong name>:192.168.1.0/24
5. Verify create failed for wrong vsphere network name or gateway is not routable

#Test Steps
1. Create with container network mapping
2. Set container ip range as <wrong name>:192.168.2.1-192.168.2.100
3. Set container network gateway as <dpg name>:192.168.1.1/24, and container ip range as <dpg name>:192.168.2.1-192.168.2.100
4. Verify create failed for wrong vsphere network name or ip range is wrong

#Test Steps
1. Create with container network mapping
2. Set container DNS as <wrong name>:8.8.8.8
3. Set container DNS as <dpg name>:abcdefg
4. Verify create failed for wrong vsphere name or wrong dns format

#Test Steps
1. Create with container network mapping <dpg name>:net1
2. Set container network gateway as <dpg name>:192.168.1.1/24
3. Set container ip range as <dpg name>:192.168.1.2-192.168.1.100
4. Set container DNS as <dpg name>:<correct dns>
5. Verify create passed
6. Integration test passed
7. Docker network ls show net1
8. Docker container created with network attached with net1, got ip address inside of network range
9. Docker create another container, and link to previous one, can talk to the the first container successfully
