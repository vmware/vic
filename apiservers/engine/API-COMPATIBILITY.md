# Docker Remote API Compatibility
VIC attempts to be compatible with the Docker remote API; however, there are some differences.  The API differences are documented below, grouped together by functional areas.

## Container APIs

### List containers
```
GET /containers/json

```
No differences

### Create a container
### Inspect a container
### List processes running inside a container
### Get container logs
### Inspect changes on a container’s filesystem
### Export a container
### Get container stats based on resource usage
### Resize a container TTY
### Start a container
### Stop a container
### Restart a container
### Kill a container
### Update a container
### Rename a container
### Pause a container
### Unpause a container
### Attach to a container
### Attach to a container (websocket)
### Wait a container
### Remove a container
### Copy files or folders from a container
### Retrieving information about files and folders in a container
### Get an archive of a filesystem resource in a container
### Extract an archive of files or folders to a directory in a container

## Image APIs

### List Images
### Build image from a Dockerfile
### Create an image
### Inspect an image
### Get the history of an image
### Push an image on the registry
### Tag an image into a repository
### Remove an image
### Search images

## Network APIs

### List networks
### Inspect network
### Create a network
### Connect a container to a network
### Disconnect a container from a network
### Remove a network

## Volume APIs

### List volumes
### Create a volume
### Inspect a volume
### Remove a volume

## MISC APIs

### Check auth configuration
### Display system-wide information
### Show the docker version information
### Ping the docker server
### Create a new image from a container’s changes
### Monitor Docker’s events
### Get a tarball containing all images in a repository
### Get a tarball containing all images.
### Load a tarball with a set of images and tags into docker
### Image tarball format
### Exec Create
### Exec Start
### Exec Resize
### Exec Inspect
