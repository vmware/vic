# Docker API server

This is the portion of a Virtual Container Host that provides a Docker API endpoint for users to interact with. It is also referred to on occasion as a 'personality'. The longer term design has multiple personalities running within a single VCH, such that the same endpoint can serve mixed API versions, and the same VCH can serve multiple API flavours.

[Issues relating to this component](https://github.com/vmware/vic/labels/component%2Fdocker-api-server)
