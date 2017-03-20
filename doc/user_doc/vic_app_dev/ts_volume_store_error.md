#Default Volume Store Error

When you create or run a container, the Docker operation fails with an error about a missing volume store.  

##Problem
Running the container fails with error:

	docker: Error response from daemon: No volume store named (default) exists.

## Cause

By default, vic-machine create does not create a volume store when you or the vSphere administrator deploy a VCH. To run containers from images that use volumes, you or the vSphere administrator must specify a volume store named default when deploying the VCH.

## Solution

Deploy a VCH by using the `vic-machine create --volume-store` option to create a VCH with a volume store named `default`. See `--volume-store` in VCH Deployment Options and [Specify Volume Stores](../vic_installation/vch_installer_examples.html#volume-stores) in Advanced Examples of Deploying a VCH in *vSphere Integrated Containers Installation*.

Use `docker volume inspect` to get information about the volume. 