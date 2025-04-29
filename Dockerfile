FROM alpine:3.21


ARG VERSION
ENV VERSION ${VERSION}

# Install CA certificates
RUN apk add --no-cache ca-certificates

# Creates non root user
RUN adduser -D -u 10001 user

COPY keeper /keeper/
RUN chmod +x /keeper

USER user

ENTRYPOINT ["/keeper/keeper", "server"]
