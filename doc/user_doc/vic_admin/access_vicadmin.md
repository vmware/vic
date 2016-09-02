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

The VIC Admin portal displays status information about the virtual container host, network, configuration, and license. It also displays the address of the Docker endpoint and the remaining capacity of the datastore. You can access live logs for different aspects of the virtual container host, or download log bundles.