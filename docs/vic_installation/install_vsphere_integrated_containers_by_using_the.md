# Install vSphere Integrated Containers by Using the Command Line Installer

The command line installer is the default means by which you install vSphere Integrated Containers. 

The command line installer deploys a vSphere Integrated Containers virtual container host in one of the following setups: 
* A vCenter Server with a cluster
* A vCenter Server with a single host
* A standalone ESXi host. 

The virtual container host allows you to use an ESXi host or cluster as the Docker endpoint for a Docker client.

**NOTE** The command line installer does not deploy a vSphere Integrated Containers management server. It does not add an extension in the vSphere Web Client. 

**Prerequisites**

* Verify that your vSphere infrastructure meets the requirements in [vSphere Integrated Containers Installation Prerequisites](vsphere_integrated_containers_installation_prerequ.md).

**Procedure**

1. Download the latest build of the command line installer bundle, ```bonneville.tgz```, from http://bonneville.eng.vmware.com:8080/job/bonneville-package/lastSuccessfulBuild/artifact/binary/.
2. Unpack ```bonneville.tgz``` into an appropriate folder on your local machine.
3. Open a command prompt and navigate to the folder in which you unpacked the ```bonneville.tgz``` bundle.
4. Run the appropriate installer executable for your operating system. 

   Use the appropriate options depending on the type of vSphere infrastructure onto which you are deploying vSphere Integrated Containers.  
   

  
**Examples**
To be added here.