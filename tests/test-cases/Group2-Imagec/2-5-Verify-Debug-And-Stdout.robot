*** Settings ***
Resource  ../../resources/Util.robot

*** Test Cases ***
Test
    ${result}=  Run Process  ${bin-dir}/imagec -standalone -reference photon -stdout -debug  shell=True  cwd=/
    Log  ${result.stdout}
    Log  ${result.stderr}
    Should Contain  ${result.stdout}  level=debug
    Should Be Equal As Integers  0  ${result.rc}
    Verify Checksums  /images/https/registry-1.docker.io/v2/library/photon/latest