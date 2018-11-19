# VIC Engine Appliance and Container VM ISOs

## Repo Specification

A repo defines the method for assembling the runtime base for both the endpointVM (appliance.iso) and containerVMs (bootstrap.iso). The endpointVM should be considered an appliance in the truest sense of the word - it is not intended to be modified or tweaked; as such it is always built from a specific repo (photon-2.0 as this document) and does not provide for changing that. The default bootstrap image is also built via the photon-2.0 repo.

The containerVMs have minimal requirements on the bootstrap contents which the repo definition tries to ensure are met:
1. a linux kernel
2. kernel modules compatible with the kernel (see the drivers array in `isos/base/repos/<repo>/init.sh` and `isos/bootstrap/bootstrap`)
3. variant specific filesystem structure
4. a tether binary built to support the linux kernel ABI in question
5. iptables binary and supporting libraries
6. an entropy daemon

These aspects are configured in the repo definition:
* the source of packages (yum repo definitions in the `*.repo` file and GPG key)
* the packages installed for the bootstrap (`repo-spec.json`)
* tweaks made to the filesystem after package install (`base.sh` and `staging.sh`)
* the package manager used to install the packages (`tdnf` or `yum`)
* the libraries and binaries used for iptables and entropy services running in the containerVM (`init.sh`)
* kernel modules loaded to support iptables, NFS, etc (`drivers` array in init.sh)

The kernel is used to boot the containerVM, and the supporting elements are preserved in a tmpfs mount for use by tether once the container filesystem has been mounted as `/`.

### How to add a new iso repo

To add a new repo definition, create a new directory under `isos/base/repos`. This will need to have the core set of files that make up a repo definition, detailed later in this readme. This repo is then built using the `bootstrap-custom` make target while specifying the `REPO=<reponame>` environment variable - this results in a `bootstrap[-repo].iso` which can be deployed via `vic-machine create ... -bi=bootstrap[-repo].iso`

The simplest approach to a new repo definition is to copy an existing repo and modify it to fit, tweaking the files described above to influce the final form of the bootstrap ISO. Chose a source repo definition that uses the same package manger as the desired custom repo. At this time only `tndf` and `yum` are supported. If others are needed then the `isos/base/utils.sh` functions need to be updated to support the new package manager.

