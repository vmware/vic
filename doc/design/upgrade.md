# Upgrade

After installing a VCH, a VIC engine admin will need to manage the lifecycle of the VCH. One aspect of this lifecycle management is upgrading and patching. For our purposes, patching and upgrading will be treated the same. 

The VIC engine product will release a new, complete, software download bundle for both patches and upgrades, versus, having a sparse patch or a separate upgrade model.

## Requirements

### Scope

vic-machine can upgrade VCH appliance no matter it's running or not, so vic-machine upgrade cannot rely on services running in VCH appliance.

For the VCH appliance, vic-machine upgrade will detect the version difference between the old VCH and itself, this includes detecting the guestinfo changes between the two versions. Then vic-machine upgrade will migrate the guestinfo metadata, and replace the iso files. If anything incorrect happens during this process, vic-machine upgrade will revert back to old version and status. The VCH will have downtime during upgrade, but should not be broken after this process.

The containers managed by this VCH will be left running the existing version. If the container is running, it will work well if it's not coupled with appliance through port-forwarding. But after VCH is upgraded, the container management might not be fully functional if there is configuration or communication level changes.

To update the containers, especially while there is security patch available, two options are available. 

- First is to upgrade the containers from vic-machine. We can migrate container configuration, replace iso file, and then restart the container process. This update requires container downtime as we cannot update tether without stop container process at this time.

 We could upgrade the containers sequentially or in parallel based on user preference.

- Second, leave it to customer like what docker did. User could easily destroy and recreate same containers through docker commands. The new container will be in new version automatically as container iso file is already updated during the appliance upgrade.

In theory, the first option does not add much value, as both options need downtime for the container. If user could accept container stop, they could recreate that container through their own script easily and flexibily, which could be related to their own business requirement.
Ideally, we should improve tether to separate tether and container lifecycle, so we can replace iso file and restart tether process transparently. Then it will be helpful to have vic-machine upgrade to update containers.

My preference: second option.

### Version Difference

VIC version includes three parts, "release tag"-"build id"-"github commit hash". vic-machine upgrade will rely on build id to detect which version is newer or older than another one.

User need to provide newer binary with bigger build id, to upgrade existing VCH. This also requires our build system to keep increase build id number no matter any kind of system change. 

Note: after introduced data migration between builds, we introduced one more internal used data migration version. That version is related to data migration plugin only, so will not be shown in vic-machine version. This will be described in data migration section.

### Impact

The VCH appliance (control plane) will have downtime during upgrade. Container lifecyles are not coupled with that of the appliance and, when using mapped vsphere networks (vic-machine --container-network argument) instead of port forwarding, the network data paths are not dependent on the appliance either. 

There will be impact on container interaction while the appliance is down:

- no direct access to container logs
- no attach ability
- network outage for NAT based port forwarding

All of those facets should resume normal operation after upgrade is complete. 

Note: Exception might happen if there is container configuration or communication changes. Addtional operations are required to fix the problem.

### VCH status

vic-machine ls and inspect will show the VCH versions and upgrade status.

### Roll Back

If anything wrong happens during upgrade, vic-machine will rollback VCH to original version and status. 

### Downgrade

Verified docker-engine upgrade/downgrade, which works well between 1.11.2 and 1.12. The running container will be stopped after upgrade/downgrade, but is good to start again in new version. And all images and containers information is not lost. But think of the complexity to support version downgrade, we'll not go with this at this moment, unless there is explicity user requirement.

### Internet connection
Internet connection is not required to upgrade/downgrade VCH, but newer version's binary should be available for vic-machine.

## Design/Implementation - Phase1

This section described the first simple implementation of upgrade. After this phase, user could upgrade VCH appliance from build to build, as long as there is no metadata changes, which means no changes in guestinfo, key value store and image metadata.

This code is already merged, and works for security patch update.

### Versioning

#### Freeze following attributes definition in VCH configuration and container VM configuration
- VCH configuration: ID, Name, Version
- Container VM serial port configuration
- Container VM log/debug file location

#### Configuration Version

Both VCH and container configuration will have version and the value is same with what vic-machine version command shown. After upgrade, VCH configuration version should be updated to new version, but container VM configuration version will still be old one.

New version's VCH will work with both old and new versions containers.

#### Embed iso file version

vic-machine need to identify the iso file version before and after upgrade. Two options here:
- Add version into iso file Primary Volume Descriptor, and read back from vic-machine
- Leverage iso file name including version, e.g. appliance-0.5.0.iso

vic-machine should check iso file version during deployment, and after upload to datastore, version should be appended to iso file name, to make sure mutilple iso files could co-exist in the same VCH, and used for different version's container.

Image file name should not be changed in datastore, cause the file path is used to create container VM. And then vic-machine will leverage file name for version checking during upgrade, to avoid download iso file from datastore back to where vic-machine is running.

Note: No feasible golang library found for this function, so will write our own library to read iso metadata.

