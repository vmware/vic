# Tether

The tether is an init replacement used in containerVMs that provides the command & control channel necessary to perform any operation inside the container. This includes launching of the container process, setting of environment variables, configuration of networking, etc. The tether is currently based on a modified SSH server tailed specifically for this purpose.

[Issues relating to this component](https://github.com/vmware/vic/labels/component%2Ftether)
