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
Resource  subtests/1-01-Docker-Info.robot
Resource  subtests/1-02-Docker-Pull.robot
Resource  subtests/1-03-Docker-Images.robot
Resource  subtests/1-04-Docker-Create.robot
Suite Setup  Install VIC Appliance To Test Server
Suite Teardown  Cleanup VIC Appliance On Test Server
Test Timeout  30 minutes

*** Test Cases ***
Docker Info Tests
    Docker Info - Basic Info
    Docker Info - Debug Info
    Docker Info - Correct container count
    Docker Info - Check modified resource pool CPU and memory limits
    Docker Info - Check updated resource pool CPU and memory usages

Docker Pull Tests
    Pull nginx
    Pull busybox
    Pull ubuntu
    Pull non-default tag
    Pull images based on digest
    Pull an image with the full docker registry URL
    Pull an image with all tags
    Pull non-existent image
    Pull image from non-existent repo
    Pull image with a tag that doesn't exist
    Pull image that already has been pulled
    Pull the same image concurrently
    Pull two images that share layers concurrently
    Re-pull a previously rmi'd image
    Pull image by multiple tags
    Issue docker pull on digest outputted by previous pull
    Pull images from gcr.io
    Verify image manifest digest against vanilla docker

Docker Images Tests
    Simple images
    All images
    Quiet images
    No-trunc images
    Filter images before
    Filter images since
    Tag images
    Specific images
    VIC/docker Image ID consistency

Docker Create Tests
    Simple creates
    Create with anonymous volume
    Create with named volume
    Create with a directory as a volume
    Create simple top example
    Create fakeimage image
    Create fakeImage repository
    Create and start named container
    Create linked containers that can ping
    Create a container after the last container is removed
    Create a container from an image that has not been pulled yet
    Create a container with no command specified
    Create a container with custom CPU count
    Create a container with custom amount of memory in GB
    Create a container with custom amount of memory in MB
    Create a container with custom amount of memory in KB
    Create a container with custom amount of memory in Bytes
    Create a container using rest api call without HostConfig in the form data
    Create a container and check the VM display name and datastore folder name
    Create disables VC destroy

Docker Finish
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} images
    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps
    Should Be Equal As Integers  ${rc}  0
