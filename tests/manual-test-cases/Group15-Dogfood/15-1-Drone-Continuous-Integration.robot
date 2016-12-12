*** Settings ***
Documentation  Test 15-1 - Drone Continuous Integration
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Drone CI
    ${output}=  Run  git clone https://github.com/vmware/vic.git drone-ci
    Log  ${output}
    ${result}=  Run Process  drone exec --docker-host %{VCH-IP}:2375 --trusted -e .drone.sec -yaml .drone.yml  shell=True  cwd=drone-ci
    Log  ${result.stderr}
    Log  ${result.stdout}
    Log  ${result.rc}