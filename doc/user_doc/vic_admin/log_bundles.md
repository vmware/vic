# Access vSphere Integrated Containers Engine Log Bundles #

vSphere Integrated Containers Engine provides log bundles that you can download from the Admin portal.

- The **Log Bundle** contains logs that relate specifically to the virtual container host that you created.
- The **Log Bundle with container logs** also includes the logs regarding  the containers that a virtual container host manages.
- Live logs (tail files) allow you to view the absolute current status of how components are running.
<ul>
- **Docker Personality** is the interface to Docker. When configured with client certificate security, it reports unauthorized access attempts to the Docker server web page.
- **Port Layer Service** is the interface to vSphere.
- **Initialization & watchdog** reports network configuration, component launch status for vic-admin and port layer, and records if they fail and relaunches them if they do. The binary  `vic-init` launches the components and redirects their output to the log files in /var/log/vic/ <br>
At higher debug levels, the component output is duplicated in that log file, so init.log  includes a superset of the log data.
- **Admin Server** includes logs for the admin server, may contain processes that failed, and network issues. When configured with client certificate security, it reports unauthorized access attempts to the admin server web page.

Live logs can help you see how any current changes you make might affect the logs. For example, when you try to trouble shoot an issue, you can see if your attempt worked or failed by looking at the live logs.

The non-live version of the logs are good for sharing with administrators or others to help solve issues.
