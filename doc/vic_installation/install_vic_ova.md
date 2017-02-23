# Installing a vSphere Integrated Containers Management Server by Using the OVA Deployment #

The OVA deployment is an alternative means by which you can install vSphere Integrated Containers. 

**IMPORTANT** In this technical preview, the OVA deployment is experimental.

The OVA installer deploys a vSphere Integrated Containers management server. You can use the OVA installer to deploy a vSphere Integrated Containers management server to a vApp, host, or cluster that is managed by vCenter Server.

If you have already installed vSphere Integrated Containers by using the command line installer, you also can deploy the management server OVA. If  the existing vSphere Integrated Containers command line installation is compatible with the OVA installation, you will see your existing virtual container host and any containers and container images in the vSphere Web Client extension.

**Prerequisites**

* Verify that your vSphere infrastructure meets the requirements in [Environment Prerequisites for vSphere Integrated Containers Installation](vic_installation_prereqs.md). 
* Obtain the latest build of the OVA installer, ```management-server.ova```, from your VMware representative.
* Save ```management-server.ova``` into an appropriate folder on your local machine.
* Install a Docker 1.9.1 client. 
 
   **NOTE**: The vSphere Integrated Containers technical preview is not compatible with Docker 1.10. To use vSphere Integrated Containers with a Docker 1.10 client, you must set `DOCKER_API_VERSION=1.21` as an environment variable in your Docker client. 
* Log into the vSphere Web Client for the vCenter Server instance on which to deploy the vSphere Integrated Containers management server.


**Procedure**

1. Select **Hosts and Clusters** from the **Home** menu in the vSphere Web Client.
2. Right-click the vCenter Server instance and select **Deploy OVF Template**.
3. Select **Local file**, click **Browse**, and navigate to the  ```management-server.ova``` file, and click **Next**.

  Deploying the OVA from a locally saved file rather than from a URL has a higher chance of success.  

4. Follow the prompts to review the OVA details and to accept the license.
5. On the Select name and folder page, optionally change the name of the management server VM, select a datacenter or folder in which to install the management server, and click **Next**.
6. Select a vApp, host, or cluster in which to install the management server, and click **Next**. 

  If you have already succeeded with a command line installation, you must select the correct resources when you deploy the OVA, in order for the existing installation to appear in the list of virtual container hosts in the vSphere Web Client extension. For example, if you used the command line installer to install vSphere Integrated Containers in a cluster, you must select the same cluster when you deploy the OVA. 
7. Select a virtual disk format, storage policy, and datastore in which to store the management server VM files, and click **Next**. 
8. Select the network to use for traffic between the management server and vSphere, for example, **VM Network**.
9. Optionally configure the network properties for the management server VM, or leave them blank to use DHCP, and click **Next**.
10. Review the vService bindings, and click **Finish** if the binding status is green.
11. When the OVA deployment finishes, power on the vSphere Integrated Containers managment server VM.

**What to Do Next**
Install the vSphere Web Client Extension for vSphere Integrated Containers.
