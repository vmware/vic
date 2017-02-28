This is the design proposal for implementing “docker rename” on vSphere Integrated Containers Engine (VIC Engine).

## Design Proposal

### The problem:

Rename involves the containerVM’s display name and the name of the datastore folder (and the files in the folder) of the containerVM on the vsphere UI, the container cache of the VIC Engine and the container network configuration (e.g., network alias). Currently both the containerVM's display name and the folder name are created during VM creation using containerName-containerID, which is done for matching the container information obtained from `docker ps` to the VM displayed on the vSphere UI. Renaming the VM display name on the UI can be achieved by using govc, which however does not update the datastore folder name. In this case, the vi admin would observe inconsistent VM display name and datastore folder name, and it becomes difficult the admin to reference the datastore folder based on the new VM display name.

### Proposed solution:

- We use containerName-containerShortID to assemble the VM display name. We do not use containerName-containerID in order to avoid the scenario wherein the containerName gets truncated to satisfy the maximum length of a VM display name in vSphere. In addition, we use the containerID to set the name of the datastore folder, thus there is no need to worry about the VM display name and datastore folder name being inconsistent. 

- VM Reconfiguration: Since vSAN requires the VM display name to be the same as the datastore folder name during VM creation, we set both the VM display name and the datastore folder name to containerID during VM creation, and then start a VM reconfiguration task after the VM is created to change its display name to containerName-containerShortID. 

- Docker support for rename: When a customer calls `docker rename`, we update the VM display name to the new name in both the docker persona and the portlayer. 

  - Network: 

    - Network alias should be updated.
    - If `--link` is used when creating the container, HostConfig of relevant containers should be automatically updated based on the backend data.
    - The containerName shown in `/etc/hosts` within the container should be updated. 
      - Proposal: Remove containerName from `/etc/hosts`; docker does not put the containerName in this file. 
      - Backward compatibility: If the user calls `docker rename` on an existing container (created by a VCH of an older version) after VCH upgrade, 
        - If the container is already powered off, we trigger a `reload` in tether, which should cause the network config to be inspected and changes to be applied when the container is powered on again.
        - If the container is still powered on, since it is not created with the new tether binary, `docker rename` will only update DNS but not `/etc/hosts`, and we will throw a warning message about the mismatch between the new containerName and the name in `/etc/hosts`. 
          
  - Storage: Nothing needs to be updated if we set the datastore folder name to containerID.

####Note on `reload`

We have two approaches to trigger `reload`:
- (a) The portlayer sends a `HUP` signal to the tether process (pid=1) via `startGuestProgram`, which then triggers `reload` in the signal handler
- (b) We add a new case to `startGuestProgram` as the `reload` command in the portlayer. Then the portlayer sends the `reload` command to the toolbox via the guest ProcessManager. The toolbox needs to implement its own handler for the `reload` command

The benefit of (a) is that it provides a single path for handling `reload` triggered from in-guest and out-of-guest, while the call path of (b) is simpler. We pick (a) in our implementation.

  
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
  - Add a test case in the upgrade test. The old VCH would create a container with containerName in its `/etc/hosts`. Then we upgrade the VCH and call `docker rename`. For an existing container that is powered off, `/etc/hosts` should contain the new containerName after it is powered on again. For an existing container that is still powered on, `/etc/hosts` is not updated.
