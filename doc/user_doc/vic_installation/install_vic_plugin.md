# Installing the vSphere Client Plug-Ins #

vSphere Integrated Containers Engine provides two UI plug-ins for vSphere:

- A basic Flex-based plug-in that adds information about virtual container hosts (VCHs) and container VMs in the Flex-based vSphere Web Client. You can install the plug-in either on a vCenter Server instance that runs on Windows, or on a vCenter Server Appliance. The basic plug-in works with the Flex-based vSphere Web Client for both vSphere 6.0 and 6.5.
- An HTML5 plug-in with more complete functionality for the HTML5 vSphere Client. The HTML5 vSphere Client is only available with vSphere 6.5. You can deploy the HTML5 plug-in for vSphere Integrated Containers Engine to a vCenter Server Appliance or a vCenter Server instance that runs on Windows, if that instance has access to a Web server. You cannot deploy the HTML5 plug-in to a vCenter Server instance on Windows that does not have access to a Web server.  

For information about the Flex-based vSphere Web Client and the HTML5 vSphere Client for vSphere 6.5, see [Introduction to the vSphere Client](https://pubs.vmware.com/vsphere-65/topic/com.vmware.wcsdk.pg.doc/GUID-3379D310-7802-4B62-8292-D11D928459FC.html) in the vSphere 6.5 documentation.

* [Install the HTML5 Plug-In on a vCenter Server Appliance](plugin_h5_vsca.md)
* [Install the HTML5 Plug-In on vCenter Server for Windows by Using a Web Server](plugin_h5_vc_web.md)
* [Install the Flex Plug-In on vCenter Server for Windows by Using a Web Server](plugin_vc_web.md)
* [Install the Flex Plug-In on vCenter Server for Windows Without Access to a Web Server](plugin_vc_no_web.md)
* [Install the Flex Plug-In on a vCenter Server Appliance by Using a Web Server](plugin_vcsa_web.md)
* [Install the Flex Plug-In on a vCenter Server Appliance Without Access to a Web Server](plugin_vcsa_no_web.md)
* [Verify the Deployment of the Flex Plug-In](plugin_verify_deployment.md)