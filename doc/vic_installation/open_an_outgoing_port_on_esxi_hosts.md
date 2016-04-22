# Open an Outgoing Port on ESXi Hosts

ESXi hosts communicate with the virtual container hosts via port 2377. For installation to succeed, port 2377 must be open for outgoing connections on all all ESXi hosts before you install vSphere Integrated Containers. Opening port 2377 for outgoing connections on ESXi hosts opens port 2377 for inbound connections on the virtual container hosts.

Failure to open port 2377 on all ESXi hosts affects installation differently, depending on whether you install by using the command line installer or by using the OVA deployment.
 
- Installation by using the command line installer fails at the verification stage. 
- Installation of the management server by using the OVA deployment succeeds, but deployment of a virtual container host fails.  

You can open port 2377 on the ESXi hosts either by installing a vSphere Installation Bundle (VIB) or by manually setting a firewall rule.  

## Set a Firewall Rule by Installing a VIB ##
VMware provides a VIB that automatically opens port 2377 on an ESXi host. If you install this VIB, port 2377 remains open even if you reboot the ESXi host, which is not the case if you open the port by setting a firewall rule manually.

1. Download the VIB from https://vic-vmware.socialcast.com/attachments/3505732.

 You must be a member of the vSphere Integrated Containers Technical Preview SocialCast channel to access the download.
2. Copy the VIB into the /tmp folder on your ESXi host.<pre>scp vmware-vic.vib root@<i>esxi_host_address</i>:/tmp</pre>
3. Use SSH to log into the ESXi host as `root`.
4. Install the VIB by running the following command:<pre>esxcli software vib install -v /tmp/vmware-vic.vib -f</pre>
5. Verify that the installation succeeded by displaying the list of firewall rules.<pre>esxcli network firewall ruleset rule list</pre> 
 The firewall rule `vmware-vic` should appear in the list.
5. (Optional) If you are running vSphere Integrated Containers in a cluster, repeat the procedure on all of the ESXi hosts in the cluster.

## Set a Firewall Rule Manually ##

**IMPORTANT**: Firewall rules that you set manually are not persistent. If you reboot the ESXi hosts, any firewall rules that you set are lost. You must recreate firewall rules each time you reboot a host.

To set a firewall rule manyally, log into each ESXi host via SSH and add the following rule after the last rule in the file ```/etc/vmware/firewall/service.xml```.

<pre>&lt;!--Port for VIC communication --&gt;
   &lt;service id='<i>id_number</i>'&gt;
   &lt;id&gt;vicoutgoing&lt;/id&gt;
   &lt;rule id='0000'&gt;
      &lt;direction&gt;outbound&lt;/direction&gt;
      &lt;protocol&gt;tcp&lt;/protocol&gt;
      &lt;porttype&gt;dst&lt;/porttype&gt;
      &lt;port&gt;2377&lt;/port&gt;
   &lt;/rule&gt;
   &lt;enabled&gt;true&lt;/enabled&gt;
   &lt;required&gt;true&lt;/required&gt;
   &lt;/service&gt;
</pre>

  
In this example, *id_number* is the number of the preceding rule in ```service.xml```, incremented by 1. For detailed instructions about how to add a rule to open a port on an ESXi host, see [VMware KB 2008226]( http://kb.vmware.com/selfservice/microsites/search.do?language=en_US&cmd=displayKC&externalId=2008226).



