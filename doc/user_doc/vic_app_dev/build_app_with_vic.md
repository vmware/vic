# Example of Building an Application with vSphere Integrated Containers Engine #

The example in this topic modifies the [voting app](https://github.com/docker/example-voting-app) by Docker to illustrate how to work around the constraints that the current version of vSphere Integrated Containers Engine imposes. For information about the constraints, see [Constraints of Using vSphere Integrated Containers Engine to Build Applications](constraints_using_vic.md). 

This example focuses on how to modify the Docker Compose YML file from the voting app to make it work with vSphere Integrated Containers. It does not describe the general function or makeup of the voting app.  

## Getting Started ##

1. Clone the Docker voting app repository from https://github.com/docker/example-voting-app.
2. Open the YML file for the simple Docker voting app, `docker-compose-simple.yml`.

<pre>
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
</pre>

This compose file uses two features that this version of vSphere Integrated Containers does not support:

- The `docker build` command
- Mapping of local folders to container volumes 

To allow the voting app to work with vSphere Integrated Containers, you must modify it to work around these constraints. 

## Modify the Application ##

This version of vSphere Integrated Containers Engine does not support the `docker build`, `tag`, `push` commands. Use regular Docker without vSphere Integrated Containers Engine to perform the steps in this section.

**NOTE**: It is possible to build and tag an image in one step. In this example, building and tagging are in separate steps.

1. Build images for the different components of the application.

   <pre>cd example-voting-app  
   docker build -t vote ./vote  
   docker build -t vote-worker ./worker  
   docker build -t vote-result ./result</pre>

2. Tag the the images to upload them to your private registry or to your personal account on Docker Hub. 
   
   This example uses a Docker Hub account. Replace <i>dockerhub_username</i> with your own account name in the commands below.

   <pre>docker tag vote <i>dockerhub_username</i>/vote  
   docker tag vote-worker <i>dockerhub_username</i>/vote-worker  
   docker tag vote-result <i>dockerhub_username</i>/vote-result</pre>
    
3. Push the images to the registry.

   <pre>docker login 
   [Provide credentials] 
   docker push <i>dockerhub_username</i>/vote  
   docker push <i>dockerhub_username</i>/vote-worker  
   docker push <i>dockerhub_username</i>/vote-result</pre> 

4. Open the `docker-compose-simple.yml` file in an editor and modify it to remove the operations that vSphere Integrated Containers does not support.

    - Remove local folder mapping
    - Remove all of the build directives
    - Update <i>dockerhub_username</i> to your Docker Hub account name
    - Save the modified file with the name `docker-compose.yml`.  

The example below shows the YML file after the modifications:
 
<pre>
version: "2"

services:
  vote:
    image: <i>dockerhub_username</i>/vote
    command: python app.py
    ports:
      - "5000:80"

  redis:
    image: redis:alpine
    ports: ["6379"]

  worker:
    image: <i>dockerhub_username</i>/vote-worker

  db:
    image: postgres:9.4

  result:
    image: <i>dockerhub_username</i>/vote-result
    command: nodemon --debug server.js
    ports:
      - "5001:80"
      - "5858:5858"
</pre>

You can download the modified YML file from the vSphere Integrated Containers Engine repository on Github from https://github.com/vmware/vic/blob/master/demos/compose/voting-app. 

## Deploy the Application to a VCH ##

The steps in this section make the following assumptions:

- You have deployed a virtual container host (VCH).
- You deployed the VCH with a volume store named `default` by specifying  `--volume-store datastore_name/path:default`.
- You deployed the VCH with the `--no-tls` option, to disable TLS authentication between the Docker client and the VCH.
- You are using Docker Compose 1.8.1.

In the procedure below, run the commands from the `example-voting-app`  folder that contains the modified `docker-compose.yml` file.

1. Run the `docker-compose` command.

	<pre>docker-compose -H <i>vch_address</i>:2375 up -d</pre>

2. In a browser, go to http://*vch_address*:5000 and http://*vch_address*:5001 to verify that the Docker voting application is running.
 
   You can vote on port 5000 and see the results on port 5001.
