# Access the Administration Portal for a Virtual Container Host #

vSphere Integrated Containers Engine provides a Web-based administration portal for virtual container hosts, called VCH Admin.

**Prerequisites**

- You deployed a virtual container host.
- Obtain the address of the virtual container host:
  - Copy the address from the output of `vic-machine create` or `vic-machine inspect`. 
  - If you deployed the virtual container host to vCenter Server, copy the address from the **Summary** tab for the vSphere Integrated Containers Engine endpoint VM in the vSphere Web Client.
  - If you deployed the virtual container host to an ESXi host, copy the address from the **Summary** tab for the vSphere Integrated Containers Engine endpoint VM in the vSphere Client.

**Procedure**

In a Web browser, go to https://<i>vch_address</i>:2378

**Result**

The VIC Admin portal displays information about the virtual container host and the environment in which is running:  

- Status information about the virtual container host, network,  firewall configuration, and license. 
- The address of the Docker endpoint.
- The remaining capacity of the datastore that you designated as the image store. 

You can access logs for different aspects of the virtual container host, or download log bundles.