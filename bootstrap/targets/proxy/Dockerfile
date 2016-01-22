FROM container-proxy-base

# copy in the local base code
COPY tether-core $TETHER/tether/

# install generic tether dependencies if any
RUN cd $TETHER/tether && gpm install

COPY tether $TETHER/target/

CMD ["/bin/bash", "-c", "cd target && GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -a -x -o /binary/proxy-tether && mv /binary /tmp/$BINBASE"]
