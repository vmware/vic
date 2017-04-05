# VCH Update (P0)
This document defines the VCH update command line, and the related changes in vic-machine, to make VCH update process more reasonable.

Here is the key points considerred in this Definition
- Batch Operation

  We'd like to update VCH configurations more than one item at one time to make the command easy to use
- Reuse all options used in vic-machine create

  In vic-machine create, we already have 57 options to define VCH configurations, though not all of them can be updated, quite amount of them should be able to be changed, and this number will grow up over time, so we definitely do not want to introduce 3X that number of options to update (update means add/modify/remove)

## Options
First option is to introduce update primitives

### Delta Options
#### Update Primitives

- add
- modify
- remove, rm

In one update command, only one delta primitive is acceptable, and only that kind of operation will be executed.

#### Example Command

- vic-machine update --add --volume-store ... --container-network ...

  add additional volume-store and container-network into VCH
- vic-machine update --modify --dns-server ... --thumbprint ... --momery ...

  update dns-server, target thumbprint and resource pool memory in the VCH
- vic-machine update --rm --volume-store ...  --dns-server ...

  Remove one volume-store and remove custom dns-server from VCH

We do not allow command like following
- vic-machine update --add --volume-store ... --container-network ... --rm --volume-store ...  --dns-server ...

In this way, the update command can only run one kind of change at one time

### Whole Thing Update
Other than introducing delta primitives, we could also require user to specify all configurations during update. That means the old VCH configuration will be replaced completely by this new configuration.

Anything specified in the new configuration, which does not exist in old configuration will be added, anything missing in new configuration will be removed. And for the configuration can only have one item, will be replaced.

To have a better user experience, vic-machine should notify user for what's changed, for example, list of volumes will be added, and another list will be removed.

Pros of this option is no new options introduced and all options used in vic-machine create could be reused in update command.

Cons of this option is that every time to update a VCH, user need to know all the old configuration first, otherwise some configuration will be removed by mistake.


If there are many options specified to create VCH, and only a few things need to be tuned, the delta update option is more clear, and less effort to write the update command. But if there are less options specified in first create, and lots of changes needed, the whole update option will be easier to know what will be there after update.

### Separate Update to Sub Commands
The idea of this option is to update a few configuration each time, which is similar to Delta Options, but we're not using different options to separate the update scope, instead we can use different update sub command. For example, to update volume-store, we can use following command
```
vic-machine update volume-store --rm ds://datastore/volume:default
vic-machine update volume-store --add nfs://host:port/container:nfs
```

The command format is like ```vic-machine update <update object> --<update keyword> --<update option> <update value>```

Pros:
- The whole update options are splitted to many update sub commands. Then in each sub command, only a small number of options are available, and that does not increase over time. Only the sub command number will increase.
- Consistent with existing update firewall command format

Cons:
- We'll need to reinvent option names

  For example, to update client network, user need to specify all options if they want to use static ip. Then in one update command, we'll need all those options. Here is the command looks like ```vic-machine update bridge-network --port-group value --gateway value --ip value```
  The existing option name in vic-machine create is ```--client-network value   --client-network-gateway value  --client-network-ip value```, we'll need to remove the prefix cause that is already shown as sub-command name.

Preference:

- The sub commands option

## Show Existing VCH configuration
Right now, there is no user friendly way to view all configurations of existing VCH. To support update operation, we need to make it easier, otherwise, update only makes things weird.

vic-machine inspect command is available now, just it shows only the execution endpoint of VCH, instead of the configuration. We can extend that command, to support configuration inspection as well.

- vic-machine inspect <--configuration, --conf>
  show VCH configuration only
- vic-machine inspect <--execution, -e>
  show VCH execution endpoint (default behavior)

### Configuration Inspection
For configuration inspection, two output format should be supported

- yml format

  This should print out all VCH configuration with predefined categories, e.g. general information, storage, network, etc. in user friendly format. Secret information will be hidden or skipped.

Following is sample format:

