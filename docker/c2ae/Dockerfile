#syntax = docker/dockerfile:1.0-experimental

FROM golang:1.12 as buildenv
RUN mkdir -p ~/.ssh && echo "gitlab.com,35.231.145.151 ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBFSMqzJeV9rUzU4kWitGjeR4PWSa29SPqJ1fVkhtj3Hw9xjLVXVYrU9QlYWrOLXBpQ6KWjbjTDTdDkoohFzgbEY=" >> ~/.ssh/known_hosts

# Fix gitlab protocal to use ssh instead of https
RUN git config --global url."git@gitlab.com:".insteadOf "https://gitlab.com/"

RUN mkdir -p /build/
COPY go.mod /build/
COPY go.sum /build/

# put go mod downlaod in its own layer, this speed up the build by A LOT :)
RUN --mount=type=ssh cd /build/ && go mod download

COPY . /build/

RUN cd /build/ &&  ls -lah scripts/ && ./scripts/build.sh


# c2ae-api stage
FROM debian:stable AS c2ae-api
# Alpine doesn't work even with libc6-compat

COPY --from=buildenv /build/bin/c2ae-api /opt/e4/bin/c2ae-api
COPY --from=buildenv /build/configs/ /opt/e4/configs

# Redirect log file to stdout
RUN ln -s /dev/stdout /var/log/e4_c2ae.log

WORKDIR /opt/e4

# Default GRPC
EXPOSE 5556

ENTRYPOINT /opt/e4/bin/c2ae-api

# c2ae-cli stage
FROM debian:stable AS c2ae-cli
# Alpine doesn't work even with libc6-compat

COPY --from=buildenv /build/bin/c2ae-cli /opt/e4/bin/c2ae-cli
COPY ./docker/c2ae/cli-entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

WORKDIR /opt/e4

ENV C2AE_API_ENDPOINT="localhost:5556"

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
