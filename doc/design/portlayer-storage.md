# Port Layer - Storage

This provides the storage manipulation portions of the port layer, including container image storage, layering along with volume creation and manipulation. [imagec](imagec.md) uses this component to translate registry images into a layered format the can be used directly by vSphere, namely VMDK disk chains.

[Issues relating to this component](https://github.com/vmware/vic/labels/component%2Fportlayer%2Fstorage)
