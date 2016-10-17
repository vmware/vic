# Upgrade

After installed VCH, VIC admin would want to upgrade existing instances to latest version. This document describes the upgrade requirement, and design options.

## Requirements

### Scope

While a new version is ready, new version's vic-machine upgrade can upgrade one running or powered off VCH appliance. The containers managed by this VCH will be left in existing version. To have a new version container, user need to create new container from upgraded VCH instance. If existing container is stopped and restarted, the old version should be kept.

VIC version includes three parts, "release tag"-"build id"-"github commit hash".

vic-machine upgrade will rely on build id change only. User need to provide newer binary with bigger build id, to upgrade existing VCH.

We could also use github commit hash for version diff, but that will be difficult to check if one version is newer than another one, so at this time, we'll leverage build id only. This requires our build system to keep increase build id no matter any system change. 

vic-machine upgrade cannot upgrade VCH to a version newer than itself.

### Impact

The VCH appliance (control plane) will have downtime during upgrade. Container lifecyles are not coupled with that of the appliance and, when using mapped vsphere networks (vic-machine --container-network argument) instead of port forwarding, the network data paths are not dependent on the appliance either. 

There will be impact on container interaction while the appliance is down:

- no direct access to container logs
- no attach ability
- network outage for NAT based port forwarding

All of those facets should resume normal operation after upgrade is complete

### VCH status

vic-machine ls and inspect will show the VCH versions and upgrade status.

### Roll Back

If anything wrong happens during upgrade, vic-machine will rollback VCH to original version and status. 

### Downgrade

For any running VCH, VIC admin can downgrade existing VCH to an old version. (I assume this function is only used for upgrade roll back, we might not want to recommend downgrade for anything working well, then the downgrade scope will be same to upgrade versions)

Verified docker-engine upgrade/downgrade, which works well between 1.11.2 and 1.12. The running container will be stopped after upgrade/downgrade, but is good to start again in new version. And all images and containers information is not lost. So from technical point of view, we should support same story.

### Internet connection
Internet connection is not required to upgrade/downgrade VCH, but newer version's binary should be available for vic-machine.

## Implementation/Issues

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

### Limitation
This document does not include approach for significant change which might break backward compatibility.
- if datastore structure change break current data consumption model, new version cannot read back images from old version's data, we'll not try to migrate docker data during upgrade.
- serial port is replaced by VMCI

# Plan

## Prerequisites

- Need to have versioning support before we claim we can upgrade from this version.
- Need to have VCH restart ready.
- As all version is built into binary, we need to automatically generate new build after tag release.

## Simple Upgrade

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

## Restrictions

User cannot run two upgrades for same VCH at the same time. 

vic-machine will check if there is already another upgrade snapshot is created before it start to create snapshot. But as create vsphere snapshot will take some time, e.g. one minute, if at this time, another upgrade process is started, it will start upgrade again cause the snapshot of previous task is not finished yet.
 