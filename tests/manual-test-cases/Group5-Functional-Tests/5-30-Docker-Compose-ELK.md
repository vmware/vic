Test 5-30 - Docker Compose ELK
=======

# Purpose:
To verify that VIC appliance can work when deploying the docker ELK services

# References:
[1 - Docker Compose Overview](https://docs.docker.com/compose/overview/)  
[2 - Docker compose ELK](https://blogs.vmware.com/cloudnative/2018/07/19/getting-started-with-elastic-stack-on-vsphere-integrated-containers/)

# Environment:
This test requires access to VMWare Nimbus cluster for dynamic ESXi and vCenter creation

# Test Steps:
1. Deploy a new vCenter with 2 ESXi hosts in a cluster
2. Deploy VIC appliance to the vSphere server 
3. Issue the following command in the docker elk app folder:  
```cd demos/compose/elk-app; COMPOSE_HTTP_TIMEOUT=300 DOCKER_HOST=<VCH IP> docker-compose up```

# Expected Outcome:
Docker compose should return with success and all containers in the compose yaml file are up and running.
Docker inspect data should show networks, alias, and IP address for the container.

# Possible Problems:
None
