# vSphere Integrated Containers Command Line Installer Options

The command line installer for vSphere Integrated Containers provides options. The options allow you to customize the installation to match your vSphere environment.

| **Option** | **Description** |
| -- | -- | -- |
| ```ceip``` | **Mandatory**. Enable or disable participation in the VMware Customer Experience Improvement Program.  Expected values are ```enable``` or ```disable```. |
| ```cert``` |2:3 |
| ```cidr``` | 2:4 |
| ```cluster``` |The path to the cluster on which to install vSphere Integrated Containers. Specify the path by using the vSphere ```govc``` CLI format. For example,  ```/<my_datacenter>/host/<my_cluster>/```. Specify this option if you are installing vSphere Integrated Containers on a vCenter Server instance that manages more than one cluster. Omit this option if vCenter Server only manages one cluster. |
| ```containerNetwork``` | 2:6 |
| ```datacenter``` | 2:7 |
| ```datastore``` | 2:8 |
| ```dns``` | The address of a DNS server, to allow you to assign static IP addresses by using the ```ip``` option. You can specify the ```dns``` option multiple times, to identify multiple DNS servers. If not specified, the installer assigns IP addresses by using DHCP. |
| ```dockerOpts``` | 2:10 |
| ```externalNetwork``` | 2:11 |
| ```force``` | 2:12 |
| ```host``` | The address of the ESXi host on which to install vSphere Integrated Containers. Specify this option if you are installing vSphere Integrated Containers on a vCenter Server instance that manages more than one host and the hosts are not included in a cluster. Omit this option if vCenter Server only manages one ESXi host.|
| ```ip``` | A static IPv4 address for the vSphere Integrated Containers appliance. Requires you to specify the ```dns``` option. If not specified, the installer assigns IP addresses by using DHCP. |
| ```key``` | 2:15 |
| ```logfile``` | 2:16 |
| ```memoryMB``` | The amount of RAM to assign to the virtual container host. Specify this option if you intend to run large numbers of containers in this virtual container host. If not specified, the installer assigns 2048 MB of RAM to the virtual container host.  |
| ```name``` | A name for the vSphere Integrated Containers appliance. If not specified, the installer sets the name to ```docker-appliance```. |
| ```numCPUs``` | The number of CPUs to assign to the virtual container host. Specify this option if you intend to run large numbers of containers in this virtual container host. If not specified, the installer creates the appliance with 2 CPUs.  |
| ```os``` | 2:20 |
| ```passwd``` | The password for the vCenter Server user account that you are using to install vSphere Integrated Containers, or the password for the ESXi host. If not specified, the installer prompts you to enter the password during installation. |
| ```pool``` | The path to a resource pool in which to place the vSphere Integrated Containers appliance. Specify the path by using the vSphere ```govc``` CLI format. For example,  ```/<my_datacenter>/host/<my_cluster>/Resources/<my_resource_pool>```.  |
| ```target``` | **Mandatory**. The address of the ESXi host or vCenter Server instance on which you are installing vSphere Integrated containers. If an ESXi host is managed by a vCenter Server instance, you must provide the address of vCenter Server rather than of the host. To facilitate IP address changes in your infrastructure, provide a fully qualified domain name (FQDN) whenever possible, rather than an IP address.|
| ```timeout``` | The timeout period for uploading images to the ESXi host and powering on virtual machines. Specify a value in the format ```XmYs``` if the default timeout of 3m0s is insufficient. |
| ```uninstall``` | Uninstalls vSphere Integrated Containers. Removes all virtual machines from the vCenter Server inventory and deletes all files from storage. <ul><li>Requires the <code>target</code> option.<li>If you installed vSphere Integrated Containers on a vCenter Server instance, you must specify the <code>user</code> option.<li>If you do not specify the <code>passwd</code> option, the installer prompts you to enter the password.</li><li>Specify the ```yes``` option to answer yes to all questions during the uninstallation process.|
| ```user``` | The username for the ESXi host or vCenter Server instance on which you are installing vSphere Integrated containers. <ul><li>If you are installing vSphere Integrated Containers directly on an ESXi host and you do not specify this option, the installer uses the <code>root</code> account for installation.</li><li> This option is <strong>mandatory</strong> if you are installing vSphere Integrated Containers on a vCenter Server instance.</li></ul>
| ```verify``` | 2:27 |
| ```yes``` | Automatically answer yes to all questions during uninstallation. |
