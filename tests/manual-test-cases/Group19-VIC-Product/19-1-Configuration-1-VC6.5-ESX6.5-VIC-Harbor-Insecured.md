Test 19-1 - Configuration 1 VC6.5 ESX6.5 VIC Harbor Insecured
=======

#Purpose:
To verify the VIC Product (Engine/Harbor) work using insecure registry with a VC6.5 and ESX6.5 server

#References:
[1 - VIC+Harbor Integration Test Plan](https://confluence.eng.vmware.com/pages/viewpage.action?spaceKey=corevc&title=VIC+-+Harbor+Integration+Test+Plan)

#Environment:
* This test requires access to VMWare Nimbus cluster for dynamic ESXi and vCenter creation
* Login to Harbor using an LDAP user1 (say admin role)
* Create a Project named say 'vic-harbor' (publicity off)
* Login with 2 other LDAP users (say user2, user3) and logout - Limitation (Issue#)
* Login as user1 and add other users with different roles (developer, guest) under Project 'vic-harbor'
* Prepare 3 windows client machines (3rd client machine could be a Linux machine as we couldn't figure out the VT-X issue for running docker on Windows 7/10 VM)

#Test Steps:
###pos001:
1. Create a VCH with Harbor private registry as an option (using --insecure-registry <harbor-ip>:80)
2. Do docker login to harbor using an user with admin role
3. Pull an image/application into the local docker registry using local docker in the client machine
4. Tag the image/application to push it to Harbor private registry using local docker in the client machine
5. Push the image/application from local to Harbor Private registry
6. From VCH host, Pull the image from Harbor registry
7. Run the application using docker run command
8. Stop, Start, attach, kill, rm and rmi

###pos002:
1. Reuse the VCH created in test Pos001
2. Do docker login to harbor using an user with developer role
3. Pull an image/application into the local docker registry using local docker in the client machine
4. Tag the image/application to push it to Harbor private registry using local docker in the client machine
5. Push the image/application from local to Harbor Private registry
6. From VCH host, Pull the image from Harbor registry
7. Run the application using docker run command
8. Stop, Start, attach, kill, rm and rmi

###neg001:
1. Reuse the VCH created in test Pos001
2. Do docker login to harbor using an user with guest role
3. Pull an image/application into the local docker registry using local docker in the client machine
4. Tag the image/application to push it to Harbor private registry using local docker in the client machine
5. Push the image/application from local to Harbor Private registry

###pos003:
1. Create a VCH2 with Harbor private registry as an option (using --insecure-registry <harbor-ip>:80)
2. Using VCH1, Do docker login to harbor registry using a developer user
3. Using VCH2, Do docker login to harbor registry using another developer user
4. Using VCH1, Pull an image/application into the local docker registry using local docker in the client machine
5. Using VCH1, Tag the image/application to push it to Harbor private registry using local docker in the client machine
6. Push the image/application from local to Harbor Private registry
7. From VCH host, Pull the image from Harbor registry
8. Run the application using docker run command
9. Stop, Start, attach, kill, rm and rmi

###pos004:
1. Create a VCH with Harbor private registry as an option (using --insecure-registry <harbor-ip>:80) from any client machine
2. Using different machines, connect to the same VCH (in parallel)
3. Do docker login to harbor using different users (in parallel)
4. Pull an image/application into the local docker registry using local docker in the client machine
5. Tag the image/application to push it to Harbor private registry using local docker in the client machine
6. Push the image/application from local to Harbor Private registry
7. From VCH, Pull same/different images using different users from Harbor registry
8. Run the applications using different users from Harbor registry (in parallel)
9. Stop, Start, attach, kill, rm and rmi using different users from Harbor registry (in parallel)

#Expected Outcome:
###pos001:
* VCH Creation with Harbor private registry should succeed - Limitation (Issue#)
* Validate that VCH is communicating with Harbor Registry by doing a docker login (using docker H <vchip>:2375 login <harbor-ip>:80)
* Ensure that the LDAP user with admin role is able to push and pull images to/from the repository to the local registry
* Ensure that the application runs in the containerVm and it is functional
* Ensure that the application tty gets attached and able to delete the ContainerVM and the images

###pos002:
* Validate that VCH is communicating with Harbor Registry by doing a docker login (using docker H <vchip>:2375 login <harbor-ip>:80)
* Ensure that the LDAP user with admin role is able to push and pull images to/from the repository to the local registry
* Ensure that the application runs in the containerVm and it is functional
* Ensure that the application tty gets attached and able to delete the ContainerVM and the images

###neg001:
* Validate that VCH is communicating with Harbor Registry
* Ensure that the LDAP user with guest role is not able to push images to the Harbor repository but able to pull the image

###pos003:
* Validate that VCH is communicating with Harbor Registry by doing a docker login (using docker H <vchip>:2375 login <harbor-ip>:80)
* Ensure that the LDAP user with admin role is able to push and pull images to/from the repository to the local registry
* Ensure that the application runs in the containerVm and it is functional
* Ensure that the application tty gets attached and able to delete the ContainerVM and the images

###pos004:
* Validate that VCH is communicating with Harbor Registry by doing a docker login (using docker H <vchip>:2375 login <harbor-ip>:80)
* Ensure that the LDAP user with admin role is able to push and pull images to/from the repository to the local registry
* Ensure that the application runs in the containerVm and it is functional
* Ensure that the application tty gets attached and able to delete the ContainerVM and the images

#Possible Problems:
* Not planning on using LDAP at this time, in all cases where LDAP is specified substitute users manually created within the product