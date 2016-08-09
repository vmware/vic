Test 1-25 - Docker Port Mapping
=======

#Purpose:
To verify that docker create works with the -p option

#References:
[1 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/create/)

#Environment:
This test requires that a vSphere server is running and available

#Test Steps:
1. Deploy VIC appliance to vSphere server
2. Issue docker create -it -p 10000:80 -p 10001:80 --name webserver nginx
3. Issue docker start webserver
4. Issue curl vch-ip:10000 --connect-timeout 20
5. Issue curl vch-ip:10001 --connect-timeout 20
6. Issue docker stop webserver
7. Issue curl vch-ip:10000
8. Issue curl vch-ip:10001
9. Issue docker create -it -p 8083:80 --name webserver2 nginx
10. Issue docker create -it -p 8083:80 --name webserver3 nginx
11. Issue docker start webserver2
12. Issue docker start webserver3
13. Issue docker create -it -p 8081-8088:80 --name webserver5 nginx
14. Issue docker create -it -p 10.10.10.10:8088:80 --name webserver5 nginx
15. Issue docker create -it -p 6379 --name test-redis redis:alpine
16. Issue docker start test-redis
17. Issue docker stop test-redis

#Expected Outcome:
* Steps 2-6 should all return without error
* Steps 7-8 should both return error
* Steps 9-11 should all return without error
* Steps 12-14 should return error
* Steps 15-17 should return without error

#Possible Problems:
None
