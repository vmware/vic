# Vagrant Dev Box

## Overview

This box is an Ubuntu 15.04 VM with the following setup by default:

* Docker daemon with port forwarded to localhost:12375

* Go toolchain

* Additional tools (lsof, strace, etc)

## Requirements

* Vagrant (https://www.vagrantup.com/downloads.html)

* VMware Fusion or Workstation

* Vagrant Fusion or Workstation license (https://www.vagrantup.com/vmware)

## Provisioning

All files matching _provision*.sh_ in this directory will be applied by the Vagrantfile, you can symlink custom scripts
if needed.  The scripts are not Vagrant specific and can be applied to a VM running on ESX for example.

## Usage

### Shared Folders

By default your *GOPATH* is shared with the same path as the host.  This is useful if your editor runs
on the host, then errors on the guest with filename:line info have the same path.  For example, when running the
following command within the top-level project directory:

``` shell
vagrant ssh -- make -C $PWD all
```

### Create the VM

``` shell
vagrant up
```

### SSH Access

``` shell
vagrant ssh
```

### Docker Access

``` shell
DOCKER_HOST=localhost:12375 docker ps
```

### Stop the VM

``` shell
vagrant halt
```

### Restart the VM

``` shell
vagrant reload
```

### Provision

After you've done a `vagrant up`, the provisioning can be applied without reloading via:

``` shell
vagrant provision
```

### Delete the VM

``` shell
vagrant destroy
```

### Deploy to ESX

``` shell
./deploy-esx.sh
```
