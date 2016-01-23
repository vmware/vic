### vSphere Integrated Containers Architecture

#### Overview

VIC is a product designed to tightly integrate container workflow, lifecycle and provisioning with the vSphere SDDC. In VIC, a container is a first-class citizen provisioned into a _Virtual Container Host_ (VCH) and able to directly integrate with vSphere infrastructure capabilities, such as networking and storage features.

[Learn more about the differences between the VIC model and a traditional software-virtualized container](docs/arch/vic-container-abstraction.md)

The architecture of VIC is designed to allow for significant modularity and flexibility and includes the following key components:

##### Port Layer Abstractions

vSphere currently lacks the notion of container primitives and abstractions through which they can be manipulated. It has a rich API with bindings for various languages (Eg. [govmomi](https://github.com/vmware/govmomi)) but these are all necessarily oriented around the notion of a VM. 

While it would be possible to write a rudimentary VIC-like container engine by driving the vSphere APIs directly from within a daemon of some kind, the tight coupling between the low-level vSphere calls and the high-level daemon API would result in very little re-usable code and monolith that's potentially difficult to maintain. An API layer that encapsulates low-level container primitives that is both container engine and operating system agnostic would be preferable.

A secondary benefit of such an API is that it could easily be extended for compatibility with emerging standards which operate at a similar layer, such as [runc](https://github.com/opencontainers/runc).








