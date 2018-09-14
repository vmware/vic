---
name: "Defect Report"
about: Report something that isn't working as expected

---

<!--
This repository is for VIC Engine. Please use it to report issues related to Virtual Container Hosts, Container VMs, and their lifecycles.

To help use keep things organized, please file issues in the most appropriate repository:
 * vSphere Client Plugins: https://github.com/vmware/vic-ui/issues
 * VIC Appliance (OVA) and User Documentation: https://github.com/vmware/vic-product/issues
 * Container Management Portal (Admiral): https://github.com/vmware/admiral/issues
 * Container Registry (Harbor): https://github.com/goharbor/harbor/issues
-->

#### Summary
<!-- Explain the problem briefly. -->


#### Environment information
<!-- Describe the environment where the issue occurred. -->

##### vSphere and vCenter Server version
<!-- Indicate the vSphere and vCenter Server version(s) being used. -->

##### VIC version
<!-- Indicate the full VIC version being used (e.g., vX.Y.Z-tag-NNNN-abcdef0). -->

##### VCH configuration
<!-- Provide the settings used to deploy the VCH (e.g., a vic-machine create command). -->


#### Details
<!-- Provide additional details. -->

##### Steps to reproduce

##### Actual behavior

##### Expected behavior


#### Logs
<!--
For issues related to a deployed VCH, please attach a log bundle.
 * If you can access the VCH Admin portal, please download and attach the log bundle(s). See https://vmware.github.io/vic/assets/files/html/vic_admin/log_bundles.html for details.
 * If the VCH Admin portal is inaccessible, you can enable SSH to the VCH endpoint VM to obtain logs manually. See https://vmware.github.io/vic/assets/files/html/vic_admin/vch_ssh_access.html for details. The VCH logs will be under /var/log/vic/ on the VM.

For issues deploying or managing a VCH, please *also* attach the vic-machine log file.
 * When using the vic-machine CLI, the vic-machine.log file is written to your current working directory.
 * When using VCH Management API, please download and attach the vic-machine-server.log file. See https://vmware.github.io/vic-product/assets/files/html/1.4/vic_vsphere_admin/appliance_logs.html for details.
-->


#### See also
<!-- Provide references to relevant resources, such as documentation or related issues. -->


#### Troubleshooting attempted
<!-- Use this section to indicate steps you've already taken to troubleshoot the issue. -->

- [ ] Searched [GitHub][issues] for existing issues. (Mention any similar issues under "See also", above.)
- [ ] Searched the [documentation][docs] for relevant troubleshooting guidance.
- [ ] Searched for a relevant [VMware KB article][kb].

<!-- Reference-style links used above; removing these will break the links. -->
[issues]:https://github.com/vmware/vic/issues
[docs]:https://vmware.github.io/vic-product/#documentation
[kb]:https://kb.vmware.com/s/global-search/%40uri#t=Knowledge&sort=relevancy&f:@commonproduct=[vSphere%20Integrated%20Containers]
