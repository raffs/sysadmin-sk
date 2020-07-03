# syntax = docker/dockerfile:1.0-experimental
FROM golang:1.13.5

# Set of labels recommended by the label-schema project
# Reference: http://label-schema.org/rc1/
LABEL org.label-schema.name="sysadmin-sk"
LABEL org.label-schema.description="Sysadmin Sidekick"
LABEL org.label-schema.url="https://github.com/raffs/sysadmin-sk"
LABEL org.label-schema.vcs-url="https://github.com/raffs/sysadmin-sk"
LABEL org.label-schema.version="0.0.1-alpha"
LABEL org.label-schema.docker.cmd="docker run -it sysadmin-sk --help"

ENV GOPATH "/usr/share/sysadmin-sk"

COPY . "/usr/share/sysadmin-sk/src/github.com/raffs/sysadmin-sk"
WORKDIR "/usr/share/sysadmin-sk/src/github.com/raffs/sysadmin-sk"

RUN export BIN_OUTPUT=/usr/local/bin/sysadmin-sk && \
    bash build/build.sh && \
    rm -rf "${GOPATH}"

ENTRYPOINT ["/usr/local/bin/sysadmin-sk"]
