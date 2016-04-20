# Tether

The tether provides two distinct sets of function:

1. the channel necessary to present the process input/output as if it were bound directly to the users local console.
2. configuration of the operating system - where regular docker can manipulate elements such as network interfaces directly, or use bind mounting of files to inject data and configuration into a container, this approach is not available when the container is running as an isolated VM; the tether has to take on all responsibility for correct configuration of the underlying OS before launching the container process.

## Operating System configuration

### Hostname

### Name resolution

### Filesystem mounts

### Network interfaces


## Management behaviours
This is a somewhat arbitrary divide, but these are essentially ongoing concerns rather than one-offs on initial start of the containerVM

### Signal handling

### Process reaping

### Secrets

### Forking
This relates specifically to VMfork, aka _Instant Clone_ - the ability to freeze a VM and spin off arbitrary numbers of children that inherit that parent VMs memory state and configuration. This requires cooperation between ESX, GuestOS, and the application processes, to handle changes to configuration such as time, MAC and IP addresses, ARP caches, open network connections, etc.  
More directly, VMfork requires that the fork be triggered from within the GuestOS as a means of ensuring it is in a suitable state, meaning the tether has to handle triggering of forking and recovery post-fork


## External communication
The vast bulk of container operations can be performed without a synchronous connection to containerVM, however `attach` is a core element of the docker command set and the one that is probably most heavily used outside of production deployment. This requires that we be able to present the user with a console for the container process.

The initial communication will be via plain network serial port, configured in client mode:

Overall flow:

1. ContainerVM is powered on
2. ContainerVM is configured so com1 is a client network serial port targeting the appliance VM
2. ESX initiates a serial-over-TCP connection to the applianceVM - this link is treated as a reliable bytestream
3. ApplianceVM accepts the connection
4. ApplianceVM acts as an SSH client over the new socket - authentication is negligble for now
5. ApplianceVM uses a custom global request type to retrieve list of containers running in containerVM (exec presents as a separate container); containerVM replies with list
6. ApplianceVM opens a session, requesting a specific container; containerVM acknowledges and glues the streams together


From the Personality to the Portlayer:
* request for container X from Personality to Interaction component
* Interaction configures network serial port to point at it's IP and connects it
* Interaction blocks waiting for SSH session for X, recording the request for a channel to X
* Interaction starts stream copying from HTTP websocket to SSH channel

Incoming TCP connection to Interaction component:
* request set of IDs on the executor
* create entries for each of the IDs in connection map
* if there is a recorded request for channel to an ID
  * establish channel
  * notify waiters
