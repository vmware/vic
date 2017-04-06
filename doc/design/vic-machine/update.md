# VCH Reconfigure (P0)
This document defines the VCH reconfigure command line, and the related changes in vic-machine, to make VCH reconfigure process more reasonable.

Here is the key points considerred in this Definition
- Batch Operation

  We'd like to change VCH configurations more than one item at one time to make the command easy to use
- Reuse all options used in vic-machine create

  In vic-machine create, we already have 57 options to define VCH configurations, though not all of them can be updated, quite amount of them should be able to be changed, and this number will grow up over time, so we definitely do not want to introduce 3X that number of options to reconfigure (update means add/modify/remove)

## Options
First option is to reconfigure all in one command

### Implicit Add/Modify and Explicit Delete
All existing vic-machine create options are supported in vic-machine reconfigure command. The meaning of those options are as following:
- for elements with labels (e.g. container-network, volume-store), if there's a value provided against that label in the reconfigure, it's a modify operation
- for elements without labels (e.g. --cert), if there's a value provided it's a replacement
- if need to remove existing configuration (e.g. static ip configuration for client-network, or registry-ca), define new options as --clear-registry-cas, --clear-client-network. Then existing configuration for registry-ca and/or client-network will be cleared. If not new configuration e.g. for client-network specified, it will share with other networks with same logic in vic-machine create.

The benefit of this option is it can update all things need to be updated in one command, but will introduce a bit more new delete options.

### Separate Reconfigure to Sub Commands
The idea of this option is to update a few configuration each time, which is one delta update option. For example, to update volume-store, we can use following command
```
vic-machine reconfigure volume-store --rm ds://datastore/volume:default
vic-machine reconfigure volume-store --add nfs://host:port/container:nfs
```

The command format is like ```vic-machine reconfigure <reconfigure object> --<reconfigure keyword> --<reconfigure option> <configuration value>```

Pros:
- The whole configuration options are splitted to many sub commands. Then in each sub command, only a small number of options are available, and that does not increase over time. Only the sub command number will increase.
- Consistent with existing update firewall command format (though the command name might be different)

Cons:
- We'll need to reinvent option names

  For example, to update client network, user need to specify all options if they want to use static ip. Then in one update command, we'll need all those options. Here is the command looks like ```vic-machine reconfigure bridge-network --port-group value --gateway value --ip value```
  The existing option name in vic-machine create is ```--client-network value   --client-network-gateway value  --client-network-ip value```, we'll need to remove the prefix cause that is already shown as sub-command name.

Preference:

- Looks like the first option is better, just leave the sub command option here, in case someone could find out more benefit from that option

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

- Command Option Format

  The yml format output is more readable for VIC admin, but if users want to recreate their vic-machine create or reconfigure command, it's hard to do based on that format, because there is no one on one mapping between the configuration and command options. 

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

Based on this output, it's easy to modify existing configuration to reconfigure command options

Note: Special character might be escaped based on current OS platform. But even with this, the command options cannot be copied to other platform.

- Subcommand Inspect (corresponding to reconfigure subcommand option)
  The reconfigure subcommand option might change command option based on each different sub command requirement, to get current configuration, we'll need to have command inspect for each subcommand. So the option is to extend vic-machine inspect command to include sub command as well.
  For example, ```vic-machine inspect bridge-network``` will print out current configuration in the command format, to make sure user can easily generate their reconfigure command based on that.

## Update In Progress Status
Same to vic-machine upgrade, concurrent reconfigure requests for same VCH is not allowed. vic-machine should be able to detect if one VCH is running reconfigure/upgrade. If true, return useful error message. And also should not leave that reconfigure flag in there after reconfigure is stopped, or even interrupted.

Issue https://github.com/vmware/vic/issues/4069 is tracking this problem.

## Backward compatible
It's possible that after vic engine released new version, user is trying to reconfigure older version's VCH with newer vic-machine.

There are few options:
- Run reconfigure using same version's vic-machine (this option should technically always work)

  For any customer that does not want to change vic-machine version, they need to run reconfigure using same version's vic-machine.

- Combine upgrade process into reconfigure

  Upgrade is actually one special kind of reconfigure, technically it's easy to have upgrade run together with reconfigure. User might want to reconfigure VCH, and also want to upgrade it after or before reconfigure. In that case, a combined reconfigure and upgrade command can provide a better user experience.

 Here is the combined command ```vic-machine reconfigure --upgrade <all reconfigure and upgrade options supported here>```

- Support backward compatibility in vic-machine

  This will need to add backward compatible code in vic-machine. Even initially there is only few changes, that will expand massively over time, which is not consistent with our data migration design, so this is not our option.

# VCH reconfigure Rollback (P0.5)
Similar to VCH upgrade workflow, before reconfigure VCH, snapshot should be took to keep existing status. If the reconfigure failed eventually, VCH should be reverted to old status automatically.

The Reconfigure Rollback is talking about after a successful reconfigure, vic-machine is still able to rollback VCH's old configuration. The command is like vic-machine reconfigure --rollback

This feature is nice to have, but as vic-machine upgrade already have this capability, mark it as P0.5.

# Manifest File
vic-machine create already has more than 50 options, and the number will increase over time. With this reconfigure feature added, most of those options need to be updated as well. To improve user experience, in my opinion, we need Manifest file to simplify vic-machine command generation.

[vic-machine design](vic-machine.md#vch-manifest) mentioned manifest idea, here we'll start from a simple implementation. The manifest file format will be same to the yml file format, mentioned in [configuration inspection](#configuration-inspection)

- vic-machine inspect --conf --manifest filename will dump all configurations in yml file format to the manifest file.
- vic-machine create --manifest will load the manifest file and create VCH based on it.
- vic-machine reconfigure --manifest will load the manifest file and update VCH configuration based on it.

This manifest file approach will simplify vic-machine create/reconfigure workflow, and leverage user readable configuration format. But it will be an additional work for existing vic-machine create.
