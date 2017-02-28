# Investigation on NSX micro-segmentation via Docker network create
This document is for this issue: https://github.com/vmware/vic/issues/3936


## A Working Installation of NSX-v
I documented how to install and configure NSX with Vsphere: How to set up VM-to-VM communication on NSX logical switch with Nimbus
https://confluence.eng.vmware.com/pages/viewpage.action?pageId=209964551

## An Overview of Workflow

### Prerequisite
VSphere admin already has NSX installed and configured on VCenter. That means there is an NSX logical switch created by the admin, which automatically creates a port group on distributed switch.
So this port group can be used as a bridge network for `vic-machine create`

![Distributed Switch with logic switch added] (pics/dswitch.png)

### Create VCH 
specifying 'bridge-network' as the port group 'vxw-dvs-55-virtualwire-1-sid-5000-logical-switch-1':
  * I use this VCH to create containers, they are automatically added to this port group (i.e. logical switch)
  * containerVMs get their IPs in the internal network
  * containerVMs can ping each other

### docker network create network-a
Our goal:
  * VMs can communicate with VMs in the same security group.
  * VMs cannot talk to VMs in a differrent security group.
  * one VM to be added to more than one security group.

To achieve our goal, use NSX security/policy groups:
  * NSX manager create three policy groups:
     * `<vch name>-PolicyGroup-rule-1`: weight 8300, Policy's Security Group to Policy's Security Group :  allow
     * `<vch name>-PolicyGroup-rule-2`: weight 7300, Any to Policy's Security Group : block
     * `<vch name>-PolicyGroup-rule-3`: weight 6300, Policy's Security Group to Any : block
  * NSX manager create a security tag `<vch name>-SecurityTag-network-a`
  * NSX manager create a security group `<vch name>-SecurityTag-network-a` (with members: security tag == `<vch name>-SecurityTag-network-a`))
  * NSX manager apply the policy groups to this security group

### docker run -net=network-a
  * create a containerVM 
  * NSX manager assign security tag `<vch name>-SecurityTag-network-a` to this container

![how security groups and policy groups can be used to isolate container networks] (pics/security-group.png)

## NSX Manager API:
![Code flow to use NSX Manager API] (pics/with-nsx.png)

### API endpoint
https://{nsxmanager's hostname or IP}/api 

### API calls
(details in http://pubs.vmware.com/nsx-63/topic/com.vmware.ICbase/PDF/nsx_63_api.pdf , note that many NSX API methods reference vCenter object IDs in URI parameters, query parameters, request bodies, and response bodies, also it takes XML as request body (sad) )
* docker network create
  * create security tag
     * POST /2.0/services/securitytags/tag

  * create security group withe security tag as member
     * POST /2.0/services/securitygroup/bulk/{scopeId} (For the scopeId use globalroot-0 for non-universal security groups and universalroot-0 for universal security groups.)

  * create security policy and apply it to security group
     * POST /2.0/services/policy/securitypolicy (specify securityGroupBinding in the xml request body)

* docker network delete 
  * check if any VMs attached to this security group
     * GET /2.0/services/securitygroup/{objectId}/translation/virtualmachine 
  * delete the network (security group, policy group, security tag)
     * DELETE /2.0/services/securitytags/tag/{tagId}
     * DELETE /2.0/services/securitygroup/{objectId} 
     * DELETE /2.0/services/policy/securitypolicy/{ID} 

* connect container to the network 
  * apply security tag to a virtual machine (containerVM)  
    * POST /2.0/services/securitytags/tag/{tagId}/vm 

* docker network list
  * list all the security groups
    * GET /2.0/services/securitygroup/internal/scope/{scopeId}
  * filter security group names by `<vch name>-SecurityGroup` prefix
  * or if we keep KV pairs of each created network, then we just need to use that to list all the networks

## State Storage
  * we can map a docker network name to a security group name as `<vch name>-SecurityGroup-<network name>`
  * we can store KV pairs for each created network and then everything else will be in the NSX control plane

## Unclear / need discussion:

### Security policy does not work with current containerVMs 
  * currently if I create VCH with 'logical-switch' as bridge network, and then create containerVMs. Even if put them in different security groups, they can still ping each other. Need to investigate this. Is this related to IP management? Need to understand more on how NSX policy group works.
  * I tried using 'ifconfig eth0 192.168.10.90 netmask 255.255.255.0 up' to change the containerVMs' IP, but security policy still does not work. (the containerVMs in different security groups can still poing each other)

### Details of integrating with current Port Layer API 
  * does the current Port Layer API map perfectly to our goals in using NSX ?
  * if not, what changes do we need to make?

### Bridged containers with exposed port
  * I have not looked into how/if this can work with NSX yet

### To support both cases: with or without NSX?
  * need to find out how to verify if a bridge network specified is a logic switch?

### User management
  * NSX allows RBAC, we can add Vcenter users and assign roles for them
  * to my understanding, in terms of using NSX for 'docker network create', we will not need RBAC controls? Unless if we want to restrict users from creating networks, otherwise we may not need RBAC at this point?



