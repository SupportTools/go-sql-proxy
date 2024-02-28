# Use a full-featured base image for building
FROM golang:1.21.6-alpine3.18 AS builder

# Install git if your project requires
RUN apk update && apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /src

# Copy the source code into the container
COPY . .

# Fetch dependencies using go mod if your project uses Go modules
RUN go mod download

# Version and Git Commit build arguments
ARG VERSION=
ARG GIT_COMMIT
ARG BUILD_DATE

# Build the Go app with versioning information
RUN GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/supporttools/go-sql-proxy/pkg/health.version=$VERSION -X github.com/supporttools/go-sql-proxy/pkg/health.GitCommit=$GIT_COMMIT -X github.com/supporttools/go-sql-proxy/pkg/health.BuildTime=$BUILD_DATE" -o /bin/go-sql-proxy

# Start from scratch for the runtime stage
FROM scratch

# Set the working directory to /app
WORKDIR /app

# Copy the built binary and config file from the builder stage
COPY --from=builder /bin/go-sql-proxy /app/

# Set the binary as the entrypoint of the container
ENTRYPOINT ["/app/go-sql-proxy"]
