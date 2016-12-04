# Deploy a VCH to a Basic vCenter Server Cluster

This topic provides instructions for deploying a virtual container host (VCH) in a very basic vCenter Server environment. This basic deployment allows you to test vSphere Integrated Containers Engine with vCenter Server before attempting a more complex deployment that corresponds to your real vSphere environment.

The vCenter Server instance to which you deploy the VCH must match the specifications listed in the prerequisites.

**Prerequisites**
* Download and unpack the vSphere Integrated Containers Engine bundle. For information about where to obtain vSphere Integrated Containers Engine, see [Download vSphere Integrated Containers Engine](download_vic.md).
* Create or obtain a vCenter Server instance with the following configuration:
  * One datacenter
  * One cluster with two ESXi hosts and DRS enabled. You can use nested ESXi hosts for this example.
  * A shared datastore, that is accessible by both of the ESXi hosts.
  * The VM Network is present
  * One distributed virtual switch with one port group named `vic-bridge`
* Verify that your vCenter Server instance and both of the ESXi hosts in the cluster meet the requirements in [Environment Prerequisites for vSphere Integrated Containers Engine Installation](vic_installation_prereqs.md).
* Familiarize yourself with the vSphere Integrated Containers Engine binaries, as described in [Contents of the vSphere Integrated Containers Engine Binaries](contents_of_vic_binaries.md). 
* Familiarize yourself with the options of the `vic-machine create` command described in [VCH Deployment Options](vch_installer_options.md).
 

**Procedure**

1. Open a terminal on the system on which you downloaded and unpacked the vSphere Integrated Containers Engine binary bundle.
2. Navigate to the directory that contains the `vic-machine` utility:
3. Run the `vic-machine create` command.

   Wrap any option arguments that include spaces or special characters in quotes. Use single quotes if you are using `vic-machine` on a Linux or Mac OS system and double quotes on a Windows system. In these examples, the user name is wrapped in quotes because it contains `@`.

   - Linux OS:
      <pre>$ vic-machine-linux create
     --target <i>vcenter_server_address</i>
     --user 'Administrator@vsphere.local'
     --password <i>vcenter_server_password</i>
     --bridge-network vic-bridge
     --image-store <i>shared_datastore_name</i>
     --no-tlsverify
     --force
     </pre>  
   - Windows:
      <pre>$ vic-machine-windows create
     --target <i>vcenter_server_address</i>
     --user "Administrator@vsphere.local"
     --password <i>vcenter_server_password</i>
     --bridge-network vic-bridge
     --image-store <i>shared_datastore_name</i>
     --no-tlsverify
     --force
     </pre> 
   - Mac OS:
       <pre>$ vic-machine-darwin create
     --target <i>vcenter_server_address</i>
     --user 'Administrator@vsphere.local'
     --password <i>vcenter_server_password</i>
     --bridge-network vic-bridge
     --image-store <i>shared_datastore_name</i>
     --no-tlsverify
     --force
     </pre> 

The `vic-machine create` command in this example specifies the minimum information required to deploy a VCH to vCenter Server:

- The address of the vCenter Server instance on which to deploy the VCH, in the `--target` option.  
- The vCenter Single Sign-On user and password in the `--user` and `--password` options. 
- The port group named `vic-bridge`, for use as the container bridge network. 
- The name of the shared datastore to use as the image store, in which to store container images.
- Disables the verification of clients that connect to this VCH by specifying the `--no-tlsverify` option.
- Disables the verification of the vCenter Server certificate by specifying the `--force` option.
   
Because the vCenter Server instance only has one datacenter and one cluster, and uses the VM Network network, `vic-machine create` automatically detects and uses these resources.

This example deploys a VCH with the default name `virtual-container-host`.

**Result**

At the end of a successful installation, `vic-machine` displays information about the new VCH:
   
<pre>Initialization of appliance successful
VCH Admin Portal:
https://<i>vch_address</i>:2378
Published ports can be reached at:
<i>vch_address</i>
Docker environment variables:
DOCKER_HOST=<i>vch_address</i>:2376
Environment saved in virtual-container-host/virtual-container-host.env
Connect to docker:
docker -H <i>vch_address</i>:2376 --tls info
Installer completed successfully</pre>

**What to Do Next** 

To test your VCH, see [Verify the Deployment of a VCH](verify_vch_deployment.md).
    
For examples of commands to deploy a VCH in various other vSphere configurations, see [Advanced Examples of Deploying a VCH](vch_installer_examples.md). 
