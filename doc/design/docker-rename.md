This is the design proposal for implementing “docker rename” on vSphere Integrated Containers Engine (VIC Engine).

## Design Proposal

### The problem:

Rename involves the containerVM’s display name and the name of the storage folder (and the files in the folder) of the containerVM on the vsphere UI, the container cache of the VIC Engine and the container network configuration (e.g., network alias). Currently both the containerVM's display name and the folder name are created during VM creation using containerName-containerID, which is done for matching the container information obtained from `docker ps` to the VM displayed on the vSphere UI. Renaming the VM display name on the UI can be easily achieved by using govc, which however does not update the storage folder name. In this case, the vi admin would observe inconsistent VM display name and storage folder name, and it becomes difficult the admin to reference the storage folder based on the new VM display name.

### Proposed solution:

- We still use containerName-containerID to assemble the VM name. However, we only use the containerID to set the name of the VM storage folder; therefore, there is no need to worry about the VM display name and storage folder name being inconsistent. 

- VM Reconfiguration: Since vSAN requires the VM display name to be the same as the storage folder name during VM creation, we set both the VM display name and the storage folder name to containerID during VM creation, and then start a VM reconfiguration task after the VM is created to change its display name to containerName-containerID. 

- Docker support for rename: When a customer calls `docker rename`, we update the VM display name to the new name in both the docker persona and the portlayer. 

  - Network: 

    - Network alias should be updated.
    - If `--link` is used when creating the container, HostConfig of relevant containers should be updated.
    - The containerName shown in `/etc/host` within the container should be updated. Proposal: Remove containerName from `/etc/host`; docker does not put the containerName in this file. For backward compatibility, if we rename a container which is created with containerName in its `/etc/host`, we use the VMware RPC API's guestinfo facility.
  
  - Storage: Nothing needs to be updated if we set the storage folder name to containerID.
  
## Testing and Acceptance Criteria

Robot scripts will be written to test the following:

1. VM reconfiguration:
  - After a containerVM is created, use govc to check the display name (containerName-containerID) and storage folder name (containerID).

2. Docker support for rename:
  - The basic functionality of `docker rename`
  - Check validity of network alias and HostConfig 
  - `docker-compose up` when there are existing containers for the same service but the configuration or image has been changed
  - `docker-compose up –force-recreate` when there are existing containers for the same service even if the configuration or image has not been changed
  
3. Backward compatibility
  - Add a test case in the upgrade test. The old VCH would create a container with containerName in its `/etc/host`. After upgrading the VCH, check whether `docker rename` updates `/etc/host` with the new name
