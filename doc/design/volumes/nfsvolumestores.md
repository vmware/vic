# Specification to support containers with NFS based shared volumes in VIC

Container users want to be able to access shared storage between their containers programatically.  Docker solves this by way of adding [NFS volumes](https://docs.docker.com/engine/reference/commandline/volume_create/#/driver-specific-options).  VIC can streamline this functionality for the user by abstracting this away from the container user and allowing the VI admin to configure the VCH with NFS based volumes access by way of our `VolumeStore`.  This way, the VI admin can add an NFS based VolumeStore to the VCH, and the VI admin need only create volumes on it without needing to know the details of the NFS target.

## Requirements

Allow the VI admin to
 1. add an NFSv3 based `VolumeStore`

Allow the container user to
 1. create volumes on the NFS based `VolumeStore`
 1. create 1 or greater containers with NFS based volumes at the given location in the container filesystem namespace
 1. validate the volume is no longer in use and delete it

