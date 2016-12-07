# Setting Up a Private Registry Using vSphere Integrated Containers Registry (Harbor)

The example in this topic illustrates how to run a private Docker registry on vSphere Integrated Containers Engine and push and pull images.  An enterprise ready registry named [vSphere Integrated Containers Registry (Harbor)](https://github.com/vmware/harbor) can be used in place of the base Docker registry.

## Prerequisites

The following examples use Linux and do not differ much for Windows and Mac users.

vSphere Integrated Containers Engine requires you to install Docker locally, so you get both a client (CLI or command line interface) and a daemon that handles all local container operations. 

Local containers run on your local machine instead of on a VMWare vSphere/ESXi environment.  You need the CLI for working with containers on vSphere Integrated Containers Engine and on your local system. 

The distinction between using the CLI against the two environments is very important in this example. By default, the CLI uses the local Docker daemon. After setting two environment variables DOCKER_HOST and DOCKER_API_VERSION, you can instruct the CLI to send all operations to vSphere Integrated Containers Engine instead of the local Docker daemon.

This example shows how to deploy an insecure registry with no authentication for simplicity.  We will also be targeting an ESXi environment.

## Workflow

In terminal #1: (local Docker)

1. Open a terminal and make sure that it uses the local Docker daemon. At the command prompt, type:  
`$> unset DOCKER_HOST`
2. Install a Virtual Container Host (VCH) for a private registry using vic-machine.
3. Run [Docker's registry](https://docs.docker.com/registry/) on the first VCH.
4. Install a second VCH for running applications and specify `--insecure-registry` to ensure it can pull images from the insecure registry in the first VCH.
5. At the terminal command prompt, using regular Docker, tag the images to be destined for the registry.
6. Modify the docker systemd config file to allow pushing to an insecure registry.
7. Restart the docker daemon.
8. Push the image using the full tagged name including host IP and port.

In terminal #2: (VIC VCH)

1. Open a terminal and make sure it is using the second VCH. At the command prompt, type:  
`$> export DOCKER_HOST=tcp://<VCH_IP>:<VCH_PORT>`
<br>
`$> export DOCKER_API_VERSION=1.23`
2. Pull the image from the registry VCH.

### Example run

Terminal 1:
<pre>
$> unset DOCKER_HOST
$> ./vic-machine-linux create --target=192.168.218.207 --image-store=datastore1 --name=vic-registry --user=root --password=vagrant --compute-resource="/ha-datacenter/host/esxbox.localdomain/Resources" --bridge-network=vic-network --no-tls --volume-store=datastore1/registry:default --force

...
INFO[2016-10-08T17:31:06-07:00] Initialization of appliance successful
INFO[2016-10-08T17:31:06-07:00]
INFO[2016-10-08T17:31:06-07:00] vic-admin portal:      
INFO[2016-10-08T17:31:06-07:00] http://192.168.218.138:2378
INFO[2016-10-08T17:31:06-07:00]
INFO[2016-10-08T17:31:06-07:00] Docker environment variables:
INFO[2016-10-08T17:31:06-07:00]   DOCKER_HOST=192.168.218.138:2375
INFO[2016-10-08T17:31:06-07:00]
INFO[2016-10-08T17:31:06-07:00]
INFO[2016-10-08T17:31:06-07:00] Connect to docker:
INFO[2016-10-08T17:31:06-07:00] docker -H 192.168.218.138:2375 info
INFO[2016-10-08T17:31:06-07:00] Installer completed successfully

$> DOCKER_HOST=tcp://192.168.218.138:2375 DOCKER_API_VERSION=1.23 docker run -d -p 5000:5000 --name registry registry:2

$> ./vic-machine-linux create --target=192.168.218.207 --image-store=datastore1 --name=vic-app --user=root --password=vagrant --compute-resource="/ha-datacenter/host/esxbox.localdomain/Resources" --bridge-network=vic-network --no-tls --volume-store=datastore1/vic-app:default --force --insecure-registry 192.168.218.138

...
INFO[2016-10-08T17:31:06-07:00] Initialization of appliance successful       
INFO[2016-10-08T17:31:06-07:00]                                              
INFO[2016-10-08T17:31:06-07:00] vic-admin portal:                            
INFO[2016-10-08T17:31:06-07:00] http://192.168.218.131:2378                  
INFO[2016-10-08T17:31:06-07:00]                                              
INFO[2016-10-08T17:31:06-07:00] Docker environment variables:                
INFO[2016-10-08T17:31:06-07:00]   DOCKER_HOST=192.168.218.131:2375           
INFO[2016-10-08T17:31:06-07:00]                                              
INFO[2016-10-08T17:31:06-07:00]                                              
INFO[2016-10-08T17:31:06-07:00] Connect to docker:                           
INFO[2016-10-08T17:31:06-07:00] docker -H 192.168.218.131:2375 info          
INFO[2016-10-08T17:31:06-07:00] Installer completed successfully             

$> sudo vi /lib/systemd/system/docker.service
$> sudo systemctl daemon-reload
$> sudo systemctl restart docker
$> docker tag busybox 192.168.218.138:5000/test/busybox
$> docker push 192.168.218.138:5000/test/busybox
</pre>

Terminal 2:

`$> export DOCKER_HOST=tcp://192.168.218.131:2375`<br>
`$> export DOCKER_API_VERSION=1.23`<br>
`$> docker pull 192.168.218.138:5000/test/busybox`

In this example, TLS is disabled for simplicity. The example does not show the modifications in /lib/systemd/system/docker.service. Example of TLS disabled:

<pre>
[Service]
Type=notify
# the default is not to use systemd for cgroups because the delegate issues still
# exists and systemd currently does not support the cgroup feature set required
# for containers run by docker
ExecStart=/usr/bin/dockerd --tls=false -H fd:// --insecure-registry 192.168.218.138:5000
</pre>

In the second step the example specifies the necessary environment variables before the docker run command.  On the Linux machine, it sets the variables for the duration of the operation, and after the docker run completes, those variables are reverted to their previous values. The example uses the registry:2 image. It is important not to specify registry:2.0 as Registry 2.0 has issues that prevents the example from working.