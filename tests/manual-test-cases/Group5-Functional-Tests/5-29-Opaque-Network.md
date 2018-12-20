Test 5-29 - NSXT Logical Switch
=======

# Purpose:
To verify the VIC appliance works with NSXT Logical Switches

# References:
[1 - VMware NSX Transformer](https://docs.vmware.com/en/VMware-NSX-T-Data-Center/index.html)

# Environment:
This test requires to use VDNET 3.0 for vCenter, ESXi and NSX-T creation and customization
1. Deploy a new topology which includes vCenter and NSXT in Nimbus using VDNET 3.0
```A new vCenter```
```A new NSXT Manager, NSXT Controller and NSXT Edge```
```Three ESXi Hosts and put two of them into os-computer-cluster-1 cluster and one of them into management cluster```
```Create an overlay transport zone in NSXT```

# Case 1 Test Steps:
1. Create two overlay logical switches for bridge and container network
2. Install VIC appliance and deploy a VCH
3. Run a variety of docker commands on the VCH appliance
4. Destroy the created VIC appliance and VCH

# Case 2 Test Steps:
1. Create 2 new logical switches for bridge and 1 logical switch for container network
2. Install VIC appliance and deploy 2 VCHs. Each VCH uses different logical switch as bridge network.
3. Deploy a selenium grid hub and 8 selenium node in each VCH
4. Verify each of the selenium node is deployed properly and connect to the hub
5. Delete VCH


# Expected Outcome:
The VCH appliance should deploy without error and each of the docker commands executed against it should return without error

# Possible Problems:
* Your testbed is deployed failed due to VDNET 3.0 is not stable every time
* Hit failure in deleting VCH which has selenium grid related container
