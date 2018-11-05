# VIC Engine Appliance and Container VM ISOs

## How to add a new iso repo

### Repo Specification

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

ISOs can be build using the *custom* iso Make targets. The appliance does not currently support a custom repo - it will always be built with the default, which is photon-2.0.

### Use sample repos to build custom iso
The source code has already provided several sample repos including centos-6.9, photon-1.0 and photon-2.0.
Use the REPO environment variable to choose with iso distro to build:

```sh
# Running against a local docker
rm -f bin/tether-linux bin/unpack  #NOTE: we have to delete the tether-linux and unpack binaries for each bootstrap ISO

# builds both appliance.iso and bootstrap.iso by default photon-2.0 
docker run -v $(pwd):/go/src/github.com/vmware/vic  gcr.io/eminent-nation-87317/vic-build-image:tdnf  make isos

# builds both appliance.iso and bootstrap.iso by photon-1.0 
docker run -v $(pwd):/go/src/github.com/vmware/vic  -e REPO=photon-1.0  gcr.io/eminent-nation-87317/vic-build-image:tdnf  make appliance bootstrap-custom

# builds a bootstrap iso with a centos-6.9 kernel and package layout
docker run -v $(pwd):/go/src/github.com/vmware/vic  -e REPO=centos-6.9  gcr.io/eminent-nation-87317/vic-build-image:yum  make bootstrap-custom
```

The output from these files will be in `bin/` titled `bootstrap[-distro].iso`. The default bootstrap build from the `bootstrap` make target will produce the `bootstrap.iso` file with no suffix.

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
