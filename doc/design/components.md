# vSphere Integrated Containers - Architecture and Design

## Component Architecture

This is a component architecture for the system, encompassing some of the main control flows. It is intended to be used as a visual map for which components interact and how those interaction occur; it does not indicate what those interaction are.

![system component architecture](https://github.com/vmware/vic/blob/master/doc/design/component_architecture.svg)


## Container
### Container Base

ContainerVMs are bootstrapped from a PhotonOS based liveCD, containing just enough to boot linux and set up the container filesystem, before performing a switch_root into the container root filesystem. The end result is a VM:
* running the PhotonOS kernel, with appropriate kernel modules for demand loading
* the specified container filesystem mounted as `/`
* [a custom init binary](#tether) that provides the command & control channel for container interaction

The contianer process runs as root with full privileges, however there's no way to make persistent changes to anything except the container filesystem - the core operating files are read-only on the ISO and refresh each time the container is started.

[Issues relating to this component](https://github.com/vmware/vic/labels/component%2Fcontainer-base)


### Tether

The tether is an init replacement used in containerVMs that provides the command & control channel necessary to perform any operation inside the container. This includes launching of the container process, setting of environment variables, configuration of networking, etc. The tether is currently based on a modified SSH server tailed specifically for this purpose.

[Issues relating to this component](https://github.com/vmware/vic/labels/component%2Ftether)


### Container Logging

Container logging, as distinct from other logging mechanisms, is the capture mechanism for container output. This fills the same niche as [Docker log drivers](https://docs.docker.com/engine/admin/logging/overview/) and has the potential to be a direct consumption of docker log drivers.

As of the v0.1 release logs are being persisted on the datastore along with the VM.

[Issues relating to this component](https://github.com/vmware/vic/labels/component%2Fcontainer-logging)


## Appliance
### Appliance Base

The appliance VM is bootstrapped from a PhotonOS based liveCD that embeds elements of VIC relevent to Virtual Container Host functions. The appliance is diskless _in essence_, but may well use a non-persistent disk to cache transient data such as docker images in-flight from [Docker Hub](https://hub.docker.com/) but not yet persisted to a datastore.

The ISO used to bootstrap the appliance is generated from the following make targets - only the first must be specified as the others are pre-requisites:
```
make appliance appliance-staging iso-base
```

The intent behind having the appliance be diskless and bootstrap off an ISO each time is:
* improved security - coupled with [the configuration mechanism](configuration.md#Configuration-persistence-mechanism) there's no way to make persistent changes without vSphere API access
* robustness - a reboot will return the appliance to a known good state, assuming an administrator has not altered the configuration in the meantime
* simplicity of update - using [vic-machine](#vic-machine) to update a VCH should be as simple as pointing the appliance VM at a new ISO version and rebooting, so long as there's no metadata migration needed; in that case the migration should be viable as a background task prior to VCH update.

[Issues relating to this component](https://github.com/vmware/vic/labels/component%2Fappliance-base)


### vicadmin

This is more speculative than any of the other components at this point. We fully expect there to be a need for user level inspection/administration of a deployed Virtual Container Host, however we've not yet identified the functions this should provide.

Current list of functions:
* log collection

Speculative list of functions (via docker-machine as a client?):
* docker API user management
* reconfigure operations (e.g. add --insecure-registry)

- [ ] Add authentication around the server - local system or full PAM
- [ ] Retrieve client certificate from VCH when using TLS

[Issues relating to this component](https://github.com/vmware/vic/labels/component%2Fvicadmin)


### Docker API server

This is the portion of a Virtual Container Host that provides a Docker API endpoint for users to interact with; it is also referred to as a 'personality'. The longer term design has multiple personalities running within a single VCH, such that the same endpoint can serve mixed API versions, and the same VCH can serve multiple API flavours.

As of the v0.1 release this makes use of the [docker engine-api](https://github.com/docker/engine-api) project to ensure API compatibility with docker.

[Issues relating to this component](https://github.com/vmware/vic/labels/component%2Fdocker-api-server)


### imagec

The component follows the naming pattern introduced by OCI with 'runc' and is a docker registry client library, purely concerned with the pull/push aspects of the registry with login becoming a transitive dependency.

[Current open issues relating to this component](https://github.com/vmware/vic/labels/component%2Fimagec)


### Port Layer
#### Port Layer - Execution

This component deals handles management of containers such as create, start, stop, kill, etc, and is broken distinct from interaction primarily because the uptime requirements may be different, [see Interaction](#portlayer-interaction).

[Issues relating to this component](https://github.com/vmware/vic/labels/component%2Fportlayer%2Fexecution)


#### Port Layer - Interaction

This component deals handles interaction with a running container and is broken distinct from execution primarily because the uptime requirements may be different.

If the execution portion of the port layer is unavailable then only container management operations around creating, starting, stopping, et al are impacted, and then only for as long as the VCH is unavailable.

If the interaction portions are unavailable it impacts ongoing use of interactive sessions and potentially loses chunks of the container output (unless serialized to vSphere infrastructure as an intermediate step - [container logging](#container-logging), [container base](#container-base), and [tether](#tether)) are the other components that factor into the log persistence discussion).

[Issues relating to this component](https://github.com/vmware/vic/labels/component%2Fportlayer%2Finteraction)


#### Port Layer - Networking

A Virtual Container Host presents a cluster of resource as a single endpoint, so while we must also supply transparent cross host networking, it must be there without requiring any user configuration.
In conjunction with [the provisioning workflows](#vic-machine) it should also allow mapping of specific vSphere/NSX networks into the docker network namespace, and mapping of existing network entities (e.g. database servers) into the docker container namespace with defined aliases.

Initial design and implementation details for MVP are [here](networking/MVPnetworking.md).

[Issues relating to this component](https://github.com/vmware/vic/labels/component%2Fportlayer%2Fnetwork)


#### Port Layer - Storage

This provides the storage manipulation portions of the port layer, including container image storage, layering along with volume creation and manipulation. [imagec](#imagec) uses this component to translate registry images into a layered format the can be used directly by vSphere, namely VMDK disk chains.

[Issues relating to this component](https://github.com/vmware/vic/labels/component%2Fportlayer%2Fstorage)


## Install and management

### vic-machine

The _docker-machine_ mechanism is something of a de-facto standard for provisioning docker hosts.
There is a significant amount of vSphere specific behaviour that needs to be expose, and that may well go beyond what docker-machine provides - hence supplying a vic-machine binary.

Ideally we'll provide vic-machine in a form that makes it viable as a docker-machine plugin, allowing some reuse of existing knowledge, with VIC specific options and behaviours presented to the user via plugin options.
It is possible that the value provided by keeping the _docker-machine_ naming is overshadowed by the flexibility that changing the name provides (perhaps _vic-machine_) - lacking concrete knowledge one way or another, this component is currently named vic-machine so as to avoid confusion with the docker binary.

While deployment of a Virtual Container Host is relatively simple if performed by someone with vSphere administrative credentials, conversations with customers have shown that the self-provisioning facet of docker is a significant portion of it's value. This component, in conjunction with [the validating proxy](#validating_proxy), provides self-provisioning capabilities and the requisite delegation of authority and specification of restrictions.

[Issues relating to this component](https://github.com/vmware/vic/labels/component%2Fvic-machine)


### Validating Proxy

The self-provisioning workflow for vSphere Integrated Containers is an authority delegation and resource restriction model. That requires that there be an active endpoint accessible to both the user and the viadmin that is capable of generating delgation tokens to be passed to the user, and validating those that are received from a user. The validating proxy fills this niche; the described proxy is very, very simple and does not include elements such as directory services integration.

[Issues relating to this component](https://github.com/vmware/vic/labels/component%2Fvalidating-proxy)



## ESX Agents

### VMOMI Authenticating Agent

As with [access to NSX management](#nsx-authenticating-agent), vSphere orchestration requires credentials and therefore credential management. The authenticating agent is indended to move the credential management problem out of the VCH and into the vSphere infrastructure. Authorizing the VCH to perform specific operations by virtue of being a specific VM (and potentially cryptographically verified), rather than generating and embedding credentials into the guest, means:
* a VCH can be safely considered untrusted
* IDS style inspection and validation can be performed on each infrastructure operation performed by a VCH
* no access to management networks is required for infrastructure orchestration


[Issues relating to this component](https://github.com/vmware/vic/labels/component%2Fvmomi-authenticating-agent)


# vSocket Relay Agent

Network serial ports as a communication channel have several drawbacks:
* serial is not intended for high bandwidth, high frequency data
* inhibits forking & vMotion without vSPC, and a vSPC requires an appliance in FT/HA configuration
* requires a VCH have a presence on the management networks
* requires opening a port on the ESX firewall

The alternative we're looking at is vSocket (uses PIO based VMCI communication), however that it Host<->VM only so we need a mechanism to relay that communication to the VCH. Initially it's expected that the Host->VCH communication still be a TCP connection for a staged delivery approach, with the longer term being an agent<->agent relay between the two hosts.

[Issues relating to this component](https://github.com/vmware/vic/labels/component%2Fvsocket-relay-agent)


### NSX Authenticating Agent

As with [access to vSphere](#vmomi-authenticating-agent), NSX management interaction requires credentials and therefore credential management. The authenticating agent is indended to move the credential management problem out of the VCH and into the vSphere infrastructure. Authorizing the VCH to perform specific operations by virtue of being a specific VM (and potentially cryptographically verified), rather than generating and embedding credentials into the guest, means:
* a VCH can be safely considered untrusted
* IDS style inspection and validation can be performed on each infrastructure operation performed by a VCH
* no access to management networks is required for infrastructure orchestration

[Issues relating to this component](https://github.com/vmware/vic/labels/component%2Fnsx-authenticating-agent)
