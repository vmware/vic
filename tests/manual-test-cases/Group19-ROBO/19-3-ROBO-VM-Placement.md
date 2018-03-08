Test 19-3 - ROBO - VM Placement
=======

# Purpose:
To verify that the VM placement feature specified works as expected in a vSphere ROBO Advanced environment without DRS.
The current placement strategy is to avoid bad host selection, instead of selecting the "best" possible host.

# References:
1. [vSphere Remote Office and Branch Office](http://www.vmware.com/products/vsphere/remote-office-branch-office.html)
2. [VM Placement without DRS](https://github.com/vmware/vic/issues/7282)

# Environment:
This test requires access to VMware Nimbus cluster for dynamic ESXi and vCenter creation. This test should be executed in the following topologies and should have vSAN enabled.
* 1 vCenter host with 3 clusters, where 1 cluster has 1 ESXi host and the other 2 clusters have 3 ESXi hosts each
* 2 vCenter hosts connected with ELM, where each vCenter host has a cluster/host/datacenter topology that emulates a customer environment (exact topology TBD)

In addition, this test should be run in multi-ESX-host and single-ESX-host cluster topologies.

See https://confluence.eng.vmware.com/display/CNA/VIC+ROBO for more details.

# Test Steps:
1. Deploy a ROBO Advanced vCenter testbed for both environments above
2. Run a few dummy VMs on the host(s) in the vCenter cluster meant for VCH installation. This ensures that the hosts have varying resource utilization.
3. Install a VCH on a particular cluster on vCenter - see note in [Environment](#environment)
4. Deploy containers that will consume resources predictably (e.g. the `progrium/stress` image)
5. Measure cluster metrics and gather resource consumption
6. Create and run regular containers such as `busybox`
7. Create and run enough containers to consume all available cluster resources
8. Attempt to create and run more containers
9. Delete some containers
10. Create and run a few containers
11. Delete the VCH

# Expected Outcome:
* Step 1 should succeed
* Step 2 should succeed
* Step 3 should succeed and the VCH should be placed on a host that satisfies the license and feature checks. The VCH's host should also meet the criteria defined in point 2 of [References](#references)
* Steps 4-5 should succeed and containers should be placed on ESX hosts in the cluster according to the criteria defined in point 2 of [References](#references)
* Step 6 should succeed and containers should be placed on ESX hosts in the cluster that have available resources according to the criteria defined in point 2 of [References](#references). In the multi-host cluster environment, the cluster resource utilization level should be as expected given containerVM sizes, cluster capacity and placement logic.
* Step 7 should succeed
* Step 8 should fail since the available resources are exhausted
* Steps 9-11 should succeed

# Possible Problems:
None
