## Introduction
The Docker Engine-API personality server is what VIC calls the server daemon that responds to Docker remote API calls.  The primary caller is most likely the Docker CLI or a user using curl.

The server, itself, builds on top of Docker's Engine-API project.  This allows this component to be REST compatible with the Docker Daemon.  Once the Engine-API rest server unmarshals the requests into golang structure, execution is handed off to a set of backend code.  These backend code validates the request inputs and calls VIC's port layer server.

VIC calls this a personality server because it translates Docker requests to VIC operations.  The port layer server should not know anything about Docker.  This allows VIC to add other personality servers in the future.
