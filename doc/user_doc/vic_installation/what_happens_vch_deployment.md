# What Happens During VCH Deployment #

When you deploy a virtual container host (VCH), `vic-machine` performs different actions depending on the configuration of the vSphere environment to which you are deploying.

## Deployment to an ESXi Host ##

When you deploy a VCH to an ESXi host that is not managed by vCenter Server, `vic-machine` performs the following actions:

- Generates TLS certificate and key files for you provide to Docker clients so that they can authenticate with the VCH.
- Creates a virtual switch and port group, each with the name `docker-machine`.
- Creates a resource pool with the name `docker-machine`.
- Creates the VCH `docker-machine` in the `docker-machine` resource pool.
- Uploads the `appliance.iso` file to the image store on the target host, and boots the VCH from that image.
- Uploads the `bootstrap.iso` file to the image store on the target host, to use when booting container VMs.
 
## Deployment to a vCenter Server Cluster  ##

- Verifies that DRS is correctly configured on the cluster.
- Verifies that a distributed virtual switch exists 

## Deployment to a Standalone Host on vCenter Server ##

dfsfd