### Configuration Data Migration

To support VCH upgrade, any changes in VCH configuration and container VM configuration will introduce additional effort to migrate data from old version to newer version. And if we support a long upgrade path, for example from 2000->2500->3000->3500->4000,  there will have a list of utilies recognize old configuration structure to help data migration.

#### Refactor current VCH configuration and container VM configuration structure
To make data migration easier, we'll need to refactor current VCH and container VM configuration structure.

VirtualContainerHostConfigSpec is too big, which includes everything in one structure. We're not be able to update part of configuration structure, it will help to update part of data with separated structures based on functionality, e.g. network, storage, and common attributes.

For ContainerVM configuration, there are a few attributes not used for container VM setup, e.g. in ExecutorConfig, Key, LayerID, RepoName, are all not container VM configuration related. Move out irrelated attributes can help the structure stability in the future. 

### Transactional Operation/Roll Back
VM snapshot is one good option for upgrade roll back. Before upgrade, vic-machine will create one snapshot of VCH. If anything wrong happens, vic-machine can switch back to the pre-upgrade snapshot.

Note: VM snapshot does not persist data in serial port, as we will use datastore files through serial port for vch logs, after roll back, log files will have all error message during upgrade if VCH appliance VM is ever started with new configuration.

### Upgrade Workflow

At this point, we'd like to start from a simple solution, which does not include data migration. 

Here is the upgrade workflow
- Find existing VCH appliance
- Query VCH configuration
- Verify new version and existing VCH's version
- Generate new VCH configuration from old one
- Valiate VCH configuration with new version's validator. If anything incompatible found, stop the upgrade. As no data migration yet.
- Upload iso files
- Prepare combined applianceVM configuration, both hardware spec (ISO path) and the extraconfig portion (VCH version)

--- everything until this portion can be done without interrupting ongoing operation ---
- Snapshot VCH appliance VM.
- Poweroff VCH appliance VM if it's not
- Update container VM configuration with new config (it should be viable to do this in a sigle reconfigure)
- Update VCH configuration with the prepared configuration change
- Power on VCH, and wait VCH initialization
- Anything wrong in the above steps, roll back to upgrade snapshot
- Ensure appliance is powered on after rollback to snapshot
- Cleanup env after upgrade (remove upgrade snapshot, remove uploaded iso files if upgrade failed)

This ensures that a failure to upload ISOs for whatever reason is detected before we take down the existing version. It limits the failure modes after shutting down the appliance to:

1. failed to update appliance configuration
2. failed to power on the appliance (e.g. system resource constraints)
3. failed to boot appliance
4. failed to acquire network addresses (this has been seen in real world examples - we may wish to attempt to preserve/reuse the IPs the appliance had prior to update, which should still be present in the extraconfig)
5. failed to rebuild system state from vsphere infrastructure

## Design - Phase2

Although we did lots of code refactor for VCH appliance configuration, it's still unavoidable to continue changing that structure. To make sure we don't break  upgrade after GA, we need a solution for the data migration anyway.

- Appliance guestinfo, which is used to persist VCH appliance configuration
- Container guestinfo, which is used to persist container configuration
- KeyValue Store, which is used to persist image metadata and network portlayer data
- Appliance log files
- vSphere object management logic change

### Guestinfo Migration

- One VirtualContainerHost structure per one version, this will be used by extraconfig package to read back old vch's configuration
- Write plugin to migrate data from from previous version, and register that plugin with correct release version
- Each version should have a separate directory, which contains all configuration files and plugins for that version
- vic-machine upgrade framework will calculate the plugin chain based on old VCH's version and vic-machine version, and then call those plugins sequentially. This will look like r1->r2->r3->f4->latest or r3->r4->r5
- Run guestinfo validation to figure out if there is any field missing (one possible option here is to translate configuration back to vic-machine input data, and then rerun validator to make sure all vsphere related objects are still exist and correct)
- Prompt missing parameters if there is any (for new parameters added and cannot be empty)
- Write new configuration back to appliance guestinfo

With this design, developer is able to add/delete/repurpose existing configuration items, except VCH version. And we'll need one guidline for what to do while one configuration item is changed.

Besides of this solution, we also dicussed another option, that is to revert configuration back to vic-machine create option, and try to recreate a new VCH appliance, and then migrate containers from the old appliance to new appliance. This solution can reduce VCH downtime, and even migrate alive.

We have two concerns for this option. 

First and most important, is that we cannot actually migrate containers from appliance1 to appliance2. This will introduce too many production refactor. For example, image store directory is created uniquely with appliance UUID, network will have conflict, connection over serial port need to be switched, etc.

Second concern is that we'll need to version more objects and logics, includes appliance guestinfo, vic-machine create options, and the logic to revert guestinfo back to vic-machine create options, those are all version sensitive.

We didn't try to figure out how to do it at this time, cause it looks overdesigned.

