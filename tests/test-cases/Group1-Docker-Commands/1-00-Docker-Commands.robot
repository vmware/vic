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
Resource  subtests/1-05-Docker-Start.robot
Resource  subtests/1-06-Docker-Run.robot
Resource  subtests/1-07-Docker-Stop.robot
Resource  subtests/1-08-Docker-Logs.robot
Resource  subtests/1-09-Docker-Attach.robot
Resource  subtests/1-10-Docker-PS.robot
Resource  subtests/1-11-Docker-RM.robot
Resource  subtests/1-12-Docker-RMI.robot
Resource  subtests/1-13-Docker-Version.robot
Resource  subtests/1-14-Docker-Kill.robot
Resource  subtests/1-15-Docker-Network-Create.robot
Resource  subtests/1-16-Docker-Network-LS.robot
Resource  subtests/1-17-Docker-Network-Connect.robot
Resource  subtests/1-18-Docker-Network-RM.robot
Resource  subtests/1-19-Docker-Volume-Create.robot
Resource  subtests/1-20-Docker-Volume-Inspect.robot
Resource  subtests/1-21-Docker-Volume-LS.robot
Resource  subtests/1-22-Docker-Volume-RM.robot
Resource  subtests/1-23-Docker-Inspect.robot
Resource  subtests/1-24-Docker-Link.robot
Resource  subtests/1-25-Docker-Port-Map.robot
Resource  subtests/1-26-Docker-Hello-World.robot
Resource  subtests/1-27-Docker-Login.robot
Resource  subtests/1-36-Docker-Rename.robot
Resource  subtests/1-37-Docker-USER.robot
Resource  subtests/1-38-Docker-Exec.robot
Resource  subtests/1-39-Docker-Stats.robot
Resource  subtests/1-40-Docker-Restart.robot
Resource  subtests/1-41-Docker-Commit.robot
Resource  subtests/1-42-Docker-Diff.robot
Resource  subtests/1-43-Docker-CP-Offline.robot
Resource  subtests/1-44-Docker-CP-Online.robot
Resource  subtests/Docker-Command-Utils.robot
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

Docker Start Tests
    Simple start
    Start from image that has no PATH
    Start non-existent container
    Start with no ethernet card
    Serially start 5 long running containers
    Parallel start 5 long running containers
    Start a container with removed network
    Simple start with attach

Docker Run Tests
    Simple docker run
    Docker run with -t
    Simple docker run with app that doesn't exit
    Docker run fake command
    Docker run fake image
    Docker run named container
    Docker run linked containers
    Docker run -d unspecified host port
    Docker run check exit codes
    Docker run ps password check
    Docker run immediate exit
    Docker run verify container start and stop time
    Docker run verify name and id are not conflated
    Docker run and auto remove
    Docker run mysql container
    Docker run mariadb container
    Docker run postgres container
    Docker run --hostname to set hostname and domainname

Docker Stop Tests
    Stop an already stopped container
    Basic docker container stop
    Basic docker stop w/ unclean exit from running process
    Stop a container with SIGKILL using default grace period
    Stop a container with SIGKILL using specific stop signal
    Stop a container with SIGKILL using specific grace period
    Stop a non-existent container
    Attempt to stop a container that has been started out of band
    Restart a stopped container
    Stop a container with Docker 1.13 CLI

Docker Logs Tests
    #Docker logs backward compatibility
    Docker logs with tail
    Docker logs with follow
    Docker logs with follow and tail
    Docker logs follow shutdown
    Docker binary logs
    Docker text logs
    Docker logs with timestamps and since certain time
    Docker logs with no flags
    Docker logs non-existent container

Docker Attach Tests
    Basic attach
    Attach to stopped container
    Attach to fake container
    Attach with short input
    Attach with short output
    Attach with short output with tty

Docker PS Tests
    Remove All Containers
    Empty docker ps command
    Docker ps only running containers
    Docker ps all containers
    Docker ps powerOn container OOB
    Docker ps powerOff container OOB
    Docker ps ports output
    Docker ps Remove container OOB
    Docker ps last container
    Docker ps two containers
    Docker ps last container with size
    Docker ps all containers with only IDs
    Docker ps with status filter
    Docker ps with label and name filter
    Docker ps with volume filter
    Docker ps with network filter
    Docker ps with volume and network filters

Docker RM Tests
    Simple volume rm
    Volume rm when in use
    Volume rm invalid volume
    Volume rm freed up volume

Docker Version Tests
    Simple Docker Version
    Docker Version Format Client Version
    Docker1.11 Version Format Client API Version
    Docker1.13 Version Format Client API Version
    Docker Version Format Client Go Version
    Docker Version Format Server Version
    Docker1.11 Version Format Server API Version
    Docker1.13 Version Format Server API Version
    Docker1.13 Version Format Server Minimum API Version
    Docker Version Format Server Go Version
    Docker Version Format Bad Field

Docker Kill Tests
    Signal a container with default kill signal
    Signal a container with SIGHUP
    Signal a non-existent container
    Signal a tough to kill container - nginx

