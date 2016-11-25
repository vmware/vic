# Contents of the vSphere Integrated Containers Engine Binaries 

After you download and unpack a vSphere Integrated Containers Engine binary bundle, you obtain following files:

| **File** | **Description** |
| --- | --- |
|`appliance.iso` | The ISO from which a virtual container host (VCH) appliance boots.|
|`bootstrap.iso` | A Photon OS kernel from which container VMs boot.|
|`ui/` | A folder that contains the files and scripts for the deployment of the vSphere Web Client Plug-in for vSphere Integrated Containers Engine.| 
|`vic-machine-darwin` | The Mac OS command line utility for the installation and management of VCHs.| 
|`vic-machine-linux` | The Linux command line utility for the installation and management of VCHs.| 
|`vic-machine-windows.exe` | The Windows command line utility for the installation and management of VCHs.| 
|`vic-ui-darwin` | The Mac OS executable for the deployment of the vSphere Web Client Plug-in for vSphere Integrated Containers Engine. <br><br> **NOTE**: Do not run this executable directly.<sup>(1)</sup>| 
|`vic-ui-linux` | The Linux executable for the deployment of the vSphere Web Client Plug-in for vSphere Integrated Containers Engine.  <br><br> **NOTE**: Do not run this executable directly.<sup>(1)</sup>| 
|`vic-ui-windows.exe` | The Windows executable for the deployment of the vSphere Web Client Plug-in for vSphere Integrated Containers Engine.  <br><br> **NOTE**: Do not run this executable directly.<sup>(1)</sup>| 
|`README`|Contains a link to the vSphere Integrated Containers Engine repository on GitHub.|
|`LICENSE`|The license file for vSphere Integrated Containers Engine|

If you build the vSphere Integrated Containers Engine binaries manually, you find the ISO files and the `vic_machine` utility in the `<git_installation_dir>/vic/bin` folder.

<sup>(1)</sup> For information about how to install the vSphere Integrated Containers Engine client plug-in, see [Installing the vSphere Web Client Plug-in for vSphere Integrated Containers Engine](install_vic_plugin.md).