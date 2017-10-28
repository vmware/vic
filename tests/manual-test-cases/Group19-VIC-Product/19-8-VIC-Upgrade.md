Test 19-8 VIC Upgrade
=======

# Purpose:
To verify VIC upgrade works and provides access to all of the vSphere Integrated Containers components.

# References:
[Official OVA download](https://my.vmware.com/en/web/vmware/info/slug/datacenter_cloud_infrastructure/vmware_vsphere_integrated_containers/1_2)

[GCP Releases](https://console.cloud.google.com/storage/browser/vic-product-ova-releases?project=eminent-nation-87317&authuser=0)

[VIC Engine Releases](https://console.cloud.google.com/storage/browser/vic-engine-releases?project=eminent-nation-87317&authuser=0)

[OVA Smoke Test](https://confluence.eng.vmware.com/pages/viewpage.action?pageId=242300713)

# Common Environment Setup:
* These tests requires that a vCenter server is running and available with DRS enabled.

# Test Template
### Test Steps:
1. Deploy **[Deploy Old OVA Link]** previous version **[OVA Old Version]** of OVA of the vSphere Integrated Containers appliance
2. Go to URL http://<appliance_IP> in browser
3. Login to management portal using the credentials set during OVA deployment
4. Go to **[Root Cert Path]** and download root certificate file
5. Get respective version of VIC engine binaries by downloading from the link in the Getting Started page (`https://<appliance_IP>:9443`)
6. Copy downloaded root cert to the same dir as vic engine binaries
7. On vic engine binaries dir, run **[VCH Create Command]** to create VCH
8. Go back to management portal, add newly created VCH as a new host
9. Add harbor registry to management portal
    - Go to **[Add Registry Path]**
    - Provide appliance IP as "Address"
    - Create new credentials for harbor and use it
10. Copy root cert to docker dir
    - `sudo mkdir -p /etc/docker/certs.d/<appliance_ip>`
    - `sudo mv ca.crt /etc/docker/certs.d/<appliance_ip>/ca.crt`
11. Tag and Push docker image to harbor registry
    - `docker pull busybox`
    - `docker tag busybox <appliance_ip>/<default_harbor_project>/busybox:test`
    - `docker login <appliance_ip> --username <username_set_during_deployment> --password <password_set_during_deployment>`
    - `docker push <appliance_ip>/<default_harbor_project>/busybox:test`
12. Create containers from the image in harbor
    - Go to **[Create Container Path]**
    - Search for recently pushed image in harbor "<appliance_ip>/<default_harbor_project>/busybox:test" and create container using it
    - Make sure container is running
13. Create container from Docker hub
    - Go to **[Create Container Path]**
    - Search for photon image and create container using it
    - Make sure container is running
14. Shut down guest OS (don't power off) deployed version **[OVA Old Version]** of OVA
15. Deploy **[Deploy New OVA Link]** new version **[OVA New Version]** of OVA of the vSphere Integrated Containers appliance (Don't power on yet)
16. Follow procedure steps in **[Upgrade OVA Link]** to swap disks and upgrade OVA
17. Power on new version of deployed OVA appliance and note down the IP
18. Go to URL http://<appliance_IP> in browser
19. Verify that the VCH added as a host in step 8 should be present in the upgraded OVA.
20. Verify that harbor registry added in step 9 is available
21. Verify that the image pushed in step 11 is available in harbor registry
22. Verify that the containers created in steps 12 and 13 are still available and running
23. Cleanup
    - Remove each of the containers created
    - Delete VCH
    - Delete new and old OVA VM

### Expected Outcome:
* All steps should succeed without error

### Possible Problems:
* None


Test - Upgrade from OVA 1.1.1 to 1.2.0
======================================
**[Deploy Old OVA Link]** [Deploy Old OVA Link](https://vmware.github.io/vic-product/assets/files/html/1.1/vic_vsphere_admin/deploy_vic_appliance.html)

**[OVA Old Version]** 1.1.1

**[Root Cert Path]** Admin > Download Root Cert

**[VCH Create Command]** `./vic-machine-linux create -t <vcenter_IP> -u <vcenter_username> -p <vcenter_password> --image-store <datastore> --bridge-network <bridge_network> --public-network <network> --registry-ca ca.crt --no-tlsverify --force`

**[Add Registry Path]** Templates > Registries > Add

**[Create Container Path]** Templates > Containers > Create

**[Deploy New OVA Link]** [Deploy New OVA Link](https://vmware.github.io/vic-product/assets/files/html/1.2/vic_vsphere_admin/deploy_vic_appliance.html)

**[OVA New Version]** 1.2.0

**[Upgrade OVA Link]** [Upgrade OVA Link](https://vmware.github.io/vic-product/assets/files/html/1.2/vic_vsphere_admin/upgrade_appliance.html)


Test - Upgrade from OVA 1.1.1 to 1.2.1
======================================
**[Deploy Old OVA Link]** [Deploy Old OVA Link](https://vmware.github.io/vic-product/assets/files/html/1.1/vic_vsphere_admin/deploy_vic_appliance.html)

**[OVA Old Version]** 1.1.1

**[Root Cert Path]** Admin > Download Root Cert

**[VCH Create Command]** `./vic-machine-linux create -t <vcenter_IP> -u <vcenter_username> -p <vcenter_password> --image-store <datastore> --bridge-network <bridge_network> --public-network <network> --registry-ca ca.crt --no-tlsverify --force`

**[Add Registry Path]** Templates > Registries > Add

**[Create Container Path]** Templates > Containers > Create

**[Deploy New OVA Link]** [Deploy New OVA Link](https://vmware.github.io/vic-product/assets/files/html/1.2/vic_vsphere_admin/deploy_vic_appliance.html)

**[OVA New Version]** 1.2.1

**[Upgrade OVA Link]** [Upgrade OVA Link](https://vmware.github.io/vic-product/assets/files/html/1.2/vic_vsphere_admin/upgrade_appliance.html)


Test - Upgrade from OVA 1.2.0 to 1.2.1
======================================

> Note: Ignore step 9 (do not apply here)

**[Deploy Old OVA Link]** [Deploy Old OVA Link](https://vmware.github.io/vic-product/assets/files/html/1.1/vic_vsphere_admin/deploy_vic_appliance.html)

**[OVA Old Version]** 1.2.0

**[Root Cert Path]** Administration > Configuration > Registry Root Certificate

**[VCH Create Command]** `./vic-machine-linux create -t <vcenter_IP> -u <vcenter_username> -p <vcenter_password> --image-store <datastore> --bridge-network <bridge_network> --public-network <network> --registry-ca ca.crt --no-tlsverify --force`

**[Add Registry Path]** Administration > Registries > +Registry

**[Create Container Path]** Home > Containers > +Container

**[Deploy New OVA Link]** [Deploy New OVA Link](https://vmware.github.io/vic-product/assets/files/html/1.2/vic_vsphere_admin/deploy_vic_appliance.html)

**[OVA New Version]** 1.2.1

**[Upgrade OVA Link]** [Upgrade OVA Link](https://vmware.github.io/vic-product/assets/files/html/1.2/vic_vsphere_admin/upgrade_appliance.html)


Test - Upgrade from OVA 1.2.0 to 1.3.0
======================================

> Note: Ignore step 9 (do not apply here)

**[Deploy Old OVA Link]** [Deploy Old OVA Link](https://vmware.github.io/vic-product/assets/files/html/1.1/vic_vsphere_admin/deploy_vic_appliance.html)

**[OVA Old Version]** 1.2.0

**[Root Cert Path]** Administration > Configuration > Registry Root Certificate

**[VCH Create Command]** `./vic-machine-linux create -t <vcenter_IP> -u <vcenter_username> -p <vcenter_password> --image-store <datastore> --bridge-network <bridge_network> --public-network <network> --registry-ca ca.crt --no-tlsverify --force`

**[Add Registry Path]** Administration > Registries > +Registry

**[Create Container Path]** Home > Containers > +Container

**[Deploy New OVA Link]** [Deploy New OVA Link](TBD - waiting on official docs to link)

**[OVA New Version]** 1.3.0

**[Upgrade OVA Link]** [Upgrade OVA Link](TBD - waiting on official docs to link)


Test - Upgrade from OVA 1.2.1 to 1.3.0
======================================

> Note: Ignore step 9 (do not apply here)

**[Deploy Old OVA Link]** [Deploy Old OVA Link](https://vmware.github.io/vic-product/assets/files/html/1.1/vic_vsphere_admin/deploy_vic_appliance.html)

**[OVA Old Version]** 1.2.1

**[Root Cert Path]** Administration > Configuration > Registry Root Certificate

**[VCH Create Command]** `./vic-machine-linux create -t <vcenter_IP> -u <vcenter_username> -p <vcenter_password> --image-store <datastore> --bridge-network <bridge_network> --public-network <network> --registry-ca ca.crt --no-tlsverify --force`

**[Add Registry Path]** Administration > Registries > +Registry

**[Create Container Path]** Home > Containers > +Container

**[Deploy New OVA Link]** [Deploy New OVA Link](TBD - waiting on official docs to link)

**[OVA New Version]** 1.3.0

**[Upgrade OVA Link]** [Upgrade OVA Link](TBD - waiting on official docs to link)


Test - Upgrade from OVA 1.1.1 to 1.2.0 to 1.2.1 to 1.3.0
========================================================

> Note: Ignore step 23 (do not run it for this test) and continue with additional test steps in this test.

**[Deploy Old OVA Link]** [Deploy Old OVA Link](https://vmware.github.io/vic-product/assets/files/html/1.1/vic_vsphere_admin/deploy_vic_appliance.html)

**[OVA Old Version]** 1.1.1

**[Root Cert Path]** Admin > Download Root Cert

**[VCH Create Command]** `./vic-machine-linux create -t <vcenter_IP> -u <vcenter_username> -p <vcenter_password> --image-store <datastore> --bridge-network <bridge_network> --public-network <network> --registry-ca ca.crt --no-tlsverify --force`

**[Add Registry Path]** Templates > Registries > Add

**[Create Container Path]** Templates > Containers > Create

**[Deploy New OVA Link]** [Deploy New OVA Link](https://vmware.github.io/vic-product/assets/files/html/1.2/vic_vsphere_admin/deploy_vic_appliance.html)

**[OVA New Version]** 1.2.0

**[Upgrade OVA Link]** [Upgrade OVA Link](https://vmware.github.io/vic-product/assets/files/html/1.2/vic_vsphere_admin/upgrade_appliance.html)

### Additional Test Steps:
a.  Shut down previous version 1.2.0 of OVA appliance

b.  Rerun steps 15 to 22 with the following variables:

- **[Deploy New OVA Link]** [Deploy New OVA Link](https://vmware.github.io/vic-product/assets/files/html/1.2/vic_vsphere_admin/deploy_vic_appliance.htm)

- **[OVA New Version]** 1.2.1

- **[Upgrade OVA Link]** [Upgrade OVA Link](https://vmware.github.io/vic-product/assets/files/html/1.2/vic_vsphere_admin/upgrade_appliance.html)

c.  Shut down previous version 1.2.1 of OVA appliance

b.  Rerun steps 15 to 23 with the following variables:

- **[Deploy New OVA Link]** [Deploy New OVA Link](TBD - waiting on official docs to link)

- **[OVA New Version]** 1.3.0

- **[Upgrade OVA Link]** [Upgrade OVA Link](TBD - waiting on official docs to link)