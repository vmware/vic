Test 22-01 - redis
=======

# Purpose:
To verify that the redis application on docker hub works as expected on VIC

# References:
[1 - Docker Hub redis Official Repository](https://hub.docker.com/_/redis/)

# Environment:
This test requires that a vSphere server is running and available

# Test Steps:
1. Deploy VIC appliance to the vSphere server


$ docker run --name some-redis -d redis
$ docker run --name some-redis -d redis redis-server --appendonly yes





# Expected Outcome:
* 

# Possible Problems:
None
