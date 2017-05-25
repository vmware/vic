Test 19-5 - Configuration 1 VC6.5 ESX6.5 VIC OVA Secured SelfSigned
=======

#Purpose:
To verify the VIC Product (Engine/Harbor) work using secure self-signed registry with a VC6.5 and ESX6.5 server

#References:
[1 - VIC+Harbor Integration Test Plan](https://confluence.eng.vmware.com/pages/viewpage.action?spaceKey=corevc&title=VIC+-+Harbor+Integration+Test+Plan)

#Environment:
* This test requires access to VMWare Nimbus cluster for dynamic ESXi and vCenter creation
* Setup the Unified OVA using the OVFTool on the Nimbus cluster
* Login to Harbor using an LDAP user1 (say admin role)
* Create a Project named say 'vic-harbor' (publicity off)
* Login with 2 other LDAP users (say user2, user3) and logout - Limitation (Issue#)
* Login as user1 and add other users with different roles (developer, guest) under Project 'vic-harbor'
* Prepare 3 windows client machines (3rd client machine could be a Linux machine as we couldn't figure out the VT-X issue for running docker on Windows 7/10 VM)

#Test Steps:
###Test Pos001 Admin Operations:
1. Create a VCH with Harbor private registry as an option (using --registry-ca <harbor-ip>:80)
2. Do docker login to harbor using an user with admin role
3. Pull an image/application into the local docker registry using local docker in the client machine
4. Tag the image/application to push it to Harbor private registry using local docker in the client machine
5. Push the image/application from local to Harbor Private registry
6. From VCH host, Pull the image from Harbor registry
7. Run the application using docker run command
8. Stop, Start, attach, kill, rm and rmi

#Expected Outcome:
###Test Pos001 Admin Operations:
* VCH Creation with Harbor private registry should succeed - Limitation (Issue#)
* Validate that VCH is communicating with Harbor Registry by doing a docker login (using docker H <vchip>:2375 login <harbor-ip>:80)
* Ensure that the LDAP user with admin role is able to push and pull images to/from the repository to the local registry
* Ensure that the application runs in the containerVm and it is functional
* Ensure that the application tty gets attached and able to delete the ContainerVM and the images

###Test Pos002 Developer Operations:
1. Reuse the VCH created in test Pos001
2. Do docker login to harbor using an user with developer role
3. Pull an image/application into the local docker registry using local docker in the client machine
4. Tag the image/application to push it to Harbor private registry using local docker in the client machine
5. Push the image/application from local to Harbor Private registry
6. From VCH host, Pull the image from Harbor registry
7. Run the application using docker run command
8. Stop, Start, attach, kill, rm and rmi

#Expected Outcome:
###Test Pos002 Developer Operations:
* Validate that VCH is communicating with Harbor Registry by doing a docker login (using docker H <vchip>:2375 login <harbor-ip>:80)
* Ensure that the LDAP user with admin role is able to push and pull images to/from the repository to the local registry
* Ensure that the application runs in the containerVm and it is functional
* Ensure that the application tty gets attached and able to delete the ContainerVM and the images

###Test Neg001 Developer Operations:
1. Reuse the VCH created in test Pos001
2. Do docker login to harbor using an user with guest role
3. Pull an image/application into the local docker registry using local docker in the client machine
4. Tag the image/application to push it to Harbor private registry using local docker in the client machine
5. Push the image/application from local to Harbor Private registry

#Expected Outcome:
###Test Neg001 Developer Operations:
* Validate that VCH is communicating with Harbor Registry
* Ensure that the LDAP user with guest role is not able to push images to the Harbor repository but able to pull the image

###Test Pos003 Two VCH With One Harbor:
1. Create a VCH2 with Harbor private registry as an option (using --registry-ca <harbor-ip>:80)
2. Using VCH1, Do docker login to harbor registry using a developer user
3. Using VCH2, Do docker login to harbor registry using another developer user
4. Using VCH1, Pull an image/application into the local docker registry using local docker in the client machine
5. Using VCH1, Tag the image/application to push it to Harbor private registry using local docker in the client machine
6. Push the image/application from local to Harbor Private registry
7. From VCH host, Pull the image from Harbor registry
8. Run the application using docker run command
9. Stop, Start, attach, kill, rm and rmi

#Expected Outcome:
###Test Pos003 Two VCH With One Harbor:
* Validate that VCH is communicating with Harbor Registry by doing a docker login (using docker H <vchip>:2375 login <harbor-ip>:80)
* Ensure that the LDAP user with admin role is able to push and pull images to/from the repository to the local registry
* Ensure that the application runs in the containerVm and it is functional
* Ensure that the application tty gets attached and able to delete the ContainerVM and the images

###Test Pos004 Three Client Machines With One Harbor:
1. Create a VCH with Harbor private registry as an option (using --registry-ca <harbor-ip>:80) from any client machine
2. Using different machines, connect to the same VCH (in parallel)
3. Do docker login to harbor using different users (in parallel)
4. Pull an image/application into the local docker registry using local docker in the client machine
5. Tag the image/application to push it to Harbor private registry using local docker in the client machine
6. Push the image/application from local to Harbor Private registry
7. From VCH, Pull same/different images using different users from Harbor registry
8. Run the applications using different users from Harbor registry (in parallel)
9. Stop, Start, attach, kill, rm and rmi using different users from Harbor registry (in parallel)

#Expected Outcome:
###Test Pos004 Three Client Machines With One Harbor:
* Validate that VCH is communicating with Harbor Registry by doing a docker login (using docker H <vchip>:2375 login <harbor-ip>:80)
* Ensure that the LDAP user with admin role is able to push and pull images to/from the repository to the local registry
* Ensure that the application runs in the containerVm and it is functional
* Ensure that the application tty gets attached and able to delete the ContainerVM and the images

#Possible Problems:
* Not planning on using LDAP at this time, in all cases where LDAP is specified substitute users manually created within the product
