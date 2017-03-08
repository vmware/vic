# Investigation on NSX micro-segmentation via Docker network create
This document is for this issue: https://github.com/vmware/vic/issues/3936


## A Working Installation of NSX-v
For VMWare engineers using Nimbus, here is a document about how to install and configure NSX with Vsphere: How to set up VM-to-VM communication on NSX logical switch with Nimbus
https://confluence.eng.vmware.com/pages/viewpage.action?pageId=209964551

## About NSX Security Groups and Security Policies
This section is a brief introduction about NSX security groups and security policies (some words are taken from this document http://www.vmware.com/content/dam/digitalmarketing/vmware/en/pdf/products/nsx/vmw-nsx-network-virtualization-design-guide.pdf).  
  * Security Groups
    NSX provides various grouping mechanisms :
      * vCenter Objects: VMs, Distributed Switches, Clusters, etc.
      * VM Properties: vNICs, VM names, VM operating Systems, etc.
      * NSX Objects: Logical Switches, Security Tags, Logical Routers, etc. 
    In this document, we propose to use vNICs as the grouping creteria. We will dicuss more on this later.

  * Security Policies
    NSX provides security policies as a way to group rules for security controls that will be applied to one or more security groups.
    Each security policy contains the following :
      * Firewall rules: 
        * NSX built in distributed firewall.
      * weight
        * Weight of a policy determines the rank of the policy versus other policies in the NSX eco-system. Higher weight rules are executed first.

## An Overview of Workflow

Current VIC implementation may have the container bridge networks all reside on a single port group and only be separated by IP space. Our container networks (created using `docker network create`) are not isolated from each other. For example, a containerVM in one container network can change its IP and be able to reach containerVMs in another container network.

Our research shows that we can use NSX security groups and policies on top of current VIC implementation and achieve micro-segmentation. 

Our goals :
  * VMs can communicate with VMs in the same network (security group).
  * VMs cannot talk to VMs in a differrent network (security group).
  * one VM to be added to more than one network (security group).

In the following example, we see that : 
  * containerVM1 and containerVM2 are in default bridge network
  * containerVM3 , containerVM4 and containerVM5 are in bridge network `network-a` 
  * containerVM5 and containerVM6 are in bridge network `network-b`
    (note that adding containerVM5 to both network `network-a` and network `network-b` are not supported on current VIC yet.)

![An Example of How We Apply Security Group on Container Networks] (pics/example-SG-vs-network.png)

Let us use this example to explain the workflow.

### Prerequisite
We assume VSphere admin already has NSX installed and configured on VCenter. That means there is an NSX logical switch created by the admin, which automatically creates a port group on distributed switch.
So this port group can be used as a bridge network for `vic-machine create`

![Distributed Switch with logic switch added] (pics/dswitch.png)

### Create VCH 
Specifying 'bridge-network' as the port group 'vxw-dvs-55-virtualwire-1-sid-5000-logical-switch-1'.
  * VIC creates VCH
  * NSX manager creates a security group `<vch name>-SG-VCH` with dynamic membership: NIC (of VCH) 
  * NSX manager creates these security policies shown in the following picture
  * NSX manager applies security policy `<vch name>-SP-0` on security group `<vch name>-SG-VCH`

  ![Security policies] (pics/security-policies)

### docker run --name=containerVM1 <image name>
  * VIC create containerVM1
  * NSX manager creates security group `<vch name>-SG-bridge` with dynamic membership: NIC-1 (of containerVM1)
  * NSX manager applies all the defined policy groups to this security group

### docker network create network-a
  * VIC creates network scopes for this container bridge network
  * NSX manager creates a security group `<vch name>-SG-network-a` 
  * NSX manager applies all the defined policy groups to this security group  
  (creating network-b is similar are similar to this)

### docker run -net=network-a --name=containerVM3 <image name>
  * VIC creates containerVM3
  * NSX manager update security group `<vch name>-SG-network-a` with dynamic mebership: NIC-3 (of containerVM3)
  (creating other containerVMs are similar to this)

### Summary of security groups and security policies of this example
We associate vNICs to security groups. So the memberships of security groups look like this :
  * `<vch name>-SG-VCH`  dynamic membership: NIC (of VCH) 
  * `<vch name>-SG-bridge`    dynamic membership: NIC-1 (of containerVM1), NIC-2 (of containerVM2)
  * `<vch name>-SG-network-a`    dynamic membership: NIC-3 (of containerVM3), NIC-4 (of containerVM4), NIC-52 (of containerVM5)
  * `<vch name>-SG-network-b`    dynamic membership: NIC-5 (of containerVM5), NIC-6 (of containerVM6)

We apply the following security policies to our security groups:
  * `<vch name>-SP-0` is applied to all the security groups
  * `<vch name>-SP-1` is applied to security groups `<vch name>-SG-bridge`, `<vch name>-SG-network-a` and `<vch name>-SG-network-b` 
  * `<vch name>-SP-2` is applied to security groups `<vch name>-SG-bridge`, `<vch name>-SG-network-a` and `<vch name>-SG-network-b`
  * `<vch name>-SP-3` is applied to security groups `<vch name>-SG-bridge`, `<vch name>-SG-network-a` and `<vch name>-SG-network-b`

The result we get :
  * all the networks are isolated to each other
  * containerVMs in the network can connect to their gateway at VCH

## NSX Manager API:
![Code flow to use NSX Manager API] (pics/with-nsx.png)

### API endpoint
https://{nsxmanager's hostname or IP}/api 

### API calls
(details in http://pubs.vmware.com/nsx-63/topic/com.vmware.ICbase/PDF/nsx_63_api.pdf , note that many NSX API methods reference vCenter object IDs in URI parameters, query parameters, request bodies, and response bodies, also it takes XML as request body (sad) )
* VCH create
  * create security group 
    * POST /2.0/services/securitygroup/bulk/{scopeId}
  * create security policies and apply them to the security group
    * POST /2.0/services/policy/securitypolicy (specify securityGroupBinding in the xml request body)

* docker network create
  * create security group `<vch name>-SG-<network name>`
     * POST /2.0/services/securitygroup/bulk/{scopeId} (For the scopeId use globalroot-0 for non-universal security groups and universalroot-0 for universal security groups.)
  * update security policies to apply them to the security group
     * PUT /2.0/services/policy/securitypolicy/{ID}

* docker network delete 
  * check if any VMs attached to this security group
     * GET /2.0/services/securitygroup/{objectId}/translation/virtualmachine 
  * delete the network (security group, policy group)
     * DELETE /2.0/services/securitytags/tag/{tagId}
     * DELETE /2.0/services/securitygroup/{objectId} 

* docker network list
  * list all the security groups
    * GET /2.0/services/securitygroup/internal/scope/{scopeId}
  * filter security group names by `<vch name>-SG-` prefix

* docker run --net=<network name> <image name>
  * update security group `<vch name>-SG-<network name>` dynamic membership 
    * PUT /2.0/services/securitygroup/bulk/{objectId}


## Discussion:
### Multi-tenants
NSX-v supports multi-tenancy, and here is an example of a NSX network topology:

![NSX-v multi-tenancy] (pics/nsx-multi-tenants.png)

It worths to mention that we would expect only one VCH is created for one tenant. Otherwise there IP overlapping would issue. 

### Associating to Security Group: Security Tag versus vNIC 
  * Security Tag: 
    * Pros:
      * We only need to tag a VM then it is added to a security group. 
      * we do not need to update security groups everytime a new containerVM is created. 
    * Cons:
      * It does not allow a containerVM to be added in multiple networks.  

  * vNIC:
    * Pros:
      * It serves our goals.
    * Cons:
      * Everytime a containerVM is created, we need to update related security groups' by adding one more membership criterias. Not sure the maximum number of membership criterias NSX-v supports. 

### Bridged Containers with Exposed Port or Directly Connected to Public network
  * all the user cases mentioned in the networking design document should still be valid 

## Unclear
### Adding a ContainerVM to Multiple Networks
This is not supported in current VIC implementation. 
  * Is adding another network adapter to a containerVM the right way to do it?
  * Was the challenge not able to add a new network adapter to a running containerVM?
  * if so, maybe after the VSPC work is committed we would be unblocked.

### Details of Integrating with Current Port Layer API 
  * does the current Port Layer API map perfectly to our goals in using NSX ?
  * if not, what changes do we need to make?

### With or Without NSX?
Do we surpport both cases :
  1. User does not have NSX setup, then we continue providing with current VIC
  2. User has NSX setup, then we provide VIC with NSX integration
  * need to find out how to verify if a bridge network specified is a logic switch?

### User Management
  * NSX allows RBAC, we can add Vcenter users and assign roles for them
  * Only security Admin or enterprise admin can operate on NSX security groups and policies. 
  * Need to discuss more with Geroge

