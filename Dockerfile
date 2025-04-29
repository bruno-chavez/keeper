FROM golang:1.24-alpine as builder

WORKDIR /keeper

# Creates non root user
RUN adduser -D -u 10001 user

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux GOAMD64=v2 go build -ldflags="-s -w" -o keeper cmd/keeper/main.go && \
    apk add --no-cache ca-certificates && \
    update-ca-certificates

FROM scratch

ARG VERSION
ENV VERSION ${VERSION}

# Non root user info
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Certs for making https requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /keeper/keeper /keeper

# Running as keeper
USER 10001

ENTRYPOINT ["/keeper"]