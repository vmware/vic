This is the design proposal for implementing “docker rename” on vSphere Integrated Containers Engine (VIC Engine).

## Design Proposal

The problem:
Rename involves changing the containerVM’s name on the vsphere UI and the name of the storage folder (and the files in the folder) of the containerVM.
Currently both the VM name and the folder name are created using containerName-containerID, which is done for better user experience (the vi admin would be able to find the containerVM and the storage folder easily from the vsphere UI).
Renaming the VM display name on the UI can be easily achieved by using govc.
However, updating the VM folder name with new container name requires powering off the VM (see here https://kb.vmware.com/selfservice/microsites/search.do?language=en_US&cmd=displayKC&externalId=1029513).
Risks might also arise in the underlying vSAN if the storage folder is renamed.

### Proposed strategies:

|     |  S1  |  S2  |  S3  |
|  ---  |  ---  |  ---  |  ---  |
|  VM name on vSphere UI  |  Update using govmomi  |
|  VM folder name  |  Do not update. The vi admin can easily find the VM folder path/name from the UI of the containerVM.  |  (1) Add the new name to the container metadata while keeping the old name. (2) docker inspect should show both the original name and the new name so that the vi admin/developers can easily for the VM folder using the original name. (3) When the container exits,  shuts down or restart, the portlayer commits the new name so that the VM folder name can be updated.  |  When creating the containerVM, only use containerID as the folder name. Therefore, there is no need to update the VM folder name during docker rename.  |

### Pros and Cons:

|     |  S1  |  S2  |  S3  |
|  ---  |  ---  |  ---  |  ---  |
|  Pros  |  Easy to implement and maintain  |  (1) After portlayer commits the new name, all the names would be consistent. (2) Easier for the developers to find out the VM folder name (simply checking the original VM name via docker inspect)  |  Easy to implement and maintain  |
|  Cons  |  Difficult for the developers to figure out the VM folder |  (1) Before portlayer commits the new name, the vi admin/developers still see inconsistent VM name and folder name. (2) Complex workflow (still need to figure out how to update the VM folder name). (3) The output of docker inspect on vic would be different from that on regular docker. (4) vSAN risk.  |  (1) It will not be straightforward for the vi admin to find the VM folder; need to either (a) look up the containerID from cli or (b) check the VM information from the vsphere UI. (2) Difficult for the developers to figure out the VM folder  |

** Note that the additional steps in implementing docker rename are not discussed here, e.g.,  checking the validity of the new name, updating the network aliases (if any) and HostConfig. 
Details of these steps can be found here https://github.com/vmware/vic/issues/380#issuecomment-279069180.


## Testing and Acceptance Criteria

1.	Robot scripts will be written to test 

  - the basic functionality of docker rename, 
  - docker-compose up when there are existing containers for the same service but the configuration and image have not been changed, 
  -  docker-compose up –force-recreate when there are existing containers for the same service even if the configuration and image have not been changed
  
2.	A sanity check will also be performed by querying the VM display name using govc.
