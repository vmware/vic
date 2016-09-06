# Using Volumes with vSphere Integrated Containers Engine #

vSphere Integrated Containers Engine supports the use of container volumes. When you create or the vSphere Administrator creates a virtual container host, you or the Administrator specify the datastore to use to store container volumes in the `vic-machine create --volume-store` option. For information about how to use the `vic-machine create --volume-store` option, see the section on `volume-store` in [Virtual Container Host Deployment Options](../vic_installation/vch_installer_options.html#volume-store) in *vSphere Integrated Containers Engine Installation and Configuration*.   

## Obtain the List of Available Volume Stores ##

To obtain the list of volume stores that are available on a virtual container host, run `docker info`.

<pre>docker -H <i>virtual_container_host_address</i>:2376 --tls info</pre>

The list of available volume stores for this virtual container host appears in the `docker info` output under `VolumeStores`.

<pre>[...]
Storage Driver: vSphere Integrated Containers Backend Engine
VolumeStores: <i>volume_store_1</i> <i>volume_store_2</i> ... <i>volume_store_n</i>
vSphere Integrated Containers Backend Engine: RUNNING
[...]</pre>

## Create a Volume in a Volume Store ##

When you use the `docker volume create` command to create a volume, you can optionally provide a name for the volume by specifying the `--name` option. If you do not specify `--name`, `docker volume create` assigns a random UUID to the volume.

- If the volume store label is anything other than `default`, you must specify the `--opt VolumeStore` option and pass the name of an existing volume store to it. If you do not specify `--opt VolumeStore`, `docker volume create` searches for a volume store named `default`, and returns an error if no such volume store exists. 

  <pre>docker -H <i>virtual_container_host_address</i>:2376 --tls volume create 
--opt VolumeStore=<i>volume_store_label</i> 
--name <i>volume_name</i></pre>

- If you or the vSphere Administrator set the volume store label to `default` when running `vic-machine create`, you do not need to specify `--opt VolumeStore`.

  <pre>docker -H <i>virtual_container_host_address</i>:2376 --tls volume create 
--name <i>volume_name</i></pre>

- If you intend to create anonymous volumes by using `docker create -v`, a volume store named `default` must exist. In this case, you include the path to the destination at which you want to mount an anonymous volume in the `docker create -v` command. Docker creates the volume in the `default` volume store, if it exists.

  <pre>docker -H <i>virtual_container_host_address</i>:2376 --tls create 
-v <i>destination_path_for_anonymous_volume</i> busybox</pre>

  **NOTE**: If you use `docker create -v`, vSphere Integrated Containers Engine only supports the `-r` and `-rw` options.

- You can optionally set the capacity of a volume by specifying the `--opt Capacity` option when you run `docker volume create`. If you do not specify the `--opt Capacity` option, the volume is created with the default capacity of 1024MB. 

  If you do not specify a unit for the capacity, the default unit will be in Megabytes.
  <pre>docker -H <i>virtual_container_host_address</i>:2376 --tls volume create 
--opt VolumeStore=<i>volume_store_label</i> 
--opt Capacity=2048
--name <i>volume_name</i></pre>
- To create a volume with a capacity in megabytes, gigabytes, or terabytes, include `MB`, `GB`, or `TB` in the value that you pass to `--opt Capacity`. The unit is case insensitive.

  <pre>docker -H <i>virtual_container_host_address</i>:2376 --tls volume create 
--opt VolumeStore=<i>volume_store_label</i> 
--opt Capacity=10GB
--name <i>volume_name</i></pre>


**NOTE**: When using a vSphere Integrated Containers Engine virtual container host as your Docker endpoint, the storage driver is always the vSphere Integrated Containers Engine Backend Engine. If you specify the `docker volume create --driver` option, it is ignored.  

## Obtain the List of Available Volumes ##

To obtain a list of volumes that are available on a virtual container host, run `docker volume ls`.

<pre>docker -H <i>virtual_container_host_address</i>:2376 --tls volume ls

DRIVER         VOLUME NAME
vsphere        <i>volume_1</i>
vsphere        <i>volume_2</i>
[...]          [...]
vsphere        <i>volume_n</i></pre>

## Delete a Named Volume from a Volume Store ##
To delete a volume, run `docker volume rm` and specify the name of the volume to delete.
<pre>docker -H <i>virtual_container_host_address</i>:2376 --tls 
volume rm <i>volume_name</i></pre>

**NOTE**: In the current builds, `docker volume rm` is not yet supported.
