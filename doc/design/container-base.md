# Container Base

ContainerVMs are bootstrapped from a PhotonOS based liveCD, containing just enough to boot linux and set up the container filesystem, before performing a switch_root into the container root filesystem. The end result is a VM:
* running the PhotonOS kernel, with appropriate kernel modules for demand loading
* the specified container filesystem mounted as `/`
* [a custom init binary](tether.md) that provides the command & control channel for container interaction

[Issues relating to this component](https://github.com/vmware/vic/labels/component%2Fcontainer-base)
