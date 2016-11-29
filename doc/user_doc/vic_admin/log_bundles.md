# Access vSphere Integrated Containers Engine Log Bundles #

vSphere Integrated Containers Engine provides log bundles that you can download from the VCH Admin portal for a virtual container host (VCH).

If the VCH is unable to connect to vSphere, logs that require a vSphere connection are disabled, and you see an error message. You can download the log bundle to troubleshoot the error.

- The **Log Bundle** contains logs that relate specifically to the VCH that you created. 
- The **Log Bundle with container logs** contains the logs for the VCH and also includes the logs regarding  the containers that the VCH manages.
- Live logs (tail files) allow you to view the current status of how components are running.
  - **Docker Personality** is the interface to Docker. When configured with client certificate security, it reports unauthorized access attempts to the Docker server web page.
  - **Port Layer Service** is the interface to vSphere.
  - **Initialization & watchdog** reports network configuration, component launch status for the VCH Admin portal and the port layer,  records if they fail, and relaunches them if they do. The binary  `vic-init` launches the components and redirects their output to the log files in `/var/log/vic/`. At higher debug levels, the component output is duplicated in that log file, so `init.log`  includes a superset of the log data.
  - **Admin Server** includes logs for the VCH admin server, may contain processes that failed, and network issues. When configured with client certificate security, it reports unauthorized access attempts to the admin server web page.

Live logs can help you see information about current commands and changes as you perform them. For example, when you are troubleshooting an issue, you can see if your command worked or failed by looking at the live logs.

You can share the non-live version of the logs with administrators or VMware Support to help you resolve issues.

## Collecting Logs Manually
If the VCH Admin portal is offline, use `vic-machine debug` to enable SSH on the VCH and use `scp -r` to capture the logs from `/var/log/vic/`.
