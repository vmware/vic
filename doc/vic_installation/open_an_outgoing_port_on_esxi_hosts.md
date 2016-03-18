# Open an Outgoing Port on ESXi Hosts

ESXi hosts communicate with the virtual container hosts via port 2377. For installation to succeed, port 2377 must be open for outgoing connections on all all ESXi hosts before you install vSphere Integrated Containers. Opening port 2377 for outgoing connections on ESXi hosts opens port 2377 for inbound connections on the virtual container hosts.

Failure to open port 2377 on all ESXi hosts affects installation differently, depending on whether you install by using the command line installer or by using the OVA deployment.
 
- Installation by using the command line installer fails at the verification stage. 
- Installation of the management server by using the OVA deployment succeeds, but deployment of a virtual container fails.  

To open port 2377 on the ESXi hosts, log into each ESXi host via SSH and add the following rule after the last rule in the file ```/etc/vmware/firewall/service.xml```.

**ESXi 6.x:**

<pre>&lt;!--Port for VIC communication --&gt;
   &lt;service id='<i>id_number</i>'&gt;
   &lt;id&gt;vicoutgoing&lt;/id&gt;
   &lt;rule id='0000'&gt;
      &lt;direction&gt;outbound&lt;/direction&gt;
      &lt;protocol&gt;tcp&lt;/protocol&gt;
      &lt;port type='dst'&gt;2377&lt;/port type&gt;
   &lt;/rule&gt;
   &lt;enabled&gt;true&lt;/enabled&gt;
   &lt;required&gt;true&lt;/required&gt;
   &lt;/service&gt;
</pre>

**ESXi 5.x:**
  
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

  
In these examples, *id_number* is the number of the preceding rule in ```service.xml```, incremented by 1. 

**IMPORTANT**: If you reboot the ESXi hosts, any firewall rules that you set are lost. You must recreate firewall rules after a reboot.

For detailed instructions about how to add a rule to open a port on an ESXi host, see [VMware KB 2008226]( http://kb.vmware.com/selfservice/microsites/search.do?language=en_US&cmd=displayKC&externalId=2008226). 
