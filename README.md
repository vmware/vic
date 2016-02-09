# vSphere Integrated Containers

VIC is a container runtime for vSphere, allowing developers familiar with Docker to develop in containers and deploy them alongside traditional VM-based workloads on vSphere clusters, and allowing for these workloads to be managed through the vSphere UI in a way familiar to existing vSphere admins.

## Building binaries for development

To build the bootstrap binaries, ensure `GOPATH` is set, then issue the following.
```
$ make gvt vendor bootstrap
```
This will install the [gvt](https://github.com/FiloSottile/gvt) utility, retrieve the build dependencies via `gvt restore`, then build the bootstrap binaries `tether-windows`, `tether-linux`, and `rpctool`.  The binaries will be created in the `./binaries` directory.

## Building the bootstrap container iso

The bootstrap container iso is the bootable cd image used to boot container VM templates.  To build the image using [docker](https://www.docker.com/), ensure `GOPATH` is set and `docker` is installed, then issue the following.

```
cd bootstrap && ./build.sh docker
```

Alternatively, the iso image can be built locally.  Again, ensure `GOPATH` is set, but also ensure the following packages are installed.
```
	curl \
	cpio \
	ssh-client \
	squashfs-tools \
	xorriso \
	git \
	gcc \
	make \
	vim \
	patch \
	rpm \
	busybox-static \
	isolinux \
	syslinux
```

Package names may vary depending on the distribution being used.  Once installed, issue the following.

```
cd bootstrap && ./build.sh docker
```

The iso image will be created in `./binaries`

## License

VIC is available under the [Apache 2 license](LICENSE).
