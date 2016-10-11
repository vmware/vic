# Virtual Container Host Status Reference #

The Web-based administration portal for virtual container hosts, VCH Admin, presents status information about a virtual container host.

If the vSphere environment in which you are deploying a virtual container host does not meet the requirements, the deployment does not succeed. However, a successfully deployed virtual container host can stop functioning if the vSphere environment changes after the deployment. If environment changes adversely affect the virtual container host, the status of the affected component turns yellow.

## Virtual Container Host (VCH) ##

**Error message**

`xx`.

**Cause** 

xx

**Solution** 

xx

## Network Connectivity ##

**Error message**

`xx`.

**Cause** 

xx

**Solution** 

xx

## Firewall ##

**Error message**

`Firewall must permit 2377/tcp outbound to use VIC`.

**Cause**

The firewall on the ESXi host on which the virtual container host is running no longer allows outbound connections on port 2377.

- The firewall was switched off when the virtual container host was deployed. The firewall has been switched on since the deployment of the virtual container host.
- A firewall ruleset was applied to the ESXi host to allow outbound connections on port 2377. The ESXi host has been rebooted since the deployment of the virtual container host. Firewall rulesets are not retained when an ESXi host reboots.

**Solution**

For information about how to reconfigure the firewall on ESXi hosts, see [VCH Deployment Fails with Firewall Validation Error](../vic_installation/ts_firewall_error.html) in *vSphere Integrated Containers Engine Installation*.

## License ##

**Error message**

`License does not meet minimum requirements to use VIC`.

**Cause** 

The license for the ESXi host or for one or more of the hosts in a vCenter Server cluster on which the virtual container host is deployed has been removed, downgraded, or has expired since the deployment of the virtual container host.

**Solution** 

- If the virtual container host is running on an ESXi host that is not managed by vCenter Server, replace the ESXi host license with a valid vSphere Enterprise license.
- If the virtual container host is running on a standalone ESXi host in vCenter Server, replace the ESXi host license with a valid vSphere Enterprise Plus license.
- If the virtual container host is running in a vCenter Server cluster, check that all of the hosts in the cluster have a valid vSphere Enterprise Plus license, and replace any licenses that have been removed, downgraded, or have expired.