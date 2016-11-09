*** Settings ***
Documentation  Test 9-03 - VICAdmin Log Failed Attempts
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Verify Temporary Redirect
    ${out}=  Run  wget --tries=3 --connect-timeout=10 --certificate=%{DOCKER_CERT_PATH}/fakeCert.pem --private-key=%{DOCKER_CERT_PATH}/fakeKey.pem --no-check-certificate ${vic-admin}/logs/vicadmin.log -O failure.log
    Should Contain  ${out}  HTTP request sent, awaiting response... 307 Temporary Redirect

Verify Failed Log Attempts
    #Save the first appliance certs and cleanup the first appliance
    ${old-certs}=  Set Variable  %{DOCKER_CERT_PATH}
    Cleanup VIC Appliance On Test Server
    
    #Install a second appliance
    Install VIC Appliance To Test Server
    
    ${out}=  Run  wget --tries=3 --connect-timeout=10 --certificate=${old-certs}/cert.pem --private-key=${old-certs}/key.pem --no-check-certificate ${vic-admin}/logs/vicadmin.log -O failure.log
    Log  ${out}

    ${out}=  Run  wget --tries=3 --connect-timeout=10 --certificate=%{DOCKER_CERT_PATH}/cert.pem --private-key=%{DOCKER_CERT_PATH}/key.pem --no-check-certificate ${vic-admin}/logs/vicadmin.log -O success.log
    Log  ${out}

    ${out}=  Run  grep -i fail success.log
    Log  ${out}
    Should Contain  ${out}  tls: failed to verify client's certificate: x509: certificate signed by unknown authority