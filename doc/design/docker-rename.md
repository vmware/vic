This is the design proposal for implementing “docker rename” on vSphere Integrated Containers Engine (VIC Engine).

## Design Proposal

### The problem:

Rename involves both the containerVM’s display name and the name of the storage folder (and the files in the folder) of the containerVM on the vsphere UI, as well as in the container cache of the VIC Engine. Currently both the containerVM's display name and the folder name are created during VM creation using containerName-containerID, which is done for better user experience (the vi admin would be able to find the containerVM and the storage folder easily from the vsphere UI). Renaming the VM display name on the UI can be easily achieved by using govc, which However does not update the storage folder name. In this case, the vi admin would observe inconsistent VM display name and storage folder name, and it becomes difficult the admin to reference the storage folder based on the new VM display name.

### Proposed solution:

Since vSAN requires the VM display name to be the same as the storage folder name during VM creation, we can do VM reconfigure after the VM is created and change its display name to the customer's desired name. When a customer calls **docker rename**, the VIC Engine would change the VM display name in vsphere while leaving the storage folder name untouched. In addition, to improve user experience, we provide a way for the customer to specify a naming template with vic-machine, which can be used to customize the VM display name and folder name so that the customer could reference these names easily even if the new display name is inconsistent with the folder name after **docker rename**. It is thus the customer's responsibility to specify an appropriate naming template.

This involves three steps:

1. Naming template support: Add command line argument in **vic-machine create** which allows the customer to specify a naming template.
2. Reconfigure the VM to change the display name and folder name after VM creation: If a naming template is specified, this step would adjust the display name and folder name based on the template after VM creation. Otherwise, we provide default settings: The default display name would be the containerName, and the default folder name would be the containerID.
3. docker support for rename: When a customer calls **docker rename**, we update the VM display name to the new name in both the docker persona and the portlayer. Network aliases and HostConfig should also be updated.


## Testing and Acceptance Criteria

Robot scripts will be written to test the following:

1. Naming template support and VM reconfiguration:
  - If a naming template is provided during **vic-machine create**, check whether the template is applied successfully after launching a containerVM.
  - Otherwise, check the display name and folder name of a containerVM and see if these match the default settings.

2. docker support for rename:
  - The basic functionality of docker rename
  - docker-compose up when there are existing containers for the same service but the configuration and image have not been changed
  - docker-compose up –force-recreate when there are existing containers for the same service even if the configuration and image have not been changed
