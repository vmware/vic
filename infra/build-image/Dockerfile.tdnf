# This is an image used for building bootstrap.iso and application.iso with photon 2.0 operate system.
#
# To use this image:
# docker run -v $(pwd):/go/src/github.com/vmware/vic gcr.io/eminent-nation-87317/vic-build-image:tdnf make most
#
# To build this image:
# docker build -t vic-build-image-tdnf -f infra/build-image/Dockerfile.tdnf infra/build-image/
# docker tag vic-build-image-tdnf gcr.io/eminent-nation-87317/vic-build-image:tdnf
# gcloud auth login
# gcloud docker -- push gcr.io/eminent-nation-87317/vic-build-image:tdnf

FROM photon:2.0

ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH $PATH:${GOPATH}/bin:/${GOROOT}/bin
ENV SRCDIR ${GOPATH}/src/github.com/vmware/vic
ENV TERM linux
ENV VIC_CACHE_DEPS=1

WORKDIR ${SRCDIR}

# gawk added purely for tolower function used in Makefile
# rpm used for initialize_bundle - tdnf seems to require rpm -initdb to work    
# kmod is needed for depmod command
# Go compilation step requires the following
#   sed
#   gcc
#   binutils
#   glibc-devel
#   linux-api-headers - like kernel-headers but named this
# findutils - for xargs for go-deps script (also in toybox if not using coreutils)
# cpio - this ensures that we don't need to change the vic scripts based on toybox or not-toybox
# it seems that erasing toybox automatically installs coreutils
# which is used to check for other utilities
# jq is used for parse repo-spec.json
# git is used for retrieving tag and commit information
# tar gzip xz are used for compression
# make for Makefile
# xorriso is used for building iso file
RUN tdnf erase -y toybox && \
    tdnf install -y jq cpio git tar gzip make xz gawk rpm kmod sed which \
    gcc binutils glibc-devel linux-api-headers findutils xorriso

ADD https://dl.google.com/go/go1.8.6.linux-amd64.tar.gz /tmp/go.tgz
RUN cd /usr/local && tar -zxf /tmp/go.tgz 

RUN mkdir -p ${SRCDIR}

COPY setup-repo.sh /usr/local/bin
RUN chmod a+x /usr/local/bin/setup-repo.sh

ENTRYPOINT [ "/usr/local/bin/setup-repo.sh" ]