*** Settings ***
Resource  ../../resources/Util.robot

*** Test Cases ***
Upgrade VCH with UpdateInProgress
    Start Process  bin/vic-machine-linux create --debug 1 --name vch-2 --target 10.192.190.169 --user Administrator@vsphere.local --password Admin!23 --force true --compute-resource --compute-resource /ha-datacenter/host/10.160.130.21 --timeout 20m  shell=True  alias=UpgradeVCH
    Wait For Process  UpgradeVCH

