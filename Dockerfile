# Use golang alpine image as the builder stage
FROM golang:1.22.4-alpine3.20 AS builder

# Install git and other necessary tools
RUN apk update && apk add --no-cache git bash

# Set the Current Working Directory inside the container
WORKDIR /src

# Copy the source code into the container
COPY . .

# Fetch dependencies using go mod if your project uses Go modules
RUN go mod download

# Version and Git Commit build arguments
ARG VERSION
ARG GIT_COMMIT
ARG BUILD_DATE

# Build the Go app with versioning information
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo \
-ldflags "-X github.com/supporttools/go-sql-proxy/pkg/health.version=$VERSION -X github.com/supporttools/go-sql-proxy/pkg/health.GitCommit=$GIT_COMMIT -X github.com/supporttools/go-sql-proxy/pkg/health.BuildTime=$BUILD_DATE" \
-o /bin/go-sql-proxy

# Start from scratch for the runtime stage
FROM scratch

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the built binary and config file from the builder stage
COPY --from=builder /bin/go-sql-proxy /bin/

# Set the binary as the entrypoint of the container
ENTRYPOINT ["/bin/go-sql-proxy"]