Docker Inspect Tests
    Simple docker inspect of image
    Docker inspect image specifying type
    Docker inspect image specifying incorrect type
    Simple docker inspect of container
    Docker inspect container specifying type
    Docker inspect container check cmd and image name
    Docker inspect container specifying incorrect type
    Docker inspect container with multiple networks
    Docker inspect invalid object
    Docker inspect non-nil volume
    Inspect RepoDigest is valid
    Docker inspect mount data

Docker Link Tests
    Link and alias

Docker Port Map Tests
    Create container with port mappings
    Create container with conflicting port mapping
    Create container with port range
    Create container with host ip
    Create container with host ip equal to 0.0.0.0
    Create container with host ip equal to public IP
    Create container without specifying host port
    Run after exit remapping mapped ports
    Remap mapped ports after OOB Stop
    Remap mapped ports after OOB Stop and Remove
    Container to container traffic via VCH public interface
    Remap mapped port after stop container, and then remove stopped container

Docker Hello-World Tests
    Hello world
    Hello world with -t

Docker Login Tests
    Docker login and pull from docker.io

Docker Rename Tests
    Remove All Containers
    Rename a non-existent container
    Rename a created container
    Rename a running container
    Rename a stopped container
    Rename a container with an empty name
    Rename a container with a claimed name
    Name resolution for a created container after renaming+starting it
    Name resolution for a running container after renaming+restarting it
    Name resolution for a running container after renaming it

Docker USER Tests
    Run Image Specifying NewUser in NewGroup
    Run Image Specifying UID 2000
    Run Specifying UID 2000 With -u
    Run Image Specifying UID:GID 2000:2000
    Run Specifying UID:GID 2000:2000 With -u
    Run as Nonexistent User With -u
    Run as Root with Nonexistent User With -u
    Run as uid 0 group 0 With -u

Docker Exec Tests
    Exec -d
    Exec Echo
    Exec Echo -i
    Exec Echo -t
    Exec Sort
    Exec Sort -i
    Exec NonExisting

Docker Stats Tests
    Stats No Stream
    Stats No Stream All Containers
    Stats API Memory Validation
    Stats API CPU Validation
    Stats No Stream Non-Existent Container
    Stats No Stream Specific Stopped Container
    Stats API Disk and Network Validation

Docker Restart Tests
    Restart Running Container
    Restart Created Container
    Restart Stopped Container
    Restart with start-stop stress

Docker Commit Tests
    Commit nano to image
    Commit env variable to image
    Unsupported commit command
    Commit with author and message
    Commit to nonexistent container

Docker Diff Tests
    Make changes to busybox image

Docker CP Offfline Tests
    Set up test files for CP Offline tests
    Copy a file from host to offline container root dir
    Copy a directory from offline container to host cwd
    Copy a directory from host to offline container, dst path doesn't exist
    Copy a non-existent file out of an offline container
    Copy a non-existent directory out of an offline container
    Copy a non-existent directory into an offline container
    Copy a large file that exceeds the container volume into an offline container
    Copy a file from host to offline container, dst is a volume
    Copy a file from host to offline container, dst is a nested volume with 2 levels
    Copy a file from host to offline container, dst is a nested volume with 3 levels
    Concurrent copy: create processes to copy a small file from host to offline container
    Concurrent copy: repeat copy a large file from host to offline container several times
    Concurrent copy: repeat copy a large file from offline container to host several times
    Sub volumes: copy from host to offline container
    Sub volumes: copy from offline container to host
    Clean up test files for CP Offline tests

Docker CP Online Tests
    Set up test files for CP Online Tests
    Copy a directory from online container to host, dst path doesn't exist
    Copy the content of a directory from online container to host
    Copy a file from online container to host, overwrite dst file
    Copy a file from host to online container, dst directory doesn't exist
    Copy a file and directory from host to online container
    Copy a directory from host to online container, dst is a volume
    Copy a file from host to offline container, dst is a volume shared with an online container
    Copy a directory from offline container to host, dst is a volume shared with an online container
    Copy a large file to an online container, dst is a volume
    Copy a non-existent file out of an online container
    Copy a non-existent directory out of an online container
    Concurrent copy: create processes to copy a small file from host to online container
    Concurrent copy: repeat copy a large file from host to online container several times
    Concurrent copy: repeat copy a large file from online container to host several times
    Sub volumes: copy from host to an online container, dst includes several volumes
    Sub volumes: copy from online container to host, src includes several volumes
    Sub volumes: copy from host to an offline container, dst includes a shared vol with an online container
    Sub volumes: copy from an offline container to host, src includes a shared vol with an online container
    Clean up test files for CP Online Tests

Docker RMI Tests
    Remove All Containers
    Basic docker pull, restart, and remove image
    Remove image with a removed container
    Remove image with a container
    Remove a fake image
    Remove an image pulled by digest

#Docker Finish
#    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} images
#    ${rc}  ${output}=  Run And Return Rc And Output  docker %{VCH-PARAMS} ps
#    Should Be Equal As Integers  ${rc}  0
