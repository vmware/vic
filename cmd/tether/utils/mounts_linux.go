// Copyright 2016 VMware, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"

	"golang.org/x/net/context"
)

// MountLabel performs a mount with the source treated as a disk label
// This assumes that /dev/disk/by-label is being populated, probably by udev
func MountLabel(label, target string, ctx context.Context) error {
	if err := os.MkdirAll(target, 0600); err != nil {
		return fmt.Errorf("unable to create mount point %s: %s", target, err)
	}

	volumes := "/dev/disk/by-label"
	source := volumes + "/" + label

	// do..while ! timedout
	var timeout bool
	for timeout = false; !timeout; {
		_, err := os.Stat(source)
		if err == nil || !os.IsNotExist(err) {
			break
		}

		deadline, ok := ctx.Deadline()
		timeout = ok && time.Now().After(deadline)
	}

	if timeout {
		detail := fmt.Sprintf("timed out waiting for %s to appear", source)
		return errors.New(detail)
	}

	if err := syscall.Mount(source, target, "ext4", syscall.MS_NOATIME, ""); err != nil {
		detail := fmt.Sprintf("mounting %s on %s failed: %s", source, target, err)
		return errors.New(detail)
	}

	return nil
}
