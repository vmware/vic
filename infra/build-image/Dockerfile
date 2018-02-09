# docker build -t vic-build-image -f infra/build-image/Dockerfile .
# docker tag vic-build-image gcr.io/eminent-nation-87317/vic-build-image:1.x
# gcloud auth login
# gcloud docker -- push gcr.io/eminent-nation-87317/vic-build-image:1.x

FROM vmware/photon:2.0

ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH $PATH:${GOPATH}/bin:/${GOROOT}/bin
ENV SRCDIR ${GOPATH}/src/github.com/vmware/vic
WORKDIR ${SRCDIR}

# gawk added purely for tolower function used in Makefile
# rpm used for initialize_bundle - tdnf seems to require rpm -initdb to work    
# kmod is needed for depmod command
# need the following for configure
#   expr for xorriso configure which toybox doesn't have
#   sed
#   gcc
#       binutils
#       glibc-devel
#       linux-api-headers - like kernel-headers but named this
# findutils - for xargs for go-deps script (also in toybox if not using coreutils)
# cpio - this ensures that we don't need to change the vic scripts based on toybox or not-toybox
# it seems that erasing toybox automatically installs coreutils
RUN tdnf erase -y toybox && tdnf install -y jq cpio git tar gzip make xz gawk rpm kmod sed gcc binutils glibc-devel linux-api-headers findutils

ADD https://dl.google.com/go/go1.8.7.linux-amd64.tar.gz /tmp/go.tgz
RUN cd /usr/local && tar -zxf /tmp/go.tgz 

ADD https://www.gnu.org/software/xorriso/xorriso-1.4.8.tar.gz /usr/local/src/xorriso.tgz
RUN cd /usr/local/src && tar -zxf xorriso.tgz
RUN cd /usr/local/src/xorriso-1.4.8 && ./configure && make && make install && make distclean

RUN mkdir -p ${SRCDIR}

COPY infra/build-image/setup-repo.sh /usr/local/bin
RUN chmod a+x /usr/local/bin/setup-repo.sh

# ENTRYPOINT [ "/usr/local/bin/setup-repo.sh" ]