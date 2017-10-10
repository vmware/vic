# Copyright 2017 VMware, Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#	http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License

*** Settings ***
Documentation  Test 1-00 - Docker Commands
Resource  ../../resources/Util.robot
Resource  subtests/1-28-Docker-Secret.robot
Resource  subtests/1-29-Docker-Checkpoint.robot
Resource  subtests/1-30-Docker-Deploy.robot
Resource  subtests/1-31-Docker-Node.robot
Resource  subtests/1-32-Docker-Plugin.robot
Resource  subtests/1-33-Docker-Service.robot
Resource  subtests/1-34-Docker-Stack.robot
Resource  subtests/1-35-Docker-Swarm.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server
Test Timeout  30 minutes

*** Test Cases ***
Docker Secret Tests
    Docker secret ls
    Docker secret create
    Docker secret inspect
    Docker secret rm

Docker Checkpoint Tests
    Docker checkpoint create
    Docker checkpoint ls
    Docker checkpoint rm

Docker Deploy Tests
    Docker deploy

Docker Node Tests
    Docker node demote
    Docker node ls
    Docker node promote
    Docker node rm
    Docker node update
    Docker node ps
    Docker node inspect

Docker Plugin Tests
    Docker plugin install
    Docker plugin create
    Docker plugin enable
    Docker plugin disable
    Docker plugin inspect
    Docker plugin ls
    Docker plugin push
    Docker plugin rm
    Docker plugin set

Docker Service Tests
    Docker service create 
    Docker service ls
    Docker service ps
    Docker serivce rm
    Docker service scale
    Docker service update
    Docker service logs

Docker Stack Tests
    Docker stack ls
    Docker stack ps
    Docker stack rm
    Docker stack services

Docker Swarm Tests
    Docker swarm init
    Docker swarm join
    Docker swarm join-token
    Docker swarm leave
    Docker swarm unlock-key
    Docker swarm update
