# Add Pre- and Post-Initialization Scripts to the vSphere Integrated Containers Appliance #

The options that the vSphere Integrated Containers installer provides do not allow for all possible vSphere configurations. After installation, your vSphere environment might not allow you to use vSphere Integrated Containers in the way that you require. 

For example, you might require the vSphere Integrated Containers appliance to have an additional network interface with a static IP address, so that you can expose the virtual container host to users who do not have access to the vCenter Server management network. This is not possible with the current installer options, which do not allow for multiple network interfaces. Achieving such a configuration requires you to reconfigure the vSphere Integrated Containers appliance after deployment. 

If the reconfiguration must be persistent, you can provide scripts that perform additional configuration on the vSphere Integrated Containers appliance during its initialization process. The scripts can be any arbitrary executable that can run on the appliance, for example UNIX shell scripts.

You can provide scripts that run either before or after the initialization of the appliance, or both. 

- Pre-initialization scripts run immediately after the mount of the Docker metadata disk, before the initialization and setup of the appliance.
- Post-initialization scripts run after the initialization and setup of the appliance, but before the Docker daemon starts.

**NOTE** Any pre- or post-initialization configuration that you perform via scripting is unsupported. This functionality is provided purely so that you can test the vSphere Integrated Containers technical preview in your specific vSphere environment.

**Prequisites**

- Deploy the vSphere Integrated Containers appliance. 
- Verify that the deployment was successful and obtain the address of the appliance.

**Procedure**

1. Log in to the vSphere Integrated Containers appliance.
2. Navigate to `/var/lib/docker/`.
3. Create folders named `pre` and `post` under `/var/lib/docker/`.
4. Copy script files into the `pre` and `post` folders, or create the scripts directly in those folders.

 - Place pre-initialization scripts in the `pre` folder.
 - Place post-initialization scripts in the `post` folder.
5. Ensure that the scripts are executable.
6. Restart the vSphere Integrated Containers appliance.

