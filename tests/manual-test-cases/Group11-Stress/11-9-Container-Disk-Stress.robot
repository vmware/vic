*** Settings ***
Documentation  Test 11-9-Container-Disk-Stress
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Container Disk Stress
    ${out}=  Run  docker %{VCH-PARAMS} pull busybox

    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} run ubuntu bash -c "apt-get update; apt-get install bonnie++; bonnie++ -u root;"
    Should Be Equal As Integers  ${rc}  0
    Should Contain  ${output}  Delete files in random order...done.
    
    Run Regression Tests