Note: extraconfig package should always be backward compatible. If it breaks this assumption, upgrade does not work for both options.

### Secret Configuration

In appliance configuration, there are few items are encrypted for security reason, including user password and certificate so far. The encryption key is readable from in guest only. vic-machine is reading configuration through vsphere API, so not be able to decrypt secret information from old appliance configuration.

Ideally, it's ok to leave them there, but if in the future there are new secret items added, it's not possible to reencrypt with old encryption key, and no way to get back old encrypted value.

For this issue, one option is to get input from user for all existing secret value, and together with new secret filed, encrypt with new key. Currently, only operation user password and certificate key are required.

This feature can be delayed until there is new secret key added.

### Container Guestinfo

VIC should not force user to upgrade containers after upgrade appliance, which means appliance will need to talk with old version's containers. Read back container configuration will be similiar to appliance configuration upgrade process. Corresponding data migration plugin should be prepared for each change, and registered to upgrade framework with release version. The process is listed below.

- One ExecutorConfig structure per one version, where the container configuration is changed. If no changes, that version will be skipped.
- Portlayer will check configuration version first while reading configuration. If it's older version, run in memory data migration through upgrade framework.
- Together with upgrade configuration, portlayer need to know if configuration is migrated or not.
- If no data migration, continue same to new version's container
- If has data migration, do not write back configuration (old version's container is running in readonly mode)

If container data migration happens, some of docker functions, e.g. start/attach/log, and docker ps/inspect information might be impacted, but docker image/volume/stop/rm should still work well, to make sure user can destroy and recreate container easily.

During appliance upgrade, vic-machine should have clear information to mention the container functional limitation after upgrade, and the solution for it.

### KeyValue Store Migration
Different to guestinfo object, VIC will not change configuration from guestinfo to anywhere else anytime soon. But from start, keyvalue store persistent position is in argument. Right now, vic has a few datastore files for keyvalue store, but those information is not configurable, which means hardcoded in portlayer.

In the future, if there is any change on assumption no matter for what reason, performance, easy to management, data migration will be hard to do.

So the idea is that, data migration framework does not assume where to load keyvalue store. If there is keyvalue store changes, plugin is the owner to read old version, update to new version and then write back to datastore files, or to anywhere defined in new version, suppose the configuration can be retrieved from guestinfo. So the input of keyvalue store plugin will be same to other appliance configuration update plugin.

The problem of this solution is that each single plugin will persist its own change, not like the guestinfo update, which is migrated in memory and persisted by migration framework. So the rollback for guestinfo is easy, but not possible for keyvalue store plugin, unless we have data roll back plugin mechanism.

The workaround for this idea is to restrict keyvalue store migration functions, e.g. only new key is supported. In the plugin, do not remove old keyvalues, so in the above failure case, old binary can still work with old keyvalues.

We have lots of other options to better support keyvalue store upgrade, but at this time, I don't want to predicate our changes in the future, that's easy to be over-engineering. While that happens, we could think about to support plugin roll back, or upgrade from portlayer, etc.

### Appliance log files

Appliance log files configuration will not be touched at this time. So upgraded VCH will use old log files if there is any such kind of change. We can revisit this while there is customer requiement on this.

### vSphere Object Management Logic

There will have new vSphere API come up, so the logic to manage vsphere objects will be enhanced, for example, as Caglar said the vmdk can be managed directly in vSphere 6.5, instead of through vm operations in vSphere 6.0. We'll switch to these new interfaces for image and volume management sometime later.

vic-machine is not supposed to migrate old image data or volume data to new vmdk files, so portlayer will need to be backward compatible. (Suppose vSphere will be backward compatibile, there should be nothing to do in vic, but need to do some research on it)

### How to Support Upgrade Between Builds

We need to keep old version's guestinfo structure definition, to make sure we can read back old configuration, and need plugin to migrate data from old structure to new structure. With current version mechanism, developer do not know which build version will be generated for his change, so we introduced another version, data migration version.

If the value is not set, by default it is 0.

Data migration framework will detect version difference between old version and latest and then run corresponding data migration plugins sequentially to migrate data.

Notes:
- Developers who change guestinfo, keyvalue store will be the owner to develop migration plugin, and be responsible to increase data migration version. 
- Each migration version should have one and only one corresponding plugin.
- If appliance configuration and container configuration are changed at the same time, two different plugins should be added, and registered to different plugin category.
- If both appliance configuration and keyvalue store are changed, two plugins are recommended as well.

### Limitation
This document does not include approach for tether communication change, that means from serial port to VMCI or something else. Need to think about how to support old version's container while that change is made.

## Restrictions

User cannot run two upgrades for same VCH at the same time. 

vic-machine will check if there is already another upgrade snapshot is created before it starts to create snapshot. But as create vsphere snapshot will take some time, e.g. one minute, if at this time, another upgrade process is started, it will start upgrade again cause the snapshot of previous task is not finished yet.

