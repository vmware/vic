This is the design proposal for implementing “docker rename” on vSphere Integrated Containers Engine (VIC Engine).

## Design Proposal

### The problem:

Rename involves the containerVM’s display name and the name of the datastore folder (and the files in the folder) of the containerVM on the vsphere UI, the container cache of the VIC Engine and the container network configuration (e.g., network alias). Currently both the containerVM's display name and the folder name are created during VM creation using containerName-containerID, which is done for matching the container information obtained from `docker ps` to the VM displayed on the vSphere UI. Renaming the VM display name on the UI can be achieved by using govc, which however does not update the datastore folder name. In this case, the vi admin would observe inconsistent VM display name and datastore folder name, and it becomes difficult the admin to reference the datastore folder based on the new VM display name.

### Proposed solution:

- We use **containerName-containerShortID** to assemble the VM display name. We do not use containerName-containerID in order to avoid the scenario wherein the containerName gets truncated to satisfy the maximum length of a VM display name in vSphere. In addition, we use the **containerID** to set the name of the datastore folder, thus there is no need to worry about the VM display name and datastore folder name being inconsistent. 

- VM Reconfiguration: Since vSAN requires the VM display name to be the same as the datastore folder name during VM creation, we set both the VM display name and the datastore folder name to containerID during VM creation, and then start a VM reconfiguration task after the VM is created to change its display name to containerName-containerShortID. 

- Docker support for rename: When a customer calls `docker rename`, we update the VM display name to the new name in both the docker persona and the portlayer. 

  - Network: 

    - Network alias should be updated.
    - If `--link` is used when creating the container, HostConfig of relevant containers should be automatically updated based on the backend data.
    - The containerName shown in `/etc/hosts` within the container should be updated. 
      - Proposal: Remove containerName from `/etc/hosts`
      - Backward compatibility for containers created by a VCH of an older version: No matter the container is powered off or still powered on, change the network setting if needed and also use govmomi to update the name in the containerVM's `guestinfo.vice./common/name`. `/etc/hosts` will be automatically regenerated with the new containerName when the containerVM is rebooted. 
          
  - Storage: Nothing needs to be updated if we set the datastore folder name to containerID.


## Testing and Acceptance Criteria

Robot scripts will be written to test the following:

1. VM reconfiguration:
  - After a containerVM is created, use govc to check the display name (containerName-containerShortID) and datastore folder name (containerID).

2. Docker support for rename:
  - The basic functionality of `docker rename`
  - Check validity of network alias and HostConfig 
  - `docker-compose up` when there are existing containers for the same service but the configuration or image has been changed
  - `docker-compose up –force-recreate` when there are existing containers for the same service even if the configuration or image has not been changed
  
3. Backward compatibility
  - Add a test case in the upgrade test. The old VCH would create a container with containerName in its `/etc/hosts`. Then we upgrade the VCH and call `docker rename` on an existing container. `/etc/hosts` should contain the new containerName after the container is rebooted. 
