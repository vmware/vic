[![Build Status](https://ci.vcna.io/api/badges/vmware/vic/status.svg)](https://ci.vcna.io/vmware/vic) [![codecov](https://codecov.io/gh/vmware/vic/branch/master/graph/badge.svg)](https://codecov.io/gh/vmware/vic) [![Download](https://img.shields.io/badge/download-latest-blue.svg)](https://github.com/vmware/vic/releases/latest) [![Go Report Card](https://goreportcard.com/badge/github.com/vmware/vic)](https://goreportcard.com/report/github.com/vmware/vic)

# vSphere Integrated Containers Engine

vSphere Integrated Containers Engine (VIC Engine) is a container runtime for vSphere, allowing developers familiar with Docker to develop in containers and deploy them alongside traditional VM-based workloads on vSphere clusters, and allowing for these workloads to be managed through the vSphere UI in a way familiar to existing vSphere admins.

See [VIC Engine Architecture](doc/design/arch/arch.md) for a high level overview.


## Project Status

VIC Engine now provides:
* basic function for most of the core lifecycle operations: pull, create, start, attach, run, kill, stop, wait, rm, logs.
* early vCenter support, leveraging DRS for initial placement.
* volume support [--volume-store](doc/user/usage.md#configuring-volumes-in-a-virtual-container-host) - SIOC is not integrated but can be set as normal.
* direct mapping of vSphere networks [--container-network](doc/user/usage.md#exposing-vsphere-networks-within-a-virtual-container-host) - NIOC is not integrated but can be set as normal.
* dual-mode management - IP addresses are reported as normal via vSphere UI, guest shutdown via the UI will trigger delivery of container STOPSIGNAL, restart will relaunch container process.
* logs command - follow is available however timestamps and corresponding filtering are not
* client authentication - basic authentication via client certificates known as _tlsverify_

We are working hard to add functionality while building out our [foundation](doc/design/arch/arch.md#port-layer-abstractions) so continue to watch the repo for new features. Initial focus is on the production end of the CI pipeline, building backwards towards developer laptop scenarios.

This limited set of current capabilities may come as a surprise to people who are familiar with [Project Bonneville](http://blogs.vmware.com/cloudnative/introducing-project-bonneville/) that was [reasonably fully featured](https://www.youtube.com/watch?v=XkFQw8ueT1w) when demonstrated at VMworld in 2015.
Project Bonneville was research aimed at determining best approaches to enabling container workflows in a vSphere environment and therefore enabled a broad set of features, but not in a manner that made it a viable product for large scale consumption. Building on top of research code is a great shortcut for fast time-to-market, but does not provide a good foundation for an enterprise quality product. vSphere Integrated Containers Engine is a full re-architecture and re-write, building off the knowledge gained during Project Bonneville while keeping almost zero code.


## Installing

Once built, pick up the correct binary based on your OS, and then the result can be installed with the following command.

```
bin/vic-machine-linux create --target <target-host>[/datacenter] --image-store <datastore name> --name <vch-name> --user <username> --password <password> --thumbprint <certificate thumbprint> --compute-resource <cluster or resource pool name> --tls-cname <FQDN, *.wildcard.domain, or static IP>
```

See `vic-machine-XXX create --help` for usage information.

A more indepth example can be found [here](doc/user/usage.md)

Container output is found in a log file on the datastore ([datastore]/containerid/containerid.log).


## Deleting

The installed Virtual Container Host can be deleted through vic-machine-XXX delete.

See `vic-machine-XXX delete --help` for usage information.

A more indepth example can be found [here](doc/user/usage.md)


## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details on submitting changes and the contribution workflow.


## Building

Building the project is done with a combination of make and containers, with golang:1.8 being the common container base. This is done so that it's possible to build directly, without a functional docker, if using a Debian based system with the Go 1.8 toolchain and Drone.io installed.

To build as closely as possible to the formal build:
```
drone exec --repo.trusted
```

To build without modifying the local system:
```
docker run -v $(pwd):/go/src/github.com/vmware/vic -w /go/src/github.com/vmware/vic golang:1.8 make all
```

To build directly:
```
make all
```

There are three primary components generated by a full build, found in `$BIN` (the ./bin directory by default). The make targets used are the following:
 1. vic-machine - `make vic-machine`
 2. appliance.iso - `make appliance`
 3. bootstrap.iso - `make bootstrap`


## Building binaries for development

Some of the project binaries can only be built on Linux.  If you are developing on a Mac or Windows OS, then the easiest way to facilitate a build is by utilizing the project's Vagrantfile.  The Vagrantfile will share the directory where the file is executed and set the GOPATH based on that share.

To build the component binaries, ensure `GOPATH` is set, then issue the following command in the root directory:
```
$ make components
```
This will install required tools and build the component binaries `tether-linux`, `rpctool` and server binaries `docker-engine-server`, `port-layer-server`.  The binaries will be created in the `$BIN` directory, ./bin by default.

To run tests after a successful build, issue the following:
```
$ make test
```

## Managing vendor/ directory

To build the VIC Engine dependencies, ensure `GOPATH` is set, then issue the following.
``
$ make gvt vendor
``

This will install the [gvt](https://github.com/FiloSottile/gvt) utility and retrieve the build dependencies via `gvt restore`


## Building the ISOs

The component binaries above are packaged into ISO files, appliance.iso and bootstrap.iso, that are used by the installer. The generation of the ISOs is split into the following targets:
iso-base, appliance-staging, bootstrap-staging, appliance, and bootstrap. Generation of the ISOs involves authoring a new root filesystem, meaning running a package manager (currently yum) and packing/unpacking archives. To install packages and preserve file permissions while unpacking these steps should be run as root, whether directly or in a container. To generate the ISOs:

```
$ make isos
```

The appliance and bootstrap ISOs are bootable CD images used to start the VMs that make up VIC Engine. To build the image using [docker](https://www.docker.com/), ensure `GOPATH` is set and `docker` is installed, then issue the following.

```
docker run -v $(pwd):/go/src/github.com/vmware/vic -w /go/src/github.com/vmware/vic golang:1.8 make isos
```

Alternatively, the iso image can be built locally.  Again, ensure `GOPATH` is set, but also ensure the following packages are installed. This will attempt to install the following packages if not present using apt-get:

```
apt-get install \
	curl \
	cpio \
	tar \
	xorriso \
	rpm \
	ca-certificates \
	yum
```

Package names may vary depending on the distribution being used.  Once installed, issue the following (the targets listed here are those executed when using the `iso` target.

```
make iso-base appliance-staging appliance bootstrap-staging bootstrap
```

The iso image will be created in `$BIN`


## Building with CI

[dronevic]:https://ci.vcna.io/vmware/vic
[dronesrc]:https://github.com/drone/drone
[dronecli]:http://readme.drone.io/0.5/install/cli/

PRs to this repository will trigger builds on our [Drone CI][dronevic].

To build locally with Drone:

Ensure that you have Docker 1.6 or higher installed.
Install the [Drone command line tools][dronecli].
From the root directory of the `vic` repository run `drone exec --repo.trusted`

## Common Build Problems
1. Builds may fail when building either the appliance.iso or bootstrap.iso with the error: `cap_set_file failed - Operation not supported`

   *Cause:* Some Ubuntu and Debian based systems ship with a defective `aufs` driver, which Docker uses as its default backing store.  This driver does not support extended file capabilities such as `cap_set_file`

   *Solution:* Edit the `/etc/default/docker` file, add the option `--storage-driver=overlay` to the `DOCKER_OPTS` settings, and restart Docker.

2. `go vet` fails when doing a `make all`

    *Cause:* Apparently some caching takes place in `$GOPATH/pkg/linux_amd64/github.com/vmware/vic` and can cause `go vet` to fail when evaluating outdated files in this cache.

    *Solution:* Delete everything under `$GOPATH/pkg/linux_amd64/github.com/vmware/vic` and re-run `make all`.

3.  `vic-machine upgrade` integration tests fail due to `BUILD_NUMBER` being set incorrectly when building locally

    *Cause:* `vic-machine` checks the build number of its binary to determine upgrade status and a locally-built `vic-machine` binary may not have the `BUILD_NUMBER` set correctly. Upon running `vic-machine upgrade`, it may fail with the message `foo-VCH has same or newer version x than installer version y. No upgrade is available.`

    *Solution:* Set `BUILD_NUMBER` to a high number at the top of the `Makefile` - `BUILD_NUMBER ?= 9999999999`. Then, re-build binaries - `sudo make distclean && sudo make clean && sudo make all` and run `vic-machine upgrade` with the new binary.

## Integration Tests

[VIC Engine Integration Test Suite](tests/README.md) includes instructions to run locally.


## License

VIC Engine is available under the [Apache 2 license](LICENSE).
