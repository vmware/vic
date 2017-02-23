# Use Different User Accounts for VCH Deployment and Operation #

A virtual container host (VCH) appliance requires the appropriate permissions in vSphere to perform various tasks during VCH operation. 

During deployment of a VCH, `vic-machine` uses the vSphere account that you specify in either of the `vic-machine create --user` or `--target` options for all deployment operations. Deployment of a VCH requires a user account with vSphere Administrator privileges. 

Day-to-day operation of a VCH requires fewer permissions than  deployment. You can configure a VCH so that it uses different user accounts for deployment and for operation by using the `vic-machine create --ops-user` and `--ops-password` options when you deploy the VCH. By specifying `--ops-user`, you can limit the post-deployment permissions of the VCH to only those vSphere operations that it needs. If you do not specify `--ops-user`, the VCH runs with full vSphere Administrator privileges.

After deployment, a VCH must have permission to perform the following operations:

- Create, modify, and delete VMs within its resource pool
- Reconfigure the endpoint VM
- Validate host firewall configuration and system licenses

When you deploy a VCH, a user account that you specify in `--ops-user` must have the correct privileges to allow the VCH to perform these operations. vSphere Integrated Containers Engine does not currently create the required roles, so to assign privileges to the `--ops-user` user account, you must manually create user roles in vSphere before you deploy the VCH. You assign  privileges to those roles, and assign the roles to the user account to use in `--ops-user`. 

**Procedure**

1. In the vSphere Web Client, create a user group, for example `VIC Ops Users`, and add the appropriate user accounts to the user group.

   The best practice when assigning roles in vSphere is to assign the roles to user groups and then to add users to those groups, rather than assigning roles to the users directly.

2. Go to **Administration** > **Roles** and create one role for each type of inventory object that VCHs need to access.

   It is possible to create a single role, but by creating multiple roles you keep the privileges of the VCH as granular as possible.

   <table>
<thead>
<tr>
<th><strong>Role to Create</strong></th>
<th><strong>Required Permissions</strong></th>
</tr>
</thead>
<tbody>
<tr>
<td><code>VCH - vcenter</code></td>
<td>Datastore &gt; Configure datastore</td>
</tr>
<tr>
<td><code>VCH - datacenter</code></td>
<td>Datastore &gt; Configure datastore<br>Datastore &gt; Low level file operations</td>
</tr>
<tr>
<td><code>VCH - datastore</code></td>
<td>Datastore &gt; AllocateSpace<br>Datastore &gt; Browse datastore <br>Datastore &gt; Configure datastore<br>Datastore &gt; Remove file<br>Datastore &gt; Low level file operations<br>Host &gt; Configuration &gt; System management</td>
</tr>
<tr>
<td><code>VCH - network</code></td>
<td>Network &gt; Assign network</td>
</tr>
<tr>
<td><code>VCH - endpoint</code></td>
<td>dvPort group &gt; Modify<br>dvPort group &gt; Policy operation<br>dvPort group &gt; Scope operation<br>vApp &gt; Add virtual machine<br>VirtualMachine &gt; Configuration &gt; Add new disk<br>VirtualMachine &gt; Configuration &gt; Advanced<br>VirtualMachine &gt; Configuration &gt; Add or Remove Device<br>VirtualMachine &gt; Configuration &gt; Remove disk<br>VirtualMachine &gt; Guest operations &gt; Guest operation program execution<br>VirtualMachine &gt; Interaction &gt; Device connection<br>VirtualMachine &gt; Interaction &gt; Power off<br>VirtualMachine &gt; Interaction &gt; Power on<br>VirtualMachine &gt; Inventory &gt; Create new<br>VirtualMachine &gt; Inventory &gt; Remove<br>VirtualMachine &gt; Inventory &gt; Register<br>VirtualMachine &gt; Inventory &gt; Unregister</td>
</tr></tbody></table>

   For information about how to create vSphere roles, see [vSphere Permissions and User Management Tasks](https://pubs.vmware.com/vsphere-65/topic/com.vmware.vsphere.security.doc/GUID-5372F580-5C23-4E9C-8A4E-EF1B4DD9033E.html) in the vSphere documentation. 

3. Go to **Networking**, create a network folder, and place the distributed virtual switches that the VCHs will use for the bridge network and any container networks into that folder.

   The parent object of distributed virtual switches that the VCH uses  as the bridge network and container networks must be set to `Read-Only`, with **Propagate to Children** enabled. By placing distributed virtual switches in a network folder, you avoid setting an entire datacenter to `Read-Only`. This restriction only applies to the bridge network and container networks. When you specify the `vic-machine create --bridge-network` and `--container-network` options, include the full inventory path to the networks in the following format:<pre><i>datacenter</i>/network/<i>network_folder</i>/<i>port_group_name</i></pre>

2. (Optional) Go to **Hosts and Clusters** and create a resource pool in which to deploy VCHs.

   By creating a resource pool for VCHs, you can set the correct permissons on just that resource pool rather than on an entire host or cluster. You specify this resource pool in the `vic-machine create --compute-resource` option when you deploy the VCH. For a more granular application of privileges, you can also apply the permissions directly to VCH vApps after deployment, rather than to a resource pool.

5. In each of the **Hosts and Clusters**, **Storage**, and **Networking** views, select inventory objects and assign the user group and the appropriate role to each one.

   <table>
<thead>
<tr>
<th>Inventory Object</th>
<th>Role to Assign</th>
<th>Propagate</th>
</tr>
</thead>
<tbody>
<tr>
<td>Top-level vCenter Server instance</td>
<td><code>VCH - vcenter</code></td>
<td>No</td>
</tr>
<tr>
<td>Datacenters</td>
<td><code>VCH - datacenter</code></td>
<td>No</td>
</tr>
<tr>
<td>Clusters. All datastores in the cluster inherit permissions from the cluster.</td>
<td><code>VCH - datastore</code></td>
<td>Yes</td>
</tr>
<tr>
<td>Standalone VMware vSAN datastores</td>
<td><code>VCH - datastore</code></td>
<td>No</td>
</tr>
<tr>
<td>Standalone datastores</td>
<td><code>VCH - datastore</code></td>
<td>No</td>
</tr>
<tr>
<td>Network folders</td>
<td><code>Read-only</code></td>
<td>Yes</td>
</tr>
<tr>
<td>Port groups</td>
<td><code>VCH - network</code></td>
<td>No</td>
</tr>
<tr>
<td>Resource pools for VCHs</td>
<td><code>VCH - endpoint</code></td>
<td>Yes</td>
</tr>
<tr>
<td>VCH vApps, for a very granular application of privileges</td>
<td><code>VCH - endpoint</code></td>
<td>Yes</td>
</tr></tbody></table>

   For information about how to assign permissions to objects in the vSphere Inventory, see [Add a Permission to an Inventory Object](https://pubs.vmware.com/vsphere-65/topic/com.vmware.vsphere.security.doc/GUID-A0F6D9C2-CE72-4FE5-BAFC-309CFC519EC8.html) in the vSphere documentation.

**What to do next**

Use `vic-machine create --ops-user=<user_account>` to deploy VCHs that operate with restricted privileges. Ensure that the various vSphere inventory objects that you specify as arguments have the user group with the appropriate role. For an example of a `vic-machine` command with the `--ops-user` option, see the section *Specify Different User Accounts for VCH Deployment and Operation* [Advanced Examples of Deploying a VCH](vch_installer_examples.md#ops-user).