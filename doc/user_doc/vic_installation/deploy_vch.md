# Deploy VCHs #

You install vSphere Integrated Containers Engine by deploying vSphere Integrated Containers Engine virtual container hosts (VCHs). You use the `vic-machine` utility to deploy a VCH. 

The `vic-machine` utility can deploy a VCH in one of the following setups: 
* vCenter Server with a cluster
* vCenter Server with one or more standalone ESXi hosts
* A standalone ESXi host

The VCH allows you to use an ESXi host or vCenter Server instance as the Docker endpoint for a Docker client. The containers that Docker developers pull or create by using a Docker client are stored and managed in the vSphere environment.

When you deploy a VCH, `vic-machine` registers the VCH as a vSphere extension. Authentication between the VCH and vSphere is handled via key pair authentication against the vSphere extension.

* [Environment Prerequisites for VCH Deployment](vic_installation_prereqs.md)
* [Deploy a VCH to an ESXi Host](deploy_vch_esxi.md)
* [Deploy a VCH to a vCenter Server Cluster](deploy_vch_vcenter.md)
* [Verify the Deployment of a VCH](verify_vch_deployment.md)
* [VCH Deployment Options](vch_installer_options.md)
* [Examples of Deploying a VCH](vch_installer_examples.md)