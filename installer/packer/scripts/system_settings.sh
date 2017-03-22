#!/bin/sh
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
systemctl enable ovf-network.service harbor_startup admiral


# Clean up temporary directories
rm -rf /var/tmp/*

# seal the template
> /etc/machine-id
rm /etc/hostname
rm /etc/ssh/ssh_host_*
umount /data
