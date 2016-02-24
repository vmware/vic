# vSphere Integrated Containers Command Line Installer Options

The command line installer for vSphere Integrated Containers provides options. The options allow you to customize the installation to match your vSphere environment.



<table> 

<thead>
    <tr>
      <th><strong>Option</strong></th>
      <th><strong>Description</strong></th>
      <th><strong>Example</strong></th>
    </tr>
	</thead>
	<tbody>
    <tr>
      <td><code>ceip</code></td>
      <td><strong>Mandatory</strong>. Enable or disable participation in the VMware Customer Experience Improvement Program.</td>
      <td><code>-ceip=enable</code>
      <br> 
	   <code>-ceip=disable</code></td>
    </tr>
    <tr>
      <td><code>cert</code></td>
      <td>The path to the X.509 certificate for the vCenter Server instance or ESXi host on which you are installing vSphere Integrated Containers. Set this option  if your vSphere environment uses SSL certificates that have been signed by a Certificate Authority. </td>
      <td>&nbsp;</td>
    </tr>
    <tr>
      <td><code>cidr</code></td>
      <td>2:4</td>
      <td>&nbsp;</td>
    </tr>
    <tr>
      <td><code>cluster</code></td>
      <td><p>The path to the cluster on which to install vSphere Integrated Containers. Specify this option if you are installing vSphere Integrated Containers in a datacenter that contains more than one cluster. Specify the path by using the vSphere&nbsp;<code>govc</code>&nbsp;CLI format, including the leading and trailing forward slashes. Omit this option if vCenter Server only manages one cluster.</p>
      <p><b>NOTE</b>: If your datacenter includes clusters and also includes standalone hosts that are not members of any of the clusters, and if you want to install vSphere Integrated Containers on one of the standalone hosts, you must specify the host address in the <code>-cluster</code> option. </p></td>
      <td><code>-cluster=/&lt;<em>my_datacenter</em>&gt; /host/&lt;<em>my_cluster</em>&gt;/</code></td>
    </tr>
    <tr>
      <td><code>containerNetwork</code></td>
      <td>2:6</td>
      <td>&nbsp;</td>
    </tr>
    <tr>
      <td><code>datacenter</code></td>
      <td>The name of the datacenter in which to install vSphere Integrated Containers. Specify this option if you are installing vSphere Integrated Containers on a vCenter Server instance that manages more than one datacenter. Omit this option if vCenter Server only manages one datacenter.</td>
      <td><code>-datacenter=&lt;<em>my_datacenter</em>&gt;</code></td>
    </tr>
    <tr>
      <td><code>datastore</code></td>
      <td>The name of the datastore in which to store the files of vSphere Integrated Containers appliance. vSphere Integrated Containers uses this datastore to store container images and the files of container virtual machines. Specify this option to install vSphere Integrated Containers on an ESXi host that contains more than one datastore. A vCenter Server cluster</li>
      </ul>
      Omit this option if the ESXi host or vCenter Server only manages one datastore.</td>
      <td><code>-datastore=&lt;<em>datastore_name</em>&gt;</code></td>
    </tr>
    <tr>
      <td><code>dns</code></td>
      <td>The address of a DNS server, to allow you to assign static IP addresses by using the&nbsp;ip&nbsp;option. You can specify the&nbsp;<code>dns</code>&nbsp;option multiple times, to identify multiple DNS servers. If not specified, the installer assigns IP addresses by using DHCP.</td>
      <td>&nbsp;</td>
    </tr>
    <tr>
      <td><code>dockerOpts</code></td>
      <td>2:10</td>
      <td>&nbsp;</td>
    </tr>
    <tr>
      <td><code>externalNetwork</code></td>
      <td>2:11</td>
      <td>&nbsp;</td>
    </tr>
    <tr>
      <td><code>force</code></td>
      <td>2:12</td>
      <td>&nbsp;</td>
    </tr>
    <tr>
      <td><code>host</code></td>
      <td><p>The address of the ESXi host on which to install vSphere Integrated Containers. Specify this option if you are installing vSphere Integrated Containers on a vCenter Server instance that manages more than one ESXi host and the hosts are not included in a cluster. Omit this option if vCenter Server only manages one ESXi host.</p>
      <p><b>NOTE</b>: If your datacenter includes clusters and also includes standalone hosts that are not members of any of the clusters, and if you want to install vSphere Integrated Containers on one of the standalone hosts, you must specify the host address in the <code>-cluster</code> option. </p></td>
      <td>&nbsp;</td>
    </tr>
    <tr>
      <td><code>ip</code></td>
      <td>A static IPv4 address for the vSphere Integrated Containers appliance. Requires you to specify the&nbsp;<code>dns</code>&nbsp;option. If not specified, the installer assigns IP addresses by using DHCP.</td>
      <td>&nbsp;</td>
    </tr>
    <tr>
      <td><code>key</code></td>
      <td>2:15</td>
      <td>&nbsp;</td>
    </tr>
    <tr>
      <td><code>logfile</code></td>
      <td>2:16</td>
      <td>&nbsp;</td>
    </tr>
    <tr>
      <td><code>memoryMB</code></td>
      <td>The amount of RAM to assign to the virtual container host. Specify this option if you intend to run large numbers of containers in this virtual container host. If not specified, the installer assigns 2048 MB of RAM to the virtual container host.</td>
      <td>&nbsp;</td>
    </tr>
    <tr>
      <td><code>name</code></td>
      <td>A name for the vSphere Integrated Containers appliance. If not specified, the installer sets the name to&nbsp;docker-appliance.</td>
      <td>&nbsp;</td>
    </tr>
    <tr>
      <td><code>numCPUs</code></td>
      <td>The number of CPUs to assign to the virtual container host. Specify this option if you intend to run large numbers of containers in this virtual container host. If not specified, the installer creates the appliance with 2 CPUs.</td>
      <td>&nbsp;</td>
    </tr>
    <tr>
      <td><code>os</code></td>
      <td>2:20</td>
      <td>&nbsp;</td>
    </tr>
    <tr>
      <td><code>passwd</code></td>
      <td>The password for the vCenter Server user account that you are using to install vSphere Integrated Containers, or the password for the ESXi host. If not specified, the installer prompts you to enter the password during installation.</td>
      <td>&nbsp;</td>
    </tr>
    <tr>
      <td><code>pool</code></td>
      <td>The path to a resource pool in which to place the vSphere Integrated Containers appliance. Specify the path by using the vSphere&nbsp;<code>govc</code>&nbsp;CLI format, including the leading and trailing forward slashes.</td>
      <td><code>-pool=/<em>&lt;my_datacenter&gt;</em> /host/<em>&lt;my_cluster&gt;</em> /Resources/<em>&lt;my_resource_pool&gt;</em>/</code></td>
    </tr>
    <tr>
      <td><code>target</code></td>
      <td><strong>Mandatory</strong>. The address of the ESXi host or vCenter Server instance on which you are installing vSphere Integrated containers. If an ESXi host is managed by a vCenter Server instance, you must provide the address of vCenter Server rather than of the host. To facilitate IP address changes in your infrastructure, provide a fully qualified domain name (FQDN) whenever possible, rather than an IP address.</td>
      <td>&nbsp;</td>
    </tr>
    <tr>
      <td><code>timeout</code></td>
      <td>The timeout period for uploading images to the ESXi host and powering on virtual machines. Specify a value in the format&nbsp;<code>XmYs</code>&nbsp;if the default timeout of 3m0s is insufficient.</td>
      <td>&nbsp;</td>
    </tr>
    <tr>
      <td><code>uninstall</code></td>
      <td>Uninstalls vSphere Integrated Containers. Removes the vSphere Integrated Containers vApp and virtual machines from the vCenter Server inventory. Deletes all of the vAPP and VM files from storage.
          <ul>
            <li>Requires the&nbsp;<code>-target</code>&nbsp;option.</li>
            <li>If you installed vSphere Integrated Containers on a vCenter Server instance, you must specify the&nbsp;<code>-user</code>&nbsp;option.</li>
            <li>If you do not specify the&nbsp;<code>-passwd</code>&nbsp;option, the installer prompts you to enter the password.</li>
            <li>Specify the&nbsp;<code>-yes</code>&nbsp;option to answer yes to all questions during the uninstallation process.</li>
          </ul>
          <p><strong>NOTE</strong>: If you do not specify the <code>-yes</code> option, the installer prompts you to confirm that you want to uninstall vSphere Integrated Containers. Enter the word <code>yes</code> to confirm. If you enter <code>y</code>, the uninstall operation quits. </p></td>
      <td>&nbsp;</td>
    </tr>
    <tr>
      <td><code>user</code></td>
      <td>The username for the ESXi host or vCenter Server instance on which you are installing vSphere Integrated containers.
          <ul>
            <li>If you are installing vSphere Integrated Containers directly on an ESXi host and you do not specify this option, the installer uses theroot&nbsp;account for installation.</li>
            <li>This option is&nbsp;<strong>mandatory</strong>&nbsp;if you are installing vSphere Integrated Containers on a vCenter Server instance.</li>
          </ul></td>
      <td>&nbsp;</td>
    </tr>
    <tr>
      <td><code>verify</code></td>
      <td>2:27</td>
      <td>&nbsp;</td>
    </tr>
    <tr>
      <td><code>yes</code></td>
      <td>Automatically answer yes to all questions during uninstallation.</td>
      <td>&nbsp;</td>
    </tr>
	</tbody>
</table>