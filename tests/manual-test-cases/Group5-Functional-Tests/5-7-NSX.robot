*** Settings ***
Documentation  Test 5-7 - NSX
Resource  ../../resources/Util.robot

*** Test Cases ***
Test
    Log To Console  TODO
    #${out}=  Deploy Nimbus Testbed  --noSupportBundles --vcvaBuild 3634791 --esxBuild 3620759 --testbedName test-vpx-4esx-virtual-fullInstall-vcva-8gbmem-nsx1m1c --runName VIC-NSX-Test --build nsx-transformers:beta:ob-3586094:master