# Deploy a VCH to an ESXi Host #

This topic provides instructions for deploying a virtual container host (VCH) to an ESXi host that is not managed by vCenter Server. This is the most straightforward way to deploy a VCH, and is ideal for testing.

**Prerequisites**
* Download and unpack the vSphere Integrated Containers Engine bundle. For information about where to obtain vSphere Integrated Containers Engine, see [Download vSphere Integrated Containers Engine](download_vic_engine.md).
* Create or obtain an ESXi host with the following configuration:
  * One datastore
  * The VM Network is present
  * You can use a nested ESXi host for this example
* Verify that the ESXi host meets the requirements in [Environment Prerequisites for vSphere Integrated Containers Engine Installation](vic_installation_prereqs.md).
* Familiarize yourself with the vSphere Integrated Containers Engine binaries, as described in [Contents of the vSphere Integrated Containers Engine Binaries](contents_of_vic_binaries.md). 
* Familiarize yourself with the options of the `vic-machine create` command described in [VCH Deployment Options](vch_installer_options.md).

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
     --force
     </pre>  
   - Windows:
      <pre>$ vic-machine-windows create
     --target <i>esxi_host_address</i>
     --user root
     --password <i>esxi_host_password</i>
     --no-tlsverify
     --force
     </pre> 
   - Mac OS:
       <pre>$ vic-machine-darwin create
     --target <i>esxi_host_address</i>
     --user root
     --password <i>esxi_host_password</i>
     --no-tlsverify
     --force
     </pre> 

The `vic-machine create` command in this example specifies the minimum information required to deploy a VCH to an ESXi host:

- The address of the ESXi host on which to deploy the VCH, in the `--target` option. 
- The ESXi host `root` user and password in the `--user` and `--password` options. If the password contains special characters, wrap it in quotes. Use single quotes if you are using `vic-machine` on a Linux or Mac OS system and double quotes on a Windows system. 
- Disables the verification of clients that connect to this VCH by specifying the `--no-tlsverify` option.
- Disables the verification of the ESXi host certificate by specifying the `--force` option.
   
Because the ESXi host only has only one datastore and uses the default VM Network network, `vic-machine create` automatically detects and uses those resources. 

When deploying to an ESXi host, `vic-machine create` creates a standard virtual switch and a distributed port group for use as the container bridge network, so you do not need to specify any network options if you do not have specific network requirements.

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
    
For examples of commands to deploy a VCH in various other vSphere configurations, see [Examples of Deploying a VCH](vch_installer_examples.md). 
