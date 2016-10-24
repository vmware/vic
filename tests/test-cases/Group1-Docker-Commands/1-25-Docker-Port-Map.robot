*** Settings ***
Documentation  Test 1-25 - Docker Port Map
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${false}
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Create container with port mappings
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -it -p 10000:80 -p 10001:80 --name webserver nginx
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start webserver
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    Wait Until Keyword Succeeds  20x  5 seconds  Hit Nginx Endpoint  ${ext-ip}  10000
    Wait Until Keyword Succeeds  20x  5 seconds  Hit Nginx Endpoint  ${ext-ip}  10001

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} stop webserver
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  curl ${ext-ip}:10000 --connect-timeout 5
    Should Not Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  curl ${ext-ip}:10001 --connect-timeout 5
    Should Not Be Equal As Integers  ${rc}  0

Create container with conflicting port mapping
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -it -p 8083:80 --name webserver2 nginx
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -it -p 8083:80 --name webserver3 nginx
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start webserver2
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start webserver3
    Should Not Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  port 8083 is not available

Create container with port range
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -it -p 8081-8088:80 --name webserver5 nginx
    Should Not Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  host port ranges are not supported for port bindings

Create container with host ip
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -it -p 10.10.10.10:8088:80 --name webserver5 nginx
    Should Not Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  host IP for port bindings is only supported for 0.0.0.0 and the external interface IP address

Create container with host ip equal to 0.0.0.0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -it -p 0.0.0.0:8088:80 --name webserver5 nginx
    Should Be Equal As Integers  ${rc}  0

Create container with host ip equal to external IP
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -it -p ${ext-ip}:8089:80 --name webserver6 nginx
    Should Be Equal As Integers  ${rc}  0

Create container without specifying host port
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -it -p 6379 --name test-redis redis:alpine
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start test-redis
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} stop test-redis
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

Run and exit with mapped ports
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rm -f $(docker ${params} ps -aq)
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  mkfifo /tmp/fifo1
    ${result}=  Start Process  docker ${params} run -i --name ctr1 -p 1900:9999 -p 2200:2222 busybox /bin/top < /tmp/fifo1  shell=True  alias=sh1
    Sleep  5
    ${rc}  ${output}=  Run And Return Rc And Output  echo q > /tmp/fifo1
    ${result}=  Wait for process  sh1
    Log  ${result.stdout}
    Log  ${result.stderr}

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -a
    Log  ${output}
    Should Contain X Times  ${output}  Exited  1

    ${rc}  ${output}=  Run And Return Rc And Output  mkfifo /tmp/fifo2
    ${result}=  Start Process  docker ${params} run -i --name ctr2 -p 1900:9999 -p 3300:3333 busybox /bin/top < /tmp/fifo2  shell=True  alias=sh2
    Sleep  5
    ${rc}  ${output}=  Run And Return Rc And Output  echo q > /tmp/fifo2
    ${result}=  Wait for process  sh2
    Log  ${result.stdout}
    Log  ${result.stderr}
    Should Be Equal As Integers  ${result.rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -a
    Log  ${output}
    Should Contain X Times  ${output}  Exited  2

Remap mapped ports after OOB Stop
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} rm -f $(docker ${params} ps -aq)
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -it -p 10000:80 -p 10001:80 --name ctr3 busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ctr3
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error

    # PowerOff VM out-of-band
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'VC'  Run And Return Rc And Output  govc vm.power -off ${vch-name}/"ctr3*"
    Run Keyword If  '%{HOST_TYPE}' == 'VC'  Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Run And Return Rc And Output  govc vm.power -off "ctr3*"
    Run Keyword If  '%{HOST_TYPE}' == 'ESXi'  Should Be Equal As Integers  ${rc}  0
    Wait Until VM Powers Off  "ctr3*"

    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -it -p 10000:80 -p 20000:22222 --name ctr4 busybox
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start ctr4
    Should Be Equal As Integers  ${rc}  0
    Should Not Contain  ${output}  Error
