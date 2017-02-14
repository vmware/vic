This is the design proposal for implementing “docker rename” on vSphere Integrated Containers Engine (VIC Engine).

Design Proposal

The problem:
Rename involves changing the containerVM’s name on the vsphere UI and the name of the storage folder (and the files in the folder) of the containerVM.
Currently both the VM name and the folder name are created using containerName-containerID, which is done for better user experience (the vi admin would be able to find the containerVM and the storage folder easily from the vsphere UI).
Renaming the VM display name on the UI can be easily achieved by using govc.
However, updating the VM folder name with new container name requires powering off the VM (see here https://kb.vmware.com/selfservice/microsites/search.do?language=en_US&cmd=displayKC&externalId=1029513).
Risks might also arise in the underlying vSAN if the storage folder is renamed.

Proposed strategies:

|     |  S1  |  S2  |  S3  |
|  ---  |  ---  |  ---  |  ---  |


Pros and Cons:

|     |  S1  |  S2  |  S3  |
|  ---  |  ---  |  ---  |  ---  |


** Note that the additional steps in implementing docker rename are not discussed here, e.g.,  checking the validity of the new name, updating the network aliases (if any) and HostConfig. 
Details of these steps can be found here https://github.com/vmware/vic/issues/380#issuecomment-279069180.


Testing and Acceptance Criteria
1.	Robot scripts will be written to test 
  - the basic functionality of docker rename, 
  - docker-compose up when there are existing containers for the same service but the configuration and image have not been changed, 
  -  docker-compose up –force-recreate when there are existing containers for the same service even if the configuration and image have not been changed
  
2.	A sanity check will also be performed by querying the VM display name using govc.
