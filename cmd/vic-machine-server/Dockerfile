# Building:
# from vic root directory
# docker build --no-cache -t vic-machine-server -f cmd/vic-machine-server/Dockerfile .
# docker tag vic-test gcr.io/eminent-nation-87317/vic-machine-server:1.x
# gcloud auth login
# gcloud docker -- push gcr.io/eminent-nation-87317/vic-machine-server:1.x

FROM vmware/photon:1.0

ENV HOST 0.0.0.0
ENV PORT 80
ENV TLS_PORT 443

# Default location for TLS - Specify `-v /host/cert/path:/certs` to use defaults
# Override by providing a volume and values for `-e TLS_CERTIFICATE` and `-e TLS_PRIVATE_KEY`
ENV TLS_CERTIFICATE=/certs/server.cert.pem
ENV TLS_PRIVATE_KEY=/certs/server.key.pem

EXPOSE 80
EXPOSE 443

COPY bin/vic-machine-server /usr/local/bin/

ENTRYPOINT /usr/local/bin/vic-machine-server
