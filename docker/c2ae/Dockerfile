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

RUN mkdir -p /opt/e4/db /opt/e4/configs

# Redirect log file to stdout
RUN ln -s /dev/stdout /var/log/e4_c2ae.log

WORKDIR /opt/e4

# Default GRPC
EXPOSE 5556

ENV C2AE_C2_ENDPOINT="localhost:5555"
ENV C2AE_C2CERT_PATH="/opt/e4/configs/c2-cert.pem"
ENV C2AE_DB_PATH="/opt/e4/db/c2ae.db"
ENV C2AE_LISTEN_ADDR="0.0.0.0:5556"
ENV C2AE_OC_ENDPOINT="localhost:55678"

ENTRYPOINT /opt/e4/bin/c2ae-api -c2 ${C2AE_C2_ENDPOINT} -c2cert ${C2AE_C2CERT_PATH} -db ${C2AE_DB_PATH} -ocSampleAll -addr ${C2AE_LISTEN_ADDR} -ocAddr ${C2AE_OC_ENDPOINT}

# c2ae-cli stage
FROM debian:stable AS c2ae-cli
# Alpine doesn't work even with libc6-compat

COPY --from=buildenv /build/bin/c2ae-cli /opt/e4/bin/c2ae-cli
COPY ./docker/c2ae/cli-entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

WORKDIR /opt/e4

ENV C2AE_API_ENDPOINT="localhost:5556"

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
