# BUILD ENVIRONMENT

FROM golang:alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/github.com/ropenttd/docker_openttd-bans-sidecar

COPY ../../go.mod .
COPY ../../go.sum .

# Get dependencies - will also be cached if we won't change mod/sum
RUN go mod download

# Then copy the rest of this source code
COPY ../../. .

# And build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/sidecar ./cmd/sidecar

# END BUILD ENVIRONMENT
# DEPLOY ENVIRONMENT

FROM scratch
MAINTAINER duck. <me@duck.me.uk>

# Copy the executable
COPY --from=builder /go/bin/sidecar /usr/local/bin/sidecar

VOLUME /config
WORKDIR /config

# Finally, let's run the script!
ENTRYPOINT ["/usr/local/bin/sidecar"]
