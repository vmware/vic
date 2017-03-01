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
  * NSX manager create an IP sets  `<vch name>-IPSets-network-a` with the subnet of this network (or Security Tag `<vch name>-SecurityTag-network-a` and add it to containerVMs)
  * NSX manager create a security group `<vch name>-SecurityGroup-network-a` (with members: IP Sets == `<vch name>-IPSets-network-a` or Security Tag == `<vch name>-SecurityTag-network-a`))
  * NSX manager apply the policy groups to this security group

### docker run -net=network-a
  * create a containerVM 

![how security groups (with IP Sets as membership condition) and policy groups can be used to isolate container networks] (pics/security-group-with-ip-sets.png)

![how security groups (with Security Tag as membership condition) and policy groups can be used to isolate container networks] (pics/security-group-with-ip-sets.png)

## NSX Manager API:
![Code flow to use NSX Manager API] (pics/with-nsx.png)

### API endpoint
https://{nsxmanager's hostname or IP}/api 

### API calls
(details in http://pubs.vmware.com/nsx-63/topic/com.vmware.ICbase/PDF/nsx_63_api.pdf , note that many NSX API methods reference vCenter object IDs in URI parameters, query parameters, request bodies, and response bodies, also it takes XML as request body (sad) )
* docker network create
  * create IP Sets
     * POST /2.0/services/ipset/{scopeMoref}
  * or create Security Tag
     * POST /2.0/services/securitytags/tag
  * create security group withe IP Sets or Security Tag as members
     * POST /2.0/services/securitygroup/bulk/{scopeId} (For the scopeId use globalroot-0 for non-universal security groups and universalroot-0 for universal security groups.)
  * create security policy and apply it to security group
     * POST /2.0/services/policy/securitypolicy (specify securityGroupBinding in the xml request body)

* docker network delete 
  * check if any VMs attached to this security group
     * GET /2.0/services/securitygroup/{objectId}/translation/virtualmachine 
  * delete the network (security group, policy group, IP Sets)
     * DELETE /2.0/services/securitytags/tag/{tagId}
     * DELETE /2.0/services/securitygroup/{objectId} 
     * DELETE /2.0/services/ipset/{ipsetId}

* docker network list
  * list all the security groups
    * GET /2.0/services/securitygroup/internal/scope/{scopeId}
  * filter security group names by `<vch name>-SecurityGroup` prefix
  * or if we keep KV pairs of each created network, then we just need to use that to list all the networks

* docker run <image> --net=<network-name>
  * create a containerVM
  * if using security tag as the membership condition, need to assign security tag to the containerVM
    * POST /2.0/services/securitytags/vm/{vmId}

## State Storage
  * we can map a docker network name to a security group name as `<vch name>-SecurityGroup-<network name>`
  * we can store KV pairs for each created network and then everything else will be in the NSX control plane

## Unclear / need discussion:
### Security Tag v.s IP Sets
  * Security Tag we only need to add the tag to a VM then it is added to a security group, which makes adding VMs to multiple security groups convinient. We do not need to update security groups everytime we create a new containerVM in a network.
  * With IP Sets, we need to add the VM's IP to an IP Set. The downside is that 1. it requires the VM's IP to be static. 2. if we have more than one vch in the environment, containerVMs' IPs overlap, we cannot allow one VM in multiple security groups. 3. it may limit the number of VMs added to an IP Set (how many IPs can be added one IP set?). 4. everytime adding a new containerVM in a network, we need to add its IP to the IP set.
  * I think Security tag is more convinient than IP Sets, if we can figure out why security tag does not work on current containerVMs. Otherwise, using IP Sets allows using what we already have and moving fast as long as we are ok with its limitation. 
  * I tried MAC Sets or Virtual Machine as security group members. They do not work for current containerVM either.

### Details of integrating with current Port Layer API 
  * does the current Port Layer API map perfectly to our goals in using NSX ?
  * if not, what changes do we need to make?

### Bridged containers with exposed port or directly connected to public network
  * I have not looked into how/if this can work with NSX yet
  * all the user cases mentioned in the networking design document should still be valid or improved

### To support both cases: with or without NSX?
  * need to find out how to verify if a bridge network specified is a logic switch?

### User management
  * NSX allows RBAC, we can add Vcenter users and assign roles for them
  * to my understanding, in terms of using NSX for 'docker network create', we will not need RBAC controls? Unless if we want to restrict users from creating networks, otherwise we may not need RBAC at this point?



