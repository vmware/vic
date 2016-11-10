# Deploy a Virtual Container Host

You install vSphere Integrated Containers Engine by deploying a vSphere Integrated Containers Engine virtual container host. You use the `vic-machine create` command to deploy a virtual container host. 

The `vic-machine` utility can deploy a virtual container host in one of the following setups: 
* vCenter Server with a cluster
* vCenter Server with one or more standalone ESXi hosts
* A standalone ESXi host

When you deploy a virtual container host, `vic-machine` registers the virtual container host as a vSphere extension. Authentication between the virtual container host and vSphere is handled via key pair authentication against the vSphere extension.

The virtual container host allows you to use an ESXi host or vCenter Server instance as the Docker endpoint for a Docker client. The containers that you pull or create by using your Docker client are stored and managed in the vSphere environment.

**Prerequisites**

* Verify that your vSphere infrastructure meets the requirements in [Environment Prerequisites for vSphere Integrated Containers Engine Installation](vic_installation_prereqs.md).
* In a vCenter Server environment, before you deploy a virtual container host, you must create a distributed virtual switch and a distributed port group for use as the bridge network for container VMs. For information about how to create a distributed virtual switch and port group, see *Network Requirements* in [Environment Prerequisites for vSphere Integrated Containers Engine Installation](vic_installation_prereqs.md#networkreqs).
* Obtain either a verified build, the latest daily build, or the source code of vSphere Integrated Containers Engine: 
 * Download the most recent verified build of vSphere Integrated Containers Engine from https://github.com/vmware/vic/releases and unpack it. This version has been tested and approved, but it does not reflect the most up-to-date version of the code.
 * Download the latest daily build of vSphere Integrated Containers Engine from https://bintray.com/vmware/vic-repo/build/view#files and unpack it. This version reflects the version of the code as it was at the last daily build. It has not been tested or approved.
 * For the very latest version, for example to include changes that you have made since the last daily build, build the vSphere Integrated Containers Engine binaries from the source code.
* Add the folder that contains the vSphere Integrated Containers Engine binaries to the `PATH` environment variable on the machine on which you are running `vic-machine`. 
* Familiarize yourself with the vSphere Integrated Containers Engine binaries, as described in [Contents of the vSphere Integrated Containers Engine Binaries](contents_of_vic_binaries.md). 
* Familiarize yourself with the options of the `create` command of the `vic-machine` utility described in [Virtual Container Host Deployment Options](vch_installer_options.md).
* For examples of commands to deploy a virtual container host in various vSphere configurations, see [Examples of Deploying a Virtual Container Host](vch_installer_examples.md).
 

**Procedure**

1. Open a terminal on the system on which you downloaded and unpacked the vSphere Integrated Containers Engine binary bundle.
2. Navigate to the directory that contains the `vic-machine` utility:
 * If you downloaded the most recent verified build or the latest daily build, go to the location in which you unpacked the vSphere Integrated Containers Engine bundle.
 * If you built the vSphere Integrated Containers Engine binaries, go to <code><i>installation_dir</i>/vic/bin</code>.
2. Run the `create` command of the `vic-machine` utility. 

   The following examples include the fewest possible options for installation in a simple vCenter Server environment with a cluster. For simplicity, these examples assume that the vSphere environment uses trusted certificates signed by a known Certificate Authority (CA), so the `--thumbprint option` is not specified, and verification of client TLS certificates is disabled.

   Deploy a virtual container host from a Mac OS system:

   <pre>$ vic-machine-darwin create
--target <i>vcenter_server_address</i>
--image-store <i>datastore_name</i> 
--user <i>username</i>
--bridge-network <i>network_name</i>
--no-tlsverify</pre>  

   Deploy a virtual container host from a Linux OS system:

   <pre>$ vic-machine-linux create
--target <i>vcenter_server_address</i>
--image-store <i>datastore_name</i> 
--user <i>username</i>
--bridge-network <i>network_name</i>
--no-tlsverify</pre> 

   Deploy a virtual container host from a Windows system:

   <pre>$ vic-machine-windows create
--target <i>vcenter_server_address</i>
--image-store <i>datastore_name</i> 
--user <i>username</i>
--bridge-network <i>network_name</i>
--no-tlsverify</pre> 


**Result**

At the end of a successful installation, `vic-machine` displays information about the new virtual container host:
   
<pre>Initialization of appliance successful
vic-admin portal:
https://<i>vch_address</i>:2378
Published ports can be reached at:
<i>vch_address</i>
Docker environment variables:
DOCKER_HOST=<i>vch_address</i>:2376
Environment saved in <i>vch_name</i>/<i>vch_name</i>.env
Connect to docker:
docker -H <i>vch_address</i>:2376 --tls info
Installer completed successfully</pre>

**What to Do Next** 

If you did not explicitly disable TLS certificate generation by using the `no-tls` option, and if your Docker client is not on the same system as the one that you used to run `vic-machine`, you must copy  all of the generated <code>*.pem</code> certificate files to the Docker client system.

To test your virtual container host, see [Verify the Deployment of a Virtual Container Host](verify_vch_deployment.md).
    
    
