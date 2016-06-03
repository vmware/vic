# Contents of the vSphere Integrated Containers Binaries 

After you download and unpack a vSphere Integrated Containers binary bundle from https://bintray.com/vmware/vic-repo/build/view#files, you obtain following files:

| **File** | **Description** |
| --- | --- |
|```appliance.iso``` | The ISO from which a virtual container host appliance boots.|
|```bootstrap.iso``` | A Photon OS kernel from which container VMs boot.|
|```vic_machine``` | The command line installation and management utility for virtual container hosts.| 
|```vic_machine-darwin``` | The Mac OS command line installation and management utility for virtual container hosts.| 
|```vic_machine-linux``` | The Linux command line installation and management utility for virtual container hosts.| 
|```vic_machine-windows``` | The Windows command line installation and management utility for virtual container hosts.| 
|`README`|Contains a link to the vSphere Integrated Containers repository on GitHub.|
|`LICENSE`|The license file for vSphere Integrated Containers|

If you build the vSphere Integrated Containers binaries manually, you find the ISO files and the ```vic_machine``` utility in the `<git_installation_dir>/vic/bin` folder.