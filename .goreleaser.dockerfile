FROM busybox:stable

WORKDIR /keeper

ARG VERSION
ENV VERSION ${VERSION}

# Creates non root user
ENV USER=user
ENV UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

COPY keeper keeper

# Running as keeper
USER user:user

ENTRYPOINT ["/keeper"]
