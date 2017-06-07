# Installation Fails with Firewall Ruleset Validation Error #
When you use either the command line installer or the OVA deployment, installation fails due to missing a firewall ruleset.
## Problem ##
Installation fails with an error during the validation phase: 

<pre>
0 enabled firewall rulesets match outbound tcp dst 2377, 
[...] 
Validating supplied configuration failed. Exiting...
</pre>

## Cause ##
ESXi hosts communicate with the virtual container hosts via port 2377. For installation to succeed, port 2377 must be open for outbound connections on all all ESXi hosts

## Solution ##
- Open port 2377 for outbound connections on all ESXi hosts. For information about how to open port 2377 on ESXi hosts, see [Open an Outgoing Port on ESXi Hosts](open_an_outgoing_port_on_esxi_hosts.md).
- In test environments, you can disable the firewall on ESXi hosts. 