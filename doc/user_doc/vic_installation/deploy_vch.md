# Deploy Virtual Container Hosts #

You install vSphere Integrated Containers Engine by deploying vSphere Integrated Containers Engine virtual container hosts. You use the `vic-machine` utility to deploy a virtual container host. 

The `vic-machine` utility can deploy a virtual container host in one of the following setups: 
* vCenter Server with a cluster
* vCenter Server with one or more standalone ESXi hosts
* A standalone ESXi host

The virtual container host allows you to use an ESXi host or vCenter Server instance as the Docker endpoint for a Docker client. The containers that Docker developers pull or create by using a Docker client are stored and managed in the vSphere environment.

When you deploy a virtual container host, `vic-machine` registers the virtual container host as a vSphere extension. Authentication between the virtual container host and vSphere is handled via key pair authentication against the vSphere extension.

* [Environment Prerequisites for Virtual Container Host Deployment](vic_installation_prereqs.md)
* [Deploy a Virtual Container Host to an ESXi Host](deploy_vch_esxi.md)
* [Deploy a Virtual Container Host to a vCenter Server Cluster](deploy_vch_vcenter.md)
* [Verify the Deployment of a Virtual Container Host](verify_vch_deployment.md)
* [Virtual Container Host Deployment Options](vch_installer_options.md)
* [Examples of Deploying a Virtual Container Host](vch_installer_examples.md)