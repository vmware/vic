## Installing Virtual Integrated Containers

The intent is that vSphere Integrated Containers (VIC) should not _require_ an installation step - deploying a [Virtual Container Host](doc/design/arch/vic-container-abstraction.md#virtual-container-host) (VCH) directly without any prior steps should always be possible. At the current time this is the only approach available.

Installation will be required for capabilities such as [self-provisioning](doc/design/validating-proxy.md) and management network isolation via [vmomi proxy](doc/design/vmomi-authenticating-agent.md).

## Deploying a Virtual Container Host

### Requirements

- ESXi/vCenter - the target virtualization environment.
   - ESXi - Enterprise license
   - vCenter - Enterprise plus license, only very simple configurations have been tested. 
- DHCP - the VCH currently requires there be DHCP on the external network (-external-network flag if not "VM Network")
- Bridge network - when installed in a vCenter environment vic-machine does not automatically create a bridge network. An existing vSwitch or Distributed Portgroup should be specified via the -bridge-network flag, and should not be the same as the external network.

Replace the `<fields>` in the example with values specific to your environment - this will install VCH to the specified resource pool of ESXi or vCenter, and the container VMs will be created under that resource pool.

- -compute-resource is the resource pool where VCH will be deployed to, which should be in govc format. Here is one resource pool path sample: `/ha-datacenter/host/localhost/Resources/test`. For users not familar with govc, please check [govc](https://github.com/vmware/govmomi/blob/master/govc) document.
- -generate-cert flag is to generate certificates and configure TLS. 
- -force flag is to remove an existing datastore folder or VM with the same name.

```
vic-machine -target target-host -image-store <datastore name> -name <vch-name> -user root -passwd <password> -compute-resource <resource pool path> -generate-cert
```
This will, if successful, produce output similar to the following:
```
INFO[2016-04-29T20:17:21-05:00] ### Installing VCH ####                      
INFO[2016-04-29T20:17:21-05:00] Generating certificate/key pair - private key in ./vch-name-key.pem 
INFO[2016-04-29T20:17:21-05:00] Validating supplied configuration            
INFO[2016-04-29T20:17:34-05:00] Creating a Resource Pool                     
INFO[2016-04-29T20:17:36-05:00] Creating VirtualSwitch                       
INFO[2016-04-29T20:17:36-05:00] Creating Portgroup                           
INFO[2016-04-29T20:17:37-05:00] Creating appliance on target                 
INFO[2016-04-29T20:17:41-05:00] Uploading images for container               
INFO[2016-04-29T20:17:41-05:00] 	bootstrap.iso 
INFO[2016-04-29T20:17:41-05:00] 	appliance.iso 
INFO[2016-04-29T20:19:15-05:00] Waiting for IP information                   
INFO[2016-04-29T20:19:33-05:00] Initialization of appliance successful       
INFO[2016-04-29T20:19:33-05:00]                                              
INFO[2016-04-29T20:19:33-05:00] SSH to appliance (default=root:password)     
INFO[2016-04-29T20:19:33-05:00] ssh root@x.x.x.x                        
INFO[2016-04-29T20:19:33-05:00]                                              
INFO[2016-04-29T20:19:33-05:00] Log server:                                  
INFO[2016-04-29T20:19:33-05:00] https://x.x.x.x:2378                    
INFO[2016-04-29T20:19:33-05:00]                                              
INFO[2016-04-29T20:19:33-05:00] Connect to docker:                           
INFO[2016-04-29T20:19:33-05:00] docker -H x.x.x.x:2376 --tls --tlscert='./vch-name-cert.pem' --tlskey='./vch-name-key.pem' info 
INFO[2016-04-29T20:19:33-05:00] Installer completed successfully...          
```



[Issues relating to Virtual Container Host deployment](https://github.com/vmware/vic/labels/component%2Fvic-machine)
