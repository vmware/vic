# Port Layer - Interaction

This component handles interaction with a running container and is distinct from execution primarily because the uptime requirements may be different.

If the execution portion of the port layer is unavailable then only container management operations around creating, starting, stopping, et al are impacted, and then only for as long as the VCH is unavailable.

If the interaction portions are unavailable it impacts ongoing use of interactive sessions and potentially loses chunks of the container output (unless serialized to vSphere infrastructure as an intermediate step - [container logging](container-logging.md), [container base](container-base.md), and [tether](tether.md)) are the other components that factor into the log persistence discussion).

[Issues relating to this component](https://github.com/vmware/vic/labels/component%2Fportlayer%2Finteraction)
