# Deploy a Virtual Container Host to a vCenter Server Cluster

This topic provides instructions for deploying a virtual container host in a very basic vCenter Server environment. This basic deployment allows you to test vSphere Integrated Containers Engine with vCenter Server before attempting a more complex deployment that corresponds to your real vSphere environment.

The vCenter Server instance to which you deploy the virtual container host must match the specifications listed in the prerequisites.

**Prerequisites**
* Download and unpack the vSphere Integrated Containers Engine bundle. For information about where to obtain vSphere Integrated Containers Engine, see [Download vSphere Integrated Containers Engine](download_vic_engine.md).
* Add the folder that contains the vSphere Integrated Containers Engine binaries to the `PATH` environment variable on the machine on which you are running `vic-machine`. 
* Create or obtain a vCenter Server instance with the following configuration:
  * One datacenter
  * One cluster with two ESXi hosts and DRS enabled. You can use nested ESXi hosts for this example.
  * One datastore
  * One network, for example the default VM Network
  * One distributed virtual switch with one distributed port group named `vic-bridge`.
* Verify that your vCenter Server instance and both of the ESXi hosts in the cluster meet the requirements in [Environment Prerequisites for vSphere Integrated Containers Engine Installation](vic_installation_prereqs.md).
* If your vSphere environment does not use trusted certificates that have been signed by a Certificate Authority (CA), obtain the thumbprint of the vCenter Server or ESXi host certificate.
* Familiarize yourself with the vSphere Integrated Containers Engine binaries, as described in [Contents of the vSphere Integrated Containers Engine Binaries](contents_of_vic_binaries.md). 
* Familiarize yourself with the options of the `vic-machine create` command described in [Virtual Container Host Deployment Options](vch_installer_options.md).
 

**Procedure**

1. Open a terminal on the system on which you downloaded and unpacked the vSphere Integrated Containers Engine binary bundle.
2. Navigate to the directory that contains the `vic-machine` utility:
3. Run the `vic-machine create` command.

   - Linux OS:
      <pre>$ vic-machine-linux create
     --target <i>vcenter_server_address</i>
     --user 'Administrator@vsphere.local'
     --password <i>vcenter_server_password</i>
     --bridge-network vic-bridge
     --no-tlsverify
     </pre>  
   - Windows:
      <pre>$ vic-machine-windows create
     --target <i>vcenter_server_address</i>
     --user "Administrator@vsphere.local"
     --password <i>vcenter_server_password</i>
     --bridge-network vic-bridge
     --no-tlsverify
     </pre> 

   - Mac OS:
       <pre>$ vic-machine-darwin create
     --target <i>vcenter_server_address</i>
     --user 'Administrator@vsphere.local'
     --password <i>vcenter_server_password</i>
     --bridge-network vic-bridge
     --no-tlsverify
     </pre> 

The `vic-machine create` command in this example specifies the minimum information required to deploy a virtual container host to vCenter Server:

- The address of the vCenter Server instance on which to deploy the virtual container host, in the `--target` option. 
- The vCenter Single Sign-On user and password in the `--user` and `--password` options. Note that the user name is wrapped in quotes, because it contains the `@` character. Use single quotes if you are using `vic-machine` on a Linux or Mac OS system and double quotes on a Windows system. 
- The distributed port group named `vic-bridge`, for use as the container bridge network. 
- Disables the verification of clients that connect to this virtual container host by specifying the `--no-tlsverify` option.
   
Because the vCenter Server instance only has one datacenter, one cluster, one network, and one datastore, `vic-machine create` automatically detects and uses these resources.

This example deploys a virtual container host with the default name `virtual-container-host`.

**Result**

If you see the error `Failed to verify certificate for target`, see [VCH Deployment Fails with a Certificate Verification Error](ts_thumbprint_error.md).

At the end of a successful installation, `vic-machine` displays information about the new virtual container host:
   
<pre>Initialization of appliance successful
vic-admin portal:
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

To test your virtual container host, see [Verify the Deployment of a Virtual Container Host](verify_vch_deployment.md).
    
For examples of commands to deploy a virtual container host in various other vSphere configurations, see [Examples of Deploying a Virtual Container Host](vch_installer_examples.md). 
