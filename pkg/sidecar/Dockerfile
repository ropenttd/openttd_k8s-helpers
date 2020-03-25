# BUILD ENVIRONMENT

FROM golang:alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/github.com/ropenttd/docker_openttd-bans-sidecar/pkg/sidecar
COPY . .

# Fetch dependencies
RUN go get -d -v

# Build the binary
RUN go build -o /go/bin/banrunner

# END BUILD ENVIRONMENT
# DEPLOY ENVIRONMENT

FROM scratch
MAINTAINER duck. <me@duck.me.uk>

# Copy the executable
COPY --from=builder /go/bin/sidecar /usr/local/bin/sidecar

VOLUME /config

# Finally, let's run the script!
USER openttd
CMD /usr/local/bin/banrunner
