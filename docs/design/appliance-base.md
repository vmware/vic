# Appliance Base

The appliance VM is bootstrapped from a PhotonOS based liveCD that embeds elements of VIC relevent to Virtual Container Host functions. The appliance is diskless _in essence_, but may well use a non-persistent disk to cache transient data such as docker images in-flight from [Docker Hub](https://hub.docker.com/) but not yet persisted to a datastore.

[Issues relating to this component](https://github.com/vmware/vic/labels/component%2Fappliance-base)
