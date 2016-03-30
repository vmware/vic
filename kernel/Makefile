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

SPECFILE = linux-esx.spec
CONFIG = config-esx-4.2.0

srpm:
	sed -i -e "s:^CONFIG_SERIAL_8250=m:CONFIG_SERIAL_8250=y:g" $(CONFIG)

	cp $(CONFIG) /usr/src/photon/SOURCES/
	cp ??-*.patch /usr/src/photon/SOURCES/
	rpmbuild --define "_topdir /usr/src/photon" --nodeps -bs $(SPECFILE)

TARGET ?= $(shell uname -m)
rpm:
	rpmbuild --define "_topdir /usr/src/photon" --nodeps --target $(TARGET) -bb $(SPECFILE)

clean:
	rm -f *.rpm *.tgz *.bz2 *.gz

++%: varname=$(subst +,,$@)
++%:
	@echo "$(varname)=$($(varname))"
+%: varname=$(subst +,,$@)
+%:
	@echo "$($(varname))"
