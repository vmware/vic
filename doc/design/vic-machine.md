On the care and feeding of VCHs - vic-machine

vic-machine is both the installer and management client for Virtual Container Hosts.


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

VCH user (existing or new) **
resource lists:
 pool **
 imagestore datastore paths **
 contianer datastore paths **
 volume datastore paths (restriction)
 network mappings:
  one network minimum for VCH comms**
  other network mappings
resource allotments:
 cpu
 memory
 network IO
 disk IO
 datastore quotas (per datastore path)
certificates
 users - for access to VCH
 hosts - for container access to external hosts
 network - for VCH/container access to networks (gating proxies)
registry lists
 whitelist
 blacklist
default container resource reservations and limits *
containerVM naming convention (displayName for vSphere) *

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

## Managing a VCH

* report VCH status and information (API endpoint, log server, et al)
* update VCH configuration - implies possible restart of component
* shutdown/reboot VCH
* upgrade VCH - should have an entirely separate doc for this