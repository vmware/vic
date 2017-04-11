#!/usr/bin/bash
# Copyright 2017 VMware, Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
set -euf -o pipefail

# Make sure we're using the traditional naming scheme for network interfaces
sed -i '/linux/ s/$/ net.ifnames=0/' /boot/grub2/grub.cfg
echo 'GRUB_CMDLINE_LINUX=\"net.ifnames=0\"' >> /etc/default/grub

# Disable console blanking
sed -i '/linux/ s/$/ consoleblank=0/' /boot/grub2/grub.cfg


# Enable systemd services
systemctl daemon-reload
systemctl enable docker.service
systemctl enable data.mount repartition.service resizefs.service getty@tty2.service
systemctl enable chrootpwd.service sshd_permitrootlogin.service vic-appliance.target
systemctl enable ovf-network.service
systemctl enable ova-firewall.service
systemctl enable harbor_startup.service harbor.service
systemctl enable admiral_startup.service admiral
systemctl enable fileserver_startup.service fileserver.service

# Clean up temporary directories
rm -rf /var/tmp/harbor

# seal the template
> /etc/machine-id
rm /etc/hostname
ls -al /etc/ssh
rm /etc/ssh/{ssh_host_dsa_key,ssh_host_dsa_key.pub,ssh_host_ecdsa_key,ssh_host_ecdsa_key.pub,ssh_host_ed25519_key,ssh_host_ed25519_key.pub,ssh_host_rsa_key,ssh_host_rsa_key.pub}
ls -al /etc/ssh
umount /data
