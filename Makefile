# Copyright (c) 2024 FajarLaksono. All Rights Reserved.

# use .env file
-include .env

# use .env.local file
-include .env.local

SERVICE_TAG=laksono-api-service
SERVICE_PATH=fajarlaksono.github.io/$(SERVICE_TAG)
BUILDER_IMAGE=golang:1.19.3-alpine

export REVISION_ID?=unknown
export BUILD_DATE?=unknown
export GIT_HASH?=unknown

RUN=docker run --rm \
	-v $(CURDIR):/opt/go/src/$(SVC) \
	-v $(GOPATH)/pkg/mod:/opt/go/pkg/mod \
	-w /opt/go/src/$(SVC) \
	-e GO111MODULE=on

# Building Images
build: 
# building the service
	$(RUN) -e CGO_ENABLE=0 -e GOOS=linux $(BUILDER_IMAGE) \
		go build -buildvcs=false -o service \
		-ldflags "-s -X main.revisionID=$(REVISION_ID) -X main.buildDate=$(BUILD_DATE) -X main.gitHash=$(GIT_HASH)" \
		./cmd/service-httpserver/...

# building the worker
	$(RUN) -e CGO_ENABLED=0 -e GOOS=linux $(BUILDER_IMAGE) \
		go build -buildvcs=false -o worker \
		-ldflags "-s -X main.revisionID=$(REVISION_ID) -X main.buildDate=$(BUILD_DATE) -X main.gitHash=$(GIT_HASH)" \
		./cmd/service-worker/...

# building the websocket
	$(RUN) -e CGO_ENABLED=0 -e GOOS=linux $(BUILDER_IMAGE) \
		go build -buildvcs=false -o websocket \
		-ldflags "-s -X main.revisionID=$(REVISION_ID) -X main.buildDate=$(BUILD_DATE) -X main.gitHash=$(GIT_HASH)" \
		./cmd/service-websocket/...

# building docker images
	docker build --tag="$(SERVICE_TAG):$(REVISION_ID)" --tag="$(SERVICE_TAG):latest" .
	docker build --file=Dockerfile.worker --tag="$(SERVICE_TAG)-worker:$(REVISION_ID)" --tag="$(SERVICE_TAG)-worker:latest" .
	docker build --file=Dockerfile.websocket --tag="$(SERVICE_TAG)-websocket:$(REVISION_ID)" --tag="$(SERVICE_TAG)-websocket:latest" .

# Running
run:
	docker-compose -f docker-compose.yaml up

# Debug start
debug-start:   
	$(RUN) -e CGO_ENABLED=0 -e GOOS=linux $(BUILDER_IMAGE) go build -buildvcs=false -o service \
			-gcflags "all=-N -1" ./cmd/service-httpserver/...
	
	$(RUN) -e CGO_ENABLED=0 -e GOOS=linux $(BUILDER_IMAGE) go build -buildvcs=false -o worker \
			-gcflags "all=-N -1" ./cmd/service-worker/...

	$(RUN) -e CGO_ENABLED=0 -e GOOS=linux $(BUILDER_IMAGE) go build -buildvcs=false -o websocket \
			-gcflags "all=-N -1" ./cmd/service-websocket/...

	docker-compose -f docker-compose-debug.yaml up

# Debug Stop
debug-stop:
	docker-compose down
	docker-compose -f docker-compose-debug.yaml down

# Clean the builds 
clean: 
	-$(RUN) $(BUILDER_IMAGE) rm -rf vendor vendor.* *.log service service.sha256 worker worker.sha256 test.xml dependencies.txt websocket
	-$(RUN) $(BUILDER_IMAGE) go clean -i
	-$(RUN) $(BUILDER_IMAGE) find . -type f -name 'coverage.xml' -delete
	-docker rmi -f $(SERVICE_TAG):$(REVISION_ID) $(SERVICE_TAG):latest $(BUILDER_IMAGE)
	-docker-compose --no-ansi -f docker-compose-test.yaml down
	-docker-compose rm -f -s -v

COMPOSE_TEST=docker-compose --no-ansi -p $(SERVICE)-test-$(TEST_ID) -f docker-compose-test.yaml
COMPOSE_TEST_CLEAN_UP=$(COMPOSE_TEST) down --remove-orphans

# Run automated test
test:
	$(RUN) $(BUILDER_IMAGE) go mod vendor
	$(COMPOSE_TEST) up -d
	$(COMPOSE_TEST) run -e CGO_ENABLED=0 -e PWD=$(CURDIR) test bash coverage.sh \
	 || ($(COMPOSE_TEST_CLEAN_UP);exit 1)
	$(COMPOSE_TEST_CLEAN_UP)

# Rebuilding Images
rebuild:
	docker-compose --no-ansi stop service

	$(RUN) -e CGO_ENABLED=0 -e GOOS=linux $(BUILDER_IMAGE) go build -buildvcs=false -o service \
		-ldflags "-s -X main.revisionID=$(REVISION_ID) -X main.buildDate=$(BUILD_DATE) -X main.gitHash=$(GIT_HASH)" \
		./cmd/service-httpserver/...

	$(RUN) -e CGO_ENABLED=0 -e GOOS=linux $(BUILDER_IMAGE) go build -buildvcs=false -o worker \
		-ldflags "-s -X main.revisionID=$(REVISION_ID) -X main.buildDate=$(BUILD_DATE) -X main.gitHash=$(GIT_HASH)" \
		./cmd/service-worker/...

	$(RUN) -e CGO_ENABLED=0 -e GOOS=linux $(BUILDER_IMAGE) go build -buildvcs=false -o websocket \
		-ldflags "-s -X main.revisionID=$(REVISION_ID) -X main.buildDate=$(BUILD_DATE) -X main.gitHash=$(GIT_HASH)" \
		./cmd/service-websocket/...

	docker-compose --no-ansi up service



