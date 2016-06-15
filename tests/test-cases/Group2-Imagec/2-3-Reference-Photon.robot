*** Settings ***
Resource  ../../resources/Util.robot

*** Test Cases ***
Test
    ${result}=  Run Process  ${bin-dir}/imagec -standalone -reference photon  shell=True  cwd=/
    Should Be Equal As Integers  0  ${result.rc}
    OperatingSystem.Directory Should Exist  /images/https/registry-1.docker.io/v2/library/photon/latest
    OperatingSystem.File Should Exist  /imagec.log
    File Should Not Be Empty  /imagec.log
    Verify Checksums  /images/https/registry-1.docker.io/v2/library/photon/latest
