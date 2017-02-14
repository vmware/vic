# Using vSphere Integrated Containers Engine with vSphere Integrated Containers Registry (Harbor)

This example illustrates using a deployed Virtual Container Host (VCH) with Harbor as a private registry with the assumption that a VCH has been set up using either static IP or FQDN.  It also assumes there is access to standard Docker that has been updated with the CA certificate used to sign the Harbor instance's server certificate and server private key.

## Workflow

1. Develop or obtain a docker container image on a computer or terminal using standard docker. Tag the image for Harbor and push the image to the server.
2. Pull down the image from Harbor to a deployed VCH and use it.

## Push a Container Image to Harbor Using Standard Docker

1. Pull the busybox container image from the docker hub to your machine, which you  have updated with the CA certificate earlier. See [Deploy a VCH with vSphere Integrated Containers Registry (Harbor)](../vic_installation/deploy_vch_harbor.html) for more information on updating certificates.
2. Tag the image for uploading to your Harbor registry and push the image up to it. 

**Important** You must log onto the Harbor server before pushing the image up to it.

    user@Devbox:~/mycerts$ docker pull busybox
    Using default tag: latest
    latest: Pulling from library/busybox

    56bec22e3559: Pull complete 
    Digest: sha256:digest
    Status: Downloaded newer image for busybox:latest
    user@Devbox:~/mycerts$ 
    user@Devbox:~/mycerts$ docker tag busybox <Harbor FQDN or static
    IP>/test/busybox

    user@Devbox:~/mycerts$ docker login <Harbor FQDN or static IP>
    Username: user
    Password: 
    Login Succeeded

    user@Devbox:~/mycerts$ docker push <Harbor FQDN or static IP>/test/busybox
    The push refers to a repository [<Harbor FQDN or static IP>/test/busybox]
    e88b3f82283b: Pushed 
    latest: digest: sha256:digest size: 527

## Pull the Image from Harbor to the VCH
In another terminal, pull the image from Harbor to the VCH.

    user@Devbox:~$ export DOCKER_HOST=tcp://<Deployed VCH IP>:2375
    user@Devbox:~$ export DOCKER_API_VERSION=1.23
    user@Devbox:~$ docker images
    REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE

    user@Devbox:~$ docker pull <Harbor FQDN or static IP>/test/busybox
    Using default tag: latest
    Pulling from test/busybox
    Error: image test/busybox not found

    user@Devbox:~$ docker login <Harbor FQDN or static IP>
    Username: user
    Password: 
    Login Succeeded

    user@Devbox:~$ docker pull <Harbor FQDN or static IP>/test/busybox
    Using default tag: latest
    Pulling from test/busybox
    56bec22e3559: Pull complete 
    a3ed95caeb02: Pull complete 
    Digest: sha256:digest
    Status: Downloaded newer image for test/busybox:latest

    user@Devbox:~$ docker images
    REPOSITORY                                TAG        IMAGE ID      CREATED            SIZE
    <Harbor FQDN or static IP>/test/busybox   latest     e292aa76ad3b        5 weeks ago         1.093 MB
    user@Devbox:~$ 

Note that the first attempt to pull the image fails with a 'not found' error message. After you log into the Harbor server, the pull attempt succeeds.