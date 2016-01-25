FROM photon

# https://github.com/vmware/photon/blob/e3b980c13d1ed93665724cefb7c3f417f8d63e5c/support/package-builder/constants.py
RUN tdnf install -y linux-api-headers ; \
    tdnf install -y glibc ; \
    tdnf install -y glibc-devel ; \
    tdnf install -y zlib ; \
    tdnf install -y zlib-devel ; \
    tdnf install -y file ; \
    tdnf install -y binutils ; \
    tdnf install -y binutils-devel ; \
    tdnf install -y gmp ; \
    tdnf install -y gmp-devel ; \
    tdnf install -y mpfr ; \
    tdnf install -y mpfr-devel ; \
    tdnf install -y mpc ; \
    tdnf install -y libgcc ; \
    tdnf install -y libgcc-devel ; \
    tdnf install -y libstdc++ ; \
    tdnf install -y libstdc++-devel ; \
    tdnf install -y libgomp ; \
    tdnf install -y libgomp-devel ; \
    tdnf install -y gcc ; \
    tdnf install -y pkg-config ; \
    tdnf install -y ncurses ; \
    tdnf install -y bash ; \
    tdnf install -y bzip2 ; \
    tdnf install -y sed ; \
    tdnf install -y procps-ng ; \
    tdnf install -y coreutils ; \
    tdnf install -y m4 ; \
    tdnf install -y grep ; \
    tdnf install -y readline ; \
    tdnf install -y diffutils ; \
    tdnf install -y gawk ; \
    tdnf install -y findutils ; \
    tdnf install -y gettext ; \
    tdnf install -y gzip ; \
    tdnf install -y make ; \
    tdnf install -y patch ; \
    tdnf install -y util-linux ; \
    tdnf install -y tar ; \
    tdnf install -y xz ; \
    tdnf install -y libtool ; \
    tdnf install -y flex ; \
    tdnf install -y bison ; \
    tdnf install -y lua ; \
    tdnf install -y popt ; \
    tdnf install -y nspr ; \
    tdnf install -y sqlite-autoconf ; \
    tdnf install -y nss ; \
    tdnf install -y elfutils-libelf ; \
    tdnf install -y libpipeline ; \
    tdnf install -y gdbm ; \
    tdnf install -y perl ; \
    tdnf install -y texinfo ; \
    tdnf install -y rpm ; \
    tdnf install -y rpm-build ; \
    tdnf install -y rpm-devel ; \
    tdnf install -y autoconf ; \
    tdnf install -y automake ; \
    tdnf install -y groff ; \
    tdnf install -y man-db ; \
    tdnf install -y man-pages ; \
    tdnf install -y elfutils ; \
    tdnf install -y cpio

# linux.spec BuildRequires
RUN tdnf install -y bc ; \
    tdnf install -y kbd ; \
    tdnf install -y kmod ; \
    tdnf install -y glib-devel ; \
    tdnf install -y xerces-c-devel ; \
    tdnf install -y xml-security-c-devel ; \
    tdnf install -y libdnet ; \
    tdnf install -y libmspack ; \
    tdnf install -y Linux-PAM ; \
    tdnf install -y openssl-devel ; \
    tdnf install -y procps-ng-devel

RUN tdnf install -y git

RUN git clone https://github.com/vmware/photon.git --branch dev \
        && (cd /photon && git checkout -q 7e3fa1f59f7a804dce985efdcab0dac0e6af7783)

ADD https://dl.bintray.com/vmware/photon_sources/1.0/linux-4.2.tar.xz /usr/src/photon/SOURCES/

COPY Makefile /photon/SPECS/linux/

RUN echo "cd /photon/SPECS/linux/ && make srpm && make rpm" > /tmp/generate.sh && chmod 755 /tmp/generate.sh

CMD ["/bin/bash", "/tmp/generate.sh"]
