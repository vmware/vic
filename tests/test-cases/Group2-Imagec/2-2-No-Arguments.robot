*** Settings ***
Resource  ../../resources/Util.robot

*** Test Cases ***
Test
    ${rc}=  Run And Return Rc  ${bin-dir}/imagec -standalone
    Should Be Equal As Integers  1  ${rc}