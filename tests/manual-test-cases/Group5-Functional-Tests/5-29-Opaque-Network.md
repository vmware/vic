Test 5-29 - NSXT Logical Switch
=======

# Purpose:
To verify the VIC appliance works with NSXT Logical Switches

# References:
[1 - VMware NSX Transformer](https://docs.vmware.com/en/VMware-NSX-T-Data-Center/index.html)

# Environment:
This test requires to use VDNET 3.0 for vCenter, ESXi and NSX-T creation and customization

# Test Steps:
1. Deploy a new topology which includes vCenter and NSXT in Nimbus using VDNET 3.0
```A new vCenter```
```A new NSXT Manager, NSXT Controller and NSXT Edge```
```Three ESXi Hosts and put two of them into os-computer-cluster-1 cluster and one of them into management cluster```
```Create an overlay transport zone in NSXT```
2. Create two overlay logical switches for bridge and container network
3. Install VIC appliance and deploy a VCH
4. Run a variety of docker commands on the VCH appliance
5. Destroy testbeds

# Expected Outcome:
The VCH appliance should deploy without error and each of the docker commands executed against it should return without error

# Possible Problems:
* Your testbed is deployed failed due to VDNET 3.0 is not stable every time
