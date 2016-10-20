Test 3-02 - Docker Compose Voting App
=======

#Purpose:
To verify that VIC appliance can work when deploying the example docker voting app

#References:
[1 - Docker Compose Overview](https://docs.docker.com/compose/overview/)  
[2 - Docker Example Voting App](https://github.com/docker/example-voting-app)

#Environment:
This test requires that a vSphere server is running and available

#Test Steps:
1. Deploy VIC appliance to the vSphere server
2. Log into the docker hub
3. Issue the following command in the docker voting app folder:  
```cd demos/compose/voting-app; COMPOSE_HTTP_TIMEOUT=300 DOCKER_HOST=<VCH IP> docker-compose up```

#Expected Outcome:
Docker compose should return with success and the voting and results servers are up and running.

#Possible Problems:
None
