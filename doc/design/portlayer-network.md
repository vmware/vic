# Port Layer - Networking

A Virtual Container Host presents a cluster of resource as a single endpoint, so while we must also supply transparent cross host networking, it must be there without requiring any user configuration.
In conjunction with [the provisioning workflows](docker-machine.md) it should also allow mapping of specific vSphere/NSX networks into the docker network namespace, and mapping of existing network entities (e.g. database servers) into the docker container namespace with defined aliases.

[Issues relating to this component](https://github.com/vmware/vic/labels/component%2Fportlayer%2Fnetwork)
