*** Settings ***
Resource  ../../resources/Util.robot

*** Test Cases ***
Test
    ${result}=  Run Process  ${bin-dir}/imagec -standalone -reference photon -logfile foo.log  shell=True  cwd=/
    Should Be Equal As Integers  0  ${result.rc}
    File Should Exist  /foo.log
    File Should Not Be Empty  /foo.log
    Verify Checksums  /images/https/registry-1.docker.io/v2/library/photon/latest
