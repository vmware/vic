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
 - `base.sh`            - An extra script to run during the base filesystem layout if needed.
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

Use the REPO environment variable to choose with iso distro to build:

```sh
make isos # builds the isos using the default photon-2.0 repo
REPO=photon-1.0 make appliance bootstrap-custom # builds the photon-2.0 appliance and a photon-1.0 bootstrap
REPO=centos-6.9 make bootstrap-custom # builds a bootstrap iso with a centos-6.9 kernel and package layout
```

The output from these files will be in `bin/` titled `bootstrap[-distro].iso`. The default bootstrap build from the `bootstrap` make target will produce the `bootstrap.iso` file with no suffix.

To use the custom bootstrap iso on your VCH, please supply the `--bootstrap-iso <iso>` flag to `vic-machine` create.
