# VMOMI Authenticating Agent

As with [access to NSX management](nsx-authenticating-agent.md), vSphere orchestration requires credentials and therefore credential management. The authenticating agent is intended to move the credential management problem out of the VCH and into the vSphere infrastructure. Authorizing the VCH to perform specific operations by virtue of being a specific VM (and potentially cryptographically verified), rather than generating and embedding credentials into the guest, means:

* a VCH can be safely considered untrusted
* IDS style inspection and validation can be performed on each infrastructure operation performed by a VCH
* no access to management networks is required for infrastructure orchestration


[Issues relating to this component](https://github.com/vmware/vic/labels/component%2Fvmomi-authenticating-agent)
