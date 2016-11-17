*** Settings ***
Documentation  Test 6-12 - Verify vic-machine version command
Resource  ../../resources/Util.robot

*** Test Cases ***
VIC-machine - Version check
    Set Test Environment Variables

    ${output}=  Run  bin/vic-machine-linux version
    @{gotVersion}=  Split String  ${output}  ${SPACE}
    ${version}=  Remove String  @{gotVersion}[2]
    Log To Console  VIC machine version: ${version}
    
    ${result}=  Run  git rev-parse HEAD
    @{gotVersion}=  Split String  ${result}  ${SPACE}
    ${commithash}=  Remove String  @{gotVersion}[0]
    
    Log To Console  Last commit hash from git: ${commithash}

    ${hash_result} =    Fetch From Right  ${version}  -
    Log To Console  Commit Hash from vic-machine version: ${hash_result}

    Should Contain  ${commithash}  ${hash_result}