```
General:
	Target: "vc.domain.com"
	Thumbprint: "98:97:D1:7B:DF:..."
	OpsUser: "Operator"
	Name: "VCH-1"
	Proxies: 
		- "https://host:port"
		- "http://host:port"
Compute:
	Path: " myCluster/Resources/myRP"
	CpuLimit: 0 MHz
	MemoryLimit: 0 MHZ
Storages:
	ImageStore: "datastore"
	VolumeStore:
		- "ds://datastore/volumes/default": "default"
		- "nfs://userid@host:port/volumes/": "nfs"
Networks:
	BridgeNetwork:
		PortGroup: "bridge"
		IpRange: "172.16.0.0/12"
	ClientNetwork:
		PortGroup: "client"
		Gateway: "10.1.0.0/16,10.2.0.0/16:10.0.0.1"
		...
	DNSServer:
		- "8.8.8.8"
		- "8.8.4.4"
Certificate:
	...
Registries:
	InsecureRegistries:
		- "..."
	RegistryCA:
		- "ca file"
Endpoint:
	Cpu: 1
	Memory: 2048 MHZ
	Debug: 3
	SSHEnabled: true
```

- Command Option Format (corresponding to Delta Update option)

  The yml format output is more readable for VIC admin, but if users want to recreate their vic-machine create or update command, it's hard to do based on that format, because there is no one on one mapping between the configuration and command options. 

  Here we'd support another option as "vic-machine inspect --cmd --conf", this will print out the vic-machine create command options minus secret information. 

  The sample output:
```
--target                vc.domain.com
--name                  VCH-1
--compute-resource      /vcqaDC/host/cls/Resources
--ops-user              operator
--image-store           datastore1
--volume-store          ds://datastore/volumes/default:default
--volume-store          nfs://userid@host:port/volumes/:nfs
--bridge-network        bridge
--bridge-network-range  172.16.0.0/12
...
```

Based on this output, it's easy to modify existing configuration to update command options

Note: Special character might be escaped based on current OS platform. But even with this, the command options cannot be copied to other platform.

- Subcommand Inspect (corresponding to update subcommand option)
  The update subcommand option might change command option based on each different sub command requirement, to get current configuration, we'll need to have command inspect for each subcommand. So the option is to extend vic-machine inspect command to include sub command as well.
  For example, ```vic-machine inspect bridge-network``` will print out current configuration in the command format, to make sure user can easily generate their update command based on that.

## Update In Progress Status
Same to vic-machine upgrade, concurrent update requests for same VCH is not allowed. vic-machine should be able to detect if one VCH is running update/upgrade. If true, return useful error message. And also should not leave that update flag in there after update is stopped, or even interrupted.

Issue https://github.com/vmware/vic/issues/4069 is tracking this problem.

## Backward compatible
It's possible that after vic engine released new version, user is trying to update older version's VCH with newer vic-machine.

There are few options:
- Run update using same version's vic-machine (this option should technically always work)

  For any customers that do not want to change vic-machine version, they need to run update using same version's vic-machine.

- Combine upgrade process into update

  Upgrade is actually one special kind of update, technically it's easy to have upgrade run together with update. User might want to update VCH configuration, and also want to upgrade it after or before update. In that case, a combined update and upgrade command can provide a better user experience.

 Here is the combined command ```vic-machine update --upgrade <all update and upgrade options supported here>```

- Support backward compatibility in vic-machine

  This will need to add backward compatible code in vic-machine. Even initially there is only few changes, that will expand massively over time, which is not consistent with our data migration design, so this is not our option.

# VCH Update Rollback (P0.5)
Similar to VCH upgrade workflow, before update VCH, snapshot should be took to keep existing status. If the update failed eventually, VCH should be reverted to old status automatically.

The Update Rollback is talking about after a successful update, vic-machine is still able to rollback VCH's old configuration. The command is like vic-machine update --rollback

This feature is nice to have, but as vic-machine upgrade already have this capability, mark it as P0.5.

# Manifest File (P1)
vic-machine create already has more than 50 options, and the number will increase over time. With this update feature added, most of those options need to be updated as well. To improve user experience, in my opinion, we need Manifest file to simplify vic-machine command generation.

[vic-machine design](vic-machine.md#vch-manifest) mentioned manifest idea, here we'll start from a simple implementation. The manifest file format will be same to the yml file format, mentioned in [configuration inspection](#configuration-inspection)

- vic-machine inspect --conf --manifest filename will dump all configurations in yml file format to the manifest file.
- vic-machine create --manifest will load the manifest file and create VCH based on it.
- vic-machine update --manifest will load the manifest file and update VCH configuration based on it.

This manifest file approach will simplify vic-machine create/update workflow, and leverage user readable configuration format. But it will be an additional work for existing vic-machine create, so just put it as P1.
