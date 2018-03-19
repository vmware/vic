Background
==========

Customers may wish to restrict the set of hosts a set of VMs are executed on.

This may be necessary for software licensing reasons (e.g., if they are billed
based on the number of physical hosts, sockets, or cores running a piece of
software), compliance reasons, or due to latency-sensitive workloads running in
an environment with stretched clusters.

Each of these use cases can be addressed by a similar pattern:
1. Create a DRS Host Group for the set of hosts which may run the workloads.
2. Create a DRS VM Group for the workloads themselves.
3. Create a VM-Host Affinity Rule to express that the identified workloads may
   only by run on the identified hosts.

Notes:
 * vSphere supports expressing the VM-Host Affinity Rule as either a requirement
   ("must") or preference ("should") [1]. If "must" rules are used, DRS will not
   allow the VMs to be run on other hosts, even in extreme circumstances (e.g.,
   HA would not perform a failover to a host not in the DRS Host Group). If
   "should" rules are used, violations will produce a log event and be reported
   as Faults on the cluster's DRS tab [2].
 * Because these rules are cluster-based, all hosts included in a DRS Host Group
   must reside in the same cluster [1].


Design
======

## Responsibility

It seems undesirable to require administrators to manage DRS Host Groups through
VIC, and it seems intractable for VIC to support such management without
duplicating significant vSphere functionality. Similarly, it seems that VM-Host
Affinity Rules should be managed directly in vSphere.

Given this, it is also reasonable for the Administrator to create the DRS VM
Group and supply its name (or an identifier) as a part of the VCH creation or
configuration process.

This design provides significant flexibility: we can support both "must" and
"should" rules as well as both Affinity and Anti-Affinity rules with no added
implementation effort (testing may be desirable). (This is insufficient to
support VM-VM Affinity, as those rules do not operate on DRS VM Groups [3].)

Because the administrator creates these vSphere entities, it is reasonable to
expect them to manage the rest of their lifecycle. As such, deletion of these
entities need not be tied to deletion of a VCH.

**Note:** Creation of an empty DRS VM Group may not be supported. This is not
          clear from the documentation [4]. One alternative would be to require
          that the administrator create a VCH, configure DRS as desired, and
          then supply the DRS VM Group as a part of a VCH configure operation.

          If creation of an empty DRS VM Group is not supported, we should
          create one automatically during the creation of each VCH. (Worth
          considering for a future release regardless.)

Impacted areas: VCH create, configure, and inspect (via CLI, API, and UI); documentation


## Permissions

DRS VM Groups are expressed as a `ClusterVmGroup` [5]  in a cluster's
configuration, which will need to be updated to include both the VCH Endpoint VM
when it is created and the VCH's container VMs as they are created.

To update this portion of the cluster's configuration, the user creating the VCH
and the operations users will both require the `Host.Inventory.EditCluster` [6].

Impacted areas: VCH create (validation and operations user grant); documentation


## Execution

To avoid placing each VM twice, it is desirable to associate each VM (including
both the VCH Endpoint VM and container VMs) with the DRS VM Group between
creation and initial power-on.

Impacted areas: portlayer


## Upgrade

