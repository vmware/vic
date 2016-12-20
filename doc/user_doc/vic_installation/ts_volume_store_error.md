#Default Volume Store Error

When you install a virtual container host (VCH) with vic-machine, the VCH fails to create a default volume store.  

##Problem
Running the application fails with error:

	docker: Error response from daemon: No volume store named (default) exists.

## Cause

By default, a volume store is not required when you install a VCH with vic-machine. To use containers that require volumes, you must specify a volume store with vic-machine create so that it is available for the containers to use.

## Solution

Install the VCH and specify a volume store. See [Specify One or More Volume Stores](vch_installer_examples.md#volume-stores) in *Advanced Examples of Deploying a VCH*.

Use `docker volume inspect` to get information about the volume. 