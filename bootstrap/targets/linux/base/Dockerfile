FROM golang:1.5.3

RUN apt-get update && apt-get -y install curl cpio ssh-client squashfs-tools xorriso git gcc make vim patch rpm busybox-static isolinux syslinux
RUN cd /tmp && git clone https://github.com/pote/gpm.git && cd gpm && ./configure && make install

# Prepare kernel
ENV TCL_BASE http://tinycorelinux.net/6.x
ENV TCL_BASE http://ftp.nluug.nl/os/Linux/distr/tinycorelinux/6.x
ENV TCL_CORE $TCL_BASE/x86_64/release/distribution_files
ENV TCL_TCZ $TCL_BASE/x86_64/tcz
ENV TCL_TCZ_X86 $TCL_BASE/x86/tcz

ENV PHOTON_RPM_REPO http://bonneville.eng.vmware.com:8080/job/bonneville-kernel/lastSuccessfulBuild/artifact/binary/
ENV PHOTON_KERNEL_RPM linux-esx-4.2.0-10.x86_64.rpm

ENV ROOTFS /rootfs
ENV BOOTFS /bootfs/boot
ENV KERNEL_TREE /linux_tree
ENV ISOOUT stdio:/dev/fd/1

# Extract the kernel rpm contents
RUN mkdir -p $KERNEL_TREE && cd $KERNEL_TREE && \
	curl -L $PHOTON_RPM_REPO/$PHOTON_KERNEL_RPM -o /tmp/$PHOTON_KERNEL_RPM && \
	rpm2cpio /tmp/$PHOTON_KERNEL_RPM | cpio -id

# Copy the kernel to the boot folder
RUN mkdir -p $BOOTFS && cd $BOOTFS && \
	 mv $KERNEL_TREE/boot/vmlinuz-esx-4.2.0 vmlinuz64

# Install the TCZ dependencies
ENV TCZ_DEPS iproute2.tcz libtirpc.tcz

# Install the TCZ dependencies
RUN mkdir -p $ROOTFS && cd $ROOTFS && \
    curl -L $TCL_CORE/corepure64.gz | gunzip | cpio -id && \
    for dep in $TCZ_DEPS; do \
        echo "Download $TCL_TCZ/$dep" &&\
        (curl -L -o /tmp/$dep $TCL_TCZ/$dep || exit 1) && \
        (unsquashfs -f -d $ROOTFS /tmp/$dep || exit 1) && \
        rm -f /tmp/$dep ;\
    done

RUN ssh-keygen -t rsa -b 1024 -f $ROOTFS/id_rsa -N ""

# Replace the modules with what came from the photon kernel rpm.
RUN rm -rf $ROOTFS/lib/modules/ $ROOTFS/lib/firmware/ && \
    mkdir -p $ROOTFS/lib/modules/ && \
    mv $KERNEL_TREE/lib/modules/ $ROOTFS/lib/

# Copy static busybox binary into rootfs
RUN cp /bin/busybox $ROOTFS/bin/busybox.static
COPY default.script $ROOTFS/etc/udhcpc/default.script

# Location of the tether build
ENV TETHER /go/src/github.com/vmware/vic/bootstrap

RUN mkdir -p /binary

# install the target specific dependencies
COPY Godeps $TETHER/
RUN cd $TETHER && gpm install
