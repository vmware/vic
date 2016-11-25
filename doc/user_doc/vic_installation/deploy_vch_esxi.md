# Deploy a Virtual Container Host to an ESXi Host #

This topic provides instructions for deploying a virtual container host to an ESXi host that is not managed by vCenter Server. This is the most straightforward way to deploy a virtual container host, and is ideal for testing.

**Prerequisites**
* Download and unpack the vSphere Integrated Containers Engine bundle. For information about where to obtain vSphere Integrated Containers Engine, see [Download vSphere Integrated Containers Engine](download_vic_engine.md).
* Add the folder that contains the vSphere Integrated Containers Engine binaries to the `PATH` environment variable on the machine on which you are running `vic-machine`. 
* Create or obtain an ESXi host with the following configuration:
  * One datastore
  * One network, for example the default VM Network
  * You can use a nested ESXi host for this example
* Verify that the ESXi host meets the requirements in [Environment Prerequisites for vSphere Integrated Containers Engine Installation](vic_installation_prereqs.md).
* If your vSphere environment does not use trusted certificates that have been signed by a Certificate Authority (CA), obtain the thumbprint of the vCenter Server or ESXi host certificate.
* Familiarize yourself with the vSphere Integrated Containers Engine binaries, as described in [Contents of the vSphere Integrated Containers Engine Binaries](contents_of_vic_binaries.md). 
* Familiarize yourself with the options of the `vic-machine create` command described in [Virtual Container Host Deployment Options](vch_installer_options.md).

**Procedure**

1. Open a terminal on the system on which you downloaded and unpacked the vSphere Integrated Containers Engine binary bundle.
2. Navigate to the directory that contains the `vic-machine` utility:
3. Run the `vic-machine create` command.

   - Linux OS:
      <pre>$ vic-machine-linux create
     --target <i>esxi_host_address</i>
     --user root
     --password <i>esxi_host_password</i>
     --no-tlsverify
     </pre>  
   - Windows:
      <pre>$ vic-machine-windows create
     --target <i>esxi_host_address</i>
     --user root
     --password <i>esxi_host_password</i>
     --no-tlsverify
     </pre> 

   - Mac OS:
       <pre>$ vic-machine-darwin create
     --target <i>esxi_host_address</i>
     --user root
     --password <i>esxi_host_password</i>
     --no-tlsverify
     </pre> 

The `vic-machine create` command in this example specifies the minimum information required to deploy a virtual container host to an ESXi host:

- The address of the ESXi host on which to deploy the virtual container host, in the `--target` option. 
- The ESXi host `root` user and password in the `--user` and `--password` options. If the password contains special characters, wrap it in quotes. Use single quotes if you are using `vic-machine` on a Linux or Mac OS system and double quotes on a Windows system. 
- Disables the verification of clients that connect to this virtual container host by specifying the `--no-tlsverify` option.
   
Because the ESXi host only has one one network, and one datastore, `vic-machine create` automatically detects and uses those resources. 

When deploying to an ESXi host, `vic-machine create` creates a standard virtual switch and a distributed port group for use as the container bridge network, so you do not need to specify any network options if you do not have specific network requirements.

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
