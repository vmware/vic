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
