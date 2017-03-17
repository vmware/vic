# Access the vSphere Integrated Containers View in the HTML5 vSphere Client #

If you have installed the HTML5 plug-in for vSphere Integrated Containers, you can find information about your vSphere Integrated Containers deployment in the HTML5 vSphere Client.

**IMPORTANT**: Do not use the vSphere Client or to perform operations on virtual container host (VCH) appliances or container VMs. Specifically, using the vSphere Client to power off, power on, or delete VCH appliances or container VMs can cause vSphere Integrated Containers Engine to not function correctly. Always use `vic-machine` to perform operations on VCHs. Always use Docker commands to perform operations on containers.

**Prerequisites**

- You are running vCenter Server 6.5.
- You installed the HTML5 plug-in for vSphere Integrated Containers.

**Procedure**

1. Log in to the HTML5 vSphere Client and go to the **Home** page.
2. Click **vSphere Integrated Containers**.

**Result**

The vSphere Integrated Containers view presents the number of VCHs and container VMs that you have deployed.

**NOTE**: More functionality will be added to the vSphere Integrated Containers view in future releases.
