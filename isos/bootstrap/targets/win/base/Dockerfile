FROM golang:1.6.0

RUN apt-get update && apt-get -y install git gcc make vim file
RUN cd /tmp && git clone https://github.com/pote/gpm.git && cd gpm && ./configure && make install && rm -fr /tmp/gpm

# Prepare tether and init
ENV TETHER /go/src/github.com/vmware/vic/bootstrap
# this needs to be set for gpm
ENV GOOS windows

#RUN git clone https://ghicken@enatai-gerrit.eng.vmware.com/bonneville-container $TETHER && cd $TETHER && git checkout proxy
COPY Godeps $TETHER/
RUN cat $TETHER/Godeps
RUN cd $TETHER && gpm install

RUN mkdir -p /binary

# so we can dump easily
ENTRYPOINT ["/bin/bash"]

