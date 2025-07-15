# Dockerfile vars
TAG=0.0.0

# vars
IMAGENAME=go-sql-proxy
REPO=docker.io/supporttools
IMAGEFULLNAME=${REPO}/${IMAGENAME}:${TAG}

.PHONY: help test build push bump all install-tools lint security deps validate ci

help:
	@echo "Makefile arguments:"
	@echo ""
	@echo "tag - Docker Tag"
	@echo ""
	@echo "Makefile commands:"
	@echo "install-tools - Install required static analysis tools"
	@echo "lint         - Run all linting tools (golangci-lint, staticcheck, go vet, deadcode)"
	@echo "test         - Run tests with race detection"
	@echo "security     - Run security scanning with gosec"
	@echo "deps         - Verify and tidy dependencies"
	@echo "validate     - Run all validation steps (lint, test, security, deps)"
	@echo "ci           - Run full CI pipeline locally (install-tools + validate)"
	@echo "build        - Build the Docker image"
	@echo "push         - Push the Docker image to the repository"
	@echo "bump         - Build and push a new image"
	@echo "all          - Run tests, build, and push"

.DEFAULT_GOAL := all

# Install required tools
install-tools:
	@echo "Installing static analysis tools..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.62.2
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@go install github.com/psampaz/go-mod-outdated@latest
	@go install github.com/remyoudompheng/go-misc/deadcode@latest
	@echo "Tools installed successfully"

# Run all linting
lint:
	@echo "Running linting tools..."
	@golangci-lint run ./...
	@staticcheck ./...
	@go vet ./...
	@deadcode .
	@echo "All linting passed!"

# Run tests
test:
	@echo "Running tests with race detection..."
	@go test -v -race ./...
	@echo "Tests completed!"

# Run security scanning
security:
	@echo "Running security scan with gosec..."
	@gosec ./...
	@echo "Security scan completed!"

# Dependency management
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@echo "Verifying dependencies..."
	@go mod verify
	@echo "Tidying dependencies..."
	@go mod tidy -v
	@echo "Checking for uncommitted changes..."
	@git diff --exit-code go.mod go.sum || (echo "ERROR: go.mod or go.sum was modified by go mod tidy - please run 'go mod tidy' locally and commit the changes" && exit 1)
	@echo "Dependencies verified!"

# Run all validation steps (mirrors CI pipeline)
validate: lint test security deps
	@echo "All validation steps passed!"

# Full CI simulation
ci: install-tools validate
	@echo "CI pipeline simulation completed successfully!"

build:
	@docker buildx build --platform linux/amd64 --pull \
		--build-arg GIT_COMMIT=`git rev-parse HEAD` \
		--build-arg VERSION=${TAG} \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		-t ${IMAGEFULLNAME} . --push

push:
	@docker push ${IMAGEFULLNAME}
	@docker tag ${IMAGEFULLNAME} ${REPO}/${IMAGENAME}:latest
	@docker push ${REPO}/${IMAGENAME}:latest

bump: build push

all: validate build push
