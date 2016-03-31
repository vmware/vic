#### Port Layer Abstractions

The Port Layer abstractions in VIC are designed to augment the vSphere APIs with low-level container primitives from which a simple container engine could be implemented. The design criteria of the Port Layer is as follows:

* The Port Layer should be primarily oriented around the notion of _isolation domains_. It should provide the means to easily express rich and flexible criteria for isolating containers and their resources, without being explicit about the mechanism through which this should be achieved.

* The Port Layer is designed to be invoked by higher-level software abstraction. It is not designed to be exposed directly to users.

* The Port Layer should be developed as Open Source Software to allow for 3rd party integration

* The Port Layer should be container engine and operating system agnostic

* The Port Layer should be designed in such a way as to optimize control plane performance

* The Port Layer should ensure a single source of truth for all state. Eg. VM power-off == container stop

From an architectural perspective, the Port Layer should be considered functionally equivalent to a project like https://github.com/docker/libcontainer in as much as it provides low-level platform-specific primitives. It is easy to see how such an abstraction could be container engine agnostic since it provides capabilities at a much lower layer. Our goal however is that it should also be operating system agnostic, which is a more challenging goal at such a low layer. 

### Operating System Independence

VMs are already completely operating system agnostic, since they virtualize at the hardware layer and all control plane operations through the vSphere APIs are therefore also necessarily OS agnostic. Guest differences are encapsulated in different builds of "VMware Tools" which is an optional in-guest agent that mediates between the guest and the hypervisor.

The Port Layer in VIC will function in exactly the same way. Control plane operations will be expressed through an OS agnostic API and distinct differences between operating system implementations will be encapuslated in the _Tether_ process that runs in each containerVM.

### The Tether Process

A traditional container runtime, such as Linux/LXC, allows the control plane and the containers to share a kernel within a common address space. Each container gets its own private namespace, but the shared kernel allows the control plane to have visibility into the containers and also allows for processes to be started and stopped inside them. 

A containerVM by contrast uses completely separate isolated kernels for the control plane and containers. The control plane can either run in the hypervisor kernel or in a distinct guest OS kernel in a separate VM, possibly even on a separate physical host. This isolation is by design: the job of a containerVM is to run only the container process in its own kernel with as minimal a guest OS stack as feasibly possible while ensuring the same strong degree of isolation as any other VM. Even the hypervisor doesn't have visibility inside the guest without an in-guest agent installed.

As such, in order for the container control plane to provide a shell into a container, to start and stop processes or to provide monitoring statistics, there must be some kind of guest agent in the containerVM. We call this guest agent a _Tether_ process. This is not the same agent as VMware Tools, but a minimal agent designed specifically for VIC.

The Tether API and Tether codebase is where all OS differences will be encapsulated. As such, the Tether API should be considered private to the Port Layer - it exists exclusively for the benefit of the internal control plane operations, not to be invoked direclty by anything that implements the Port Layer. 






