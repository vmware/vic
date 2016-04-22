# Add Pre- and Post-Initialization Scripts to the vSphere Integrated Containers Appliance #

The options that the vSphere Integrated Containers installer provides do not allow for all possible vSphere configurations. After installation, your vSphere environment might not allow you to use vSphere Integrated Containers in the way that you require. 

To allow vSphere Integrated Containers to run correctly in your environment, you might have to manually reconfigure the vSphere Integrated Containers appliance. Because the appliance is stateless, any manual reconfiguration that you perform is lost if you reboot the appliance. To make a reconfiguration persistent, you can provide scripts that configure the appliance during its initialization process. The scripts can be any arbitrary executable that can run on the appliance, for example UNIX shell scripts. The scripts can run either before or after the initialization of the appliance, or both.

- Pre-initialization scripts run immediately after the mount of the Docker metadata disk, before the initialization and setup of the appliance.
- Post-initialization scripts run after the initialization and setup of the appliance, but before the Docker daemon starts.

**NOTE** Any pre- or post-initialization configuration that you perform via scripting is unsupported. This functionality is provided purely so that you can test the vSphere Integrated Containers technical preview in your specific vSphere environment.

For example, you might require the vSphere Integrated Containers appliance to have an additional network interface with a static IP address, so that you can expose the virtual container host to users who do not have access to the vCenter Server management network. This is not possible with the current installer options, which do not allow for multiple network interfaces. Achieving such a configuration requires you to reconfigure the vSphere Integrated Containers appliance after deployment. 

Similarly, your vSphere environment might also require TLS authentication to be persistent across reboots of the vSphere Integrated Containers appliance.   

**Prequisites**

- Deploy the vSphere Integrated Containers appliance.

 To use pre- and post-initialization scripts, you must use build 58 or later of the vSphere Integrated Containers command line installer.
- Verify that the deployment was successful and obtain the address of the appliance.

**Procedure**

1. Use SSH to log in to the vSphere Integrated Containers appliance as `root`.
2. Create executable files named `pre`, or `post`, or both, under `/var/lib/docker/`.

 - Write pre-initialization scripting in the `pre` executable.
 - Write post-initialization scripting in the `post` executable.
 - Use `#!/bin/ash` as the script interpreter
 - Ensure that the scripts are executable.
3. Restart the vSphere Integrated Containers appliance.

For an example of how to use post-initialization scripting to implement persistent TLS authentication, see [Appendix: Example of Implementing Persistent TLS Authentication](appendix_persistent_tls.md).