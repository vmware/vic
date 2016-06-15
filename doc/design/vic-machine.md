# On the care and feeding of VCHs - vic-machine

vic-machine is both the management client for Virtual Container Hosts and the mechanism by which they are initially deployed.

## Roles, responsibilities, and multi-tenacy models
This document notes three separate roles in the course of deploying and managing a VCH:
* vSphere Administrator (_viadmin_)
* VCH Administrator (_admin_)
* VCH User (_user_)

This workflow is the bridge between the infrastructure administration portion of an organisation and the users. The specific intent here is to allow each of the roles to operate at a level of detail appropriate to that role, with minimal impact or dependency on the others. The:
* _viadmin_ - identified by having administrative access to vSphere
  - operates at a business decision level, mapping relative priority of projects and teams to accessible resources and permissible resource limits. Delegates usage authority of those resources to the _admin_ within specific limits enforced by vSphere.
* _admin_ - identified by delegated authority in the form of access to a signed VCH manifest file
  - controls sub-division, if any, of assigned resources among finer-grained projects and teams
  - deploys a VCH to manage a set of assigned resources. Delegates consumption decisions about those resources to _user_
* _user_ - identified by granted API access to a specific VCH
  - controls the specifics of _what_ is done with the available resources, by way of an API client such as the Docker client

This model provides for a form of multi-tenancy, with the _viadmin_ able to specify a service account that will be configured with appropriate RBAC rules. It's unclear at this time if it's viable to create sub-users with further restricted RBAC rulesets, so the working assumption is that all sub-division performed by _admin_ operate with the same service account, and the sub-division is enforced by configuration of vSphere constructs such as resource pools, but without the authority isolation provided by different vSphere users. How a VCH prevents manipulation of those sub-division limits is covered in [the security architecture](security.md).
It is not intended that there be RBAC within a VCH at this time.

How this multi-tenacy model is used is left up to the business, but there are two primary models that we consider during development:
* team based - the VCH is assigned to a team for their use, potentially running a mix of independent workloads
* application based - a VCH is used as the management construct for a given appliancation, with all containers being portions of that app


## VCH manifest

Certain information is required to deploy a VCH, with the _viadmin_ and _admin_ each contributing portions of that data. The VCH manifest is the mechanism by which that cooperation happens, and is the token via which authority delegation from _viadmin_ to _admin_ occurs. Creation of a manifest is conceptually a compositing process with the following inputs and an immutable result:

1. input manifest (if omitted, this creates a new manifest from scratch)
2. a restriction set of some kind ( e.g. compute resource path, datastore prefix, registry whitelist)
  - as a note, user credentials are also considered a restriction as they control permissible operations


The _viadmin_ uses vic-machine to create a _base manifest_ containing the following restrictions:
* target vSphere environment
* vSphere user and credential for created VCHs:
  - must already exist, or
  - must be created by vic-machine during manifest creation, or
  - demand created during VCH deployment, requiring stored admin credentials:
