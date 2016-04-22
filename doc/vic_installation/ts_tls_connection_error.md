# Connection to Docker Client Fails with a TLS Connection Error #
After a successful installation of vSphere Integrated Containers, connecting a Docker client to the virtual container host fails with a TLS connection error.

## Problem ##
After you have set the `DOCKER_HOST` variable to point to your virtual container host, when you attempt a Docker operation in your Docker client, the connection fails with the error `An error occurred trying to connect: Get https://<vic_appliance_address>:2376/v1.22/info: tls: oversized record received with length 20527`.

## Cause ##
The Docker client is attempting to authenticate the connection with the virtual container host. TLS authentication is not activated on the virtual container host, or the Docker client does not have the appropriate certificate and key.

## Solution ##

- Activate TLS authentication between the virtual container host and Docker client by following the instructions in [Using TLS Authentication with vSphere Integrated Containers](using_tls_with_vic.md).
- Specify the `-tls` and `-tlsverify` options when you connect to the virtual container host. 

Alternatively, disable TLS authentication in the Docker client. 

1. Open a Docker client terminal.
2. Prevent Docker from looking for certificates by disabling the docker certificate path variable.

 `unset DOCKER_CERT_PATH` 
3. Disable TLS authentication by disabling the docker TLS variable. 

 `unset DOCKER_TLS_VERIFY` 
4. Check that your Docker client can now connect to the virtual container host by running a Docker command. 

 `docker info` 

 You should see information about the virtual container host that is running in vSphere.