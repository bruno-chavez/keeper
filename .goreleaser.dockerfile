FROM busybox:stable


ARG VERSION
ENV VERSION ${VERSION}

RUN mkdir -p /etc/ssl/certs && \
    wget -qO /etc/ssl/certs/ca-certificates.crt https://curl.se/ca/cacert.pem

# Creates non root user
RUN adduser -D -u 10001 user
USER 10001

COPY keeper /keeper/
RUN chmod +x /keeper/keeper

ENTRYPOINT ["keeper"]
