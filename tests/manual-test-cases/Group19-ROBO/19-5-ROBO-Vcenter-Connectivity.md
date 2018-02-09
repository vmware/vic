Test 19-5 - ROBO vCenter Connectivity
=======

# Purpose:
To verify that the applications deployed in containerVMs in a ROBO Advanced environment are functional when the ESXi(s) hosting the containerVMs are disconnected from the vSphere host. This test exercises the WAN connectivity and resiliency support for a ROBO environment that could represent a customer's cluster topology. 

# References:
1. [vSphere Remote Office and Branch Office](http://www.vmware.com/products/vsphere/remote-office-branch-office.html)

# Environment:
This test requires access to VMware Nimbus cluster for dynamic ESXi and vCenter creation. This test should be executed in the following topologies and should have vSAN enabled.
* 1 vCenter host with 3 clusters, where 1 cluster has 1 ESXi host and the other 2 clusters have 3 ESXi hosts each
* 2 vCenter hosts connected with ELM, where each vCenter host has a cluster/host/datacenter topology that emulates a customer environment (exact topology TBD)

See https://confluence.eng.vmware.com/display/CNA/VIC+ROBO for more details.

# Test Steps:
1. Deploy a ROBO Advanced vCenter testbed for both environments above
2. Install the VIC appliance on vCenter
3. Create and start some container services such as nginx, wordpress or a database
4. Run a multi-container application exercising network links with docker-compose
5. To simulate a WAN link outage, _abruptly_ disconnect each ESX host from vCenter (possibly by changing firewall rules)
6. Verify that the containers/services/applications started in Steps 3 and 4 are still alive and responding
7. Create/start a container
8. For each ESXi host that hosts containerVM(s), re-connect it to vCenter
9. Create/start a container
10. Delete the VCH

# Expected Outcome:
* Steps 1-6 should succeed
* Step 7 should fail since the vCenter host is disconnected from the VCH's host
* Steps 8-10 should succeed

# Possible Problems:
None
