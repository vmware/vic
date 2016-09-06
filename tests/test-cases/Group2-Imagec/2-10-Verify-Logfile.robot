*** Settings ***
Resource  ../../resources/Util.robot

*** Test Cases ***
Test
    ${result}=  Run Process  ${bin-dir}/imagec -debug -standalone -reference photon -logfile /tmp/foo.log -destination /tmp/images  shell=True 
    Should Be Equal As Integers  0  ${result.rc}
    OperatingSystem.File Should Exist  /tmp/foo.log
    File Should Not Be Empty  /tmp/foo.log
    Verify Checksums  /tmp/images/https/registry-1.docker.io/v2/library/photon/latest
