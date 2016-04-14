# vSocket Relay Agent

Network serial ports as a communication channel have several drawbacks:

* serial is not intended for high bandwidth, high frequency data
* inhibits forking & vMotion without vSPC, and a vSPC requires an appliance in FT/HA configuration
* requires a VCH have a presence on the management networks
* requires opening a port on the ESX firewall

The alternative we're looking at is vSocket (uses PIO based VMCI communication), however that is Host<->VM only so we need a mechanism to relay that communication to the VCH. Initially it's expected that the Host->VCH communication will be a TCP connection for a staged delivery approach, with the longer term being an agent<->agent relay between the two hosts.

[Issues relating to this component](https://github.com/vmware/vic/labels/component%2Fvsocket-relay-agent)
