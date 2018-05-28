#!/bin/bash
# Copyright 2018 VMware, Inc. All Rights Reserved.
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
#

echo "base customizations for $(basename $(dirname $0))"

# cpio will return an error on open syscall if lib64 is a symlink
rm -f $(rootfs_dir $PKGDIR)/lib64

# work arounds for libgcc.x86_64 0:4.4.7-18.el6 needing /dev/null
# it looks like udev package will create this but has it's own issues
mkdir -p $(rootfs_dir $PKGDIR)/{dev,lib64}
mknod $(rootfs_dir $PKGDIR)/dev/null c 1 3
chmod 666 $(rootfs_dir $PKGDIR)/dev/null