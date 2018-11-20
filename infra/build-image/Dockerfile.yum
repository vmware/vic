# This is an image used for building bootstrap.iso and application.iso with centos 6.9 operate system.
#
# To use this image:
# docker run -v $(pwd):/go/src/github.com/vmware/vic gcr.io/eminent-nation-87317/vic-build-image:yum make most
#
# To build this image:
# docker build -t vic-build-image-yum -f infra/build-image/Dockerfile.yum infra/build-image/
# docker tag vic-build-image-yum gcr.io/eminent-nation-87317/vic-build-image:yum
# gcloud auth login
# gcloud docker -- push gcr.io/eminent-nation-87317/vic-build-image:yum

FROM centos:7

ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH $PATH:${GOPATH}/bin:/${GOROOT}/bin
ENV SRCDIR ${GOPATH}/src/github.com/vmware/vic
ENV TERM linux
ENV VIC_CACHE_DEPS=1

WORKDIR ${SRCDIR}

# rpm used for initialize_bundle - yum seems to require rpm -initdb to work
# Go compilation step requires the following
#   sed
#   gcc
#   binutils
#   glibc-devel
#   glibc-static
# cpio - this ensures that we don't need to change the vic scripts based on toybox or not-toybox
# it seems that erasing toybox automatically installs coreutils
# jq is used for parse repo-spec.json
# git is used for retrieving tag and commit information
# tar gzip xz are used for compression
# make for Makefile
# xorriso is used for building iso file
# epel-release allow yum to install packages and dependencies
RUN set -ex; \
    yum install -y epel-release; \
    yum install -y \
    which jq cpio git tar gzip make rpm sed gcc \
    binutils glibc-devel xorriso glibc-static;

ADD https://dl.google.com/go/go1.8.6.linux-amd64.tar.gz /tmp/go.tgz
RUN cd /usr/local && tar -zxf /tmp/go.tgz 

RUN mkdir -p ${SRCDIR}

COPY setup-repo.sh /usr/local/bin
RUN chmod a+x /usr/local/bin/setup-repo.sh

ENTRYPOINT [ "/usr/local/bin/setup-repo.sh" ]
