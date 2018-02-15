Test 1-38 - Docker Exec
=======

# Purpose:
To verify that docker exec command is supported by VIC appliance

# References:
[1 - Docker Command Line Reference](https://docs.docker.com/engine/reference/commandline/exec/)

# Environment:
This test requires that a vSphere server is running and available

# Test Steps:
1. Deploy VIC appliance to vSphere server
2. Issue docker run -d busybox /bin/top
3. Issue docker exec <containerID> /bin/echo ID - 5 times with incrementing ID
4. Issue docker exec -i <containerID> /bin/echo ID - 5 times with incrementing ID
5. Issue docker exec -t <containerID> /bin/echo ID - 5 times with incrementing ID
6. Issue docker exec -it <containerID> /bin/echo ID - 5 times with incrementing ID
7. Issue docker exec -it <containerID> NON_EXISTING_COMMAND

# Expected Outcome:
* Step 2-6 should echo the ID given
* Step 7 should return an error

## Exec

# Possible Problems:
None

# Exec Power Off test for long running Process
## Test Steps
1. Pull an image that contains `/bin/top`. Busybox suffices here.
2. Create a container running `/bin/top` to simulate a long running process.
3. Run the container and detach from it.
4. Start 20 simple exec operations in parallel against the detached container.
5. Stop the container while the execs are still running in order to trigger the exec power off errors.
6. collect all output from the parallel exec operations.

## Expected Outcome
* step 1 should successfully complete with an rc of 0
* step 2 should successfully complete with an rc of 0
* step 3 should successfully launch the container and detach from it. The container should remain running for 5 seconds only.
* step 4 all 20 execs should be started successfully, we expect most if not all to fail.
* step 5 container should halt successfully.
* step 6 should contain the error message for exec operations that are interrupted by a power off operation. Specifically a poweroff that was explicitly triggered.

# Exec Power Off test for short Running Process
## Test Steps
1. Pull an image that contains `sleep`. Busybox suffices here.
2. Create a container running `sleep` to simulate a long running process.
3. Run the container and detach from it.
4. Start 20 simple exec operations in parallel against the detached container.
5. Wait(`docker wait`) for the container to exit just in case it has not(20 execs should be ample stress).
6. collect all output from the parallel exec operations.

## Expected Outcome
* step 1 should successfully complete with an rc of 0
* step 2 should successfully complete with an rc of 0
* step 3 should successfully launch the container and detach from it. The container should remain running for 5 seconds only and should return with an RC of 0.
* step 4 all 20 execs should be started successfully, we expect most if not all to fail.
* step 5 container should halt successfully.
* step 6 should contain the error message for exec operations that are interrupted by a power off operation. Specifically a poweroff that was implicitly triggered.
