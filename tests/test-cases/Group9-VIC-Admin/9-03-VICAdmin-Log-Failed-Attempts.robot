*** Settings ***
Documentation  Test 9-03 - VICAdmin Log Failed Attempts
Resource  ../../resources/Util.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server

*** Test Cases ***
Verify Unable To Verify
    ${out}=  Run  wget --tries=3 --connect-timeout=10 %{VIC-ADMIN}/logs/vicadmin.log -O failure.log
    Should Contain  ${out}  ERROR: cannot verify
    Should Contain  ${out}  certificate, issued by
    Should Contain  ${out}  Unable to locally verify the issuer's authority.
    
Verify Temporary Redirect
    ${out}=  Run  wget --tries=3 --connect-timeout=10 --no-check-certificate %{VIC-ADMIN}/logs/vicadmin.log -O failure.log
    Should Contain  ${out}  HTTP request sent, awaiting response... 307 Temporary Redirect

Verify Failed Log Attempts
    #Save the first appliance certs and cleanup the first appliance
    #${old-certs}=  Set Variable  %{DOCKER_CERT_PATH}
    Run  cp -r %{DOCKER_CERT_PATH} old-certs
    Cleanup VIC Appliance On Test Server
    
    #Install a second appliance
    Install VIC Appliance To Test Server
    OperatingSystem.File Should Exist  old-certs/cert.pem
    OperatingSystem.File Should Exist  old-certs/key.pem
    ${out}=  Run  wget -v --tries=3 --connect-timeout=10 --certificate=old-certs/cert.pem --private-key=old-certs/key.pem --no-check-certificate %{VIC-ADMIN}/logs/vicadmin.log -O failure.log
    Log  ${out}
    ${out}=  Run  wget -v --tries=3 --connect-timeout=10 --certificate=%{DOCKER_CERT_PATH}/cert.pem --private-key=%{DOCKER_CERT_PATH}/key.pem --no-check-certificate %{VIC-ADMIN}/logs/vicadmin.log -O success.log
    Log  ${out}
    ${out}=  Run  cat success.log
    Log  ${out}
    ${out}=  Run  grep -i fail success.log
    Should Contain  ${out}  tls: failed to verify client's certificate: x509: certificate signed by unknown authority
    Run  rm -r old-certs
