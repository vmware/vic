FROM container-win-base

COPY tether-core $TETHER/tether/

# Enable SSH debug
#RUN sed -i -e 's/const debugMux = false/const debugMux = true/' /go/src/golang.org/x/crypto/ssh/mux.go && sed -i -e 's/const debugHandshake = false/const debugHandshake = true/' /go/src/golang.org/x/crypto/ssh/handshake.go

COPY tether $TETHER/target/
# ensure the base deps are present
RUN cd $TETHER/tether && gpm install

# construct a script to generate the rootfs and iso - we can't do it as a command as docker gets
# confused by something and it doesn't parse as a valid json array
RUN echo "cd $TETHER/target && \
	GOARCH=amd64 GOOS=windows CGO_ENABLED=0 go build -tags netgo -o tether && \
	file tether && \
	cp -f tether /binary/win-tether.exe && \
	mv /binary /tmp/\$BINBASE" > /tmp/generate.sh && chmod a+x /tmp/generate.sh

CMD ["/tmp/generate.sh"]
