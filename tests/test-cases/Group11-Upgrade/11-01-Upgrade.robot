*** Settings ***
Documentation  Test 11-01 - Upgrade
Resource  ../../resources/Util.robot
# Suite Setup  Install VIC with version to Test Server  5604
# Suite Teardown  Clean up VIC Appliance And Local Binary
Default Tags

*** Keywords ***
Install VIC with version to Test Server
    [Arguments]  ${version}=5604
    Log To Console  \nDownloading vic ${version} from bintray...
    ${rc}  ${output}=  Run And Return Rc And Output  wget https://bintray.com/vmware/vic-repo/download_file?file_path=vic_${version}.tar.gz -O vic.tar.gz
    ${rc}  ${output}=  Run And Return Rc And Output  tar zxvf vic.tar.gz
	Install VIC Appliance To Test Server  vic-machine=./vic/vic-machine-linux  appliance-iso=./vic/appliance.iso  bootstrap-iso=./vic/bootstrap.iso  certs=${false}

Clean up VIC Appliance And Local Binary
    Cleanup VIC Appliance On Test Server
    Run  rm -rf vic.tar.gz vic

Get Container IP
    [Arguments]  ${id}  ${network}=default
    ${rc}  ${ip}=  Run And Return Rc And Output  docker ${params} network inspect ${network} | jq '.[0].Containers."${id}".IPv4Address'
    Should Be Equal As Integers  ${rc}  0
    [Return]  ${ip}

Launch Container
    [Arguments]  ${name}  ${network}=default
    ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} run --name ${name} --net ${network} -itd busybox
    Should Be Equal As Integers  ${rc}  0
    ${id}=  Get Line  ${output}  -1
    ${ip}=  Get Container IP  ${id}  ${network}
    [Return]  ${id}  ${ip}

*** Test Cases ***
Upgrade VCH with containers
    Log  Disabled until PR \#2554 is merged  WARN
    # ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network create bar
    # Should Be Equal As Integers  ${rc}  0
    # Comment  Launch first container on bridge network
    # ${id1}  ${ip1}=  Launch Container  vch-restart-test1  bridge
    # ${id2}  ${ip2}=  Launch Container  vch-restart-test2  bridge

    # ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -it -p 10000:80 -p 10001:80 --name webserver nginx
    # Should Be Equal As Integers  ${rc}  0
    # Should Not Contain  ${output}  Error
    # ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start webserver
    # Should Be Equal As Integers  ${rc}  0
    # Should Not Contain  ${output}  Error
    # Wait Until Keyword Succeeds  20x  5 seconds  Hit Nginx Endpoint  ${vch-ip}  10000
    # Wait Until Keyword Succeeds  20x  5 seconds  Hit Nginx Endpoint  ${vch-ip}  10001

    # Log To Console  \nUpgrading VCH...
    # ${rc}  ${output}=  Run And Return Rc And Output  bin/vic-machine-linux upgrade --debug 1 --name=${vch-name} --target=%{TEST_URL} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --force=true --compute-resource=%{TEST_RESOURCE} --timeout %{TEST_TIMEOUT}
    # Should Contain  ${output}  Completed successfully
    # Should Not Contain  ${output}  Rolling back upgrade
    # Should Be Equal As Integers  ${rc}  0

    # ${rc}  ${output}=  Run And Return Rc And Output  bin/vic-machine-linux inspect --name=${vch-name} --target=%{TEST_URL} --thumbprint=%{TEST_THUMBPRINT} --user=%{TEST_USERNAME} --password=%{TEST_PASSWORD} --compute-resource=%{TEST_RESOURCE}
    # Should Contain  ${output}  Completed successfully
    # Should Be Equal As Integers  ${rc}  0
    # Log  ${output}
    # Get Docker Params  ${output}  ${true}

    # # wait for docker info to succeed
    # Wait Until Keyword Succeeds  20x  5 seconds  Run Docker Info  ${params}

    # ${status}=  Get State Of Github Issue  2448
    # Run Keyword If  '${status}' == 'closed'  Fail  Test 11-1-VCH-Upgrade.robot needs to be updated now that Issue #2448 has been resolved
    # Log  Issue \#2448 is blocking implementation  WARN
    # #${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network ls
    # #Should Be Equal As Integers  ${rc}  0
    # #Should Contain  ${output}  bar
    # ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} network inspect bridge
    # Should Be Equal As Integers  ${rc}  0
    # ${ip}=  Get Container IP  ${id1}  bridge
    # Should Be Equal  ${ip}  ${ip1}
    # ${ip}=  Get Container IP  ${id2}  bridge
    # Should Be Equal  ${ip}  ${ip2}
    # ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} inspect vch-restart-test1
    # Should Be Equal As Integers  ${rc}  0
    # Should Contain  ${output}  "Id"
    # ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} stop vch-restart-test1
    # Should Be Equal As Integers  ${rc}  0
    # ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -a
    # Should Be Equal As Integers  ${rc}  0
    # Should Contain  ${output}  Exited (0)
    # ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start vch-restart-test1
    # Should Be Equal As Integers  ${rc}  0
    # ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} ps -a
    # Should Be Equal As Integers  ${rc}  0
    # Should Not Contain  ${output}  Exited (0)

    # Wait Until Keyword Succeeds  20x  5 seconds  Hit Nginx Endpoint  ${vch-ip}  10000
    # Wait Until Keyword Succeeds  20x  5 seconds  Hit Nginx Endpoint  ${vch-ip}  10001

    # ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} create -it -p 10000:80 -p 10001:80 --name webserver1 nginx
    # Should Be Equal As Integers  ${rc}  0
    # Should Not Contain  ${output}  Error
    # ${rc}  ${output}=  Run And Return Rc And Output  docker ${params} start webserver1
    # Should Be Equal As Integers  ${rc}  1
    # Should Contain  ${output}  port 10000 is not available
