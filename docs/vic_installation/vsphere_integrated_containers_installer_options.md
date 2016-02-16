# vSphere Integrated Containers Command Line Installer Options

The command line installer for vSphere Integrated Containers provides options. The options allow you to customize the installation to match your vSphere environment.

<table>
  <thead>
    <tr>
      <th><strong>Option</strong></th>
      <th><strong>Description</strong></th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><span class="style1">ceip</span></td>
      <td><strong>Mandatory</strong>. Enable or disable participation in the VMware Customer Experience Improvement Program. Expected values are&nbsp;enable&nbsp;ordisable.</td>
    </tr>
    <tr>
      <td><span class="style1">cert</span></td>
      <td>2:3</td>
    </tr>
    <tr>
      <td><span class="style1">cidr</span></td>
      <td>2:4</td>
    </tr>
    <tr>
      <td><span class="style1">cluster</span></td>
      <td>The path to the cluster on which to install vSphere Integrated Containers. Specify the path by using the vSphere&nbsp;<span class="style1">govc</span>&nbsp;CLI format. For example,&nbsp;<span class="style1">/&lt;<em>my_datacenter</em>&gt;/host/&lt;<em>my_cluster</em>&gt;/</span>. Specify this option if you are installing vSphere Integrated Containers on a vCenter Server instance that manages more than one cluster. Omit this option if vCenter Server only manages one cluster.</td>
    </tr>
    <tr>
      <td><span class="style1">containerNetwork</span></td>
      <td>2:6</td>
    </tr>
    <tr>
      <td><span class="style1">datacenter</span></td>
      <td>2:7</td>
    </tr>
    <tr>
      <td><span class="style1">datastore</span></td>
      <td>2:8</td>
    </tr>
    <tr>
      <td><span class="style1">dns</span></td>
      <td>The address of a DNS server, to allow you to assign static IP addresses by using the&nbsp;ip&nbsp;option. You can specify the&nbsp;<span class="style1">dns</span>&nbsp;option multiple times, to identify multiple DNS servers. If not specified, the installer assigns IP addresses by using DHCP.</td>
    </tr>
    <tr>
      <td><span class="style1">dockerOpts</span></td>
      <td>2:10</td>
    </tr>
    <tr>
      <td><span class="style1">externalNetwork</span></td>
      <td>2:11</td>
    </tr>
    <tr>
      <td><span class="style1">force</span></td>
      <td>2:12</td>
    </tr>
    <tr>
      <td><span class="style1">host</span></td>
      <td>The address of the ESXi host on which to install vSphere Integrated Containers. Specify this option if you are installing vSphere Integrated Containers on a vCenter Server instance that manages more than one host and the hosts are not included in a cluster. Omit this option if vCenter Server only manages one ESXi host.</td>
    </tr>
    <tr>
      <td><span class="style1">ip</span></td>
      <td>A static IPv4 address for the vSphere Integrated Containers appliance. Requires you to specify the&nbsp;<span class="style1">dns</span>&nbsp;option. If not specified, the installer assigns IP addresses by using DHCP.</td>
    </tr>
    <tr>
      <td><span class="style1">key</span></td>
      <td>2:15</td>
    </tr>
    <tr>
      <td><span class="style1">logfile</span></td>
      <td>2:16</td>
    </tr>
    <tr>
      <td><span class="style1">memoryMB</span></td>
      <td>The amount of RAM to assign to the virtual container host. Specify this option if you intend to run large numbers of containers in this virtual container host. If not specified, the installer assigns 2048 MB of RAM to the virtual container host.</td>
    </tr>
    <tr>
      <td><span class="style1">name</span></td>
      <td>A name for the vSphere Integrated Containers appliance. If not specified, the installer sets the name to&nbsp;docker-appliance.</td>
    </tr>
    <tr>
      <td><span class="style1">numCPUs</span></td>
      <td>The number of CPUs to assign to the virtual container host. Specify this option if you intend to run large numbers of containers in this virtual container host. If not specified, the installer creates the appliance with 2 CPUs.</td>
    </tr>
    <tr>
      <td><span class="style1">os</span></td>
      <td>2:20</td>
    </tr>
    <tr>
      <td><span class="style1">passwd</span></td>
      <td>The password for the vCenter Server user account that you are using to install vSphere Integrated Containers, or the password for the ESXi host. If not specified, the installer prompts you to enter the password during installation.</td>
    </tr>
    <tr>
      <td><span class="style1">pool</span></td>
      <td>The path to a resource pool in which to place the vSphere Integrated Containers appliance. Specify the path by using the vSphere&nbsp;<span class="style1">govc</span>&nbsp;CLI format. For example, <span class="style1">/<em>&lt;my_datacenter&gt;</em>/host/<em>&lt;my_cluster&gt;</em>/Resources/<em>&lt;my_resource_pool&gt;</em></span>.</td>
    </tr>
    <tr>
      <td><span class="style1">target</span></td>
      <td><strong>Mandatory</strong>. The address of the ESXi host or vCenter Server instance on which you are installing vSphere Integrated containers. If an ESXi host is managed by a vCenter Server instance, you must provide the address of vCenter Server rather than of the host. To facilitate IP address changes in your infrastructure, provide a fully qualified domain name (FQDN) whenever possible, rather than an IP address.</td>
    </tr>
    <tr>
      <td><span class="style1">timeout</span></td>
      <td>The timeout period for uploading images to the ESXi host and powering on virtual machines. Specify a value in the format&nbsp;<span class="style1">XmYs</span>&nbsp;if the default timeout of 3m0s is insufficient.</td>
    </tr>
    <tr>
      <td><span class="style1">uninstall</span></td>
      <td>Uninstalls vSphere Integrated Containers. Removes all virtual machines from the vCenter Server inventory and deletes all files from storage.
          <ul>
            <li>Requires the&nbsp;<span class="style1">target</span>&nbsp;option.</li>
            <li>If you installed vSphere Integrated Containers on a vCenter Server instance, you must specify the&nbsp;<span class="style1">user</span>&nbsp;option.</li>
            <li>If you do not specify the&nbsp;<span class="style1">passwd</span>&nbsp;option, the installer prompts you to enter the password.</li>
            <li>Specify the&nbsp;<span class="style1">yes</span>&nbsp;option to answer yes to all questions during the uninstallation process.</li>
          </ul></td>
    </tr>
    <tr>
      <td><span class="style1">user</span></td>
      <td>The username for the ESXi host or vCenter Server instance on which you are installing vSphere Integrated containers.
          <ul>
            <li>If you are installing vSphere Integrated Containers directly on an ESXi host and you do not specify this option, the installer uses theroot&nbsp;account for installation.</li>
            <li>This option is&nbsp;<strong>mandatory</strong>&nbsp;if you are installing vSphere Integrated Containers on a vCenter Server instance.</li>
          </ul></td>
    </tr>
    <tr>
      <td><span class="style1">verify</span></td>
      <td>2:27</td>
    </tr>
    <tr>
      <td><span class="style1">yes</span></td>
      <td>Automatically answer yes to all questions during uninstallation.</td>
    </tr>
  </tbody>
</table>