When testing a custom bootstrap:
1. confirm that your entropy daemon is running and present in the process listing for the containerVM, e.g.
```
 # ps -ef
PID   USER     TIME  COMMAND
...
292 root      0:00 /.tether/lib64/ld-linux-x86-64.so.2 --library-path /.tether/lib64 /.tether/bin/haveged -w 1024 -v 1 -F
```
2. confirm that the iptables rules are applied as expected for a container on a bridge network to confirm all bridge traffic is allowed (`docker run -it busybox`). In this case the bridge network is the default `172.16.0.0/16`:
```
/ # /.tether/bin/iptables-wrapper -L -v
Chain INPUT (policy DROP 1 packets, 36 bytes)
 pkts bytes target     prot opt in     out     source               destination
    2    72 VIC        all  --  any    any     anywhere             anywhere
    0     0 ACCEPT     all  --  lo     any     anywhere             anywhere

Chain FORWARD (policy DROP 0 packets, 0 bytes)
 pkts bytes target     prot opt in     out     source               destination
    0     0 ACCEPT     all  --  lo     lo      anywhere             anywhere

Chain OUTPUT (policy DROP 0 packets, 0 bytes)
 pkts bytes target     prot opt in     out     source               destination
    0     0 VIC        all  --  any    any     anywhere             anywhere
    0     0 ACCEPT     all  --  any    lo      anywhere             anywhere

Chain VIC (2 references)
 pkts bytes target     prot opt in     out     source               destination
    0     0 ACCEPT     all  --  eth0   any     anywhere             anywhere             state RELATED,ESTABLISHED
    0     0 ACCEPT     all  --  any    eth0    anywhere             anywhere
    0     0 ACCEPT     all  --  eth0   any     172.16.0.0/16        anywhere
    1    36 RETURN     all  --  any    any     anywhere             anywhere
```
3. confirm that the iptables rules are applied as expected for a container on a DHCP container network with a forwarded, published port (`docker run -it --net=public -p 80:8080 tomcat:alpine /bin/ash`) - this should show rules for DHCP (ports `67` and `68`), port `8080` and NAT rules from port `80` to `8080`. Checking iptables with the advanced rules asserts that the necessary shared libraries for the functions are present and can be dynamically loaded:
```
/ # /.tether/bin/iptables-wrapper -L -v
Chain INPUT (policy DROP 1 packets, 36 bytes)
 pkts bytes target     prot opt in     out     source               destination
    5   776 VIC        all  --  any    any     anywhere             anywhere
    0     0 ACCEPT     all  --  lo     any     anywhere             anywhere

Chain FORWARD (policy DROP 0 packets, 0 bytes)
 pkts bytes target     prot opt in     out     source               destination
    0     0 ACCEPT     all  --  lo     lo      anywhere             anywhere

Chain OUTPUT (policy DROP 0 packets, 0 bytes)
 pkts bytes target     prot opt in     out     source               destination
    0     0 VIC        all  --  any    any     anywhere             anywhere
    0     0 ACCEPT     all  --  any    lo      anywhere             anywhere

Chain VIC (2 references)
 pkts bytes target     prot opt in     out     source               destination
    0     0 ACCEPT     udp  --  eth0   any     anywhere             anywhere             udp spts:67:68 dpts:67:68
    0     0 ACCEPT     udp  --  any    eth0    anywhere             anywhere             udp spts:67:68 dpts:67:68
    0     0 ACCEPT     all  --  eth0   any     anywhere             anywhere             state RELATED,ESTABLISHED
    0     0 ACCEPT     all  --  any    eth0    anywhere             anywhere
    0     0 ACCEPT     icmp --  eth0   any     anywhere             anywhere             icmp echo-request
    0     0 ACCEPT     icmp --  any    eth0    anywhere             anywhere             icmp echo-reply
    0     0 ACCEPT     tcp  --  eth0   any     anywhere             anywhere             tcp dpt:8080
    1    36 RETURN     all  --  any    any     anywhere             anywhere
/ # /.tether/bin/iptables-wrapper -L -v -t nat
Chain PREROUTING (policy ACCEPT 1 packets, 36 bytes)
 pkts bytes target     prot opt in     out     source               destination
    0     0 REDIRECT   tcp  --  eth0   any     anywhere             anywhere             tcp dpt:80 redir ports 8080

Chain INPUT (policy ACCEPT 0 packets, 0 bytes)
 pkts bytes target     prot opt in     out     source               destination

Chain OUTPUT (policy ACCEPT 0 packets, 0 bytes)
 pkts bytes target     prot opt in     out     source               destination

Chain POSTROUTING (policy ACCEPT 0 packets, 0 bytes)
 pkts bytes target     prot opt in     out     source               destination
/ #
```
4. Check that the containerVM can mount an NFS volume (if this is a necessary configuration) - this confirms that the NFS kernel modules or core support are present.

#### Examples

This is the example repo spec for photon-2.0:

```json
{
    "name": "photon-2.0",
    "packagemanager": "tdnf",
    "kernel": "linux-esx",
    "init": "/bin/init",
    "packages_script_staging": {
        "bootstrap": "",
        "appliance": "erase toybox"
    },
    "packages": {
        "base": "filesystem toybox kmod haveged iptables",
        "bootstrap": "",
        "appliance": "procps-ng e2fsprogs shadow systemd iputils iproute2 tdnf gzip lsof logrotate photon-release mingetty rpm dbus net-tools openssh ca-certificates sudo vim lsof"
    }
}
```

#### Structure:

Each repository folder under isos/base/repos consists of:

 - `init.sh`            - Functions needed by bootstrap to set the entropy and iptables config paths.
 - `base.sh`            - An extra script to run at the start of the base filesystem layout if needed.
 - `staging.sh`         - An extra script to run after the staging step if filesystem modification is needed (see packages_script_staging for package level modification).
 - `distro.repo`        - A yum-compatible specification of the desired repositories for package downloads, added to /etc/yum.repos.d/.
 - `repo-spec.json`     - The repository specification. See below.

