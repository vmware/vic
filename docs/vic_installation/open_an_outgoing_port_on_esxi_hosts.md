# Open an Outgoing Port on ESXi Hosts

ESXi hosts communicate with the virtual container hosts via port 2377. For installation to succeed, port 2377 must be open before you install vSphere Integrated Containers. 

To open port 2377 on the ESXi hosts, log into each host via SSH and add the following rule after the last rule in the file ```/etc/vmware/firewall/service.xml```.
  
     
    <!--Port for VIC communication -->
    <service id='*id_number*'>
     <id>vicoutgoing</id>
      <rule id='*id_number*'>
       <direction>outbound</direction>
       <protocol>tcp</protocol>
       <port type='dst'></port type>
       <port>2377</port>
     </rule>
      <enabled>true</enabled>
     <required>true</required>
    </service>

  
In this example, *id_number* is the number of the preceding rule in ```service.xml```, incremented by 1. 

For detailed instructions about how to add a rule to open a port on an ESXi host, see [VMware KB 2008226]( http://kb.vmware.com/selfservice/microsites/search.do?language=en_US&cmd=displayKC&externalId=2008226). 
