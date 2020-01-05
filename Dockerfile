FROM golang:1.13.5

# Set of labels recommended by the label-schema project
# Reference: http://label-schema.org/rc1/
LABEL org.label-schema.name="sysadmin-sk"
LABEL org.label-schema.description="Sysadmin Sidekick"
LABEL org.label-schema.url="https://github.com/raffs/sysadmin-sk"
LABEL org.label-schema.vcs-url="https://github.com/raffs/sysadmin-sk"
LABEL org.label-schema.version="0.0.1-alpha"
LABEL org.label-schema.docker.cmd="docker run -it sysadmin-sk --help"

# Define global env variables
ENV GOPATH "/usr/share/sysadmin-sk"
ENV PROJECT_DIR "${GOPATH}/src/github.com/raffs/sysadmin-sk"

COPY . "${PROJECT_DIR}"

WORKDIR "${PROJECT_DIR}"

# building sysadmin-sk binary from cmd/sysadmin-sk
RUN cd "${PROJECT_DIR}/cmd/sysadmin-sk" && \
    BIN_OUTPUT=/usr/bin/sysadmin-sk "${PROJECT_DIR}/build/build.sh" && \
    rm -rf "${GOPATH}"

ENTRYPOINT ["/usr/bin/sysadmin-sk"]