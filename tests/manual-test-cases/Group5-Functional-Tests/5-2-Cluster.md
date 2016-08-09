Test 5-2 - Cluster
=======

#Purpose:
To verify the VIC appliance works in when the vCenter appliance is using a cluster

#References:
[1 - VMware vCenter Server Availability Guide](http://www.vmware.com/files/pdf/techpaper/vmware-vcenter-server-availability-guide.pdf)

#Environment:
This test requires access to VMWare Nimbus cluster for dynamic ESXi and vCenter creation

#Test Steps:
1. Deploy a new vCenter in Nimbus
2. Create a new datacenter:  
```govc datacenter.create ha-datacenter```
3. Create a new cluster:  
```govc cluster.create cluster1```
4. Add 2 ESXi hosts to the cluster:  
```govc cluster.add -hostname=<ESXi IP> -username=<USER> -cluster=cluster1 -password=<PW> -noverify=true```
5. Deploy VCH Appliance to the new vCenter cluster:    
```bin/vic-machine-linux create --target=<VC IP> --user=Administrator@vsphere.local --image-store=datastore1 --appliance-iso=bin/appliance.iso --bootstrap-iso=bin/bootstrap.iso --generate-cert=false --password=Admin\!23 --force=true --bridge-network=bridge --compute-resource=/ha-datacenter/host/cluster1/<ESXi IP 1>/Resources --external-network=vm-network --name=VCH-test```
6. Run a variety of docker commands on the VCH appliance

#Expected Outcome:
The VCH appliance should deploy without error and each of the docker commands executed against it should return without error

#Possible Problems:
None
