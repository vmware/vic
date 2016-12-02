# Use and Limitations of Containers in vSphere Integrated Containers Engine

vSphere Integrated Containers Engine currently includes the following capabilities and limitations:

## Supported Features
This version of vSphere Integrated Containers Engine supports these features:

- Docker Compose (basic)
- Registry pull from docker hub and private registry
- Named Data Volumes
- Anonymous Data Volumes
- Bridged Networks
- External Networks
- Port Mapping
- Network Links/Alias

## Limitations
vSphere Integrated Containers Engine includes these limitations:

- Container VMs only support root user.
- When you do not configure a PATH environment variable, or create a container from an image that does not supply a PATH, vSphere Integrated Containers Engine provides a default PATH.
- You can resolve the symbolic names of a container from within another container, except in the following cases:
	- Aliases
	- IPv6
	- Service discovery
- Containers can acquire DHCP addresses only if they are on a network that has DHCP.

## Future Features

This version of vSphere Integrated Containers Engine does not support these features:

- Pulling images via image digest 
- Pushing a registry
- Sharing concurrent data volume between containers
- Mapping a local host folder to a container volume
- Mapping a local host file to a container
- Docker build
- Docker copy files into a container, both running and stopped
- Docker container inspect does not return all container network for a container

For limitations of using vSphere Integrated Containers with volumes, see [Using Volumes with vSphere Integrated Containers Engine](using_volumes_with_vic.md).

## Workflow Guidelines

`Docker compose` is a good baseline reference for this guideline as you can perform the same tasks with `Docker compose` manually using the Docker CLI and with scripting using the CLI. This guideline uses Docker Compose 1.8.1. The Future Features list puts constraints on types of containerized applications that can be deployed in this release.

**Note** These guidelines and recommendations exist for the current feature set in this release of vSphere Integrated Containers Engine.

###Guidelines for Building Container Images

Without the docker build and registry pushing features, you need to use regular Docker to build a container and to push it to the global hub or your corporate private registry. The example workflow using Docker's multi-tiered [voting app](https://github.com/docker/example-voting-app) illustrates the workaround the constraints.

Use the [guidelines](README.md) to modify the Docker Compose yml file to work with vSphere Integrated Containers Engine.

#### Original Compose File

    version: "2"
    
    services:
      vote:
        build: ./vote
        command: python app.py
        volumes:
         - ./vote:/app
        ports:
         - "5000:80"
    
      redis:
        image: redis:alpine
        ports: ["6379"]
    
      worker:
        build: ./worker
    
      db:
        image: postgres:9.4
    
      result:
        build: ./result
        command: nodemon --debug server.js
        volumes:
        - ./result:/app
        ports:
        - "5001:80"
        - "5858:5858"


The compose file uses two features that are not yet supported in this release: docker build and local folder mapping to container volume.  

#### Modify the App

To modify this app and deploy it onto a vSphere environment, perform these steps.

1. Clone the repository from github
2. Use regular docker to build each component that requires a build.
3. Tag the the images to upload to your private registry or private account on Docker Hub. In this example, the account is victest on Docker Hub. Create your own account and use that in place of the victest keywords.  

	**Note** These steps are performed in a terminal using regular docker as opposed to the vSphere Integrated Containers docker personality daemon. You can build and tag an image in one step.

    	**build the images:**  
    	$> cd example-voting-app  
    	$> docker build -t vote ./vote  
    	$> docker build -t vote-worker ./worker  
    	$> docker build -t vote-result ./result  
    
    	**tag the images for a registry:**  
    	$> docker tag vote victest/vote  
    	$> docker tag vote-worker victest/vote-worker  
    	$> docker tag vote-result victest/vote-result  
    
    	**push the images to the registry:**  
    	$> docker login (... and provide credentials)  
    	$> docker push victest/vote  
    	$> docker push victest/vote-worker  
    	$> docker push victest/vote-result  

4. Analyze the application. The local folder mapping and all the build directives from the yml file are removed. 

#### Updated Compose File

    version: "2"
    
    services:
      vote:
    	image: victest/vote
    	command: python app.py
    	ports:
    	  - "5000:80"
    
      redis:
    	image: redis:alpine
       		ports: ["6379"]
    
      worker:
       		image: victest/vote-worker
    
      db:
       		image: postgres:9.4
    
      result:
       		image: victest/vote-result
    	command: nodemon --debug server.js
    	ports:
    	  - "5001:80"
    	  - "5858:5858"

1. Assuming you have deployed a VCH with vic-machine and the VCH_IP is the IP address of the deployed VCH, you have the IP after the VCH was successfully installed. Using the example-voting-app folder, with the modified compose yml file:

	$> docker-compose -H VCH_IP up -d

2. Go to "http://VCH_IP:5000" and "http://VCH_IP:5001" to verify the voting app is running.

###Guidelines for Sharing Configuration

Without the data volume sharing and docker copy features, providing configuration to a containerized application has some constraints. 

An example of configuration is your web server config files. The guideline is to pass in configuration using command line arguments or environment variables. Add a script to the container image that ingests the command line argument/environment variable and passes the configuration to the contained application. A benefit of using environment variables to transfer configuration is the containerized app then closely follows the popular 12-factor app model.

With no direct support for sharing volumes between containers, processes that must share files have these options:

- Build them into the same image and run in the same container.
- Add a script to the container that mounts an NFS share where containers must be on the same network:
	- Run container with NFS server sharing a data volume.
	- Mount NFS share on the containers that need to share files.
