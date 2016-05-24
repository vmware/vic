# Application topologies

This describes the core application topologies that will be used to validate function and determine relative priority of feature implementation. This is by no means a complete list and the ordering may change over time as we gain additional feedback from users. 

## Prime topology - platform 2.5 tiered application

### Description
This is the topology currently being used as the primary touchstone for feature support and priority. It's based off [the docker voting app](https://github.com/docker/example-voting-app) as a great example of a mapping from a traditional tiered application to a containerized environment. This is not assumed to be a pure [twelve factor app](http://12factor.net/) because while that significantly simplifies the infrastructure requirements it's not a viable assumption for all workloads.

There are two scenarios described in the breakdown:
1. unmodified voting app
2. voting app using non-containerized database with direct exposure to an external network on the front-end

The second scenario is assumed to be a better reflection of probable deployment scenarios as people are moving to containerized workloads, hence the description of this as a 2.5 platform application. The breakdowns extract the docker CLI operations that correspond to the docker compose file. In the modified case it also details additional configuration elements that permit for the integration with non-containerized workloads and existing networks. 

### Breakdown - unmodified
Usage scenario: docker voting app - direct from compose file
```
	docker network create front-tier
	docker network create back-tier
	
	docker volume create voting-app
	docker volume create result-app
	docker volume create db-data
	
	docker create --volume=voting-app:/app --publish=5000:80 --link=redis --net=front-tier --name voting-app voting-app
	docker network connect back-tier voting-app  
	docker start voting-app
	docker create --volume=result-app:/app --publish=5001:80 --link=db --net=front-tier --name result-app result-app 
	docker network connect back-tier result-app  
	docker start result-app
	docker run --publish=6379 --net=back-tier --name redis redis:alpine
	docker run --link=db --link=redis --net=back-tier --name worker worker
	docker run --volume=db-data:/var/lib/postgresql/data --net=back-tier --name db postgres:9.4
```

### Breakdown - modified
Usage scenario: docker voting app - using non-containerized database and external net. Results only visible internally on front-end net
```
#	mappings supplied to vic-machine when deploying VCH
	-docker-network=vsphere-external:external
	-docker-network=corp.net:corp
	-container=votesdb.corp.net:db.corp

# docker command breakdown
	docker network create front-tier
	docker network create back-tier
	
	docker volume create voting-app
	docker volume create result-app
	docker volume create db-data
	
	docker create --volume=voting-app:/app --publish=80 --link=redis --net=external --name voting-app voting-app 
	docker network connect back-tier voting-app 
	docker start voting-app
	docker create --volume=result-app:/app --publish=80 --link=db --net=front-tier --name result-app result-app 
	docker network connect back-tier result-app 
	docker start result-app
	docker run --publish=6379 --net=back-tier --name redis redis:alpine
	docker create --link=redis -link=db.corp:db --net=corp --name worker worker
	docker network connect back-tier worker
	docker start worker
	
# non-containerized workloads
  VM running database: votesdb.corp.net
```
