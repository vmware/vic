# Docker Machine

The _docker-machine_ mechanism is something of a de-facto standard for provisioning docker hosts; ideally we'll keep the name and the bulk of the common options to allow customers to leverage their existing knowledge, and the VIC specific options and behaviour will present to the user as a plugin.
It is possible that the value provided by keeping the _docker-machine_ naming is overshadowed by the flexibility that changing the name provides (perhaps _vic-machine_) - lacking concrete knowledge one way or another, this component is currently named docker-machine as that is the functional niche it fills.

While deployment of a Virtual Container Host is relatively simple if performed by someone with vSphere administrative credentials, conversations with customers have shown that the self-provisioning facet of docker is a significant portion of it's value. This component, in conjunction with [the validating proxy](validating_proxy.md), provides self-provisioning capabilities and the requisite delegation of authority and specification of restrictions.

[Issues relating to this component](https://github.com/vmware/vic/labels/component%2Fdocker-machine)