By default, we will not attempt association with a DRS VM Group. This ensures
that existing VCHs can be safely upgraded without the need to grant existing
operations users additional [permissions](#Permissions).

Impacted areas: testing


Decomposition
=============

## 0. Investigate empty DRS VM Groups

The inability to create an empty DRS VM Group would have significant impact on
the overall workflow and user experience; it would not be possible to have the
administrator do all setup work prior to creating the VCH.

We should investigate this early as it may affect the overall design/approach.


## 1. Prototype portlayer work

As an initial step to validate this design, we should prototype the portlayer
changes using hard-coded rules and groups. This allows us to ensure that a VCH
works as expected when configured as we intend.

This should also be prioritized early, to ensure that effort isn't wasted on any
user-facing changes that may need to be changed as a result of design changes.


## 2. Finalize portlayer work

This is likely to involve moving the hard-coded portions from the prototyping
work above out of the portlayer and into the install code.

This may also involve refactoring of `Commit`, as it is "getting excessively
large and should be broken out into a separate subfunction with the addition of
this logic" [7].

Failure of the operation to update the `ClusterVmGroup` should be treated as a
fatal error.

No operations requiring `Host.Inventory.EditCluster` should be performed when no
DRS VM Group has been specified.


## 3. Update vic-machine create

This is essentially using a DRS VM Group value supplied by the user instead of
the hard-coded value used to enable the portlayer work to proceed.

This includes both tagging the VCH Endpoint VM itself during creation as well as
persisting the DRS VM Group so that container VMs are tagged as well.

This includes changes to both the CLI (which will probably take a rule name) and
API (which will probably accept either a rule name or identifier). 

This requires validation that the DRS VM Group exists and that the operations
user has the necessary `Host.Inventory.EditCluster` privilege.

This should include end-to-end testing of both the CLI and API, which will
require additional investigation to establish a pattern for configuring the
prerequisites (DRS Host Group, DRS VM Group, VM-Host Affinity Rule) from robot.


## 4 Update vic-machine inspect

The inspect CLI (including inspect config) and API should be updated to return
the configured DRS VM Group.

This should include end-to-end testing of both the CLI and API, to ensure that
the correct value is returned when a group is configured and a sane response
is returned when no group is configured.


## 5. Update vic-machine configure

Apply the changes from vic-machine create to configure as well, including
retroactively tagging the VCH Endpoint VM.

This requires validation that the operations user has the necessary right when
performing any configure which introduces a DRS VM Group.

There may be complexity around ensuring the atomicity of this operation [8].


## 6. Update VCH Management UI

As a part the VCH Creation Wizard, users should be able to select a DRS VM Group
to be supplied to the API.


## 7. Operations user grant

The operations user "grant permissions" logic will need to be updated to include
granting `Host.Inventory.EditCluster` on the VCH's cluster.

We should then leverage this grant functionality in some of the testing for the
feature as a way of validating that it works as intended.


## 8. Additional testing

We should implement an additional end-to-end test case to ensure that an old
VCH, created with an operations user that did not have this permission, can be
upgraded to a version (i.e., one with the DRS VM Group functionality) and
additional containers can be created without issue. 


## 9. Documentation

Update documentation to reflect this work.


Open Questions
==============

1. Is it possible to create an empty DRS VM Group?
2. Does deleting a VM remove it from any `ClusterVmGroup`s it was a part of?
3. Does unregistering a VM remove it from any `ClusterVmGroup`s it was a part of?
4. Does govc expose the necessary commands to create/edit host/vm-group rules?


References 
==========

1 - https://docs.vmware.com/en/VMware-vSphere/6.5/com.vmware.vsphere.resmgmt.doc/GUID-2FB90EF5-7733-4095-8B66-F10D6C57B820.html
2 - https://docs.vmware.com/en/VMware-vSphere/6.5/com.vmware.vsphere.resmgmt.doc/GUID-FF28F29C-8B67-4EFF-A2EF-63B3537E6934.html
3 - https://docs.vmware.com/en/VMware-vSphere/6.5/com.vmware.vsphere.resmgmt.doc/GUID-94FCC204-115A-4918-9533-BFC588338ECB.html
4 - https://docs.vmware.com/en/VMware-vSphere/6.5/com.vmware.vsphere.resmgmt.doc/GUID-56C0F68B-23D7-4CD3-A93A-BCF20EAA0C35.html
5 - http://pubs.vmware.com/vsphere-6-0/index.jsp?topic=%2Fcom.vmware.wssdk.apiref.doc%2Fvim.cluster.VmGroup.html
6 - http://pubs.vmware.com/vsphere-6-0/index.jsp#com.vmware.wssdk.apiref.doc/vim.ComputeResource.html?path=3_1_0_2_5_16_14#reconfigureEx
7 - https://github.com/vmware/vic/issues/6461#issuecomment-332584964
8 - https://github.com/vmware/vic/issues/6461#issuecomment-374366308
