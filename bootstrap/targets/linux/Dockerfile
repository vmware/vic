FROM container-linux-base

COPY 60-persistent-storage.rules.diff /tmp/
RUN patch --unified -d $ROOTFS -p 0 --input=/tmp/60-persistent-storage.rules.diff

COPY 99-vmware-memhotplug.rules $ROOTFS/etc/udev/rules.d/
COPY 99-vmware-cpuhotplug.rules $ROOTFS/etc/udev/rules.d/

COPY isolinux $BOOTFS/isolinux
RUN cp /usr/lib/ISOLINUX/isolinux.bin $BOOTFS/isolinux/isolinux.bin
RUN cp /usr/lib/syslinux/modules/bios/ldlinux.c32 $BOOTFS/isolinux/ldlinux.c32
COPY vmfork.sh $ROOTFS/

#RUN git clone https://ghicken@enatai-gerrit.eng.vmware.com/bonneville-container $TETHER && cd $TETHER && git checkout proxy
# install the generic tether dependencies
COPY tether $TETHER/tether/
RUN cd $TETHER/tether && gpm install

# Enable SSH debug
#RUN sed -i -e 's/const debugMux = false/const debugMux = true/' /go/src/golang.org/x/crypto/ssh/mux.go && sed -i -e 's/const debugHandshake = false/const debugHandshake = true/' /go/src/golang.org/x/crypto/ssh/handshake.go

RUN cd $TETHER/tether && \
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go test -a -x -tags netgo -installsuffix netgo -v ./...

RUN cd $TETHER/tether/cmd/tether && \
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -a -x -tags netgo -installsuffix netgo -o tether-linux && \
	ldd tether-linux | grep -q 'not a dynamic executable' && \
	cp -f tether-linux $ROOTFS/bin/tether && \
	ln -sf ../bin/tether $ROOTFS/sbin/init && return 0 || \
	echo "error: tether is not static" && return 1

COPY rpctool.go $TETHER/
RUN cd $TETHER/ && \
    GOARCH=amd64 GOOS=linux go build -a -x -o $ROOTFS/sbin/rpctool --ldflags '-extldflags "-static"' rpctool.go

COPY init $ROOTFS/init

ENTRYPOINT ["/bin/bash"]

# construct a script to generate the rootfs and iso - we can't do it as a command as docker gets
# confused by something and it doesn't parse as a valid json array
RUN echo "sed -i -e \"s!BUILD_ID!\$BUILD_ID!\" $BOOTFS/isolinux/boot.msg && \
        cd $ROOTFS && find | cpio -o -H newc | gzip --fast > $BOOTFS/core.gz && \
	xorriso -publisher 'VMware Inc.' -as mkisofs -V bonneville-tcl \
	-l -J -R -no-emul-boot -boot-load-size 4 -boot-info-table \
	-b boot/isolinux/isolinux.bin -c boot/isolinux/boot.cat \
	-o \$ISOOUT /bootfs && \
	mv /binary /tmp/\$BINBASE" > /tmp/generate.sh && chmod a+x /tmp/generate.sh

# Construct root filesystem and generate the ISO to stdout
# To write the iso to a file, set ISOOUT to the target path (directory must exist)
CMD ["/tmp/generate.sh"]
