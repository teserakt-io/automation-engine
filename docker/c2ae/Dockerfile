###########################################################
#
# c2ae-api stage
#
###########################################################
FROM alpine AS c2ae-api

# Redirect log file to stdout
RUN ln -s /dev/stdout /var/log/e4_c2ae.log

WORKDIR /opt/e4

# Default GRPC
EXPOSE 5556
# Default HTTP
EXPOSE 8886

VOLUME /opt/e4/configs

ENTRYPOINT /opt/e4/bin/c2ae-api

COPY ./bin/c2ae-api /opt/e4/bin/c2ae-api
COPY ./configs/ /opt/e4/configs

###########################################################
#
# c2ae-cli stage
#
###########################################################
FROM alpine AS c2ae-cli

COPY ./docker/c2ae/cli-entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

WORKDIR /opt/e4

VOLUME /opt/c2ae/c2ae-cert.pem

ENV C2AE_API_ENDPOINT="localhost:5556"
ENV C2AE_API_CERT="/opt/c2ae/c2ae-cert.pem"

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]

COPY ./bin/c2ae-cli /opt/e4/bin/c2ae-cli