* stored _viadmin_ credentials if necessary:
    - encrypted credentials and [validating proxy](components.md#validating-proxy) URI, or
    - unencrypted credentials (_viadmin_ and _admin_ roles are held by the same entity and manifests are stored securely)

The _base manifest_ is the minimum set of informtion necessary from the _viadmin_ role. There is additional information that will almost always be required for full function of a VCH:
* vSphere network to use for container network
  - must already exist
* vSphere switch to use for container network:
  - must already exist
  - if switch, port group or network is demand created during VCH deployment, requiring stored admin credentials
    - in an enterprise environment, creation of a portgroup without VXLAN or VLAN is unlikely to suffice, so pre-existing portgroup is highly recommended

To create a VCH, disambiguation of which resource to use from the available set is necessary. If there is no ambiguity, i.e only one of each, then this can be omitted:
* _client_ network - necessary for any access to the VCH
* default container network:
  - _private_ port group or dedicated vSwitch, or
  - IPAM config for _client_ network (defaults to DHCP)
* datastore paths for (may all be the same path):
  - images
  - containers
  - volumes


## Implementation approach

Recommended initial approach is to have several components that get rolled into both vic-machine. Initial implementation should focus on:
* CLI parsing into config structs
* validation of config struct values
* creation of VCH from config struct

![vic-machine high level logic](images/vic-machine-high-level.png)

Components:
* command line argument parsing
  - maps command line arguments to internal config data structures
  - _package: main, path: cmd/vic-machine_
* validation of config structure values
  - given an internal config data structure it's necessary to check the specifics against the target vSphere to ensure they are valid - past validity is no indicator of current validity
  - this may also translate from symbolic names to morefs, resulting in a manifest that is resilient to name changes, but fragile across vSphere instances with identical naming schemes. This should be a user selectable behaviour, I prefer it on by default.
  - _package: spec, path: lib/spec_
* creation of VCH from configuration
  - when a manifest has been validated against a vSphere, this component creates the corresponding vSphere objects - this is primarily the VCH applianceVM, but may also include port groups and other objects.
  - _package: management, path: install/management_
* [vmomi gateway](components.md#vmomi_gateway)
* manifest creation and consistency
  - logic to validate that layered restrictions don't violate prior restrictions, e.g. refining [datastore1]/a/path to [datastore1]/a/path/to/vch is acceptable, but not to [datastore1]/a/second/path
  - this should operate on the internal config data structures - probably a sliding window of two layers of the composite manifest each loaded into a config struct
  - signing and signature validation
  - _package: manifest, path: install/manifest_
* load/save manifest to/from internal config data structures
  - this is the mapping from the config held in the compositied manifest to the current end configuration, held in a serializable config structure
  - signing and signature validation of manifest layers
  - _package: manifest, path: install/manifest_

The reason for this breakdown is because some of these elements need to be duplicated in the:
* [validating proxy](components.md#validating-proxy):
  - load/save manifest
  - manifest consistency
  - validation
  - reification
  - vmomi gateway
* applianceVM
  - validation
  - vmomi gateway


## =================================================
below this point is working notes.

## Installing - per vSphere target

### Inputs

1. vSphere SDK endpoint
2. vSphere administrative credentials

### Actions

1. deploy ESX agents
2. upload ISOs to common location
3. create custom tasks, alerts, and icons
4. create VCH tag (enable filtering of VCHs)
5. install UI plugin


## Installing - per VCH

### Inputs

* VCH user (existing or new) **
* resource lists:
 - pool **
 - imagestore datastore paths **
 - container datastore paths **
 - volume datastore paths (restriction)
 - network mappings:
  - one network minimum for VCH comms**
  - other network mappings
* resource allotments:
 - cpu
 - memory
 - network IO
 - disk IO
 - datastore quotas (per datastore path)
* certificates
 - users - for access to VCH
 - hosts - for container access to external hosts
 - network - for VCH/container access to networks (gating proxies)
* registry lists
 - whitelist
 - blacklist
* default container resource reservations and limits *
* containerVM naming convention (displayName for vSphere) *

### Actions

Some of the elevated privilege operations could be delegated during self-provisioning to avoid manifestations of un-utilized authority, e.g. resource pool, user, and rbac entries for a potential but uncreated VCH. This delegation of higher authority requires additional care in the self-provisoning path.

### Requiring elevated privileges
1. create vSphere user for VCH
2. create RBAC entries for VCH resources - resource pool, datastores, networks, et al
3. obtain credentials for VCH user (e.g. SSO token)
  * should be revokable
  * should only have expiration date if no concern about clean VCH retirement
4. create and size VCH resource pool/vApp
  * if vApp then should also configure the start/stop behaviours
  * this may encompass disabling certain operations via the UI
5. place credentials in VCH applianceVM extraConfig

## Requiring VCH user privileges
1. validate supplied configuration
2. construct extraConfig/guestinfo configuration for applianceVM
3. create VCH applianceVM
  * this may encompass disabling certain operations via the UI
4. upload ISOs if not shared
5. initialize applianceVM

At this point install transitions to managing - reporting VCH status from initial install is the same as reporting that information for any VCH regardless of age.


## Deleting - per vSphere target

### Inputs

1. vSphere SDK endpoint
2. vSphere administrative credentials
3. VCH identification
 - VCH resource pool path in govc format, e.g. /ha-datacenter/vm
 - VCH name
 - VCH ID, which is the value returned by vic-machine ls

The VCH name and VCH resource pool path are identical to the value used in vic-machine create command. Which can uniquely identify one VCH instance.
The VCH ID is the value returned by vic-machine ls, which probably be the VM mob-id, query from VC or ESX.

Either VCH ID or the VCH resource pool path and VCH name should be specified.

### Actions

1. Get VCH VM
2. Read back VCH configuration from VM guestinfo
3. Delete following resources based on VCH configuration.
 - Container VMs managed by the VCH
   - Contianer datastore paths and resource pool path configured during installation will be used here, to detect if the VM belongs to this VCH. If yes, these VMs will be removed.
    - Container VMs will be removed if they are in stopped status.
     - If container VMs are in powered on state, delete will return failure if -force is not specified.
      - If container VMs are in powered on state, and -force is specified, vic-machine delete will power off and remove them those VMs.

 - volumes managed by the VCH
   - -force option is required to delete volumes together with VCH uninstallation.
    - Volume datastore path (sample: ds://datastore1/volume/vch1) configured during installation will be used here, to detect if the volumes are created by this VCH.
     - If volume directory is empty, vic-machine delete will delete this directory as well, otherwise warning it.
 - images managed by the VCH
   - Images are shared by VCHs. Currently no reference or metadata to specify which image is referenced by which VCH, so for TP3, vic-machine delete will not touch images.
 - vSphere networks created
   - vic-machine create or VCH port-layer-server will create network, so vic-machine delete should delete any new networks. (Still need to confirm for how to get created networks)
 - VCH specific metadata from vsphere objects if any
   - Cause image will not be removed at this time, image metadata will not touched as well. But for container metadata, if any, will be deleted. (Need to finalize where it is)
 - the resource pool and appliance VM

### Samples
```
vic-machine delete 
--target root:password@192.168.1.1
--vch-path /ha-datacenter/vm/vch1
--force
```
This command will delete VCH /ha-datacenter/vm/vch1 appliance, all containers VMs, networks and volumes created by this VCH.

## Managing a VCH

* report VCH status and information (API endpoint, log server, et al)
* update VCH configuration - implies possible restart of component
* shutdown/reboot VCH
* upgrade VCH - should have an entirely separate doc for this