The following fields of repo-spec are required to build a repo:

 - `name`               - The name of the repository. This should match the name of the repo directory, and probably be formatted like `repo-version`.
 - `packagemanager`     - The package manager to use during install. Currently supported managers are yum and tdnf.
 - `init`               - The init system. Appliances always use /lib/systemd/systemd, bootraps can use a slim boot script at `/bin/init` or a full init system like SYSVinit/systemd
 - `kernel`             - Determines the iso kernel. Can be an absolute path to a Kernel RPM or a repository package, e.g. linux or linux-esx.
 - `packages.base`      - Packages that are required in both the VCH and containerVM, e.g. filesystem coreutils kmod.
 - `packages.bootstrap` - Packages to be installed on the bootstrap iso but not the appliance iso.
 - `packages.appliance` - Packages to be installed on the appliance iso but not the bootstrap iso. 

The following fields of the repo-spec are optional:

 - `packages_script_staging.*` - A package manager script that runs on the after the base install but before the boostrap/appliance install. 

## How to build an custom iso

ISOs can be build using the *custom* iso Make targets. The appliance does not support a custom repo - it will always be built with the default, which is photon-2.0.

### Use sample repos to build custom iso
The source code has already provided several sample repos including centos-6.9, photon-1.0 and photon-2.0.
Use the REPO environment variable to choose with iso distro to build:

```sh
# Running against a local docker
unset DOCKER_HOST DOCKER_TLS_VERIFY DOCKER_CERT_PATH
# NOTE: we have to delete the tether-linux and unpack binaries for each bootstrap ISO
rm -f bin/tether-linux bin/unpack  

# builds vic-machine binaries and both appliance.iso and bootstrap.iso with photon-2.0
docker run -v $(pwd):/go/src/github.com/vmware/vic  gcr.io/eminent-nation-87317/vic-build-image:tdnf  rm -f bin/tether-linux bin/unpack \&\& make isos vic-machine

# builds appliance.iso with photon-2.0 and bootstrap.iso by photon-1.0
docker run -v $(pwd):/go/src/github.com/vmware/vic  -e REPO=photon-1.0  gcr.io/eminent-nation-87317/vic-build-image:tdnf  rm -f bin/tether-linux bin/unpack \&\& make appliance bootstrap-custom


# builds a bootstrap iso with a centos-6.9 kernel and package layout
docker run -v $(pwd):/go/src/github.com/vmware/vic  -e REPO=centos-6.9  gcr.io/eminent-nation-87317/vic-build-image:yum  rm -f bin/tether-linux bin/unpack \&\& make bootstrap-custom
```

The output from these files will be in `bin/` titled `bootstrap[-repo].iso`. The default bootstrap build from the `bootstrap` make target will produce the `bootstrap.iso` file with no suffix.

To use the custom bootstrap iso on your VCH, please supply the `--bootstrap-iso <iso>` flag to `vic-machine` create.

### Use custom repos to build iso

First, prepare your custom repo as the above instructions. Then you have a folder named \<REPONAME\> containing needed files.

And then, you should be able to run these commands against a VCH or a local docker, adding/removing the Capacity option from the initial volume create. If you have your own fork/branch of the vmware/vic repo in GitHub you can updated the “-b vmware/master” argument to “-b YourGithub/YourBranchName” for example to build that:

```sh 
# Running directly against a VCH
 
# NOTE: we have to delete the tether-linux and unpack binaries for each bootstrap ISO as we get "FATAL: kernel too old" messages otherwise
# This is in the commands, but noted here for clarity
 
# Volume to store the source and binaries for cross build efficiency - sized for artifacts and build caches
docker volume create --name=vic-build --opt Capacity=4G
 
# build image for the custom iso
docker create -it -e REPO=<REPONAME> -v vic-build:/go/src/github.com/vmware/vic --name=custom-builder gcr.io/eminent-nation-87317/vic-build-image:yum -b vmware/master rm -f bin/tether-linux bin/unpack \&\& make bootstrap-custom
 
# Copy your repo folder <REPONAME> into the local ./isos/base/repos directory and it gets copied from there - this path is to make life easy if you have the repo checked out. This should result in an isos/base/repos/<REPONAME> directory.
docker cp isos/base/repos/<REPONAME> custom-builder:/go/src/github.com/vmware/vic/isos/base/repos/
 
# build the custom iso
docker start -ai custom-builder
 
# copy the generated artifacts out of the container
docker cp custom-builder:/go/src/github.com/vmware/vic/bin/bootstrap-<REPONAME>.iso .
``` 
