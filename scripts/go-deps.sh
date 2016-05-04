#!/bin/bash
# Copyright 2016 VMware, Inc. All Rights Reserved.
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
# Lists the non-standard library Go packages the specified package depends
# on.
#
# Usage: script/go-deps.sh pkg
#
#     pkg       This is github.com/vmware/vic/cmd/imagec for example
#

PKG=$1

echo "Generating deps for $PKG" >&2
go list -f '{{join .Deps "\n"}}' $PKG |  xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}' | sed -e 's:github.com/vmware/vic/\(.*\)$:\1/*:'