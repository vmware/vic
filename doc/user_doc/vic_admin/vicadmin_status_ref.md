# Virtual Container Host Status Reference #

The Web-based administration portal for virtual container hosts, VCH Admin, presents status information about a virtual container host.

If the vSphere environment in which you are deploying a virtual container host does not meet the requirements, the deployment does not succeed. However, a successfully deployed virtual container host can stop functioning if the vSphere environment changes after the deployment. If environment changes adversely affect the virtual container host, the status of the affected component changes from green to yellow.

## Virtual Container Host (VCH) ##

VCH Admin checks the status of the processes that the virtual container host runs:

- The port layer server, that presents an API of low-level container primitive operations, and implements those container operations via the vSphere APIs.
- VCH Admin server, that runs the VCH Admin portal. 
- The vSphere Integrated Containers Engine initialization service and watchdog service for the other components. 
- The Docker engine server, that exposes the Docker API and semantics, translating those composite operations into port layer primitives.

### Error ###

- The **Virtual Container Host** status is yellow.
- The **Virtual Container Host** status is yellow and an error message informs you that the virtual container host cannot connect to vSphere.

#### Cause ####

- One or more of the virtual container host processes is not running correctly, or the virtual container host is unable to connect to vSphere.
- The management network connection is down and the virtual container host endpoint VM cannot connect to vSphere.

#### Solution ####

1. (Optional) If you see the error that the virtual container host is unable to connect to vSphere, check the virtual container host management network.
2. In the VCH Admin portal for the virtual container host, click the link for the **VCH Admin Server** log.
2. Search the log for references to the different virtual container host processes.

  The different processes are identified in the log by the following names:

  - `port-layer-server`
  - `vicadmin`
  - `vic-init`
  - `docker-engine-server`

3. Identify the process or processes that are not running correctly and attempt to remediate the issues as required.

## Registry and Internet Connectivity ##

VCH Admin checks connectivity on the public network by attempting to connect from the virtual container host to docker.io and google.com. VCH Admin only checks the public network connection. It does not check other networks, for example the bridge, management, client, or container networks.

### Error ###

The **Registry and Internet Connectivity** status is yellow.

#### Cause ####

The public network connection is down.

#### Solution ####

Check the **VCH Admin Server** log for references to network issues. Use the vSphere Web Client to remediate the management network issues as required.

## Firewall ##

VCH Admin checks that the firewall is correctly configured on the ESXi host or the ESXi hosts in the cluster on which the virtual container host is running.

### Error ###

- The **Firewall** status is unavailable.
- The **Firewall** status is yellow and shows the error `Firewall must permit 2377/tcp outbound to use VIC`.

#### Cause ####

- The management network connection is down and the virtual container host endpoint VM cannot connect to vSphere. 
- The firewall on the ESXi host on which the virtual container host is running no longer allows outbound connections on port 2377.

  - The firewall was switched off when the virtual container host was deployed. The firewall has been switched on since the deployment of the virtual container host.
  - A firewall ruleset was applied to the ESXi host to allow outbound connections on port 2377. The ESXi host has been rebooted since the deployment of the virtual container host. Firewall rulesets are not retained when an ESXi host reboots.

#### Solution ####

- If the **Firewall** status is unavailable: 
  - Check the **VCH Admin Server** log for references to network issues. 
  - Use the vSphere Web Client to remediate the management network issues as required.
- If you see the error about port 2377, reconfigure the firewall on the ESXi host or hosts to allow  outbound connections on port 2377. For information about how to reconfigure the firewall on ESXi hosts, see [VCH Deployment Fails with Firewall Validation Error](../vic_installation/ts_firewall_error.html) in *vSphere Integrated Containers Engine Installation*.


## License ##

VCH Admin checks that the ESXi hosts on which you deploy virtual container hosts have the appropriate licenses.

### Error ###

- The **License** status is yellow and shows the error `License does not meet minimum requirements to use VIC`.
- The **License** status is unavailable.

#### Cause ####

- The license for the ESXi host or for one or more of the hosts in a vCenter Server cluster on which the virtual container host is deployed has been removed, downgraded, or has expired since the deployment of the virtual container host.
- The management network is down, or the virtual container host endpoint VM is unable to connect to vSphere.

#### Solution ####

- If the license does not meet the requirements:
  - If the virtual container host is running on an ESXi host that is not managed by vCenter Server, replace the ESXi host license with a valid vSphere Enterprise license.
  - If the virtual container host is running on a standalone ESXi host in vCenter Server, replace the ESXi host license with a valid vSphere Enterprise Plus license.
  - If the virtual container host is running in a vCenter Server cluster, check that all of the hosts in the cluster have a valid vSphere Enterprise Plus license, and replace any licenses that have been removed, downgraded, or have expired.
- If the **License** status is unavailable: 
  - Check the **VCH Admin Server** log for references to network issues. 
  - Use the vSphere Web Client to remediate the management network issues as required.