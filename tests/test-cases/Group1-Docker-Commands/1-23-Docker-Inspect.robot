*** Settings ***
Documentation  Test 1-23 - Docker Inspect
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${false}
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Simple docker inspect of image
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} pull busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect busybox
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Evaluate  json.loads(r'''${output}''')  json
    ${id}=  Get From Dictionary  ${output[0]}  Id

Docker inspect image specifying type
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect --type=image busybox
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Evaluate  json.loads(r'''${output}''')  json
    ${id}=  Get From Dictionary  ${output[0]}  Id

Docker inspect image specifying incorrect type    
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect --type=container busybox
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error: No such container: busybox
    
Simple docker inspect of container
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create busybox
    Should Be Equal As Integers  ${rc}  0    
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect ${container}
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Evaluate  json.loads(r'''${output}''')  json
    ${id}=  Get From Dictionary  ${output[0]}  Id

Docker inspect container specifying type
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect --type=container ${container}
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Evaluate  json.loads(r'''${output}''')  json
    ${id}=  Get From Dictionary  ${output[0]}  Id

Docker inspect container check cmd and image name
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create busybox /bin/bash
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect ${container}
    Should Be Equal As Integers  ${rc}  0
    ${output}=  Evaluate  json.loads(r'''${output}''')  json
    ${config}=  Get From Dictionary  ${output[0]}  Config
    ${image}=  Get From Dictionary  ${config}  Image
    Should Contain  ${image}  busybox
    ${cmd}=  Get From Dictionary  ${config}  Cmd
    Should Be Equal As Strings  ${cmd}  [u'/bin/bash']

Docker inspect container specifying incorrect type
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect --type=image ${container}
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error: No such image: ${container}

Docker inspect container with multiple networks
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network create net-one
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network create net-two
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} create --name=two-net-test --net=net-one busybox
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} network connect net-two two-net-test
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${container}=  Run And Return Rc And Output  docker %{VCH-PARAMS} start two-net-test
    Should Be Equal As Integers  ${rc}  0
    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect -f '{{range $key, $value := .NetworkSettings.Networks}}{{$key}}{{end}}' two-net-test
    Should Contain  ${out}  net-two
    Should Contain  ${out}  net-one
    Should Be Equal As Integers  ${rc}  0
    
Docker inspect invalid object
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect fake
    Should Be Equal As Integers  ${rc}  1
    Should Contain  ${output}  Error: No such image or container: fake