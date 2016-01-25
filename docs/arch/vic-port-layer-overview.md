#### Port Layer Abstractions

The Port Layer abstractions in VIC are designed to augment the vSphere APIs with low-level container primitives from which a simple container engine could be implemented. The design criteria of the Port Layer is as follows:

* The Port Layer should be primarily oriented around the notion of _isolation domains_. It should provide the means to easily express rich and flexible criteria for isolating containers and their resources, without being explicit about the mechanism through which this should be achieved.
* The Port Layer is designed to be invoked by higher-level software abstraction. It is not designed to be exposed directly to users.
* The Port Layer should be developed as Open Source Software to allow for 3rd party integration
* The Port Layer should be container engine and operating system agnostic
* The Port Layer should be designed in such a way as to optimize control plane performance
* The Port Layer should ensure a single source of truth for all state. Eg. VM power-off == container stop

The Port Layer APIs are organized into 5 discreet domains:

**VCH**

Deals with the creation and lifecycle management of a Virtual Container Host. A VCH represents:

* An isolation domain and control plane endpoint for a single tenant
* A dynamically-configurable resource boundary for containerVMs
* A containerVM bootstrap ISO and VM/kernel configuration defaults for all provisioned containerVMs

A VCH is not explicity tied to a single image cache (see Storage), but is limited to a single operating-system type by the bootstrap ISO associated with it. Images deployed via the VCH endpoint must therefore be of a compatible filesystem type.

A VCH is limited to a single tenant, but may be connected to by multiple users.

The creation of a VCH requires access to resources typically controlled by a vSphere admin, but management of VCH creation and lifecycle should be possible by dev and ops personas. This is done at a high-level by the vSphere admin creating a binary token representing specific system resources which can then be used as input to the VCH creation. 

We plan to deliver VCH creation via a Docker Machine plugin

**Execution**



