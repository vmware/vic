# Using Volumes with vSphere Integrated Containers Engine #

vSphere Integrated Containers Engine supports the use of container volumes. When you create or the vSphere Administrator creates a virtual container host (VCH), you or the Administrator specify the datastore to use to store container volumes in the `vic-machine create --volume-store` option. For information about how to use the `vic-machine create --volume-store` option, see the section on `volume-store` in [VCH Deployment Options](../vic_installation/vch_installer_options.html#volume-store) in *vSphere Integrated Containers Engine Installation*.  

- [Obtain the List of Available Volume Stores](#list_vs) 
- [Obtain the List of Available Volumes](#list_vols)
- [Create a Volume in a Volume Store](#create_vol)
- [Creating Volumes from Images](#image_volumes)
- [Create a Container and Attach it to an Anonymous or Named Volume](#create_container)
- [Attach an Existing Volume to a Container](#attach)
- [Obtain Information About a Volume](#inspect_vol) 
- [Delete a Named Volume from a Volume Store](#delete_vol) 

For simplicity, the examples in this topic assume that the VCHs implement TLS authentication with self-signed untrusted certificates, with no client verification.

<a name="list_vs"></a>
## Obtain the List of Available Volume Stores ##

To obtain the list of volume stores that are available on a VCH, run `docker info`.

<pre>docker -H <i>virtual_container_host_address</i>:2376 --tls info</pre>

The list of available volume stores for this VCH appears in the `docker info` output under `VolumeStores`.

<pre>[...]
Storage Driver: vSphere Integrated Containers Backend Engine
VolumeStores: <i>volume_store_1</i> <i>volume_store_2</i> ... <i>volume_store_n</i>
vSphere Integrated Containers Backend Engine: RUNNING
[...]</pre>

<a name="list_vols"></a>
## Obtain the List of Available Volumes ##

To obtain a list of volumes that are available on a VCH, run `docker volume ls`.

<pre>docker -H <i>virtual_container_host_address</i>:2376 --tls volume ls

DRIVER         VOLUME NAME
vsphere        <i>volume_1</i>
vsphere        <i>volume_2</i>
[...]          [...]
vsphere        <i>volume_n</i></pre>

<a name="create_vol"></a>
## Create a Volume in a Volume Store ##

When you use the `docker volume create` command to create a volume, you can optionally provide a name for the volume by specifying the `--name` option. If you do not specify `--name`, `docker volume create` assigns a random UUID to the volume.

- If the volume store label is anything other than `default`, you must specify the `--opt VolumeStore` option and pass the name of an existing volume store to it. If you do not specify `--opt VolumeStore`, `docker volume create` searches for a volume store named `default`, and returns an error if no such volume store exists. 

  <pre>docker -H <i>virtual_container_host_address</i>:2376 --tls volume create 
--opt VolumeStore=<i>volume_store_label</i> 
--name <i>volume_name</i></pre>

- If you or the vSphere Administrator set the volume store label to `default` when running `vic-machine create`, you do not need to specify `--opt VolumeStore`.

  <pre>docker -H <i>virtual_container_host_address</i>:2376 --tls volume create 
--name <i>volume_name</i></pre>

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

After you create a volume by using docker volume create, you can attach it to a container by running either of the following commands:

<pre>docker -H <i>virtual_container_host_address</i>:2376 --tls 
create -v /<i>volume_name</i> busybox</pre>
<pre>docker -H <i>virtual_container_host_address</i>:2376 --tls 
run -v /<i>volume_name</i> busybox</pre>

**NOTE**: When using a vSphere Integrated Containers Engine VCH as your Docker endpoint, the storage driver is always the vSphere Integrated Containers Engine Backend Engine. If you specify the `docker volume create --driver` option an error stating that a bad driver has been selected will occur.

<a name="image_volumes"></a>
## Creating Volumes from Images ##

Some images, for example, `mongo` or `redis:alpine`, contain volume bind information in their metadata. vSphere Integrated Containers Engine creates such volumes with the default parameters and treats them as anonymous volumes. vSphere Integrated Containers Engine treats all volume mount paths as unique, in the same way that Docker does. This should be kept in mind if you attempt to bind other volumes to the same location as anonymous or image volumes. A specified volume always takes priority over an anonymous volume.

If you require an image volume with a different volume capacity to the default, create a named volume with the required capacity. You can mount that named volume to the location that the image metadata specifies. You can find the location by running `docker inspect image_name` and consulting the `Volumes` section of the output. The resulting container has the required storage capacity and the endpoint.  

<a name="create_container"></a>
## Create a Container and Attach it to an Anonymous or Named Volume ##

If you intend to create named or anonymous volumes by using `docker create -v` when creating containers, a volume store named `default` must exist in the VCH. In this case, you include the path to the destination at which you want to mount an anonymous volume in the `docker create -v` command. Docker creates the anonymous volume in the `default` volume store, if it exists. The VCH attaches the anonymous volume to the container.

For example, to create a busybox container that is mounted to the `volumes` folder of an anonymous volume in the default volume store, run the following command:

<pre>docker -H <i>virtual_container_host_address</i>:2376 --tls 
create -v /volumes busybox</pre>

You can create containers that are attached to named volumes by using `docker create -v` and specifying a volume name. When you create containers that are attached to named volumes, the VCH checks whether the volume exists in the volume store, and if it does not, creates it. The VCH attaches the existing or new volume to the container.

For example, to create a busybox container that is mounted to the `volumes` folder of a volume named `volume_1` in the default volume store with default capacity, run the following command:

<pre>docker -H <i>virtual_container_host_address</i>:2376 --tls 
create -v volume_1:/volumes busybox</pre>

**NOTES**: 
- vSphere Integrated Containers Engine does not support mounting directories as data volumes. A command such as <code>docker create -v /<i>folder_name</i>:/<i>folder_name</i> busybox</code> is not supported.
- If you use `docker create -v` to create containers that are attached to volumes, vSphere Integrated Containers Engine only supports the `-r` and `-rw` options.

<a name="attach"></a>
## Mount an Existing Volume on a Container ##
vSphere Integrated Containers Engine currently supports mounting a volume on only one container at a time. When you mount a volume on a container by using `docker create -v`,  that volume remains mounted on the container until you remove that container. When you have removed the container you can mount the volume on a new container.

This example performs the following operations:

- Creates a container named `container1` from the `busybox` image.
- Mounts the `myData` folder of a volume named `volume1` on that container, starts the container, and attaches to it.
- After performing operations in `volume1:/myData` then stopping and detaching `container1`, creates `container2` from the `ubuntu` image and mounts the `myData` folder of `volume1` on it.

<pre>docker -H <i>virtual_container_host_address</i>:2376 --tls 
create --name container1 -v volume1:/myData busybox
docker start container1
docker attach container1 

[Perform container operations and detach]

docker stop container1 
docker rm container1
docker create -it --name container2 -v volume1:/myData ubuntu
docker start container2 
docker attach container2 

[Perform container operations with the same volume that was 
previously mounted to container1]</pre>

<a name="inspect_vol"></a>
## Obtain Information About a Volume ##
To get information about a volume, run `docker volume inspect` and specify the name of the volume.
<pre>docker -H <i>virtual_container_host_address</i>:2376 --tls 
volume inspect <i>volume_name</i></pre>

<a name="delete_vol"></a>
## Delete a Named Volume from a Volume Store ##
To delete a volume, run `docker volume rm` and specify the name of the volume to delete.
<pre>docker -H <i>virtual_container_host_address</i>:2376 --tls 
volume rm <i>volume_name</i></pre>


