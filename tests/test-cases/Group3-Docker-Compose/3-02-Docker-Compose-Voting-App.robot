*** Settings ***
Documentation  Test 3-02 - Docker Compose Voting App
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server  certs=${false}
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Compose Voting App
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} login --username=victest --password=vmware!123
    Should Be Equal As Integers  ${rc}  0

    Set Environment Variable  COMPOSE_HTTP_TIMEOUT  300
    # must set CURL_CA_BUNDLE to work around Compose bug https://github.com/docker/compose/issues/3365
    Set Environment Variable  CURL_CA_BUNDLE  ${EMPTY}

    ${rc}  ${out}=  Run And Return Rc And Output  docker-compose --skip-hostname-check -f %{GOPATH}/src/github.com/vmware/vic/demos/compose/voting-app/docker-compose.yml %{VCH-PARAMS} up -d
    Log  ${out}
    #Log  ${out.stdout}
    #Log  ${out.stderr}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect -f {{.State.Running}} vote
    Log  ${out}
    Should Contain  ${out}  true
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect -f {{.State.Running}} result
    Log  ${out}
    Should Contain  ${out}  true
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect -f {{.State.Running}} worker
    Log  ${out}
    Should Contain  ${out}  true
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect -f {{.State.Running}} db
    Log  ${out}
    Should Contain  ${out}  true
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect -f {{.State.Running}} redis
    Log  ${out}
    Should Contain  ${out}  true
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect -f '{{range $key, $value := .NetworkSettings.Networks}}{{$key}}{{end}}' vote
    Log  ${out}
    Should Not Be Empty  ${out}
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect -f '{{range $key, $value := .NetworkSettings.Networks}}{{index $value "Aliases"}}{{end}}' vote
    Log  ${out}
    Should Contain  ${out}  vote
    Should Be Equal As Integers  ${rc}  0

    ${rc}  ${out}=  Run And Return Rc And Output  docker %{VCH-PARAMS} inspect -f '{{range $key, $value := .NetworkSettings.Networks}}{{index $value "IPAddress"}}{{end}}' vote
    Log  ${out}
    Should Not Be Empty  ${out}
    Should Be Equal As Integers  ${rc}  0
