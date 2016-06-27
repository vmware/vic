Test 3-2 - Docker Compose Voting App
=======

#Purpose:
To verify that VIC appliance can work when deploying the example docker voting app

#References:
[1 - Docker Compose Overview](https://docs.docker.com/compose/overview/)
[2 - Docker Example Voting App](https://github.com/docker/example-voting-app)

#Environment:
This test requires that a vSphere server is running and available

#Test Steps:
1. Download the voting app
2. Deploy VIC appliance to the vSphere server
3. Issue DOCKER_HOST=<VCH IP> docker-compose up in the docker voting app folder
4. Verify that the server is running on http://localhost:5000 and the results are found on http://localhost:5001

#Expected Outcome:
Docker compose should return with success and the server and results should be running.

#Possible Problems:
None