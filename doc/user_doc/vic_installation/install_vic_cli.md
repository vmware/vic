# Deploy a Virtual Container Host

You use the `create` command of the `vic-machine` command line utility to deploy a vSphere Integrated Containers virtual container host. 

The `vic-machine` utility can deploy a virtual container host in one of the following setups: 
* vCenter Server with a cluster
* vCenter Server with one or more standalone ESXi hosts
* A standalone ESXi host

The virtual container host allows you to use an ESXi host or vCenter Server instance as the Docker endpoint for a Docker client. The containers that you pull or create in your Docker client are stored and managed in the vSphere environment.

**NOTE** The `vic-machine` utility does not add an extension in the vSphere Web Client. 

**Prerequisites**

* Verify that your vSphere infrastructure meets the requirements in [Environment Prerequisites for vSphere Integrated Containers Installation](vic_installation_prereqs.md).
* If you are deploying a virtual container host in a vSphere environment with more than one ESXi host, create a private port group for container VMs to use to communicate with each other. For information about how to create a private port group, see [Create a Private Port Group for Virtual Container Hosts](create_a_private_port_group_for_vch.md).
* Obtain either a verified build, the latest daily build, or the source code of vSphere Integrated Containers: 
 * Download the most recent verified build of vSphere Integrated Containers from https://github.com/vmware/vic/releases and unpack it. This version has been tested and approved, but it does not reflect the most up-to-date version of the code.
 * Download the latest daily build of vSphere Integrated Containers from https://bintray.com/vmware/vic-repo/build/view#files and unpack it. This version reflects the version of the code as it was at the last daily build. It has not been tested or approved.
 * For the very latest version, for example to include changes that you have made since the last daily build, build the vSphere Integrated Containers binaries from the source code.
* Add the folder that contains the vSphere Integrated Containers binaries to the `PATH` environment variable. 
* Familiarize yourself with the vSphere Integrated Containers binaries, as described in [Contents of the vSphere Integrated Containers Binaries](contents_of_vic_binaries.md). 
* Familiarize yourself with the options of the `create` command of the `vic-machine` utility described in [Virtual Container Host Deployment Options](vch_installer_options.md).
* For examples of commands to deploy a virtual container host in various vSphere configurations, see [Examples of Deploying a Virtual Container Host](vch_installer_examples.md).
 

**Procedure**

1. Open a terminal on the system on which you downloaded and unpacked the vSphere Integrated Containers binary bundle.
2. Navigate to the directory that contains the `vic-machine` utility:
 * If you downloaded the most recent verified build or the latest daily build, go to the location in which you unpacked the vSphere Integrated Containers bundle.
 * If you built the vSphere Integrated Containers binaries, go to <code><i>installation_dir</i>/vic/bin</code>.
2. Run the `create` command of the `vic-machine` utility. 

   The following examples include the fewest possible options, for installation in a simple vCenter Server environment with a cluster.

   Deploy a virtual container host from a Mac OS system:

   <pre>$ vic-machine-darwin create
--target <i>vcenter_server_address</i>
--image-datastore <i>datastore_name</i> 
--user <i>username</i>
--compute-resource /<i>datacenter_name</i>/host/<i>cluster_name</i>/
--bridge-network <i>network_name</i></pre>  

   Deploy a virtual container host from a Linux OS system:

   <pre>$ vic-machine-linux create
--target <i>vcenter_server_address</i>
--image-datastore <i>datastore_name</i> 
--user <i>username</i>
--compute-resource /<i>datacenter_name</i>/host/<i>cluster_name</i>/
--bridge-network <i>network_name</i></pre> 

   Deploy a virtual container host from a Windows system:

   <pre>$ vic-machine-windows create
--target <i>vcenter_server_address</i>
--image-datastore <i>datastore_name</i> 
--user <i>username</i>
--compute-resource /<i>datacenter_name</i>/host/<i>cluster_name</i>/
--bridge-network <i>network_name</i></pre> 

   At the end of a successful installation, `vic-machine` displays a success message:
   
   <pre>Initialization of appliance successful
SSH to appliance (default=root:password)
ssh root@<i>vch_address</i>
Log server:
https://<i>vch_address</i>:2378
Connect to docker:docker -H <i>vch_address</i>:2376 
--tls --tlscert='./<i>vch_name</i>-cert.pem' 
--tlskey='./<i>vch_name</i>.pem' info</pre>

3. (Optional) Copy the generated certificate and key files to the Docker client system.

  If you did not explicitly disable TLS certificate generation by using the `generate-cert=false` option, and if your Docker client is not on the same system as the one that you used to run `vic-machine`, you must copy the <code><i>vch_name</i>-cert.pem</code> and <code><i>vch_name</i>-key.pem</code> files to the Docker client system.

**What to Do Next**

If you did not explicitly disable TLS certificate generation by using the `generate-cert=false` option, and if your Docker client is not on the same system as the one that you used to run `vic-machine`, copy the <code><i>vch_name</i>-cert.pem</code> and <code><i>vch_name</i>-key.pem</code> files to the Docker client system. 

To test your virtual container host, see [Verify the Deployment of a Virtual Container Host to vCenter Server](verify_vch_deployment.md) or [Verify the Deployment of a Virtual Container Host to an ESXi Host](verify_vch_deployment_esx.md).
    
    