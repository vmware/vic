Test 9-01 - VIC Admin ShowHTML
=======

#Purpose:
To verify that the VIC Administration appliance can display HTML

#Environment:
This test requires that a vSphere environment be running and available

#Test Steps:
1. Deploy VIC appliance to the vSphere server
2. Pull the VICadmin web page and verify that it contains valid HTML
3. Pull the Portlayer log file and verify that it contains valid data
4. Pull the VCH-Init log and verify that it contains valid data
5. Pull the Docker Personality log and verify that it contains valid data
6. Create a container via the appliance
7. Pull the container log bundle from the appliance and verify that it contains the new container's logs

#Expected Outcomes:
* VICadmin should display a web page that at a minimum includes <title>VIC Admin</title>
* VICadmin responds with a log file indicating that the portlayer sever has started
* VICadmin responds with a log file indicating VCH init has begun reaping processes
* VICadmin responds with log file indicating docker personality service has started
* VICadmin responds with a ZIP file containing at a minimum the vmware.log file from the new container
