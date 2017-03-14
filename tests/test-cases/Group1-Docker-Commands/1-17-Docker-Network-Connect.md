Test 1-17 - Docker Network Connect
=======

#Purpose:
To verify that docker network connect command is supported by VIC appliance

#References:
[1 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/network_connect/)

#Environment:
This test requires that a vSphere server is running and available

#Test Steps:
1. Deploy VIC appliance to vSphere server
2. Issue docker network create test-network
3. Issue docker create busybox ifconfig
4. Issue docker network connect test-network <containerID>
5. Issue docker start <containerID>
6. Issue docker logs <containerID>
7. Issue docker network connect test-network fakeContainer
8. Issue docker network connect fakeNetwork <containerID>

9. Issue docker network create cross1-network
10. Issue docker network create cross1-network2
11. Issue docker create --net cross1-network --name cross1-container busybox /bin/top
12. Issue docker network connect cross1-network2 <containerID>
13. Issue docker start <containerID>
14. Issue docker create --net cross1-network --name cross1-container2 debian ping -c2 cross1-container
15. Issue docker network connect cross1-network2 <containerID>
16. Issue docker start <containerID>
17. Issue docker logs --follow cross1-container2

18. Issue docker network create cross2-network
19. Issue docker network create cross1-network2
20. Issue docker run -itd --net cross2-network --name cross2-container busybox /bin/top
21. Get the above container's IP - ${ip}
22. Issue docker run --net cross2-network2 --name cross2-container2 debian ping -c2 ${ip}
23. Issue docker logs --follow cross2-container2
24. Issue docker run -d --net cross2-network -p 8080:80 nginx
25. Get the above container's IP - ${ip}
26. Issue docker run --net cross2-network2 --name cross2-container3 debian ping -c2 ${ip}
27. Issue docker logs --follow cross2-container3

28. Issue docker network create --internal internal-net
29. Issue docker run --net internal-net busybox ping -c1 www.google.com
30. Issue docker network create public-net
31. Issue docker run --net internal-net --net public-net busybox ping -c2 www.google.com
32. Issue docker run -itd --net internal-net busybox
33. Get the above container's IP - ${ip}
34. Issue docker run --net internal-net busybox ping -c2 ${ip}

#Expected Outcome:
* Step 4 should complete successfully
* Step 6 should print the results of the ifconfig command and there should be two network interfaces in the container(eth0, eth1)
* Step 7 should result in an error with the following message:  
```
Error response from daemon: No such container: fakeContainer
```
* Step 8 should result in an error with the following message:  
```
Error response from daemon: network fakeNetwork not found
```
* Steps 9-16 should return without errors
* Step 17's output should contain "2 packets transmitted, 2 packets received"
* Steps 18-22 should return without errors
* Step 23's output should contain "2 packets transmitted, 0 packets received, 100% packet loss"
* Steps 24-26 should return without errors
* Step 27's output should include "2 packets transmitted, 0 packets received, 100% packet loss"

* Step 28 should return without an error
* Step 29 should return with a non-zero exit code and the output should contain "Network is unreachable"
* Step 30 should return without an error
* Step 31's output should contain "2 packets transmitted, 2 packets received"
* Steps 32-33 should return without errors
* Step 34's output should contain "2 packets transmitted, 2 packets received"

#Possible Problems:
None