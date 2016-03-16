[![Build Status](https://ci.vmware.run/api/badges/vmware/vic/status.svg)](https://ci.vmware.run/vmware/vic)


# vSphere Integrated Containers

VIC is a container runtime for vSphere, allowing developers familiar with Docker to develop in containers and deploy them alongside traditional VM-based workloads on vSphere clusters, and allowing for these workloads to be managed through the vSphere UI in a way familiar to existing vSphere admins.

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details on submitting changes and the contribution workflow.

## Building binaries for development

To build the bootstrap binaries, ensure `GOPATH` is set, then issue the following.
```
$ make all
```
This will install required tools, build the bootstrap binaries `tether-windows`, `tether-linux`, `rpctool` and server binaries `docker-server`, `port-layer-server`.  The binaries will be created in the `./binaries` directory.

To run tests after a successfull build, issue the following.
```
$ make test
```

## Managing vendor/ directory

To build the VIC dependencies, ensure `GOPATH` is set, then issue the following.
``
$ make gvt vendor
``

This will install the [gvt](https://github.com/FiloSottile/gvt) utility and retrieve the build dependencies via `gvt restore`

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
cd bootstrap && ./build.sh local
```

The iso image will be created in `./binaries`

[dronevic]:https://ci.vmware.run/vmware/vic
[dronesrc]:https://github.com/drone/drone
[dronecli]:http://readme.drone.io/devs/cli/

## Building with CI

Merges to this repository will trigger builds with [Drone][dronevic].

To build locally with Drone:

Ensure that you have Docker 1.6 or higher installed.
Install the Drone command line tools.
From the root directory of the `vic` repository run `drone exec -trusted`

## Starting docker-server

Generate certificate pair

```
go run `go env GOROOT`/src/crypto/tls/generate_cert.go --host localhost
```

Start docker-server

```
binary/docker-server --port=2376 --tls-certificate=cert.pem --tls-key=key.pem
```

## Starting port-layer-server

```
binary/port-layer-server --port 8080 --path /tmp/
```

## Testing with docker client

Download docker client
```
https://get.docker.com/builds/Linux/x86_64/docker-1.10.2
https://get.docker.com/builds/Darwin/x86_64/docker-1.10.2
https://get.docker.com/builds/Windows/x86_64/docker-1.10.2.exe
```

Pull an image
```
$ DOCKER_HOST=tcp://127.0.0.1:2376 docker pull tomcat
```

## License

VIC is available under the [Apache 2 license](LICENSE).
