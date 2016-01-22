FROM golang:1.5.3

# install GPM
#RUN apt-get update && apt-get -y install gcc make
RUN cd /tmp && git clone https://github.com/pote/gpm.git && cd gpm && ./configure && make install && rm -fr /tmp/gpm

# Prepare tether and init
ENV TETHER /go/src/enatai-gerrit.eng.vmware.com/bonneville-container/
# seems early docker versions are broken with WORKDIR/ENV
WORKDIR /go/src/enatai-gerrit.eng.vmware.com/bonneville-container/

RUN mkdir -p target /binary $TETHER/tether

# install the target specific dependencies
COPY Godeps $TETHER/
RUN cd $TETHER && gpm install

