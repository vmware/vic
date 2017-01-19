# Specification to support containers with NFS based shared volumes in VIC

Container users want to be able to access shared storage between their containers programatically.  Docker solves this by way of adding [NFS volumes](https://docs.docker.com/engine/reference/commandline/volume_create/#/driver-specific-options).  VIC can streamline this functionality for the user by abstracting this away from the container user and allowing the VI admin to configure the VCH with NFS based volumes access by way of our `VolumeStore`.  This way, the VI admin can add an NFS based VolumeStore to the VCH, and the VI admin need only create volumes on it without needing to know the details of the NFS target.

### Requirements

Allow the VI admin to
 1. add an NFSv3 based `VolumeStore`

Allow the container user to
 1. create volumes on the NFS based `VolumeStore`
 1. create 1 or greater containers with NFS based volumes at the given location in the container filesystem namespace
 1. validate the volume is no longer in use and delete it

### Non Requirements

 1. Instantiation or provisioning of shared storage
 2. Exposing shared storage configuration via VIC (e.g. IOPS per client, storage policy, etc.)
 3. Management of shared storage via VIC (e.g. container quiesce for storage maintanence, etc.)

### Implementation

Adding shared storage to our model fits with the `VolumeStore` interface.  At install, a VI admin can specify an NFS target as a `VolumeStore` (potentially) using a `nfs://host/<path>` URI with a volume store name.  The container user only needs to pass the volume store name as one of the `volume create` driver opts to create a volume which will be backed by this shared storage target.  Then many containers can be created with the specified volume attached.

_VolumeStore_
The `VolumeStore` interface is used by the storage layer to implement the volume storage layer on different backend implementations.  The currenty (and only) implementation used by VIC is to manimpulate vsphere `.vmdk` backed block devices on the Datastore.  We intend to create a similar implementation for NFSv3.

https://github.com/vmware/vic/blob/master/lib/portlayer/storage/volume.go#L36
```
 35 // VolumeStorer is an interface to create, remove, enumerate, and get Volumes.
 36 type VolumeStorer interface {
 37 »···// Creates a volume on the given volume store, of the given size, with the given metadata.
 38 »···VolumeCreate(op trace.Operation, ID string, store *url.URL, capacityKB uint64, info map[string][]byte) (*Volume, error)
 39
 40 »···// Destroys a volume 
 41 »···VolumeDestroy(op trace.Operation, vol *Volume) error
 42
 43 »···// Lists all volumes 
 44 »···VolumesList(op trace.Operation) ([]*Volume, error)
 45 
 46 »···// List the configured volume stores 
 47 »···VolumeStoresList(op trace.Operation) (map[string]url.URL, error) 
 48 }   
```

When we create the NFS `VolumeStore`, we'll store the NFS target parameters (`host` + `path`) in the implementation's struct.  This is the only information we'll need to mount the NFS target on the container.
```
type NFSv3VolumeStore struct {
 target *url.URL
}
```

_VolumeCreate_
In the vsphere model, a volume is a `.vmdk` backed block device.  Creation of a volume entails attaching a new disk to the VCH, preparing it with a filesystem, and detaching it.  The resulting `.vmdk` lives in its own folder in the volume store directory (specified during install w/ `vic-machine`).  We're going to follow the same model except there is nothing to prepare.  Each volume will be a directory (which the container client will mount directly) and live at the top of the volume store directory (which we will prepare during install).

Some psuedo code.
```
func VolumeCreate() {
// volPath := vicVolumePath(nameOfVolume)
// mkdir volPath
// return volPath
}
```
_The implementation is still being worked on.  The open question is whether the VCH appliance will mount the target to manipulate the NFS target, or use an NFS client implementation in userspace instead.  I'd much (*MUCH*) rather do the latter.  The `linux` VFS implementation throws `sync` errors when mounts are unavailable.  And we don't want to bring down the appliance because of a network hiccup.  NFS is a simple protocol and there is a public pkg which implements most of it.  Adding the few primitives we need shouldn't be that difficult, but more evaluation of the work required is needed_
