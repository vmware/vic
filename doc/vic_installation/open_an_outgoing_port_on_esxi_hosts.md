# Open an Outgoing Port on ESXi Hosts

ESXi hosts communicate with the virtual container hosts via port 2377. For installation to succeed, port 2377 must be open on all hosts before you install vSphere Integrated Containers. 

To open port 2377 on the ESXi hosts, log into each host via SSH and add the following rule after the last rule in the file ```/etc/vmware/firewall/service.xml```.
  
<pre>&lt;!--Port for VIC communication --&gt;
   &lt;service id='*id_number*'&gt;
   &lt;id&gt;vicoutgoing&lt;/id&gt;
   &lt;rule id='*id_number*'&gt;
      &lt;direction&gt;outbound&lt;/direction&gt;
      &lt;protocol&gt;tcp&lt;/protocol&gt;
      &lt;port type='dst'&gt;&lt;/port type&gt;
      &lt;port&gt;2377&lt;/port&gt;
   &lt;/rule&gt;
   &lt;enabled&gt;true&lt;/enabled&gt;
   &lt;required&gt;true&lt;/required&gt;
   &lt;/service&gt;
</pre>
  
In this example, *id_number* is the number of the preceding rule in ```service.xml```, incremented by 1. 

For detailed instructions about how to add a rule to open a port on an ESXi host, see [VMware KB 2008226]( http://kb.vmware.com/selfservice/microsites/search.do?language=en_US&cmd=displayKC&externalId=2008226). 
