# VICadmin

This is more speculative than any of the other components at this point. We fully expect there to be a need for user level inspection/administration of a deployed Virtual Container Host, however we've not yet identified the functions this should provide.

Current list of functions:
* log collection

Speculative list of functions (via docker-machine as a client?):
* docker API user management
* reconfigure operations (e.g. add --insecure-registry)

- [ ] Add authentication around the server - local system or full PAM 
- [ ] Retrieve client certificate from VCH when using TLS

[Issues relating to this component](https://github.com/vmware/vic/labels/component%2Fvicadmin)
