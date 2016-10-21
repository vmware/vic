Test 3-03 - Docker Compose Basic
=======

#Purpose:
To verify that VIC appliance can work when deploying the most basic example from docker documentation

#References:
[1 - Docker Compose Getting Started](https://docs.docker.com/compose/gettingstarted/)

#Environment:
This test requires that a vSphere server is running and available

#Test Steps:
1. Create a compose file that includes a basic python server and redis server
2. Deploy VIC appliance to the vSphere server
3. Issue:  
```DOCKER_HOST=<VCH IP> docker-compose up```

#Expected Outcome:
* Docker compose should return with success and the server should be running.
* The server should report the following output:
```
The server is now ready to accept connections
```

#Possible Problems:
None