# Access vSphere Integrated Containers Engine Log Bundles #

vSphere Integrated Containers Engine provides log bundles that you can download from the VCH Admin portal for a virtual container host (VCH).

You access the VCH Admin Portal at https://<i>vch_address</i>:2378.

If the VCH is unable to connect to vSphere, logs that require a vSphere connection are disabled, and you see an error message. You can download the log bundle to troubleshoot the error.

- The **Log Bundle** contains logs that relate specifically to the VCH that you created. 
- The **Log Bundle with container logs** contains the logs for the VCH and also includes the logs regarding the containers that the VCH manages.
- Live logs (tail files) allow you to view the current status of how components are running.
  - **Docker Personality** is the interface to Docker. When configured with client certificate security, it reports unauthorized access attempts to the Docker server web page.
  - **Port Layer Service** is the interface to vSphere.
  - **Initialization & watchdog** reports:
  		- Network configuration
  		- Component launch status for the other components
  		- Reports component failures and restart counts

  	At higher debug levels, the component output is duplicated in the log files for those components, so `init.log`  includes a superset of the log data.

    **Note:** This log file is duplicated on the datastore in a file in the endpoint VM folder named `tether.debug`, to allow the debugging of early stage initialization and network configuration issues.

  - **Admin Server** includes logs for the VCH admin server, may contain processes that failed, and network issues. When configured with client certificate security, it reports unauthorized access attempts to the admin server web page.

Live logs can help you to see information about current commands and changes as you make them. For example, when you are troubleshooting an issue, you can see whether your command worked or failed by looking at the live logs.

You can share the non-live version of the logs with administrators or VMware Support to help you to resolve issues.

Logs also include vic-machine commands used during VCH installation to help you resolve issues.

## Collecting Logs Manually
If the VCH Admin portal is offline, use `vic-machine debug` to enable SSH on the VCH and use `scp -r` to capture the logs from `/var/log/vic/`.

## Setting the Log Size Cap
The log size cap is set at 20MB. If the size exceeds 20 MB, then vSphere Integrated Containers Engine compresses the files and saves a history of the last two rotations. These files are rotated:

`/var/log/vic/port-layer.log` <br>
`/var/log/vic/init.log` <br>
`/var/log/vic/docker-personality.log` <br>
`/var/log/vic/vicadmin.log`