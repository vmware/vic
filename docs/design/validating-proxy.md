# Validating Proxy

The self-provisioning workflow for vSphere Integrated Containers is an authority delegation and resource restriction model. That requires that there be an active endpoint accessible to both the user and the viadmin that is capable of generating delgation tokens to be passed to the user, and validating those that are received from a user. The validating proxy fills this niche; the described proxy is very, very simple and does not include elements such as directory services integration.

[Issues relating to this component](https://github.com/vmware/vic/labels/component%2Fvalidating-proxy)
