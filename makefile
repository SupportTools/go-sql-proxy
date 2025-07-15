# Dockerfile vars
TAG=0.0.0

# vars
IMAGENAME=go-sql-proxy
REPO=docker.io/supporttools
IMAGEFULLNAME=${REPO}/${IMAGENAME}:${TAG}

.PHONY: help test build push bump all

help:
	@echo "Makefile arguments:"
	@echo ""
	@echo "tag - Docker Tag"
	@echo ""
	@echo "Makefile commands:"
	@echo "test     - Run tests and static analysis"
	@echo "build    - Build the Docker image"
	@echo "push     - Push the Docker image to the repository"
	@echo "bump     - Build and push a new image"
	@echo "all      - Run tests, build, and push"

.DEFAULT_GOAL := all

test:
	@echo "Running tests and static analysis..."
	golint ./... && \
	staticcheck ./... && \
	go vet ./... && \
	go mod tidy && \
	go mod verify && \
	gosec ./... && \
	deadcode ./... && \
	go fmt ./...

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

all: test build push
