*** Settings ***
Documentation  Test 16-1 - Windows VIC Machine
Resource  ../../resources/Util.robot
Test Setup  Set Test Environment Variables

*** Variables ***
${ver}  0.6.0

*** Keywords ***
Cleanup Folders
    ${output}=  Execute Command  rm -Recurse -Force vic*
    ${output}=  Execute Command  rm *.pem

*** Test Cases ***
Install VCH With TLS
    Open Connection  %{WINDOWS_URL}  prompt=>
    Login  %{WINDOWS_USERNAME}  %{WINDOWS_PASSWORD}
    Cleanup Folders
    ${output}=  Execute Command  wget https://bintray.com/vmware/vic/download_file?file_path=vic_${ver}.tar.gz -OutFile vic.tar.gz
    ${output}=  Execute Command  7z x vic.tar.gz
    ${output}=  Execute Command  7z x vic.tar
    ${output}=  Execute Command  ./vic/vic-machine-windows.exe create --target %{TEST_URL} --user %{TEST_USERNAME} --password %{TEST_PASSWORD}
    Get Docker Params  ${output}  ${true}
    Run Regression Tests
    ${output}=  Execute Command  ./vic/vic-machine-windows.exe delete --target %{TEST_URL} --user %{TEST_USERNAME} --password %{TEST_PASSWORD}    
    Cleanup Folders
    
Install VCH Without TLS
    Open Connection  %{WINDOWS_URL}  prompt=>
    Login  %{WINDOWS_USERNAME}  %{WINDOWS_PASSWORD}
    Cleanup Folders
    ${output}=  Execute Command  wget https://bintray.com/vmware/vic/download_file?file_path=vic_${ver}.tar.gz -OutFile vic.tar.gz
    ${output}=  Execute Command  7z x vic.tar.gz
    ${output}=  Execute Command  7z x vic.tar
    ${output}=  Execute Command  ./vic/vic-machine-windows.exe create --target %{TEST_URL} --user %{TEST_USERNAME} --password %{TEST_PASSWORD} --no-tls
    Get Docker Params  ${output}  ${false}
    Run Regression Tests
    ${output}=  Execute Command  ./vic/vic-machine-windows.exe delete --target %{TEST_URL} --user %{TEST_USERNAME} --password %{TEST_PASSWORD}    
    Cleanup Folders