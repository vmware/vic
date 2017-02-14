# Set Up the `--ops-user` User Account #

A virtual container host (VCH) appliance requires the appropriate permissions in vSphere to perform various tasks during VCH operation. 

During deployment of a VCH, `vic-machine` uses the vSphere account that you specify in either of the `vic-machine create --user` or `--target` options for all deployment operations. Deployment of a VCH requires a user account with vSphere Administrator privileges. 

Day-to-day operation of a VCH requires fewer permissions than  deployment. You can configure a VCH so that it uses different user accounts for deployment and for operation by using the `vic-machine create --ops-user` and `--ops-password` options when you deploy the VCH. By specifying `--ops-user`, you can limit the post-deployment permissions of the VCH to only those vSphere operations that it needs. If you do not specify `--ops-user`, the VCH runs with full vSphere Administrator privileges.

After deployment, a VCH must have permission to perform the following operations:

- Create, modify, and delete VMs within its resource pool
- Modify itself
- Validate host firewall configuration and system licenses

When you deploy a VCH, a user account that you specify in `--ops-user` must have the correct privileges to allow the VCH to perform these operations. To assign privileges to the `--ops-user` user account,  you must create user roles in vSphere before you deploy the VCH. You assign  privileges to those roles, and assign the roles to the user account to use in `--ops-user`. 

For information about how to create vSphere roles, see [vSphere Permissions and User Management Tasks](https://pubs.vmware.com/vsphere-65/topic/com.vmware.vsphere.security.doc/GUID-5372F580-5C23-4E9C-8A4E-EF1B4DD9033E.html) in the vSphere documentation.  For information about how to assign permissions to objects in the vSphere Inventory, see [Add a Permission to an Inventory Object](https://pubs.vmware.com/vsphere-65/topic/com.vmware.vsphere.security.doc/GUID-A0F6D9C2-CE72-4FE5-BAFC-309CFC519EC8.html).

**Procedure**

1. In the vSphere Web Client, create a user group, for example `VIC Ops Users`, and add the appropriate user accounts to the user group.

   The best practice when assigning roles in vSphere is to assign the roles to user groups and then to add users to those groups, rather than assigning roles to the users directly.

2. Go to **Administration** > **Roles** and create one role for each type of inventory object that VCHs need to access.

   vSphere Integrated Containers Engine does not currently create the required roles, so you must create them manually. [Required VCH Roles](#required_roles) below lists the roles that you must create and the permissions that each role requires.

3. Go to **Networking**, create a network folder, and place any distributed virtual switches that the VCHs will use into that folder.

   The parent object of any distributed virtual switches that the VCH uses must be set to `Read-Only`, with **Propagate to Children** enabled. By placing distributed virtual switches in a network folder, you avoid setting an entire datacenter to `Read-Only`.

2. Go to **Hosts and Clusters** and create a resource pool in which to deploy VCHs.

   By creating a resource pool for VCHs, you can set the correct permissons on just that resource pool rather than on an entire host or cluster.

5. In each of the **Hosts and Clusters**, **Storage**, and **Networking** views, select inventory objects and assign the user group and the appropriate role to each one.

   [Inventory Objects to Which to Assign Permissions](#inventory_objects) below lists the inventory objects that VCHs access, and the role that you must assign to each one.

**What to do next**

Use `vic-machine create` to deploy VCHs, specifying an appropriate user account in the `--ops-user` option. Specify the inventory objects on which you have set permissions in the various `vic-machine create` options. For an example of a `vic-machine` command with the `--ops-user` option, see the section *Specify Different User Accounts for VCH Deployment and Operation* [Advanced Examples of Deploying a VCH](vch_installer_examples.md#ops-user).

<a name="required_roles"></a>
## Required VCH Roles ##

You must create the following vSphere roles for use by the `--ops-user` account.

|**Role to Create**|**Required Permissions**|
|---|---|
|`VCH - vcenter`|Datastore &gt; Configure datastore|
|`VCH - datacenter`|Datastore &gt; Configure datastore<br>Datastore &gt; Low level file operations|
|`VCH - datastore`|Datastore &gt; AllocateSpace<br>Datastore &gt; Browse datastore <br>Datastore &gt; Configure datastore<br>Datastore &gt; Remove file<br>Datastore &gt; Low level file operations<br>Host &gt; Configuration &gt; System management|
|`VCH - network`|Network &gt; Assign network|
|`VCH - endpoint`|dvPort group &gt; all permissions<br>vApp &gt; Add virtual machine<br>VirtualMachine &gt; Configuration &gt; Add new disk<br>VirtualMachine &gt; Configuration &gt; Advanced<br>VirtualMachine &gt; Configuration &gt; Add or Remove Device<br>VirtualMachine &gt; Configuration &gt; Remove disk<br>VirtualMachine &gt; Guest operations &gt; Guest operation program execution<br>VirtualMachine &gt; Interaction &gt; Device connection<br>VirtualMachine &gt; Interaction &gt; Power off<br>VirtualMachine &gt; Interaction &gt; Power on<br>VirtualMachine &gt; Inventory &gt; Create new<br>VirtualMachine &gt; Inventory &gt; Remove<br>VirtualMachine &gt; Inventory &gt; Register<br>VirtualMachine &gt; Inventory &gt; Unregister|

<a name="inventory_objects"></a>
## Inventory Objects to Which to Assign Permissions ##

You must assign permission to perform the VCH roles to the appropriate objects in the vCenter Server inventory.

|Inventory Object|Role to Assign|Propagate|
|---|---|---|
|Top-level vCenter Server instance|`VCH - vcenter`|No|
|Datacenters|`VCH - datacenter`|No|
|Clusters. All datastores in the cluster inherit permissions from the cluster.|`VCH - datastore`|Yes|
|Standalone VMware vSAN datastores|`VCH - datastore`|No|
|Standalone datastores|`VCH - datastore`|No|
|Network folders|`Read-only`|Yes|
|Port groups|`VCH - network`|No|
|Resource pools|`VCH - endpoint`|Yes|