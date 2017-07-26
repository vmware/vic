Test 22-01 - nginx
=======

# Purpose:
To verify that the nginx application on docker hub works as expected on VIC

# References:
[1 - Docker Hub nginx Official Repository](https://hub.docker.com/_/nginx/)

# Environment:
This test requires that a vSphere server is running and available

# Test Steps:
1. Deploy VIC appliance to the vSphere server


$ docker run --name some-nginx -v /some/content:/usr/share/nginx/html:ro -d nginx
$ docker run --name some-nginx -d some-content-nginx
$ docker run --name some-nginx -d -p 8080:80 some-content-nginx
$ docker run --name my-custom-nginx-container -v /host/path/nginx.conf:/etc/nginx/nginx.conf:ro -d nginx





# Expected Outcome:
* 

# Possible Problems:
None
