# VCH Deployment Fails with Firewall Validation Error #
When you use `vic-machine create` to deploy a virtual container host (VCH), deployment fails because firewall port 2377 is not open on the target ESXi host or hosts.

## Problem ##
Deployment fails with a firewall error during the validation phase: 

<pre>Firewall must permit dst 2377/tcp outbound to the VCH management interface</pre>

## Cause ##

ESXi hosts communicate with the VCHs through port 2377 via Serial Over LAN. For deployment of a VCH to succeed, port 2377 must be open for outgoing connections on all ESXi hosts before you run `vic-machine create`. Opening port 2377 for outgoing connections on ESXi hosts opens port 2377 for inbound connections on the VCHs.

## Solution ##

Set a firewall ruleset on the ESXi host or hosts. In test environments, you can disable the firewall on the hosts.

### Set a Firewall Ruleset Manually 

In production environments, if you are deploying to a standalone ESXi host, set a firewall ruleset on that ESXi host. If you are deploying to a cluster, set the firewall ruleset on all of the ESXi hosts in the cluster.

**IMPORTANT**: Firewall rulesets that you set manually are not persistent. If you reboot the ESXi hosts, any firewall rules that you set are lost. You must recreate firewall rules each time you reboot a host.

1. Use SSH to log in to each ESXi host as `root` user. 
2. Follow the instructions in [VMware KB 2008226]( http://kb.vmware.com/selfservice/microsites/search.do?language=en_US&cmd=displayKC&externalId=2008226) to add the following rule after the last rule in the file ```/etc/vmware/firewall/service.xml```.
<pre>
&lt;service id='<i>id_number</i>'&gt;
  &lt;id&gt;vicoutgoing&lt;/id&gt;
  &lt;rule id='0000'&gt;
    &lt;direction&gt;outbound&lt;/direction&gt;
    &lt;protocol&gt;tcp&lt;/protocol&gt;
    &lt;port type='dst'&gt;2377&lt;/port&gt;
  &lt;/rule&gt;
  &lt;enabled&gt;true&lt;/enabled&gt;
  &lt;required&gt;true&lt;/required&gt;
&lt;/service&gt;
</pre>

  
  In this example, *id_number* is the number of the preceding ruleset in ```service.xml```, incremented by 1.

### Disable the Firewall

In test environments, you can disable the firewalls on the ESXi hosts instead of opening port 2377. 
 
1. Use SSH to log in to each ESXi host as `root` user. 
2. Run the following command: 

  ```$ esxcli network firewall set --enabled false``` 