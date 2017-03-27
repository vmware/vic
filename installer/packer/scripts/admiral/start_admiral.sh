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

set -x -euf -o pipefail

PORT=""
ADMIRAL_TLS=""
ADMIRAL_DATA_LOCATION=""
ADMIRAL_KEY_LOCATION=""
ADMIRAL_CERT_LOCATION=""

if [ -n $ADMIRAL_TLS ]; then
	/usr/bin/docker run -d -p ${PORT}:${PORT} \
	 --name vic-admiral \
	 -e XENON_OPTS="--sandbox=${ADMIRAL_DATA_LOCATION}" \
	 --log-driver=json-file \
	 --log-opt max-size=1g \
	 --log-opt max-file=10 \
		 vmware/admiral:dev-vic
else
	/usr/bin/docker run -d -p ${PORT}:${PORT} \
	 --name vic-admiral \
	 -e ADMIRAL_PORT=-1 \
	 -e XENON_OPTS="--sandbox=${ADMIRAL_DATA_LOCATION} --securePort=${PORT} --certificateFile=/tmp/server.crt --keyFile=/tmp/server.key --port=-1" \
	 -v "$ADMIRAL_CERT_LOCATION:/tmp/server.cert" \
	 -v "$ADMIRAL_KEY_LOCATION:/tmp/server.key" \
	 --log-driver=json-file \
	 --log-opt max-size=1g \
	 --log-opt max-file=10 \
		 vmware/admiral:dev-vic
fi
