FROM container-proxy-base

COPY tether $TETHER/tether

CMD ["/bin/bash", "-c", "cd tether/cmd/proxy/ && GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -a -x -o /binary/proxy-tether && mv /binary /tmp/$BINBASE"]
