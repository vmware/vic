*** Settings ***
Documentation  Test 9-03 - VICAdmin Log Failed Attempts
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Test
    ${out}=  Run  wget --tries=3 --connect-timeout=10 --certificate=%{DOCKER_CERT_PATH}/fakeCert.pem --private-key=%{DOCKER_CERT_PATH}/fakeKey.pem --no-check-certificate ${vic-admin}/logs/vicadmin.log -O failure.log
    Log  ${out}

    ${out}=  Run  wget --tries=3 --connect-timeout=10 --certificate=%{DOCKER_CERT_PATH}/cert.pem --private-key=%{DOCKER_CERT_PATH}/key.pem --no-check-certificate ${vic-admin}/logs/vicadmin.log -O success.log
    Log  ${out}

    ${out}=  Run  grep -i fail success.log
    Log  ${out}
    Should Contain  ${out}  tls: failed to verify client's certificate: x509: certificate signed by unknown authority