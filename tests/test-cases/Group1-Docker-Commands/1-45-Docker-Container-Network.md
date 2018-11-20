Test 1-45 - Docker Container Network
====================================

# Purpose:
To verify that when containerVM is based on custom iso, the tomcat 
application on docker hub works as expected on VIC. And verify that
tomcat on vic-specific container-network works as expected.

# References:
[1 - Docker Hub tomcat Official Repository](https://hub.docker.com/_/tomcat/)

# Environment:
This test requires that a vSphere server is running and available

# Test Steps:
1. Deploy VIC appliance to the vSphere server with custom iso as containerVM
2. Run an tomcat container with a mapped port and verify the server is up and running:  
`docker run --name tomcat1 -d -p 8080:8080 tomcat:alpine`
3. Run an tomcat container on the specific container network:
`docker run --name tomcat2 -d --net=public tomcat:alpine`
4. Run an tomcat container with a mapped port on the specific container network:
`docker run --name tomcat3 -d -p 8083:8080 --net=public tomcat:alpine`

# Expected Outcome:
* Each step should succeed, tomcat should be running without error in each case

# Possible Problems:
None
