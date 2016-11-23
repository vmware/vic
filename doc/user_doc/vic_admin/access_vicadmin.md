# Access the Administration Portal for a Virtual Container Host #

vSphere Integrated Containers Engine provides a Web-based administration portal for virtual container hosts, called VCH Admin.

**Prerequisites**

- You deployed a virtual container host.
- Obtain the address of the virtual container host:
  - Copy the address from the output of `vic-machine create` or `vic-machine inspect`. 
  - If you deployed the virtual container host to vCenter Server, copy the address from the **Summary** tab for the vSphere Integrated Containers Engine endpoint VM in the vSphere Web Client.
  - If you deployed the virtual container host to an ESXi host, copy the address from the **Summary** tab for the vSphere Integrated Containers Engine endpoint VM in the vSphere Client.

**Procedure**

1. In a Web browser, go to https://<i>vch_address</i>:2378.
2. (Optional) If the virtual container host uses full TLS authentication with trusted CA certificates, enter the username and password for the vCenter Server instance or ESXi host.

**Result**

The VCH Admin portal displays information about the virtual container host and the environment in which is running:  

- Status information about the virtual container host, network,  firewall configuration, and license. For information about these statuses and how to remedy error states, see the [Virtual Container Host Status Reference](vicadmin_status_ref.md).
- The address of the Docker endpoint.
- The remaining capacity of the datastore that you designated as the image store.
- Live logs and log bundles for different aspects of the virtual container host. For information about the logs, see [Access vSphere Integrated Containers Engine Log Bundles](log_bundles.md).