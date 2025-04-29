FROM alpine:3.21


ARG VERSION
ENV VERSION ${VERSION}

# Install certs and create a non-root user
RUN apk add --no-cache ca-certificates && \
    adduser -D -u 10001 user

COPY keeper /keeper/
RUN chmod +x /keeper

USER user

ENTRYPOINT ["/keeper/keeper", "server"]
