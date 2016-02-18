# Install vSphere Integrated Containers by Using the Command Line Installer

The command line installer is the default means by which you install vSphere Integrated Containers. 

The command line installer deploys a vSphere Integrated Containers virtual container host in one of the following setups: 
* A vCenter Server with a cluster
* A vCenter Server with a single host
* A standalone ESXi host. 

The virtual container host allows you to use an ESXi host or cluster as the Docker endpoint for a Docker client.

**NOTE** The command line installer does not deploy a vSphere Integrated Containers management server. It does not add an extension in the vSphere Web Client. 

**Prerequisites**

* Verify that your vSphere infrastructure meets the requirements in [Environment Prerequisites for vSphere Integrated Containers Installation](vsphere_integrated_containers_installation_prerequ.md).
* Download the latest build of the command line installer bundle, ```bonneville.tgz```, from http://bonneville.eng.vmware.com:8080/job/bonneville-package/lastSuccessfulBuild/artifact/binary/.
* Unpack ```bonneville.tgz``` into an appropriate folder on your local machine.
* Familiarize yourself with the installer options described in [vSphere Integrated Containers Installer Options](vsphere_integrated_containers_installer_options.md).
* Familiarize yourself with the contents of the vSphere Integrated Containers installer bundle, as described in [Contents of the vSphere Integrated Containers Command Line Installer Bundle](contents_of_the_vsphere_integrated_containers_comm.md).
* Install a Docker 1.9.1 client. Note that vSphere Integrated Containers is not compatible with Docker 1.10.


**Procedure**

1. Open a command prompt and navigate to the folder in which you unpacked the ```bonneville.tgz``` bundle.
2. Run the installer executable for your operating system. 

   The following examples include the fewest possible options, for installation in a simple environment.

   ```$ install-win.exe -target=<esx_host_or_vcenter_server_address> -user=<esx_host_vcenter_server_username> -passwd="<password>" -name=<name_for_VIC_appliance> -ceip=<enable/disable>```
   
   ```$ install-linux -target=<esx_host_or_vcenter_server_address> -user=<esx_host_vcenter_server_username> -passwd="<password>" -name=<name_for_VIC_appliance> -ceip=<enable/disable>```
   
   ```$ install-osx -target=10.153.199.227-target=<esx_host_or_vcenter_server_address> -user=<esx_host_vcenter_server_username> -passwd="<password>" -name=<name_for_VIC_appliance> -ceip=<enable/disable>```  
   
3. Verify that the vSphere Integrated Containers appliance has initiated correctly. 

   At the end of a successful installation, the installer displays a success message:
   
   ```Initialization of appliance successful.```
   
   ```You can run docker commands via: set DOCKER_HOST=tcp://<virtual_container_host_address>:2376 docker```
   
   ```You can access logs via: https://<virtual_container_host_address>:2378```
   
   ```Installer completed successfully...```
4. (Optional) If you installed onto a vCenter Server instance, log into the vSphere Web Client.
5. (Optional) In the Hosts and Clusters view, navigate to the cluster on which you installed the appliance.

    vSphere Integrated Containers version 0.0.1 is displayed in the **Summary** tab for the appliance vApp. The vApp contains a virtual machine and a virtual machine template, both with the name that you provided during installation. 
6.  (Optional) Select the virtual machine running inside the vApp and click the **Summary** tab.

    The **Notes** panel includes the address of the Docker endpoint to use to run Docker commands.
    
7.  On the machine where you installed the Docker client, connect the Docker client to the virtual container host that is running in your vSphere infrastructure. 
 
  ```$ export DOCKER_HOST=tcp://<virtual_container_host_address>:2376 docker```
8.  Check that the Docker client is connected to the virtual container host.
 
  ```$ docker info```

    You should see confirmation that the Storage Driver is ```vmware-01```. If you installed onto a vCenter Server instance, you should see that the Operating System is ```VMware vCenter Server/6.0.0```. The datastore for containers should be a datastore in your vSphere environment.
9.  Pull a test Docker container image into the virtual container host.

    For example, pull in the BusyBox container.

    ```$ docker pull busybox:latest```
    
10. In the vSphere Web Client, go to **Storage**, and navigate to the datastore that the virtual container host uses to store containers.
11. Select **Manage** > **Files**, and navigate to *datastore_name* > *name_of_VIC_appliance* > **images**. 
 
    You should see a VMDK file for every container image that you have pulled into the virtual container host.

12. Run a Docker container.

    For example, run the BusyBox container.

    ```$ docker run busybox```
    
10. In the **Hosts and Clusters** view of the vSphere Web Client, navigate to the vSphere Integrated Containers appliance, click **Related Objects** > **Virtual Machines**.
 
    You should see a virtual machine for every container that you are running.

    
    