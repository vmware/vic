# MITM proxy images

This directory contains the necessary consituent parts for building the images used in Group1-Docker-Commands/1-02-Docker-Pull.robot

## Building the images
### Registry Image
The registry image is in `vendor/distribution-library-image`. It is the same as the one Docker distributes, but there is no default volume, so pushed images can be saved to the registry image itself. Preparing the image for use in the test will look something like this:

```console
$ cd vendor/distribution-library-image
$ docker build -t registry .
$ docker pull busybox:latest
$ docker run -itd --net=host registry
$ docker tag busybox localhost:5000/busybox
$ docker push localhost:5000/busybox
$ docker commit registry
$ docker tag registry victest/registry-busybox:latest
$ docker push victest/registry-busybox:latest
```
You may need to add the local registry to your insecure registry list in `/etc/docker/daemon.json` before you can push to the registry.
`docker build .` in either directory to recreate the image. The registry image is a standard docker registry image, except that the `VOLUME` 

### MITMproxy Image
This is a custom image that performs a MITM on the image being pulled from the registry container when Docker is configured to use it as an HTTP proxy. `docker build .` should be sufficient for this image and then just `docker push victest/docker-layer-injection-proxy:latest`
