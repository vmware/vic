
# Container Create w/ Volumes


### Docker Personality

Changes to this portion will involve interacting with the new planned changes to the portlayer.

    It will involve determining whether a volume exists or needs to be created and then asking the portlayer to perform a join operationon the obtained volume. The changes here will mainly be looking for the new inputs to the `docker create` command that we want to support. This call is really managing and informing a call to the VolumeCreate call in the case of anonymous volume creation, if the volume already exists then a previous `VolumeCreate` can be assumed and only a joing will be necessary after fetching the existing volume's data. 

examples:

```
docker create -v "/mnt/pnt:<some basic opts>"

and 

docker create -v "name:/mnt/pnt:<some basic opts>"
```

    We will need to look at and do some validation on the two new potential inputs. One thing to not is that calls based on this method of volume creation means that if the volume does not exist it will be created on the default datastore specified by the vic-machine. Additionally, We will have to have a default size for the volume driver args when an anonymous volume is specified possibly at VCH deployment time. see https://github.com/vmware/vic/pull/1152 for additional discussion on this.
    The anonymous volume usage will be required to use both the default Capacity and default VolumeStore since docker cli does not support passing more args into the container create call. If either of these defaults are not set then anonymous volumes by definition will not be permitted.
    
#### Inputs

+ **mount path** is the destination of where the vdmk will be mounted inside the container. _This is required_ if it is the only value set the user is specifying an anonymous volume and we must generate our own name for the volume and populate that in the volume metadata. 


+ **name** is the value that will be listed a as the name of the volume and the md5 sum of this name will be used as the label for the block device and the target of the mount(portlayer join operation). If this is specified it must be validated. Additionally, we should look out for an already exists error. 

+ **general args** are as follows [rw|ro], [z|Z], [[r]shared|[r]slave|[r]private], and [nocopy]. These should be parsed and placed into the DriverArgs that are specified to the portlayer. right now we only support rw/ro. __TODO__ we do want to research the [no copy] option.


__NOTE:__ : in MountPoint for the volume metadata(docker perspective) we need to include something that says "Mountpoint is a block device" or something along those lines.


### Join call for attaching a volume to a vm

This call, which will be implemented in the volume portion of the storage layer within the portlayer srever, will involve a config spec change. The three things needed for this call are the handle to the container, a filled volume struct, and the driver options for the device addition(such as rw/ro). We will add a value to the extraconfig->executorConfig which will append a new Mountspec for the device to be mounted. The Op type will be an "Add"


```
[]DeviceChange{
    op:Add,
    state(?):exists,
    VirtualDevice{
    file:<vmdkPath->(should come from volume struct)>
    }
}

[]Extraconfig.append{
    executorConfig:
        label:<generated on creation, should be md5 sum>
        MountPoint:<where to mount the vmdk in the container>
}

```

The function signature should look as such

```
func (v *VolumeStore) Join(container_handle *Handle, volume *Volume, diskOpts map[string]string)
```

this will be added to the VolumeStore interface that is in Faiyaz's PR # 1196 - Volumes on vSphere
