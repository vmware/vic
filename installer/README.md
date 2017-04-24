# VIC Unified Installer

This directory will host all the code that is going to be part of the VIC unified installer OVA.

It is currently under heavy development and not suitable for any use except for development, this file will be updated to reflect the status of the installer as development progresses.

### Usage

```
esxcli system settings advanced set -o /Net/GuestIPHack -i 1
esxcli network firewall set --enabled false
```

The machine that is running Packer (make ova-release) must be reachable from the launched VM and
have `ovftool` installed

#### Build bundle and OVA

First, we have to set the revisions of the components we want to bundle in the OVA:

```
export BUILD_HARBOR_REVISION=1.1.0-rc1     # Optional, defaults to dev
export BUILD_ADMIRAL_REVISION=v1.1.0-rc1   # Optional, defaults to dev
export BUILD_VICENGINE_REVISION=1.1.0-rc4  # Required
```

Then set the required env vars for the build environment and make the release:

```
export PACKER_ESX_HOST=1.1.1.1
export PACKER_USER=root
export PACKER_PASSWORD=password
export PACKER_LOG=1

make ova-release
```

Deploy OVA with ovftool in a Docker container on ESX host
```
docker run -it --net=host -v ~/go/src/github.com/vmware/vic/bin:/test-bin \
  gcr.io/eminent-nation-87317/vic-integration-test:1.27 ovftool --acceptAllEulas --X:injectOvfEnv \
  --X:enableHiddenProperties -st=OVA --powerOn --noSSLVerify=true -ds=datastore1 -dm=thin \
  --net:Network="VM Network" \
  --prop:appliance.root_pwd="VMware1\!" --prop:appliance.permit_root_login=True --prop:registry.port=443 \
  --prop:management_portal.port=8282 --prop:registry.admin_password="VMware1\!" \
  --prop:registry.db_password="VMware1\!" /test-bin/vic-1.1.0-a84985b.ova \
  vi://root:password@192.168.1.20
```

### Troubleshooting

#### ova-release failed

```
2017/03/16 10:26:25 packer: 2017/03/16 10:26:25 starting remote command: test -e
/vmfs/volumes/datastore1/vic
2017/03/16 10:26:25 ui error: ==> ova-release: Step "StepOutputDir" failed, aborting...
==> ova-release: Step "StepOutputDir" failed, aborting...
Build 'ova-release' errored: unexpected EOF

==> Some builds didn't complete successfully and had errors:
2017/03/16 10:26:25 ui error: Build 'ova-release' errored: unexpected EOF
2017/03/16 10:26:25 Builds completed. Waiting on interrupt barrier...
2017/03/16 10:26:25 machine readable: error-count []string{"1"}
2017/03/16 10:26:25 ui error:
==> Some builds didn't complete successfully and had errors:
2017/03/16 10:26:25 machine readable: ova-release,error []string{"unexpected EOF"}
2017/03/16 10:26:25 ui error: --> ova-release: unexpected EOF
2017/03/16 10:26:25 ui:
==> Builds finished but no artifacts were created.
2017/03/16 10:26:25 waiting for all plugin processes to complete...
2017/03/16 10:26:25 /usr/local/bin/packer: plugin process exited
2017/03/16 10:26:25 /usr/local/bin/packer: plugin process exited
2017/03/16 10:26:25 /usr/local/bin/packer: plugin process exited
2017/03/16 10:26:25 /usr/local/bin/packer: plugin process exited
--> ova-release: unexpected EOF

==> Builds finished but no artifacts were created.
2017/03/16 10:26:25 /usr/local/bin/packer: plugin process exited
installer/vic-unified-installer.mk:31: recipe for target 'ova-release' failed
make: *** [ova-release] Error 1
```

Solution: Cleanup datastore by removing the `vic` folder


#### Connection refused

```
2017/03/16 12:48:46 ui: ==> ova-release: Connecting to VM via VNC
==> ova-release: Connecting to VM via VNC
2017/03/16 12:49:13 ui error: ==> ova-release: Error connecting to VNC: dial tcp 10.17.109.107:5900:
getsockopt: connection refused
==> ova-release: Error connecting to VNC: dial tcp 10.17.109.107:5900: getsockopt: connection
refused
```

Solution: Disable firewall on ESX host `esxcli network firewall set --enabled false`

#### No IP address ready

```
2017/03/23 12:03:45 packer: 2017/03/23 12:03:45 opening new ssh session
2017/03/23 12:03:45 packer: 2017/03/23 12:03:45 starting remote command: esxcli --formatter csv
network vm list
2017/03/23 12:03:46 packer: 2017/03/23 12:03:46 opening new ssh session
2017/03/23 12:03:46 packer: 2017/03/23 12:03:46 starting remote command: esxcli --formatter csv
network vm port list -w 73094
2017/03/23 12:03:46 packer: 2017/03/23 12:03:46 [DEBUG] Error getting SSH address: No interface on
the VM has an IP address ready
```

Solution: Disable firewall on the build machine. The launched VM is unable to get the kickstart file
from your build machine.
