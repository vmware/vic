Test 15-1 Drone Continous Integration
=======

#Purpose:
To verify the VCH appliance can be used as the docker engine replacement for a drone continuous integration environment

#References:
[1- Drone Server Setup](http://readme.drone.io/setup/overview/)

#Environment:
This test requires that a vSphere server is running and available

#Test Steps:
1. Install a new VCH appliance into the vSphere server
2. Pull an image of the drone server:  
```sudo docker pull drone/drone:0.4```
3. Create an oauth app in github, and record the client_id and client_secret for it.  You will need the URL for your server and set the authorization callback URL to http://${url}/authorize
4. Create a dronerc file in /etc/drone/dronerc with the following contents:  
```
REMOTE_DRIVER=github  
REMOTE_CONFIG=https://github.com?client_id=%{CLIENT_ID}&client_secret=%{CLIENT_SECRET}
```
5. Start the server on the VCH appliance:  
```
sudo docker -H ${params} run --volume /var/lib/drone:/var/lib/drone --volume /var/run/docker.sock:/var/run/docker.sock --env-file /etc/drone/dronerc --restart=always --publish=80:8000 --detach=true --name=drone drone/drone:0.4
```
6. Log into the server by navigating to http://localhost
7. Configure at least 1 drone worker on the server
8. Make sure pull request and merge request hooks are configured and allow the server to perform continuous integration on the repository

#Expected Outcome:
The drone server should function properly within the VCH environment and it should be able to execute the CI tests properly on the drone worker that is configured.

#Possible Problems:
None